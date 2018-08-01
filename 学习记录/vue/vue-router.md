# vue-router

[参考链接](https://github.com/DDFE/DDFE-blog/issues/9)

和主要流程相关的代码主要是 `src` 目录下的 `components`, `history` 目录和 `create-matcher.js`, `create-route-map.js`, `inex.js` 和 `install.js`。

## 引用入口

一般在 Vue 项目的 `router.js` 或者 `app.js` 会引用并定义 vue-router，如下：

```javascript
import Vue from 'vue'
import VueRouter from 'vue-router'

// 1. 插件
// 安装 <router-view> and <router-link> 组件
// 且给当前应用下所有的组件都注入 $router and $route 对象
Vue.use(VueRouter)

// 2. 定义各个路由下使用的组件，简称路由组件
const Home = { template: '<div>home</div>' }
const Foo = { template: '<div>foo</div>' }
const Bar = { template: '<div>bar</div>' }

// 3. 创建 VueRouter 实例 router
const router = new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { path: '/', component: Home },
    { path: '/foo', component: Foo },
    { path: '/bar', component: Bar }
  ]
})

// 4. 创建 启动应用
// 一定要确认注入了 router
// 在 <router-view> 中将会渲染路由组件
new Vue({
  router,
  template: `
    <div id="app">
      <h1>Basic</h1>
      <ul>
        <li><router-link to="/">/</router-link></li>
        <li><router-link to="/foo">/foo</router-link></li>
        <li><router-link to="/bar">/bar</router-link></li>
        <router-link tag="li" to="/bar">/bar</router-link>
      </ul>
      <router-view class="view"></router-view>
    </div>
  `
}).$mount('#app')
```

## 作为插件

Vue.use(plugin) 是 Vue 插件的基本写法。这个机制会调用插件的 `install` 方法，在 `src/index.js` 中暴露：

```javascript
// src/index.js

// 赋值 install
VueRouter.install = install

// 自动使用插件
if (inBrowser && window.Vue) {
  window.Vue.use(VueRouter)
}
```

在 `src/install.js` 中定义。`install.js` 主要做了：

1. `export let _Vue` 暴露一个 Vue 的引用。

    插件在打包的时候是肯定不希望把 vue 作为一个依赖包打进去的，但是呢又希望使用 `Vue` 对象本身的一些方法，此时就可以采用上边类似的做法，在 install 的时候把参数赋值给 `Vue` ，这样就可以在其他地方使用 `Vue` 的一些方法而不必引入 vue 依赖包（前提是保证 install 后才会使用）。

2. `export function install (Vue)` 暴露 `install` 方法。

    1. `_Vue = Vue`。即处理问题 1。
    2. `Object.defineProperty(Vue.prototype, '$router', {get () { return this.$root._router }})` 来定义让组件都可以使用 `$router` 和 `$route` （`$route` 代码省略了）

    在 Vue.js 中所有的组件都是被扩展的 Vue 实例，也就意味着所有的组件都可以访问到这个实例原型上定义的属性。

    3. `Vue.mixin({ beforeCreate() {xxxx}})` 来注册 `beforeCreate` 钩子。在钩子中调用了 `this._router.init(this)` 来初始化，并且调用 `Vue.util.defineReactive(this, '_route', this._router.history.current)` 来定义响应式的 `_route` 对象。（把 `_router` 定义成响应式表示只要这个值发生变化，就会触发整个 render 过程）
    4. `Vue.component('router-view', View)` 来注册 `router-view` 和 `router-link`。
    5. 定义 router 钩子函数的合并策略为追加。主要用于 Vue.extend 和 mixin 时的同名 key 的合并策略。

    ```javascript
    const strats = Vue.config.optionMergeStrategies
    // use the same hook merging strategy for route hooks
    strats.beforeRouteEnter = strats.beforeRouteLeave = strats.beforeRouteUpdate = strats.created
    ```

## 实例化 VueRouter

```javascript
const router = new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { path: '/', component: Home },
    { path: '/foo', component: Foo },
    { path: '/bar', component: Bar }
  ]
})
```

除了 `install` 方法之外，这里讲一下 VueRouter 的构造函数，在 `src/index.js` 中。

1. 调用

  ```javascript
  this.matcher = createMatcher(options.routes || [])
  ```

  这个方法定义在 `src/create-matcher.js` 中。

  1. 首先调用 `createRouteMap(routes)` 方法，根据传入的 `routes` 配置生成对应的路由 map。这个方法在 `src/create-route-map.js` 中定义。

    在 `createRouteMap(routes)` 中，扫描构造函数参数的 `routes` 数组，拿出每一项创建一个 `RouteRecord` 对象。（这个对象还会存 parent 以保存嵌套路由关系）接下来更新两个数组：`pathMap` 用以存放 path -> record 的映射关系；`nameMap` 用以存放 name -> record 的映射关系。这里的 path 和 name 都是构造函数时 `routes` 数组的元素的属性，也就是路由配置项。

  2. 返回一个对象包含两个方法，分别是 `match` 和 `addRoutes`。（之后就可以使用 `this.matcher.match` 了）

2. 实例化 History。

  VueRouter 把几种路由模式 (hash, history, abstract) 分别做成了类，实现了一个基类 History `src/history/base.js`。根据传入的 `mode` 进行分别的实例化，并赋值到 `this.history`。

  ```javascript
  switch (mode) {
    case 'history':
      this.history = new HTML5History(this, options.base)
      break
    case 'hash':
      this.history = new HashHistory(this, options.base, this.fallback)
      break
    case 'abstract':
      this.history = new AbstractHistory(this, options.base)
      break
    default:
      if (process.env.NODE_ENV !== 'production') {
        assert(false, `invalid mode: ${mode}`)
      }
  }
  ```

## 实例化 Vue

在入口处的最后一步，创建了 VueRouter 实例之后就开始创建 Vue 实例，把 `router` 当做对象的一个 key 传入，如下：

```javascript
new Vue({
  router,
  template: `
    <div id="app">
      <h1>Basic</h1>
      <ul>
        <li><router-link to="/">/</router-link></li>
        <li><router-link to="/foo">/foo</router-link></li>
        <li><router-link to="/bar">/bar</router-link></li>
        <router-link tag="li" to="/bar">/bar</router-link>
      </ul>
      <router-view class="view"></router-view>
    </div>
  `
}).$mount('#app')
```

在 install 方法中，有一段调用 `Vue.mixin` 进行全局混入的代码，如下：

```javascript
Vue.mixin({
  beforeCreate () {
    // 判断是否有 router
    if (this.$options.router) {
      this._routerRoot = this
      // 赋值 _router
      this._router = this.$options.router
      // 初始化 init
      this._router.init(this)
      // 定义响应式的 _route 对象
      Vue.util.defineReactive(this, '_route', this._router.history.current)
    }
  }
})

```

因此每次在创建 vue 实例的时候，都会进入这个钩子。其中 `this._router` 指向之前初始化的 VueRouter 实例。接下来还要为 Vue 的原型链添加 `$route` 和 `$router` 对象，方便实例使用，如下：

```javascript
Object.defineProperty(Vue.prototype, '$router', {
  get () { return this._routerRoot._router }
})

Object.defineProperty(Vue.prototype, '$route', {
  get () { return this._routerRoot._route }
})
```

所以实际上 `this._routerRoot` 就是创建的那个 Vue 实例。而 `this._routerRoot._router` 就是 mixin 代码中的 `this._router`，也就是 RouterVue 实例；`this._routerRoot._route` 就是 mixin 代码中的定义成响应式的部分，值就是 `this._router.history.current`。

下面会分别来看这两者。

### router.init

在 `src/index.js` 中，VueRouter 的构造函数之后，就定义了 `init(app)` 方法，这里的 `app` 就是 Vue 实例。这里主要的工作有：

1. `this.app = app`，记录当前实例。（当前版本会记录多个 app，即 `this.apps.push(app)`）
2. 把构造函数中定义的 `this.history` 取出，针对 `HTML5History` 和 `HashHistory` 进行一些特殊处理。（主要是调用 `transitionTo` 方法实现跳转到当前 URL 上标示的路由，而不是每次都打开首页）

  > 因为在这两种模式下才有可能存在进入时候的不是默认页，需要根据当前浏览器地址栏里的 path 或者 hash 来激活对应的路由，此时就是通过调用 transitionTo 来达到目的；而且此时还有个注意点是针对于 HashHistory 有特殊处理，为什么不直接在初始化 HashHistory 的时候监听 hashchange 事件呢？这个是为了修复 [vuejs/vue-router#725](https://github.com/vuejs/vue-router/issues/725) 这个 bug 而这样做的。简要来说就是说如果在 beforeEnter 这样的钩子函数中是异步的话，beforeEnter 钩子就会被触发两次，原因是因为在初始化的时候如果此时的 hash 值不是以 / 开头的话就会补上 #/，这个过程会触发 hashchange 事件，所以会再走一次生命周期钩子，也就意味着会再次调用 beforeEnter 钩子函数。

  展开一下这里的 `transitionTo` 方法，定义在 `src/history/base.js` 中，作用是跳转到某个路由页面。它的步骤是：

  1. 调用 `this.router.match` 方法。在 VueRouter 的构造函数中曾经使用 `this.history = new History(this, options.base` 来实例化 History。而在 History 的构造函数 `constructor(router, base)` 中又设置了 `this.router = router`。因此实际上 `this.router` 指向 VueRouter 的实例。根据 VueRouter 类的定义，`match` 方法内部调用了 `this.matcher.match`，参数为目标路由 `location` 和当前路由 `this.current`，返回目标路由对象。（如果没有找到则创建）

    查找时，根据传入的 `location` 上有 `name` 还是 `path`，分别从之前创建的 2 个 MAP 中寻找。不论最终找到与否，都会调用 `_createRoute` 方法进行创建路由对象并返回。创建的代码位于 `src/util/route.js`，其中值得关注的是每一个路由 route 对象都对应有一个 matched 属性，它是一个数组。如果能够匹配到 map 中的某个 RouteRecord，那么数组会把当前 record 添加到最前端，然后继续访问 record 的父亲，直到根。因此最终这个数组是**从根开始**到最终匹配的 record 的数组；如果没有匹配到，那这里就是个空数组。

  2. 调用 `confirmTransition`。先检查目标和当前是否相同，相同则结束。之后检查两者的内部结构以确认哪些组件需要更新，哪些不需要。

  3. 在 `confirmTransition` 的回调中调用 `this.updateRoute`，作用是调用目标路由组件上的各类钩子。

3. `history.listen(route => {this.app._route = route})`

  当前版本会操作多个 app，即

  ```javascript
  history.listen(route => {
    this.apps.forEach(app => {
      app._route = route
    });
  });
  ```

  `history.listen(cb)` 的作用是注册 `this.cb = cb`。在每次更新路由 `this.updateRoute` 时会调用 `this.cb && this.cb(route)`。这里的作用就是更新当前实例的 `_route` 属性。作用在下面进行分析。

### 响应式的 _route

```javascript
Vue.util.defineReactive(this, '_route', this._router.history.current)
```

首先从 `src/history/base.js` 中看到，`this.current` 是在 `this.updateRoute` 时每次被设定更新的，永远指向当前路由对象。给 `_route` 定义了这么一个响应式的属性值也就意味着如果该属性值发生了变化，就会触发更新机制，继而调用应用实例的 render 重新渲染。而让这个值发生变化的，就是上面最后注册的 `history.listen` 中的 `cb`，也就是那句 `app._route = route`。

## router-link 和 router-view

这是 vue-router 提供的两个自定义组件。分别来看。

### router-view

组件代码定义在 `src/components/view.js` 中。

1. 它是一个函数式组件。重点关注渲染函数 `render (h, { props, children, parent, data })`
2. 设置 `data.routerView = true` 并观察 `parent.$vnode.data.routerView` 来确定嵌套路由，从而确定深度，存放在 `data.routerViewDepth`
3. 使用 `parent.$route.matched[depth]` 来获取当前的路由记录(RouteRecord)对象。如果没有获取到对象，则直接调用 `h()` 绘制空节点。
4. 获取 `props.name` 作为要渲染视图名称（默认 `'default'`），并从 `component = matched.components[name]` 中获取 Vue 组件。
5. 非 keepAlive 模式下，每次都设置组件的全部生命周期。最后调用 `return h(component, data, children)` 执行渲染。

### router-link

这个代码在 `src/components/link.js` 中。

1. 调用 `const { location, route, href } = router.resolve(this.to, current, this.append)` 获取目标路由的信息。
2. 根据当前的路由信息以及配置 (exact) 给当前链接设置 activeClass
3. 注册点击事件，根据配置 (replace) 决定使用 `router.push` 还是 `router.replace`。这里还要判断一些别的事件以免误触发点击，比如是否按了 ctrl 之类的。
4. 根据配置项 (tag) 决定创建 a 标签还是从儿子中寻找 a 标签（找不到就还是把事件注册在自己的节点上）
5. 调用 `return h(this.tag, data, this.$slots.default)` 执行渲染。

## 总结

总结下来，vue-router 的流程大致是这样

install 方法，执行于 `Vue.use(VueRouter)`

1. 注册两个组件 `router-view`, `router-link`。
2. 注册了 `beforeCreated` 钩子。
  1. 执行 VueRouter 实例的 `init` 方法进行初始化。
  2. 定义 `app._route` 为响应式，一旦发生变化则触发 app 重新渲染。
3. 给 `Vue.prototype` 挂载 `$router` 和 `$route` 两个对象，方便其他组件使用。两者分别指向 VueRouter 的实例以及当前路由对象。

VueRouter 实例的构造函数，执行于 `const router = new VueRouter({routes, mode, base});`

4. 创建 matcher，暴露 `match`, `addRoutes` 两个方法。
5. 根据不同的 `mode` 实例化 History

VueRouter 实例的 `init` 方法，执行于创建 Vue 实例时由 `beforeCreated` 钩子调用

6. `this.app = app` 记录当前 Vue 实例
7. 调用 `history.transitionTo` 直接跳转到目标页面，而不是统一进首页。（因为默认设置了 START 为 `'/'`）注意 `transitionTo` 的内部会调用 `updateRoute`。
8. `history.listen` 注册回调，在每次 `updateRoute` 时修改 `app._route`，因为响应式从而触发 app 重新渲染

`router-link` 部分

9. 渲染时根据当前路由信息决定是否设置 activeClass。
10. 分析 `to` 中包含的目标路由信息。
11. 渲染为 DOM 标签并绑定点击事件。

`router-view` 部分

12. 渲染时根据当前路由信息，找到配好的 Vue components。如果找不到则渲染空白；找到了就设置一些钩子函数，然后通过调用渲染函数，让 Vue 执行它的渲染即可。

实际切换页面时（也包括第一次默认的 `history.transitionTo`）

13. 点击时根据配置使用 `router.push` 或者 `router.replace`。这两个方法最终都调用了 `history.transitionTo`。所以无论如何切换页面的入口都在 `history.transitionTo`，内部调用 `updateRoute`，执行注册在 `listen` 的回调，即更新 `app._route`。又因为响应式，进入了重新渲染流程，导致 `router-view` 重新渲染，页面切换完成。

