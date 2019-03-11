# Vue SSR

## 目录结构及用途

一个 Vue SSR 项目的大致目录结构是

```
src
├── components
│   ├── Foo.vue
│   ├── Bar.vue
│   └── Baz.vue
├── App.vue
├── app.js # 通用 entry(universal entry)
├── entry-client.js # 仅运行于浏览器
└── entry-server.js # 仅运行于服务器
```

### app.js

用于暴露 `createApp()` 方法，从而供两个 entry 调用并创建新的 Vue 实例。在 SSR 模式下必须在每个请求都创建新的 Vue 实例，而不是复用同一个。

如果有使用 vue-router 或者 vuex，也会在这里引用他们，并当做参数传入 Vue 的构造函数。

另外如果有使用 vue-router 或者 vuex，还会在同级存在一个单独的文件暴露 `createRouter()` 或者 `createStore()` 方法，供这里调用。

```javascript
// app.js
import Vue from 'vue'
import App from './App.vue'
import { createRouter } from './router'

export function createApp () {
  // 创建 router 实例
  const router = createRouter()

  const app = new Vue({
    // 注入 router 到根 Vue 实例
    router,
    render: h => h(App)
  })

  // 返回 app 和 router
  return { app, router }
}
```

### entry-client.js

entry-client 负责调用 app.js 中暴露的 `createApp()` 方法，并挂载到容器节点中。其他客户端的特有逻辑也在这里。

entry-client 构建后会成为 client-bundle.js，给浏览器使用。

```javascript
import { createApp } from './app'

// 客户端特定引导逻辑……

const { app } = createApp()

// 这里假定 App.vue 模板中根元素具有 `id="app"`
app.$mount('#app')
```

### entry-server.js

entry-server 也负责调用 app.js 中暴露的 `createApp()` 方法。在使用 vue-router 的情况下，它还需要负责注册初始路由，以及匹配路由的工作。它暴露一个 `context => resolve(app)` 的方法，通常还是异步的（因为 `router.onReady()` 是异步的）。这里的 `context` 就是 express 的上下文，通常包含 url, path 等属性。

entry-server 构建后会成为 server-bundle.js，供服务端的 express 代码使用。因此它的输入就是 express 中间件的输入（context)，经过一系列操作后，把 `app` 通过 Promise 返回给 express。express 的后续操作参见下面的 server.js。

```javascript
// entry-server.js
import { createApp } from './app'

export default context => {
  // 因为有可能会是异步路由钩子函数或组件，所以我们将返回一个 Promise，
    // 以便服务器能够等待所有的内容在渲染前，
    // 就已经准备就绪。
  return new Promise((resolve, reject) => {
    const { app, router } = createApp()

    // 设置服务器端 router 的位置
    router.push(context.url)

    // 等到 router 将可能的异步组件和钩子函数解析完
    router.onReady(() => {
      const matchedComponents = router.getMatchedComponents()
      // 匹配不到的路由，执行 reject 函数，并返回 404
      if (!matchedComponents.length) {
        return reject({ code: 404 })
      }

      // Promise 应该 resolve 应用程序实例，以便它可以渲染
      resolve(app)
    }, reject)
  })
}
```

### server.js

server.js 的作用就是使用 server-bundle.js （由 entry-server 而来)，拿到 app 后调用 `renderer.renderToString()` 方法，最终不论是使用 template (`<!--vue-ssr-outlet-->`) 还是直接 `res.end`，都需要返回出去。以 express 为例：

```javascript
// server.js
const createApp = require('/path/to/built-server-bundle.js')

server.get('*', (req, res) => {
  const context = { url: req.url }

  createApp(context).then(app => {
    renderer.renderToString(app, (err, html) => {
      if (err) {
        if (err.code === 404) {
          res.status(404).end('Page not found')
        } else {
          res.status(500).end('Internal Server Error')
        }
      } else {
        res.end(html)
      }
    })
  })
})
```

## 数据预取

SSR 需要把首屏的渲染结果直接返回，因此需要在服务端进行数据预取，并把预取的数据加模板，渲染成页面。而对客户端来说，因为有一个**激活**的工作，因此需要客户端的初始状态和服务端保持一致才行，因此客户端也需要预取数据。

### 引入 vuex

涉及到数据，就会使用 vuex。因此我们添加一个 store.js，和 router 以及 app 一样，返回一个 `createStore()` 方法而不是 `store` 实例。在方法中，初始化一个实例，并且包含获取数据的方法。

这里定义了 state, actions 和 mutations，还缺少一个地方 `dispatch` 来调用 actions。这个调用就在页面组件中。

```javascript
// store.js
import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

// 假定我们有一个可以返回 Promise 的
// 通用 API（请忽略此 API 具体实现细节）
import { fetchItem } from './api'

export function createStore () {
  return new Vuex.Store({
    state: {
      items: {}
    },
    actions: {
      fetchItem ({ commit }, id) {
        // `store.dispatch()` 会返回 Promise，
        // 以便我们能够知道数据在何时更新
        return fetchItem(id).then(item => {
          commit('setItem', { id, item })
        })
      }
    },
    mutations: {
      setItem (state, { id, item }) {
        Vue.set(state.items, id, item)
      }
    }
  })
}
```

在 app.js 还需要引入这个 store.js，创建 store 实例并传递给 Vue 构造函数，这里省略了。

### 路由组件（页面组件）

即一个组件占据一个页面，在 routers 配置中 component 的值。在它们的 script 部分，暴露一个 `asyncData` 方法，调用刚才定义的 actions。

注意：这个方法执行在 vue 实例化之前，因此没有 `this`。取而代之的是通过参数直接传入 `store` 和 `route` 对象，方便调用 `dispatch` 方法和获取当前路由状态（参数）。

```html
<!-- Item.vue -->
<template>
  <div>{{ item.title }}</div>
</template>

<script>
export default {
  asyncData ({ store, route }) {
    // 触发 action 后，会返回 Promise
    return store.dispatch('fetchItem', route.params.id)
  },
  computed: {
    // 从 store 的 state 对象中的获取 item。
    item () {
      return this.$store.state.items[this.$route.params.id]
    }
  }
}
</script>
```

### 服务端预取

在 entry-server.js 中，获得匹配的路由组件 `matchedComponents` 之后，查看组件是否有 `asyncData` 方法。如有则调用。

在这些 `asyncData` 调用之后，需要最终调用一下 `context.state = store.state`，把状态挂到 context 上，方便之后输出到 `window.__INITIAL_STATE__` 上面。

```javascript
import { createApp } from './app'

export default context => {
  return new Promise((resolve, reject) => {
    const { app, router, store } = createApp()

    router.push(context.url)

    router.onReady(() => {
      const matchedComponents = router.getMatchedComponents()
      if (!matchedComponents.length) {
        return reject({ code: 404 })
      }

      // 对所有匹配的路由组件调用 `asyncData()`
      Promise.all(matchedComponents.map(Component => {
        if (Component.asyncData) {
          return Component.asyncData({
            store,
            route: router.currentRoute
          })
        }
      })).then(() => {
        // 在所有预取钩子(preFetch hook) resolve 后，
        // 我们的 store 现在已经填充入渲染应用程序所需的状态。
        // 当我们将状态附加到上下文，
        // 并且 `template` 选项用于 renderer 时，
        // 状态将自动序列化为 `window.__INITIAL_STATE__`，并注入 HTML。
        context.state = store.state

        resolve(app)
      }).catch(reject)
    }, reject)
  })
}
```

在这种模式下，客户端（浏览器端）代码应该读取这个初始状态，并设置到前端的 store 中去，因此在 entry-client.js 中，也会看到一句代码：

```javascript
const { app, router, store } = createApp()

if (window.__INITIAL_STATE__) {
  // 读取服务端传递过来的 store 状态，写入到前端的 store 中去
  store.replaceState(window.__INITIAL_STATE__)
}
```

### 客户端预取

客户端预取有两种方式：

#### 先获取数据，再导航路由

使用此策略，应用程序会等待视图所需数据全部解析之后，再传入数据并处理当前视图。好处在于，可以直接在数据准备就绪时，传入视图渲染完整内容，但是如果数据预取需要很长时间，用户在当前视图会感受到"明显卡顿"。因此，如果使用此策略，建议提供一个数据加载指示器 (data loading indicator)。

我们可以通过检查匹配的组件，并在全局路由钩子函数中执行 asyncData 函数，来在客户端实现此策略。注意，在初始路由准备就绪之后，我们应该注册此钩子，这样我们就不必再次获取服务器提取的数据。

lavas 采用的是这个方案，为此它还注册了全局的 ProgressBar。

```javascript
// entry-client.js

// ...忽略无关代码

router.onReady(() => {
  // 添加路由钩子函数，用于处理 asyncData.
  // 在初始路由 resolve 后执行，
  // 以便我们不会二次预取(double-fetch)已有的数据。
  // 使用 `router.beforeResolve()`，以便确保所有异步组件都 resolve。
  router.beforeResolve((to, from, next) => {
    const matched = router.getMatchedComponents(to)
    const prevMatched = router.getMatchedComponents(from)

    // 我们只关心非预渲染的组件
    // 所以我们对比它们，找出两个匹配列表的差异组件
    let diffed = false
    const activated = matched.filter((c, i) => {
      return diffed || (diffed = (prevMatched[i] !== c))
    })

    if (!activated.length) {
      return next()
    }

    // 这里如果有加载指示器 (loading indicator)，就触发

    Promise.all(activated.map(c => {
      if (c.asyncData) {
        return c.asyncData({ store, route: to })
      }
    })).then(() => {

      // 停止加载指示器(loading indicator)

      next()
    }).catch(next)
  })

  app.$mount('#app')
})
```

#### 先匹配并导航路由，再获取数据

此策略将客户端数据预取逻辑，放在视图组件的 `beforeMount` 函数中。当路由导航被触发时，可以立即切换视图，因此应用程序具有更快的响应速度。然而，传入视图在渲染时不会有完整的可用数据。因此，对于使用此策略的每个视图组件，都需要具有条件加载状态。

这可以通过纯客户端 (client-only) 的全局 mixin 来实现：

```javascript
// 写在 entry-client.js 中
Vue.mixin({
  beforeMount () {
    const { asyncData } = this.$options
    if (asyncData) {
      // 将获取数据操作分配给 promise
      // 以便在组件中，我们可以在数据准备就绪后
      // 通过运行 `this.dataPromise.then(...)` 来执行其他任务
      this.dataPromise = asyncData({
        store: this.$store,
        route: this.$route
      })
    }
  }
})
```

#### 处理组件重用

无论选择哪种策略，当路由组件重用（同一路由，但是 params 或 query 已更改，例如，从 user/1 到 user/2）时，也应该调用 asyncData 函数。我们也可以通过纯客户端 (client-only) 的全局 mixin 来处理这个问题：

```javascript
// 也在 entry-client.js 中
Vue.mixin({
  beforeRouteUpdate (to, from, next) {
    const { asyncData } = this.$options
    if (asyncData) {
      asyncData({
        store: this.$store,
        route: to
      }).then(next).catch(next)
    } else {
      next()
    }
  }
})
```

## 客户端激活

SSR 模式下，服务端返回的是渲染好的 HTML 结果。但是这些 HTML 都是静态的，和普通的 Vue SPA 应用的前端节点完全不同（没有响应式，没有 vdom 等等）。因此在 SSR 模式下，前端代码要负责把这些静态 HTML 转化为动态节点。（不能删掉重新渲染）

前端代码判断是需要做激活操作还是普通的（丢弃）渲染操作，取决于容器上是否存在 `data-server-rendered` 属性。此外，也可以使用 `app.$mount('#app', true)` 的第二个参数来强制激活。

在开发模式下，前端还会额外检测客户端激活的 vdom 和服务端返回的 dom 是否匹配，如不匹配将打印一条警告，并退出激活模式，而是丢弃所有服务端 DOM 重新渲染。在生产模式下，这个检测会被跳过。