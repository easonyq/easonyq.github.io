# React 细节记录

## defaultProps

使用静态成员 `defaultProps` 可以设置 react 组件属性的默认值。

```ts
interface IMyComponentProps {
    name?: string
}

// 使用类的写法
class MyComponent extends React.Component<IMyComponentProps, {}> {
    static defaultProps: IMyComponentProps = {
        name: 'eason'
    }

    render() {
        return (
            <div>Hello {this.props.name}!</div>
        )
    }
}

// 使用方法的写法
function MyComponent(props: IMyComponentProps) {
    return (
        <div>Hello {this.props.name}!</div>
    )
}
MyComponent.defaultProps = {
    name: 'eason'
}
```

## 高阶组件

High Order Component (HOC)，它是一个函数，输入参数（至少）是一个组件，返回一个新的组件。因为 react 中函数也可以作为组件（只要返回 JSX 即可)，所以这个函数本身也是组件，可以直接在 JSX 中使用，称高阶组件。举例如下：

```js
// HOC 一般需要把自己不需要的 props 透传给子组件
const HOC = Wrapped => (props => {
    <div>
        <span>HOC</span>
        <Wrapped {...props} />
    </div>
})

const MyComponent = props => <div>Hello!</div>;
const MyHOC = HOC(MyComponent);

class App extends React.Component {
    render() {
        return (
            <div className="app">
                <MyHOC/>
            </div>
        )
    }
}
```

高阶组件是实现组件复用的两种常用方式之一（另一种是渲染属性）。react-router 的 `withRouter`，redux 的 `connect` 都是通过高阶组件来实现的。

## 渲染属性 Render Props

渲染属性是 react 实现组件复用的另外一种常用方式。父组件使用一个函数作为 prop 传递到子组件，子组件把自己的而状态当做参数传入这个 prop，从而实现子组件的状态被父组件使用，但渲染逻辑依然由父组件决定。

```js
class Child extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            name: 'child'
        };
    }
    render() {
        return (
            <div>
                <span>This is Child Component</span>
                {this.props.render(this.state)}
            </div>
        )
    }
}

class Parent extends React.Component {
    render() {
        return (
            <div>
                <Child render={state => (
                    <Hello name={state} />
                )}/>
            </div>
        )
    }
}

const Hello = props => (
    <div>
        <span>Hello! {props.name}</span>
    </div>
);
```

Child 是一个可复用的组件。Parent 使用 Child 时，传入了名为 render 的属性，值为一个函数。这个函数接受一个对象为参数，返回一段 JSX（使用了 Hello 组件）。当 Child 渲染时，它调用了父组件传过来的 render 方法，以自己的状态 state 为参数。所以这时父组件的 render 函数的参数其实就是子组件的 state。

从结果上说，Child 的状态被父组件 Parent 使用了，但状态的使用（渲染）方式依然由父组件决定。所以这种模式下，渲染和数据分离，只有数据能被复用。

**注意：**如果 Child 继承自 React.PureComponent，每次 Parent 渲染都会重新生成 render 属性（渲染方法）。这个方法在浅比较中会被判断为不相等，所以 PureComponent 的作用会消失，每次都会重新渲染。

要解决这个问题，把 render 属性变成实例方法：

```js
class Parent extends React.Component {
    renderChild(state) {
        return <Hello name={state}>;
    }

    render() {
        return (
            <div>
                <Child render={this.renderChild}/>
            </div>
        );
    }
}
```

这个模式不一定非要使用名为 `render` 的属性。通常情况会使用 `children` 属性，因为这个属性可以不显式的写在属性列表中，而是直接写成子元素。

```js
// Parent render()
<Child>
    {state => (
        <Hello name={state} />
    )}
</Child>

// Child render()
<div>
    {this.props.children(this.state)}
</div>
```

## React.PureComponent

```js
class MyComponent extends React.PureComponent {
    render() {
        // xxx
    }
}
```

从继承 `Component` 改为 `PureComponent` 可以让组件成为“纯组件”，这样可以减少不必要的 render 操作次数，提升性能，也可以省略 `shouldCoponentUpdate` 函数。常用于实现纯展示组件。

具体来说，针对纯组件的情况下，react 决定是否需要重绘组件的判断标准是它的 props 和 state 是否发生了改变。但注意这里只做一次**浅比较**，也就是只比较：

1.  `Object.keys(state | props)` 的长度是否一致
2. 每个 `key` 是不是两者都有，并且是否是同一个引用

所以深层嵌套数据是不参与比较的，这也是纯组件性能高的原因。当然这也要求我们在操作纯组件的状态改变时需要谨慎。例如保持引用不变而修改深层数据（例如 array.pop），react 不会察觉它有变化，因此不会重绘，界面也就不会更新。所以**必须使用不可变数据**，即每次创建一个新的 state 才行。

另外，纯组件也可以有 `shouldComponentUpdate` 方法，会被优先调用。当不存在时才使用上述浅比较，所以两者可以共存。

## react-router withRouter

一个普通的组件，被 `withRouter(component)` 调用后，会成为一个新的组件。这个组件可以从属性中获取 `match`, `location` 和 `history`，从而进行一些路由相关的判断或操作。

```js
import React from "react";
import PropTypes from "prop-types";
import { withRouter } from "react-router";

// A simple component that shows the pathname of the current location
class ShowTheLocation extends React.Component {
  static propTypes = {
    match: PropTypes.object.isRequired,
    location: PropTypes.object.isRequired,
    history: PropTypes.object.isRequired
  };

  render() {
    const { match, location, history } = this.props;

    return <div>You are now at {location.pathname}</div>;
  }
}

// Create a new component that is "connected" (to borrow redux
// terminology) to the router.
const ShowTheLocationWithRouter = withRouter(ShowTheLocation);
```

## styled-components

使用类库 styled-components 可以快速完成两个事情：

1. 使用 HTML 标签创建一个带样式的 React 组件，这可以让 jsx 中不出现 class，比较清晰。

2. 为已有的组件套上样式，返回新的组件

另外还有几个快捷函数也比较实用：

1. createGlobalStyle 用于创建全局 CSS。常规的 `styled.div` 或者 `styled(component)` 都是互相隔离的。

2. css 创建一个CSS片段。返回值可以直接集成到其他的 styled-components 的样式代码中。

## @loadable/component

用于分割最终生成的 bundle。使用这个类库后，可以把比较大（或者独立）的 React 组件单独打包成一个 bundle，从而缩小 bundle 的尺寸。

通常还配合 webpackChunkName，例如

```js
import loadable from '@loadable/component';

const LazyDesigner = loadable(() => import(/* webpackChunkName: "designer" */ './components/Designer'), {
    fallback: <LoadingView name="设计器" />
});
```

会最终生成一个 designer.js。

## immer

通常我们需要在 react 中生成下一个新的**不可变的状态**，需要做不少克隆操作。使用 immer 可以避开这种不直观的方式，直接采用修改的写法，由类库帮助我们进行克隆等操作，生成下一个状态。

```js
import produce from "immer"

const baseState = [
    {
        todo: "Learn typescript",
        done: true
    },
    {
        todo: "Try immer",
        done: false
    }
]

const nextState = produce(baseState, draftState => {
    draftState.push({todo: "Tweet about it"})
    draftState[1].done = true
})

// baseState 并没有改变
// nextState = [
//     {
//         todo: "Learn typescript",
//         done: true
//     },
//     {
//         todo: "Try immer",
//         done: true
//     },
//     {
//         todo: "Tweet about it"
//     }
// ];
```

immer 在配合编写 redux reducer 的时候非常强大。

```js
// redux reducer
const reducer = (state, action) => {
    produce(state, draft => {
        switch (action.type) {
            case 'BIRTHDAY':
                draft.user.age += 1;
        }
    })
}

// react setState
// 在当做 setState 的参数时，produce 可以不接第一个参数 state，而直接编写 draft 的转换函数
onBirthDayClick = () => {
    this.setState(
        produce(draft => {
            draft.user.age += 1
        })
    )
}
```

immer 也可以生成 patches，供后续打补丁使用。

## react-helmet

用来设置 `<head>` 内部的标签。`<Helmet>` 内部的内容会被追加到 `<head>` 中。后设置的会覆盖之前设置的。

```ts
import {Helmet} from 'react-helmet'

let Simple = prop => (
    <div className="wrapper">
        <Helmet>
            <title>Hello</title>
            <meta charSet="utf-8"/>
            <link rel="stylesheet" href="//baidu.com/index.css">
        </Helmet>
    </div>
)
```

## getDerivedStateFromProps

从 react 16.3 开始，官网标记3个声明周期函数为*不安全的*，它们是 `componentWillMount`, `componentWillReceiveProps` 和 `componentWillUpdate`。取而代之的是 `static getDerivedStateFromProps` 和 `getSnapshotBeforeUpdate`。

如此修改的原因是，官方希望在 17 版本退出 async rendering，希望在实际DOM挂在之前，虚拟DOM构建阶段中止组件的生命周期。下次再恢复时，需要重新执行生命周期，因此这3个方法不能保证只被调用一次，因此不安全。所以添加了两个取代的生命周期（新增的两个和原来的三个不能同时使用）

static getDerivedStateFromProps:

组件每次被 render 时都会触发，是 render 函数执行之前最后一个执行的生命周期函数。根据 props 更新 state，返回值就是新的 state。如果返回 null 则表示不需要更新 state。

因为是静态方法，因此没有 `this`。如果需要上一个状态或者属性，需要额外使用一个变量记录下来，参与比较。

getSnapshotBeforeUpdate:

在 update 发生时出发，具体在 render 之后，组件 DOM 渲染之前。


