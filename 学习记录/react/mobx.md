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
import {observable} from 'mobx';
import {observer} from 'mobx-react';

class TodoList {
    @observable todos = []
    @computed get unfinishedTodoCount() {
        return this.todos.filter(todo => !todo.finished).length
    }
    @action addTodo(title) {
        this.todos.push({
            title,
            finished: false
        })
    }
}

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

// 可以在外部修改数据
setTimeout(() => {
    store.addTodo('learn mobx-react');
}, 1000);
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

### Provider & inject

mobx-react 提供了 `<Provider>` 组件来传入 store，让所有子孙组件都能够访问这个 store。子孙组件获取 store 的方式就通过 `inject`。

```js
// 在需要使用 Provider 提供属性的地方声明 inject 即可
// 如果同时使用 @inject 和 @observer，必须让 @inject 在前面
@inject("color")
@observer
class Button extends React.Component {
    render() {
        return <button style={{ background: this.props.color }}>{this.props.children}</button>
    }
}

class Message extends React.Component {
    render() {
        return (
            <div>
                {this.props.text} <Button>Delete</Button>
            </div>
        )
    }
}

class MessageList extends React.Component {
    render() {
        const children = this.props.messages.map(message => <Message text={message.text} />)
        return (
            <Provider color="red">
                <div>{children}</div>
            </Provider>
        )
    }
}
```

如果不使用 decorator 语法，也可以写作

```js
// inject 在外层，observer 在内层
// 这样能很明显地看出，inject(property)(component) 返回一个高阶组件，和 redux 的 connect 类似
// 类组件
let Button = inject('color')(observer(class Button extends React.Component {
    // render() {...}
}))

// 方法组件
let Button = inject('color')(observer(({color}) => {
    // return JSX
}))
```

一般来说 Provider 会使用在最外层的根元素。这套机制使用 React 中的 context 机制进行传递。在新版中可以使用 `React.createContext` 进行替代。

[参考1 - 如何使用 React.createContext 来解决属性层层传递](https://hackernoon.com/how-do-i-use-react-context-3eeb879169a2)

[参考2 - 官网](http://react.html.cn/docs/context.html)

### 最简版例子

三个步骤：

1. 使用 `mobx.observable(state)` 来定义状态

2. 使用 `mobx.action(newState => state = newState)` 来定义状态改变

3. 使用 `mobxReact.observer(ReactComponent)` 来定义响应式的 React 组件。在组件中调用 action 实现状态改变。

```js
// 通过 observable 定义组件的状态
const user = mobx.observable({
    name: "Jay",
     age: 22
})

// 通过 action 定义如何修改组件的状态
const changeName = mobx.action(name => user.name = name)
const changeAge = mobx.action(age => user.age = age)

// 通过 observer 定义 ReactComponent 组件。
// 也可以写作 @observer
const Hello = mobxReact.observer(class Hello extends React.Component {
    componentDidMount() {
        // 视图层通过事件触发 action
        changeName('Wang') // render Wang
    }

    render() {
        // 渲染
        console.log('render',user.name);
        return <div>Hello,{user.name}!</div>
    }
})

ReactDOM.render(<Hello />, document.getElementById('mount'));

// 非视图层事件触发，外部直接触发 action
changeName('Wang2')// render Wang2

// 重点：没有触发重新渲染
// 原因：Hello 组件并没有用到 `user.age` 这个可观察数据
changeAge('18')  // no console
```

注意：

1. 状态可以定义在组件外部，且在组件外部也可以触发状态变化，因此就实现了显示和状态的分离。

2. mobx 自行确定组件是否需要重绘（通过可观察属性），因此不需要 react 递归检查，也不需要开发者使用 `shouldComponentUpdate`，性能很高。

## mobx-state-tree

mobx 使用的是可变数据，因此可以直接改变状态，比较方便。同时回溯历史状态就不太方便。而 mobx-state-tree （简称 MST）可以解决这个问题，让 mobx 同时拥有两种好处。

```js
import { types, onSnapshot } from "mobx-state-tree"

const Todo = types
    .model("Todo", {
        title: types.string,
        done: false
    })
    .actions(self => ({
        toggle() {
            self.done = !self.done
        }
    }))

const Store = types.model("Store", {
    todos: types.array(Todo)
})

// create an instance from a snapshot
const store = Store.create({
    todos: [
        {
            title: "Get coffee"
        }
    ]
})

// listen to new snapshots
onSnapshot(store, snapshot => {
    console.dir(snapshot)
})

// invoke action that modifies the tree
store.todos[0].toggle()
// prints: `{ todos: [{ title: "Get coffee", done: true }]}`
```

MST 要求状态使用三层结构：

1. 一棵树（一个完整的状态）拥有多个 model （如例子中的 `Todo` 和 `Store`）

2. 一个 model 可以有多个节点（如例子中 `Todo` 的 `title` 和 `done`）。他们可以是数据类型(`type.string`)，也可以是具体的值（`false`）。model 的返回值可以链式调用 action，定义修改状态的方法（可变数据的特点）。

3. 定义完毕后，通过 `create` 方法，填入实际的数据，获得实例。

4. 修改状态时，可以直接调用 action，以可变数据的形式更换状态

5. 使用内置的 snapshot 功能，可以把当前状态保存下来，且以后不会再改变。这是不可变数据的特点。

[参考 - 更完整的 MST 概览](https://juejin.im/post/5c4931e451882523ea6e0c42)

### 使用 snapshot 创建 model

获取到 snapshot 之后，MST 支持从它生成新的实例，只要使用 `applySnapshot` 方法。

```js
import {applySnapshot} from 'mobx-state-tree'

// 往 store 中写入刚才获取的 snapshot，相当于恢复 store
applySnapshot(store, snapshot)
```

这样的话，回退到历史状态也可以实现

```js
import { applySnapshot, onSnapshot } from "mobx-state-tree"

var states = []
var currentFrame = -1

onSnapshot(store, snapshot => {
    if (currentFrame === states.length - 1) {
        currentFrame++
        states.push(snapshot)
    }
})

export function previousState() {
    if (currentFrame === 0) return
    currentFrame--
    applySnapshot(store, states[currentFrame])
}

export function nextState() {
    if (currentFrame === states.length - 1) return
    currentFrame++
    applySnapshot(store, states[currentFrame])
}
```

### volatile state

MST 中的所有状态都是**可持久化**的，也就是可以序列化成标准的 JSON，且类型必须和 model 声明时匹配。

但如果需要在 model 中存储无需持久化的（即获取快照等序列化操作时不计算在内的）或者数据结构类型无法预知的动态数据时可以使用 volatile state。

```js
const Todo = types.model({}).extend(self => {
    // 当 views computed value 要使用这个变量，它必须被声明为 observable
    // localState 不是一个对外暴露的属性（不在 model 里面也不在 props 里面），它是一个 volatile state，可以理解为临时变量
    const localState = observable.box(3)

    return {
        views: {
            // note this one IS a getter (computed value)
            get x() {
                return localState.get()
            }
        },
        actions: {
            setX(value) {
                localState.set(value)
            }
        }
    }
})
```

为了避免代码重复，上述例子可以简写为：

```js
const Todo = types
    .model({})
    // volatile 方法做了2个事情：定义了 localState(observable)，对外设置了 computed value 名为 localState。
    .volatile(self => ({
        localState: 3
    }))
    .actions(self => ({
        setX(value) {
            self.localState = value
        }
    }))
```

虽然外部可以访问，但 volatile state 不会被记录到 snapshot 或者 patch 中。另外对这些状态的修改也不会触发 snapshot 或者 patch。

```ts
const Model = types.model('Model')
    .volatile(self => ({
        myData: {} as any
    })
    .actions(self => ({
        run(foo: () => any) {
            foo();
        }
    }))

const model = Model.create();

autorun(() => console.log(model.myData))

model.run(() => {
    model.myData = {'name': 'eason'}
})
model.run(() => {
    model.myData.name = 'zoe'
})

// print:
// {} - 初始值
// {name: eason} - 第一次修改
// 第二次修改不打印，说明 volatile state 是 observalbe 的，但是只观察引用，不观察深度的每一个属性。
```

### snapshots

每次调用 action 后，都会生成 snapshots。

1. 通过 `getSnapshot(store, applyPostProcess)` 获取 snapshots，返回值是一个 immutable plain object，用于序列化。
2. 通过 `onSnapshot(store, callback)` 来监听变化，只要新的 snapshot 生成就会调用 callback。（mobx 事务只出发一次）
3. 可以把 snapshots 直接传入，用于恢复成状态，即 `store = Model.create(snapshot)`。
4. 通过 `applySnapshot(model, snapshot)` 来把当前 model 尽可能恢复成 snapshot 的状态。

### patches

JSON-patches 是用来描述 JSON 变动的一组数据，本身也是 JSON 格式，它的结构是

```ts
export interface IJsonPatch {
    op: "replace" | "add" | "remove"
    path: string
    value?: any
}
```

其中 path 是要修改的属性在 JSON 中的层级，例如 `/todos[1]/name`。

1. 和 snapshots 不同， 一个修改可能产生多个 patches，它不受事务的限制，而且例如 Array.splice 就是一个产生多个 patches 的例子。
2. 使用 `onPatch(model, listener)` 来监听变化。model 或者 model 的子孙节点变化都会触发，`path` 是相对于 model 所在层级的。
3. `applyPatch(model, patch)` 可以对 model 应用 patch （这里也可以传入 patch 的数组）。
4. patches 适合用来做 undo/redo

### identifiers & references

如果为一个 model 声明了 `types.identifier`，那么在 snapshot 中只需要指定这个 identifier 属性，就可以引用到实际的对象，并调用上面的方法，如下面的例子。

```js
const Todo = types.model({
    id: types.identifier,
    title: types.string
})

const TodoStore = types.model({
    todos: types.array(Todo),
    selectedTodo: types.reference(Todo)
})

// create a store with a normalized snapshot
const storeInstance = TodoStore.create({
    todos: [
        {
            id: "47",
            title: "Get coffee"
        }
    ],
    selectedTodo: "47"
})

// because `selectedTodo` is declared to be a reference, it returns the actual Todo node with the matching identifier
console.log(storeInstance.selectedTodo.title)
// prints "Get coffee"
```

注意点：

1. 每个 model 只能定义**至多一个** identifier
2. identifier 属性初始化后不可修改
3. identifier 必须全局唯一（不单单是在当前层级）
4. 可以使用 `types.refinement(types.identifier, id => /^MyApp_/.test(id))` 来限制 id 的格式
5. 使用 `type.reference` 方法时，参数类型必须有 identifier 类型的属性存在

### hooks

MST 中总共包含 6 种 hooks

1. afterCreate  在 create 方法调用之后。如果要使用到父节点的话，最好使用 afterAttach。
2. afterAttach  在当前节点挂载到父节点之后。这里可以安全地访问父节点
3. beforeDetach 在当前节点从父节点取下之前。（当前节点还没删除，例如 `detach(node)` 被调用时）
4. beforeDestroy 在当前节点被销毁之前。子节点比父节点更早触发。
5. preProcessSnapshot 在创建实例（`create(snapshot)`）或者应用快照 （`applySnapshot`） 时先进入这个 hook 进行处理，通常在这里执行数据转换和修改等。方向是 快照 -> 实例。定义这个 hook 不能在 action 中，应该使用单独的方法调用。
6. postProcessSnapshot 在创建快照之前调用，方向是 实例 -> 快照。通常和 preProcessSnapshot 的操作完全相反。这个 hook 也不能在 action 中定义。