# redux 基础

参考[官网](https://cn.redux.js.org/) 和 [阮一峰教程](http://www.ruanyifeng.com/blog/2016/09/redux_tutorial_part_one_basic_usages.html)

这里只记录一些核心知识，供快速浏览，详细的代码可以查看上面两篇文章。

## 核心设计理念

1. Web 应用是一个状态机，视图与状态是一一对应的。

2. 所有的状态，保存在一个对象里面。

## 核心概念

### store

保存数据的地方，整个应用只有一个。用它可以获取 state， 更新 state（通过 dispatch)和注册订阅(通过 subscribe)。

创建 store 的方法是调用 `createStore` 方法，它的参数是一个方法，名为 reducer，后面介绍。

### state

记录当前的状态（包含各种数据）。根据核心设计理念，state 和 view 是一一对应的。

state 本身是一个简单的对象。

### action

从状态 A 变成状态 B 只能通过发送 action 的方式来进行。一般这种状态的改变来自于用户的交互，因此 action 也是由 view 发出的（通过 store.dispatch 方法)。

action 可以理解为是一个指令。举例来说，你让一个软件播放一首特定的歌曲，你就需要告诉他【播放一首歌】以及【歌的名字】这样两个信息。

action 也是一个简单的对象，但它必须包含 `type` 字段，表明是做什么的（对应例子中的【播放一首歌】)。剩余的信息都可以自定义（对应例子中的【歌的名字】）。

#### action creator

action creator 只是大家为了方便额外制造的一个子概念，它不是 redux 的核心概念，只是约定俗成的写法而已。它的内容就是使用方法来创建 action，避免重复代码而已。

#### dispatch action

action 的角色是指令，因此需要一个发布指令的方法，这就是 `store.dispatch(action)`。

### reducer

调用 `store.dispatch` 之后，store 接到了指令，但还不知道具体要如何处理。这个处理即从状态A 如何变成状态 B（或者说状态 B 是怎样的），这就需要 reducer，也是当时 `createStore` 方法传入的参数，它本身也是个方法。

```javascript
const reducer = function (state, action) {
    // TODO
    return newState;
}
```

如果用上述的例子解释，这个函数就需要负责往 state 里面的歌单添加一首歌曲。歌曲的名字就在 action 的自定义属性里面，最终返回一个**新的** state。

reducer 不需要手动调用。当调用 `store.dispatch` 方法时，会自动被调用。

#### reducer 必须是纯函数

纯函数意味着同样的输入必须是同样的输出，所以：

1. 如随机数，当前时间等可变因素不能出现在 reducer 内部

2. 不能改写参数，主要是旧的 state

3. 不能使用系统 I/O 的 API

#### state 的默认值

reducer 还要负责设置 state 的默认值。这一般通过方法参数的默认值实现，例如 `(state = defaultState, action) => {xxx}`。

在调用 `createStore` 方法时，还可以提供第二个参数，用以指定 state 的默认值。这通常用在服务端渲染时，把服务端渲染后的状态挂在 `window` 上。因此如果有了这个值，reducer 的默认值就不生效了（相当于进入 reducer 时 state 已经有值了，所以不会再取参数默认值了）。

#### 必须返回新的 state

reducer 必须返回**新的** state，而不能仅仅在旧的 state 上处理一下就返回。因此多会用到如 `Object.assign({}, state, {change})` 之类的方法。

#### 必须兼顾未知情况

通常 reducer 里面会根据 `action.type` 来决定是什么操作，从而决定调用什么方法或者执行具体的操作。这时候如果发现是不认识的 `type` ，必须返回原 state，否则 `undefined` 就会被作为新的 state，从而发生系统错误。

#### reducer 的拆分

因为整个系统只有一个大的 state，所以需要一个大的 reducer 来处理它。但实际上不同的指令处理的是不同区块的数据，互相没有关联，因此 reducer 可以拆分成小的。这样大的 reducer 就负责根据 `type` 来分配并调用小的 reducer。这和 react 组件的嵌套关系也是一致的，所以组件和 state 是可以一一对应的。

这个过程在代码上写起来比较重复，因此 redux 提供了 `combineReducers` 方法。

通常情况，可以把所有子 reducer 放在同一个文件，使用 `export` 进行导出。这样在根 reducer 中，使用 `import * as reducers from 'xxx'`，并且 `const reducer = combineReducers(reducers)` 即可。

### store.subscribe

这个方法可以监听 state 的变化，从而执行一些代码。比较常见的是把 view 的更新（例如组件的 `render` 或者 `setState` 方法）放入监听中，就可以自动渲染。

这个方法返回一个函数，再次调用可以取消监听。

## 数据流转方向

![](http://www.ruanyifeng.com/blogimg/asset/2016/bg2016091802.jpg)

1. 用户交互后，由 view （通常是 `onClick`）构造 action，并调用 `store.dispatch(action)`

2. store 自动调用 reducer，由它返回一个新的 state

3. state 发生了变化，由 `store.subscribe` 注册的监听函数被触发

4. 通常这个监听函数会调用 `component.setState(store.getState())`，react 自动更新 view。

# 中间件和异步操作

上述基本流程中，reducer 可以立刻算出新的 state，称为同步。但如果需要发送请求等异步操作，就需要中间件的介入。

中间件的原理是在 `store.dispatch` 方法调用时，经过中间件处理后，再执行 reducer。这样 reducer 依然是同步的，也没有 I/O 操作，异步的过程在中间件中执行，对前后都是透明的。

![](http://www.ruanyifeng.com/blogimg/asset/2016/bg2016092002.jpg)

## 中间件的使用

承接常见能力的中间件基本上都已经有了，不太需要额外开发。因此只关心它的使用方法即可。主要就是 `applyMiddleware` 方法，用作 `createStore` 的第二个参数。

```javascript
import {applyMiddleware, createStore} from 'redux';
import createLogger from 'redux-logger';
const logger = createLogger();

const store = createStore(
  reducer,
  applyMiddleware(logger)
);
```

1. 如前所述，`createStore` 方法第二个参数是 state 的默认值。如果提供了这个参数，那 `applyMiddleware` 就是第三个参数。

2. `applyMiddleware` 可以接受多个中间件作为参数，以逗号间隔。这些参数中是有顺序要求的，具体要参考中间件本身的文档。

## 异步操作

同步操作只要发起一个 action，而异步操作需要发起三个，分别是操作发起时、成功时和失败时的 action。举例如下：

```javascript
// 写法一：名称相同，参数不同
{ type: 'FETCH_POSTS' }
{ type: 'FETCH_POSTS', status: 'error', error: 'Oops' }
{ type: 'FETCH_POSTS', status: 'success', response: { ... } }

// 写法二：名称不同
{ type: 'FETCH_POSTS_REQUEST' }
{ type: 'FETCH_POSTS_FAILURE', error: 'Oops' }
{ type: 'FETCH_POSTS_SUCCESS', response: { ... } }
```

异步操作时 state 也需要增加几个状态，来描述异步操作当前的情况：

```javascript
let state = {
  // ...
  isFetching: true, // 正在请求中，相当于 isLoading
  didInvalidate: true, // 数据是否已经过期
  lastUpdated: 'xxxxxxx' // 上次更新数据的时间
};
```

异步操作的流程如下：

1. 要发起异步操作时，发送一个发起操作的 action，触发 state 更新为“正在请求中”的状态（isFetching = true)，并更新视图（可能是个 Loading）

2. 异步操作返回时，根据返回结果发送操作成功或者失败的 action，出发 state 更新状态，并更新视图（可能是正常结果，也可能是错误处理）

第一个步骤和同步相同，重点在于第二个步骤，即异步操作返回时如何发送第二个 action。

### 方案一：redux-thunk

我们先考虑改造 action creator。如下的 fetchPosts 是一个 action creator，但不同于基础部分返回一个 action 对象，它返回的是一个函数。这个区别后面会讲到。

```javascript
const fetchPosts = postTitle => {
  return (dispatch, getState) => {
    dispatch(requestPosts(postTitle));
    return fetch(`/some/API/${postTitle}.json`)
      .then(response => response.json())
      .then(json => dispatch(receivePosts(postTitle, json)));
    };
  }
};
```

`fetchPosts` 返回的方法接两个参数，分别是 `dispatch` 和 `getState`。和上述操作流程一样，先发送一个发起操作的 action (`requestPosts(postTitle)`)，随后调用 `fetch` 异步操作，在返回后，发起第二个操作成功的 action （`receivePosts(postTitle, json)`)。

`store.dispatch` 方法只接受对象类型的 action。为了让这里返回函数被接受，需要一个中间件叫做 [redux-thunk](https://github.com/gaearon/redux-thunk)。经过这个中间件的强化，`store.dispatch` 就可以接受函数作为参数了。

从本质上说，当发起异步操作后，实际上还是只发送了一个 action。这第二个 action 是在中间件内部发送的，而参数 `dispatch` 和 `getState`，也是所有中间件的固定参数格式。

### 方案二：redux-promise

这个思路和上面类似，只是不同于让 dispatch 接受函数，这个方案是让它接受 Promise 对象，用到的中间件叫做 redux-promise。

```javascript
const fetchPosts = (dispatch, postTitle) => {
  return new Promise((resolve, reject) => {
    dispatch(requestPosts(postTitle));
    return fetch(`/some/API/${postTitle}.json`).then(response => {
      return {
        type: 'FETCH_POSTS',
        payload: response.json()
      };
    });
  });
});
```