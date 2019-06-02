# mobx

从宏观的角度来看，mobx和redux是平级的，都是用作整个系统状态管理的工具。它们的开发背景都源于react，但是最终都脱离了react，即可以支持其他的框架（如 Angular 和 Vue）甚至单纯的 JS 代码。

不过通常我们还是把这两者和 react 配套使用，因此它们也各自有和 react 的连接库，分别叫做 react-redux 和 mobx-react。

## 核心概念和特点

1. 任何源自应用状态的东西都应该自动地获得。（实际我还没理解这句话）

2. 常规写法使用装饰器（@observable)，因此需要环境支持 ES7。一般使用 babel-preset-mobx 或者是 typescript。

和 redux 类似，mobx 也有一些核心概念。但它们因为设计思路不同，互相没有对应关系，例如 action 这个概念两者都有，但是含义是不同的，需要注意。

下列概念的引入均脱离 react，单纯 mobx + 单纯的JS代码。配合 react 是 mobx-react 的工作。

### observable state

可观察的状态是 mobx 最基本的单元，可以类似理解为 Vue 的 data 里面的内容（也可能是 Vue 参考的 mobx）。不过 mobx 强大之处在于，它的可观察对象可以是任何东西，包括函数，类，对象，数组等等，不仅仅局限于基本类型或者可序列化类型（例如数值，字符串等）。

```javascript
import { observable } from "mobx";

class Todo {
    id = Math.random();
    @observable title = "";
    @observable finished = false;
}
```

如上代码创建了 2 个可观察状态，在之后会被使用到。

### 衍生 Derivations

任何源自状态并且不会再有任何进一步的相互作用的东西就是衍生。mobx 中衍生主要有两类

1. 计算属性

    是从 observable 中计算得到的结果（类似于 Vue 的 computed 从 data 中计算而来）。这里要求计算过程是一个纯函数，因此不能有副作用，也必须每次结果相同。

    计算属性的写法就是在 getter 之前加上 @computed，如下：

    ```javascript
    class TodoList {
        @observable todos = [];
        @computed get unfinishedTodoCount() {
            return this.todos.filter(todo => !todo.finished).length;
        }
    }
    ```

    当出现在 @computed 内部的 observable （如上代码中的 `todos`）变更后，值会被重新计算。

2. 反应 Reactions

    也从 observable 出发，但是不返回一个值（区别于计算属性），而是执行一些列带副作用的操作。例如说更新界面（把UI界面看成是 observable 的华丽展现），向后端发送新的数据等，就是反应。

    反应在 mobx 中，有三种显式的方法可供调用，分别是 autorun, when, reaction。这三者底层实现可能是相同的，只是参数略有区别。以 autorun 为例：

    ```javascript
    autorun(() => {
        console.log("Tasks left: " + todos.unfinishedTodoCount)
    })
    ```

    mobx 会检测在 autorun 函数中使用到的observalbe。当它们变化时，autorun 重新执行，内部可能也是使用观察者模式进行实现的。

    当 mobx 连接到 react 的时候，react 组件会根据状态(observable)的不同渲染成不同的样子，这个渲染工作就是 reactions 的一种。这个情况下，mobx-react 类库会自动进行这个工作，就不需要我们手动调用 autorun 等方法了，只需要在 react 组件的声明时（包括类和方法两种声明方式）设置为 observable 即可，会在之后 mobx-react 中提到。

### 动作 Actions

**注意：mobx 的 actions 定义和 redux 的 actions 是不同的。mobx 的 actions 大致等于 redux 的 reducer，是一个函数，而不是一条指令（简单对象）**

mobx 的 actions 是可选的，这是它和 redux 最大的差别之一。mobx 的状态是可以直接更改的，并不需要像 redux 的 reducer 那样每次返回新的对象，而不得修改原有状态。所以说 mobx 更加自由和宽松。

不使用 actions，直接修改状态：

```javascript
store.todos.push(
    new Todo("Get Coffee"),
    new Todo("Write simpler code")
);
store.todos[0].finished = true;
```

这些都是合法的，而且可以正常出发 reactions。不过为了代码整洁，条理清晰，也可以使用 actions，方便对状态的变化进行管理。在配置项 `enforceActions: true` 可以做出配置，要求必须通过 actions 改变状态。大型项目一般还是使用 actions。

## 数据流原则

mobx 也是单向数据流，通过动作(actions)改变状态(observable)，状态的改变衍生出视图的改变(reactions)。

## 异步 Actions

如果项目配置了必须使用 actions，而又涉及异步操作时，这里就需要注意了。

```javascript
mobx.configure({ enforceActions: true }) // 不允许在动作之外进行状态修改

class Store {
    @observable githubProjects = []
    @observable state = "pending" // "pending" / "done" / "error"

    @action
    fetchProjects() {
        this.githubProjects = []
        this.state = "pending"
        fetchGithubProjectsSomehow().then(
            projects => {
                const filteredProjects = somePreprocessing(projects)
                this.githubProjects = filteredProjects
                this.state = "done"
            },
            error => {
                this.state = "error"
            }
        )
    }
}
```

如上代码会抛出错误，因为 `@action` 修饰的仅仅是 `fetchProjects` 方法，并不包含内部的回调函数（`fetchGithubProjectsSomehow` 的 `then` 部分）。

### 把回调也变成 actions

所以原则上说，应当把这块回调函数也套在 action 里面，才能避免报错，因为这个回调也修改了状态。这就是修复方法一

```javascript
mobx.configure({ enforceActions: true })

class Store {
    @observable githubProjects = []
    @observable state = "pending" // "pending" / "done" / "error"

    @action
    fetchProjects() {
        this.githubProjects = []
        this.state = "pending"
        fetchGithubProjectsSomehow().then(this.fetchProjectsSuccess, this.fetchProjectsError)

    }

    // action.bound 是为了绑定 this
    @action.bound
    fetchProjectsSuccess(projects) {
        const filteredProjects = somePreprocessing(projects)
        this.githubProjects = filteredProjects
        this.state = "done"
    }
    @action.bound
    fetchProjectsError(error) {
        this.state = "error"
    }
}
```

或者为了代码简洁考虑，依然把这两个方法合并到 then 中，使用 `action` 方法，而不是修饰符。不过这样也要为它起个名字。

```javascript
mobx.configure({ enforceActions: true })

class Store {
    @observable githubProjects = []
    @observable state = "pending" // "pending" / "done" / "error"

    @action
    fetchProjects() {
        this.githubProjects = []
        this.state = "pending"
        fetchGithubProjectsSomehow().then(
            // 内联创建的动作
            action("fetchSuccess", projects => {
                const filteredProjects = somePreprocessing(projects)
                this.githubProjects = filteredProjects
                this.state = "done"
            }),
            // 内联创建的动作
            action("fetchError", error => {
                this.state = "error"
            })
        )
    }
}
```

因为这两个 action 只用了一次，所以还可以使用 `runInAction` 方法，它就是为了只使用一次的 action 而提供的。

```javascript
mobx.configure({ enforceActions: true })

class Store {
    @observable githubProjects = []
    @observable state = "pending" // "pending" / "done" / "error"

    @action
    fetchProjects() {
        this.githubProjects = []
        this.state = "pending"
        fetchGithubProjectsSomehow().then(
            projects => {
                const filteredProjects = somePreprocessing(projects)
                // 将‘“最终的”修改放入一个异步动作中
                runInAction(() => {
                    this.githubProjects = filteredProjects
                    this.state = "done"
                })
            },
            error => {
                // 过程的另一个结局:...
                runInAction(() => {
                    this.state = "error"
                })
            }
        )
    }
}
```

### async/await

引入了 async/await 的代码，看上去像是同步代码。但实际上它还是 Promise 的语法糖，因此在 await 之后的代码依然是独立的，不能够被 `@action` 修饰。因此在上述例子中把 `fetchGithubProjectsSomehow().then()` 改成 await 之后，`runInAction` 依然不能省略。

### flow - mobx 提供的最优解

使用 `function* ()` 代替 `async`， 使用 `yeild` 代替 `await`，就可以把整个 action 变成 flow，自动支持异步操作（也不必再使用 `@action` 修饰符）。

```javascript
mobx.configure({ enforceActions: true })

class Store {
    @observable githubProjects = []
    @observable state = "pending"

    fetchProjects = flow(function * () { // <- 注意*号，这是生成器函数！
        this.githubProjects = []
        this.state = "pending"
        try {
            const projects = yield fetchGithubProjectsSomehow() // 用 yield 代替 await
            const filteredProjects = somePreprocessing(projects)
            // 异步代码块会被自动包装成动作并修改状态
            this.state = "done"
            this.githubProjects = filteredProjects
        } catch (error) {
            this.state = "error"
        }
    })
}
```

## mobx-react

当 mobx 和 react 配套使用时，就需要引入 mobx-react 库。它提供的最大的便利之一，就是通过 `@observer` 修饰符（或者 `observer` 方法）把 react 组件套起来，这样当数据（observable)变化时，组件会自动更新，不需要手动调用 `autorun` 等方法。

```javascript
import React, {Component} from 'react';
import ReactDOM from 'react-dom';
import {observer} from 'mobx-react';

// 类声明的组件使用@observer修饰符
@observer
class TodoListView extends Component {
    render() {
        return <div>
            <ul>
                {this.props.todoList.todos.map(todo =>
                    <TodoView todo={todo} key={todo.id} />
                )}
            </ul>
            Tasks left: {this.props.todoList.unfinishedTodoCount}
        </div>
    }
}

// 方法声明的组件使用observer方法
const TodoView = observer(({todo}) =>
    <li>
        <input
            type="checkbox"
            checked={todo.finished}
            onClick={() => todo.finished = !todo.finished}
        />{todo.title}
    </li>
)

const store = new TodoList();
ReactDOM.render(<TodoListView todoList={store} />, document.getElementById('mount'));
```

在被 `observer` 包裹的组件中的 `render` 方法被使用到的变量都属于被观察的状态，只要发生变化就会由 mobx 自动触发 react 的重绘。它和普通的 react 组件相比有几个特点：

1. **性能更优秀**。传统的 react 组件内部使用 `shouldComponentUpdate` 方法对前后状态进行深度比较，从而决定是否要更新组件，这个比较开销比较大。而 mobx 的观察机制决定它能够准确得知状态是否发生了改变，因此一旦发生，它不需要进行深度比较，而是直接使用 `forceUpdate` 更新组件，性能有很大提升。

2. 组件状态在组件外部定义，不在组件之内。这样使得组件状态的维护更加清晰，兄弟组件之间复用状态也比较简单。这也是状态统一管理的思路。


### 和 redux 相比

Mobx 的优势来源于可变数据（Mutable Data）和可观察数据 (Observable Data) 。

Redux 的优势来源于不可变数据（Immutable data）。

1. 可变数据的优势：相比 redux, 组件状态的修改可以直接进行，不需要通过 action -> dispatch -> reducer 这套复杂的流程并生成一个新的状态。当然严格模式下 action 还是需要的，不过后面两部依然可以省略。

2. 不可变数据的优势：可预测性和可回溯。可以快速回到任意一个历史的状态，并且历史状态不可能被改变。

    不可变数据不一定要使用 Immutable.js 库，更重要的是一种约定。只要约定每次返回新的状态，不修改旧的，就符合了不可变数据的原则。Immutable.js 只是一种更简便的实现而已。

具体的选择可以根据业务特性，是否愿意为了可回溯的特性，牺牲代码的简便，状态改变流程的缩短。