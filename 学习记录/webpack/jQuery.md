# webpack 引用 jQuery

最近有 Lavas 开发者给我提出一个问题：他想使用 jQuery，但是 jQuery 并不是 CommonJS 规范的，他是通过暴露全局变量挂到 window 上来实现的，所以要加到 webpack 里面需要一些配置。

此类文章网上蛮多的，但都比较老，有些甚至是针对 webpack v2。这次是我自己实践操作过，针对 webpack v3 和 webpack-chain v4 亲测有效，在此纪录。

## 方法一： script 标签单独引用

这应该是最简单的做法了:

1. 在 index.html （Vue 全局的模板文件，包含 `<div id="app"></div>` 挂载点的那个）的头部添加单体的 `<script src="jquery.min.js"></script>`

2. 在 webpack 的配置中添加 `externals`

    ```javascript
    externals: {
        // key 表示 import 'jquery'，也就是包的名字
        // value 表示实际输出到代码中的全局对象的名字，即 window.jQuery
        jquery: 'jQuery'
    }
    ```
externals 的做法是建立一个假的映射关系。当发现有代码 `require('jquery')`，就返回 `jQuery` 而不是真的去 node_modules 里面找，所以就变成了读全局变量，而不会报错说找不到。

缺点也很明显，因为独立于 webpack，所以 webpack 把代码合并的时候并不理会 jquery，因此需要多发一个请求。当然这能有效的缩小 bundle 的大小，有利有弊。

## 方法二： expose-loader

略微复杂一些，但可以让 jquery 进入 webpack，一起参与优化 & 打包。

expose-loader 的作用是把指定的 js 模块 export 的变量声明为全局变量（并且支持命名）。所以这里是把 jquery export 的内容指定为全局变量 `window.jQuery` 或者 `window.$` 就完成任务了。

1. 安装 jquery 和 expose-loader （jquery 也要安装是因为要让他从 node_modules 引入，而不是单独的某个静态文件）

    ```bash
    npm i jquery expose-loader --save
    ```

2. 添加 webpack 配置，让 `jquery` 经过 `expose-loader`，让 webpack 可以认识 jquery。

    ```javascript
    // 也可以配置的时候直接写这条 rule，但核心是要写在最前面。
    config.module.rules.unshift({
      test: require.resolve('jquery'),
      // 这个 options 是表示挂载到 window 上的哪个变量。
      // 如下就是 window.$ 和 window.jQuery 都可以访问到
      use: [{
        loader: 'expose-loader',
        options: '$'
      }, {
        loader: 'expose-loader',
        options: 'jQuery'
      }]
    })
    ```

3. 也可以使用 webpack-chain 来进行这一步的配置。鉴于 webpack-chain v4 的 API 文档相当难找，摸索了很久。

    ```javascript
    config.module.rule('jquery')
      .test(require.resolve('jquery'))
      // .pre() 表示插入到最前面（进入 enforce 部分），和 unshift 效果相同
      .pre()
      .use('expose-jQuery')
          .loader('expose-loader')
          .options('jQuery')
          .end()
      .use('expose-$')
          .loader('expose-loader')
          .options('$')
          .end()
    ```

4. 在代码中进行引用。可以在某个 Vue 组件，也可以在某个 js 的头部。

    ```javascript
    // 表示要使用 window.jQuery 这种方式。
    // 不写这句的话相当于没有使用，那么代码不会被打包进去。
    import 'expose-loader?jQuery!jquery'

    // 没有 babel 环境就用 require
    require('expose-loader?jQuery!jquery')
    ```

但要特别注意：因为是 webpack 打包的，所以在开发状态这些 script 是动态添加到 `<body>` 的最后。因此如果在 index.html 里面的 `<script>` 是访问不到 `window.jQuery` 的。如果确实有类似需求，应该考虑把这个使用的代码或者类库也放到 bundle 里面来，理清依赖关系才对。

## 其他方式和插件

[这篇文章](https://array_huang.coding.me/webpack-book/chapter2/webpack-jquery-plugins.html) 还讲述了其他的内容，例如 `webpack.ProvidePlugin` 可以自动生成那句 `require('jquery')`。此外这个站点也有其他一些 webpack 的实用扫盲教程，相当不错。
