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

```javascript
import { createApp } from './app'

// 客户端特定引导逻辑……

const { app } = createApp()

// 这里假定 App.vue 模板中根元素具有 `id="app"`
app.$mount('#app')
```

### entry-server.js

entry-server 也负责调用 app.js 中暴露的 `createApp()` 方法。在使用 vue-router 的情况下，它还需要负责注册初始路由，以及匹配路由的工作。它暴露一个 `context => resolve(app)` 的方法，通常还是异步的（因为 `router.onReady()` 是异步的）。这里的 `context` 就是 express 的上下文，通常包含 url, path 等属性。

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

在编译完成后，会生成一个 server-bundle.js 和 client-bundle.js。这里先讲服务端。

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