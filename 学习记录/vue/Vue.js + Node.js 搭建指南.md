# 【新手向】Vue.js + Node.js(koa) 合体指南

## webpack 大法好

Webpack 是大家熟知的前端开发利器，它可以搭建包含热更新的开发环境，也可以生成压缩后的生产环境代码，还拥有灵活的扩展性和丰富的生态环境。但它的缺点也非常明显，那就是配置项又多又复杂，随便拿出某一个配置项（例如 `rules`， `plugins`， `devtool`等等）都够写上一篇文章来说明它的 N 种用法，对新手造成极大的困扰。Vue.js（以下简称 Vue）绝大部分情况使用 webpack 进行构建，间接地把这个问题丢给了 Vue 的新手们。不过不论是 Vue 还是 webpack，其实他们都知道配置问题的症结所在，因此他们也想了各自的办法来解决这个问题，我们先看看他们的努力。

### Vue cli 2.x - 提供开箱即用的配置

在之前一长段时间中，我们要初始化一个 Vue 项目，一般是使用 vue-cli 提供的 `vue init` 命令（这也是 Vue cli 的 v2 版本，之后简称 Vue cli 2.x)。而且通常一些比较有规模的项目都会使用 `vue init webpack my-project` 来使用 webpack 模板，那么上面提到的配置问题就来了。

为了解决这个问题， Vue 的做法是提供开箱即用的配置，即通过 `vue init` 出来的项目，默认生成的巨多的配置文件，截图如下：

![Vue cli 2 build 目录](http://boscdn.bpc.baidu.com/assets/easonyq/vue-webpack/vue-cli-2-build.png)

开箱即用是保证了，但一旦要修改，就相当于是进入了一个黑盒，开发者对于一堆文件，一堆 JSON 望洋兴叹。

### webpack 4 - 极大地简化配置

webpack 4 推出也有一年左右了，它的核心改动之一是极大地简化配置。它添加了 `mode`，把一些显而易见的配置做成内置的。因此例如 `NoEmitOnErrorsPlugin()`, `UglifyJSPlugin()` 等等都不必写了；分包用的 `CommonsChunkPlugin()` 也浓缩成了一个配置项 `optimization.splitChunks`，并且已有能适应绝大部分情况的默认值。

据说 webpack 4 构建出来的代码的体积还更小了，因此这次升级显然是必要的。

### Vue cli 3.x - 升级 webpack，还搞出了插件

大约小半年前，Vue cli 推出了 v3 版本，也是一个颠覆性的升级。它把核心精简为 `@vue/cli`，把 webpack 搞成了 `@vue/cli-service`, 把其他东西抽象为“插件”。这些插件包括 babel, eslint, Vuex, Unit Testing 等等，还允许自定义编写和发布。我不在这里介绍 Vue cli 3.x 的用法和生态，但从结果看，现在通过 `vue create` 创建的的 Vue 项目清爽了不少。

![Vue cli 3 目录](http://boscdn.bpc.baidu.com/assets/easonyq/vue-webpack/vue-cli-3.png)

### 所以现在的问题是什么？

如果我们单纯开发一个前端 Vue 项目，webpack-dev-server 能帮助我们启动一个 nodejs 服务器并支持热加载，非常好用。可如果我们要开发的是一个 nodejs + Vue 的全栈项目呢？**两者是不可能启动在同一个端口的**。那我们能做的只是让 nodejs 启动在端口 A，让 Vue (webpack-dev-server) 启动在端口 B。而如果 Vue 需要发送请求访问 nodejs 提供的 API 时，还会遇上跨域问题，虽然可以通过配置 proxy 解决，但依然非常繁琐。**而实质上，这是一整个项目的前后端而已，我们应该使用一条命令，一个端口来启动它们。**

抛开 Vue，此类需求 webpack 本身其实是支持的。因为它除了提供 webpack-dev-server 之外，还提供了 webpack-dev-middleware。它以 express middleware 的方式，同样集成了热加载的功能。因此如果我们的 nodejs 使用的是 express 作为服务框架的话，我们可以以 `app.use` 的方式引入这个中间件，就可以达成两者的融合了。

再说回 Vue cli 3。它通过 `vue-cli-service` 命令，把 webpack 和 webpack-dev-server 包裹起来，这样用户就看不到配置文件了，达成了简洁的目的。不过实质上，配置文件依然存在，只是移动到了 `node_modules/@vue/cli-service/webpack.config.js` 而已。当然为了个性化需求，它也支持用户通过配置对象 (`configureWebpack`) 或者链式调用 (`chainWebpack`) 两种间接的方式，但不再提供直接修改配置文件的方式了。

然而致命的是，即便它提供了足够的方式修改配置，**但它不能把 webpack-dev-server 变成 webpack-dev-middleware**。这表示使用 Vue cli 3 创建的 Vue 部分和 nodejs(express) 部分是不能融合的。

### 怎么解决？

说了这么多，其实这就是我最近实际碰到的问题以及分析问题的思路。鉴于 Vue cli 3 黑盒的特性，我们无法继续使用它了（可能以后有升级能解决这个问题，至少目前不行）。而使用 Vue cli 2 又因为它内置的是 webpack 3 且配置文件一大堆，也让人无所适从。这么看，唯一剩下的路就只能**自行使用并配置 webpack 4**了，这也是本文的内容所在。

## 技术栈

### nodejs 部分

目前比较主流的构建 nodejs 部分的 Web 框架是 [express](https://expressjs.com/)，且不说它的语法有多优雅，使用有多广泛等等，最主要的原因是刚才提过的 webpack-dev-middleware 就是一个 express 的中间件，因此两者可以无缝衔接。

可惜的是，在我实际的项目中，我使用了 [koa](https://koa.bootcss.com/) 作为了我的 nodejs 框架。其实要说它比 express 好在哪里我也说不上来，也不是本文的重点。可能出于尝鲜的目的，或者团队技术栈统一的目的，或者其他鬼使神差的巧合，反正我用了它，而且开始时还没意识到有这个融合的问题，直到后来发现 webpack-dev-middleware 和 koa 是不兼容的，我内心有过一丝后悔……当然这是后话了。

本文以 koa 为基准。如果您使用的是 express，其实大同小异，而且更加简单。

### Vue 部分

Vue 没什么好多说的，就一个版本，不存在 express / koa / 其他的选择。只是这里我没有使用 SSR，而是普通的 SPA 项目（单页应用，前端渲染）。

## 目录结构

既然是两个项目合体，总有一个目录结构的安排问题。这里我不谈每个项目内部需要如何组织，那是 Vue / koa 本身的问题，也是个人喜好的问题。我想谈的是这两者之间的组织方式，不外乎以下 3 种：（实际上也是个人喜好问题，见仁见智，这里只是统一一下表述，避免后续的混淆）

*以下截图中的前后端项目均为独立项目，即融合之前的，可以单独运行的那种，所以能看到两份 package.json 和 package-lock.json*

### 后端项目为基础，前端项目为子目录

![](http://boscdn.bpc.baidu.com/assets/easonyq/vue-webpack/vue-inside-koa.png)

除了红框中的 vue 目录外，其他都是 nodejs 的代码。而且因为我只是做个示意，所以 nodejs 代码其实也仅仅包含两个 index.js，public 目录和两个 package.json。实际的 nodejs 项目应该会有更多的代码，例如 actions（把每个路由处理单独到一个目录），middlewares（过所有路由的中间件）等等。

这个安排的思路是认为前端是整个项目的一部分（页面展示部分），所以 Vue 单独放在一个目录里面。**我采用的就是这种结构。**

### 前端项目为基础，后端项目为子目录

![](http://boscdn.bpc.baidu.com/assets/easonyq/vue-webpack/koa-inside-vue.png)

这就和前面一种相反，红框中的是后端代码。这么安排的理由可能是因为我们是前端开发者，所以把前端代码位于基础位置，后端提供的 API 辅助 Vue 的代码运行。

### 中立，不偏向任何人

![](http://boscdn.bpc.baidu.com/assets/easonyq/vue-webpack/vue-and-koa.png)

看了前面两种，自然能想到这第三种办法。不过我认为这种办法纯粹没事儿找事儿，因为根据 npm 的要求，package.json 是必须放在根目录的，所以实际上想把两者完全分离并公平对待是弊大于利的（例如各类调用路径都会多几层），适合强迫症患者。

## 改造 Vue 部分

Vue 部分的改造点主要是：

1. package.json 融合到根目录（nodejs) 的 package.json 里面去。这里主要包括依赖 （`dependency` 和 `devDependency`）以及执行命令（`scripts`）两部分。其余的如 `browserslist`, `engine` 等 babel 可能用到的字段，因为 nodejs 代码不需要 babel，所以可以直接复制过去，不存在融合。

2. 编写 `webpack.config.js`。（因为 Vue cli 3 是自动生成且隐藏的，这个就需要自己写）

下面详细来看。

### 融合 package.json

刚才有提到过，像 `browserslist`, `engine` 这类 babel 等使用的字段，因为 nodejs 端是不需要的，所以简单的复制过去即可。需要动脑的是依赖和命令。

依赖方面，其实前后端共用的依赖也基本不存在，所以实际上也是一个简单的复制。需要注意的是类似 `vue`, `vue-router`, `webpack`, `webpack-cli` 等等都是 `devDependency`，而不是 `dependency`。真正需要放到 `dependency` 的，其实只有 `@babel/runtime` 这一个（因为使用了 `plugin-transform-runtime`)。

命令方面，本身 Vue 必备的是“启动开发环境”和“构建”两条命令（可选的还有测试，这个我这里先不讨论）。因为开发环境需要和 nodejs 融合，所以这条我们放到 nodejs 部分说。剩下的是构建命令，常规操作是通过设置 `NODE_ENV` 为 `production` 来让 webpack 走入线上构建的情况。另外值得注意的是，因为现在 package.json 和 webpack.config.js 不在同级目录了，所以需要额外指定目录，命令如下：(`cross-env` 是一个相当好用的跨平台设置环境变量的工具)

```javascript
{
  "scripts": {
    "build": "cross-env NODE_ENV=production webpack --config vue/webpack.config.js"
  }
}
```

### 编写 webpack.config.js

*本文的重点不是 webpack 的配置方式，因此这里比较简略，不详细讲述每个配置项的含义*

webpack.config.js 本质上是一个返回 JSON 的配置文件，我们会用到其中的几个 key。如果要了解 webpack 全部的配置项，可以查看 [webpack 的中文网站](https://webpack.docschina.org/configuration)介绍。另外如果不想分段查看，你可以在[这里](https://github.com/easonyq/vue-nodejs-template/blob/master/vue/webpack.config.js)找到完整的 webpack.config.js。

#### mode

webpack 4 新增配置项，常规可选值 `'production'` 和 `'development'`。这里我们根据 `process.env.NODE_ENV` 来确定值。

```javascript
let isProd = process.env.NODE_ENV === 'production'

module.exports = {
  mode: isProd ? 'production' : 'development'
}
```

#### entry

定义 webpack 的入口。我们需要把入口设置为创建 Vue 实例的那个 JS，例如 `vue/src/main.js`。

```javascript
{
  entry: {
    "app": [path.resolve(__dirname, './src/main.js')]
  }
}
```

#### output

定义 webpack 的输出配置。在开发状态下，webpack-dev-middleware（以下简称 wdm）并不会真的去生成这个 `dist` 目录，它是通过一个内存文件系统，把文件输出到内存。所以这个目录仅仅是一个标识而已。

```javascript
{
  output: {
    filename: '[name].[hash:8].js',
    path: isProd ? resolvePath('../vue-dist') : resolvePath('dist'),
    publicPath: '/'
  }
}
```

#### resolve

主要定义两个东西：webpack 处理 `import` 时自动添加的后缀顺序和供快速访问的别名。

```javascript
{
  resolve: {
    extensions: ['.js', '.vue', '.json'],
    alias: {
      'vue$': 'vue/dist/vue.esm.js',
      '@': resolvePath('src'),
    }
  }
}
```

#### module（重点）

`module` 在 webpack 中主要确定如何处理项目中不同类型的模块。我们这里采用最常用的配法，即告诉 webpack，什么样的后缀文件用什么样的 loader 来处理。

```javascript
{
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader'
      },
      {
        test: /\.js?$/,
        loader: 'babel-loader',
        exclude: file => (
          /node_modules/.test(file) && !/\.vue\.js/.test(file)
        )
      },
      {
        test: /\.less$/,
        use: [
          isProd ? MiniCssExtractPlugin.loader : 'vue-style-loader',
          'css-loader',
          'less-loader'
        ]
      },
      {
        test: /\.css$/,
        use: [
          isProd ? MiniCssExtractPlugin.loader : 'vue-style-loader',
          'css-loader'
        ]
      }
    ]
  }
}
```

上述配置了 4 种文件的处理方式，它们分别是：

1. `/\.vue$/`

    处理 Vue 文件，使用 Vue 专门提供的 `vue-loader`。这个处理器做的事情就是**把 Vue 里面的 `<script>` 和 `<style`> 的部分独立出来，让它们可以继续由下面的 rules 分别处理**。否则一个 `.vue` 文件是不可能进入 `.js` 或者 `.css` 的处理流程的。另外如果 `<style>` 有 `lang` 属性，还可以进入例如 `.less`, `.styl` 等其他处理流程。

    它还需要专门的插件 `VueLoaderPlugin()`，之后可以看到，不要漏掉。

2. `/\.js?$/`

    表面上是处理后缀为 `.js` 的文件，但实质上在这里也用来处理 Vue 里面 `<script>` 的内容。在这里我们要做的是使用 `babel-loader` 对代码中的高级写法转译为兼容低版本浏览器的写法。具体转译规则使用 `.babelrc` 文件配置。另外这里还忽略了 `node_modules`（因为其他包在发布时都已经转码过了，不用再处理徒增时间）。不过得确保 `node_modules` 里的单体 Vue 文件依然参与转译，这也是 [Vue 官方文档](https://vue-loader.vuejs.org/zh/guide/pre-processors.html#%E6%8E%92%E9%99%A4-node-modules) 的推荐写法。

3. `/\.less$/`

    我的项目中使用 less 作为样式预处理器，因此在每个 Vue 文件都使用了 `<style lang="less">`。这样通过 `vue-loader`，就能让这条规则中配置的几个 loader 来处理 Vue 文件中的样式了。这几个 loader 分别做的事情是：

    1. `vue-style-loader`

        把样式以 `<style>` 标签的形式插入在页面头部，**只在开发状态下使用**。

        它和 `style-loader` 的差异并不大，但既然 [Vue 官方文档](https://vue-loader.vuejs.org/zh/guide/pre-processors.html)建议使用这个，那我们就用这个吧。

    2. `mini-css-extract-plugin` 的 loader

        把样式抽成一个单独的 css 文件并在 `<head>` 标签中以 `<link rel="stylesheet">` 的方式引用，取代原来 webpack 3.x 的 `extract-text-webpack-plugin`，**只在生产状态下使用**。

        它同样需要在插件中增加 `MiniCssExtractPlugin()` 以配合使用。

    3. `css-loader`

        支持通过隐式的方式加载资源。例如如果在 JS 文件中编写 `import 'style.css'`，或者在样式文件中编写 `url('image.png')`，经过 `css-loader` 可以把 `style.css` 和 `image.png` 都引入到 webpack 的处理流程中。这样针对 css 的所有 loader 就可以处理 `style.css`，而针对所有图片的 loader（例如尺寸很小就自动转为 base64 的 `url-loader`）都可以处理 `image.png` 了。

    4. `less-loader`

        使用 less 预处理器必须加载的 loader，用以把 less 语法转化为普通的 css 语法。

        基本上每个预处理器都有对应的 loader，例如 `stylus-loader`, `sass-loader` 等等，可以按需使用。

4. `/\.css$/`

    和 `.js` 规则相似，这条规则可以同时应用于 `.css` 的后缀文件以及 Vue 中的 `<style>` （且没有写 `lang` 的）部分。

#### plugins（重点）

插件和规则类似，也是对加载进入 webpack 的资源进行处理。不过和规则不同，它并不以正则（多数为后缀名）决定是否进入，而是全部进入后通过插件自身的 JS 写法来确定处理哪些，因此更加灵活。

前面提过，有些功能需要 loader 和 plugins 配合使用，因此也需要声明在这里，比如 `VueLoaderPlugin()` 和 `MiniCssExtractPlugin()`。

在我们的项目中，插件分为两类。一类是不论环境（开发还是生产）都要使用的，这一类有 2 个：

```javascript
{
  "plugins": [
    // 和 vue-loader 配合使用
    new VueLoader(),

    // 输出 index.html 到 output
    new HtmlwebpackPlugin({
      template: resolvePath('index.html')
    })
  ]
}
```

另外一类是生产环境才需要使用的，也有 2 个：

```javascript
if (isProd) {
  webpackConfig.plugins.push(
    // 每次 build 清空 output 目录
    new CleanWebpackPlugin(resolvePath('../vue-dist'))
  )
  webpackConfig.plugins.push(
    // 分离单独的 CSS 文件到 output，和 MiniCssExtractPlugin.loader 配合使用
    new MiniCssExtractPlugin({
      filename: 'style.css',
    })
  )
}
```

#### optimization

optimization 是一个 webpack 4 新增的配置项，主要处理生产环境下的各类优化（例如压缩，提取公共代码等等），所以大部分的优化都在 `mode === 'production'` 时会使用。这里我们只使用它的一个功能，即分包，以前的写法是 `new webpack.optimize.CommonChunkPlugin()`，现在只要配置就可以了，配置方法也很简单：

```javascript
{
  optimization: {
    splitChunks: {
      chunks: 'all'
    }
  }
}
```

这样，从 `node_modules` 来的代码会打包到一起，命名为 `vendors~app.[hash].js`，这里面可能包含了 Vue, Vuex, Vue Router 等等第三方的代码。这些代码不会经常修改，所以独立出来并添加长时间的强制缓存能显著提升站点访问速度。

### 看看编译完成的产物

通过运行 `npm run build`，能够调用 `webpack-cli` 来运行刚才编写的配置文件。编译成功后，会在根目录下生成一个 `vue-dist` 目录，里面存放的内容如下：(如果做了 Vue 的路由懒加载，即 `const XXX = () => import('@/XXX.vue')`，文件会根据路由分割，因此数量会更多)

![vue-dist](http://boscdn.bpc.baidu.com/assets/easonyq/vue-webpack/vue-dist.png)

总共 4 个文件

1. `index.html` 存放唯一的 HTML 入口，里面包含对各 JS, CSS 文件的引用，并定义了容器节点。使用静态服务器启动后，由于 JS 的执行，可以执行前端渲染。

2. `style.css` 存放所有的样式。这些样式都是从每个 Vue 文件的 `<style lang="less">` 部分中抽出来合并而成的。

3. `app.[hash].js` 存放所有的**自定义** JS，即每个 Vue 文件的 `<script>` 部分，以及如 `app.js`, `router.js`, `store.js` 的代码等等。

4. `vendors~app.[hash].js` 如上所述，存放所有**类库** JS，如 vue-router, vuex 本身的代码。

对 nodejs 来说，需要关心的只是这个 `index.html` 而已，其他 3 个都会由它负责引入。那么我们接下来看看如何改造 nodejs 部分。

## 改造 nodejs(koa) 部分

koa 部分需要我们改造的点主要有：

1. `package.json`

    nodejs 项目的标配，记录依赖，脚本，项目信息等等。**我们需要在这里和 Vue 端的 package.json 进行合并，尤其是 `npm run dev` 脚本的合并。**

2. `index.js`

    nodejs 的代码入口，通过命令 `node index.js` 启动整个项目。在这里可以注册路由规则，注册中间件，注册 wdm 等等，绝大部分逻辑都在这里。

    因为开发环境和生产环境的行为不尽相同（例如开发环境需要 wdm 而生产环境不需要），因此可以分为两个文件（`index.dev.js` 和 `index.prod.js`），也可以在一个文件中通过环境变量判断，这个因人而异。

    虽然 koa 的路由规则和中间件都可以写在这里，但通常只要是略有规模的项目，都会把路由处理和中间件分别独立成 `actions` 和 `middlewares` 目录分开存放（名字怎么起看自己喜好）。配置文件（例如配置启动端口号）也通常会独立成 `config.js` 或者 `config` 目录。其他的例如 `util` 目录等也都按需建立。

    **我们需要在这里统一前后路由，并使用 wdm 等**

### package.json

在“改造 Vue 部分”的 package.json 中曾经讲过，Vue 项目的依赖都直接复制到外层的 package.json 中来，还增加了一条 `npm run build` 命令。这里会再列出两条命令，达成最基本的的需求。

我为了区分运行环境，把 `index.js` 拆解为了 `index.dev.js` 和 `index.prod.js`。如上面所说，你也可以就在一个文件里用 `process.env.NODE_ENV` 来判断运行环境。

#### 启动开发服务器

常规的 koa 服务，一般我们通过 `node index.js` 来启动。但 nodejs 默认没有热加载，因此修改了 nodejs 代码需要重启服务器，比较麻烦。以前我会使用 `chokidar` 模块监听文件系统的变化，并在变化时执行 `delete require.cache[path]` 来实现简单的热加载机制。但现在有一款更方便的工具帮我们做了这个事情，那就是 `nodemon`。

```javascript
"nodemon -e js --ignore vue/ index.dev.js"
```

它的使用方式也很简单，把 `node index.js` 换成 `nodemon index.js`，他就会监听以这个入口执行的所有文件的变化，并自动重启。但我们这里还额外使用了两个配置项。`-e` 表示指定扩展名，这里我们只监听 js。 `--ignore` 指定忽略项，因为 `vue/` 目录中有 webpack 帮我们执行热加载，因此它的修改可以忽略。其他可用的配置项可以参考 nodemon 的[主页](https://github.com/remy/nodemon)。

#### 启动线上环境服务器

这个就简单了，直接执行 `node` 命令即可。所以最终的脚本部分如下：

```javascript
{
  "scripts": {
    "dev": "nodemon -e js --ignore vue/ index.dev.js",
    "build": "cross-env NODE_ENV=production webpack --config vue/webpack.config.js",
    "start": "node index.prod.js"
  }
}
```

### index.js

这个文件是 koa 的启动入口，它的大致结构如下（我使用了 koa-router 来管理路由，且只列举最最简单的骨架）：

```javascript
// 引用基本类库
const Koa = require('koa')
const Router = require('koa-router')
const koaStatic = require('koa-static')

// 初始化
const app = new Koa()
const router = new Router()

// 常规项目可能有中间件，即处理所有路由的逻辑，如验证登录，记录日志等等，这里省略

// 注册路由到 koa-router。
// 常规项目路由很多，应该独立到一个目录去一个个注册
router.get('/api/hello', ctx => {
  ctx.body = {message: 'Greeting from koa'}
})

// koa-router 以中间件的形式注册给 koa
// 就理解为固定写法
app.use(router.routes());
app.use(router.allowedMethods());
// 为 public 目录启动静态服务，可以放图片，字体等等。Vue 部分打包的资源在 vue-dist，不在这里。
app.use(koaStatic('public'));

// 实际项目可能端口还要写到配置文件，这里随意了
app.listen(8080)
```

我们从开发环境和线上环境两个方面来讨论对这个文件的改造方式。

#### 开发环境

首先是合并前后端路由。

Vue 端使用的是 history 路由模式，因此本身就需要 nodejs 来配合，[Vue 官方](https://router.vuejs.org/zh/guide/essentials/history-mode.html)推荐的是 [connect-history-api-fallback 中间件](https://github.com/bripkens/connect-history-api-fallback)，不过那是针对 express 的。我找到了一个给 koa 使用的相同功能的中间件，名为 [koa2-history-api-fallback](https://github.com/luzuoquan/koa2-history-api-fallback)。

不论是哪一个中间件，原理是一样的。因为 SPA 只生成一个 `index.html`，因此**所有的 navigate 路由**都必须定向到这个文件才行，否则例如 `/user/index` 这样的路由，浏览器会去找 `/user/index.html`，显然是找不到的。

既然 Vue 需求的是**所有的 navigate 路由**，显然它不能注册在 koa 的路由之前，否则 koa 的路由将永远无法生效。因此路由注册顺序就很自然了：**先后端再前端**。另外这里说的 navigate 路由，指的是请求 HTML 的第一个请求，静态资源的请求不在其内，因此例如上述的 `public` 静态路由和 Vue 的中间件的前后顺序就无所谓了。

所以路由部分融合后大概是这样：

```javascript
// 后端(koa)路由
// koa-router 的单个注册部分省略
app.use(router.routes());
app.use(router.allowedMethods());

// 前端(vue)路由
// 所有 navigate 请求重定向到 '/'，因为 webpack-dev-middleware 只服务这个路由
app.use(history({
  htmlAcceptHeaders: ['text/html'],
  index: '/'
}));
app.use(koaStatic('public'));
```

其次就是问题的最开端，webpack-dev-middleware 的使用。

和 Vue 的中间件类似，webpack-dev-middleware 也是只支持 express 的（这些都表明 express 的生态更好），然后我也找了个 koa 版本的替代方案，叫做 [koa-webpack](https://github.com/shellscape/koa-webpack)。

使用起来倒也不麻烦，如下：

```javascript
const koaWebpack = require('koa-webpack')
const webpackConfig = require('./vue/webpack.config.js')

// 注意这里是个异步，所以和其他 app.use 以及最终的 app.listen 必须在一起执行
// 可以使用 async/await 或者 Promise 保证这一点
koaWebpack({
  config: webpackConfig,
  devMiddleware: {
    stats: 'minimal'
  }
}).then(middleware => {
  app.use(middleware)
})
```

一个完整的 `index.dev.js` 可以查看[这里](https://github.com/easonyq/vue-nodejs-template/blob/master/index.dev.js)。

#### 线上环境

线上环境和开发环境有两处不同，我们着重讲一下这两个不同点。

首先，线上环境不使用 webpack-dev-middleware (koa-webpack)，因此这部分代码不需要了。

其次，因为构建后的 Vue 代码全部位于 vue-dist 目录，而我们需要的 HTML 入口以及其他 JS, CSS文件都在其中，因此我们需要把 vue-dist 目录添加到静态服务中可供访问，另外 history fallback 的目标也有所改变，如下：

```javascript
// 后端(koa)路由
// koa-router 的单个注册部分省略
app.use(router.routes());
app.use(router.allowedMethods());

// 前端(vue)路由
// 所有 navigate 请求重定向到 /vue-dist/index.html 这个文件，配合下面的 koaStatic('vue-dist')，这里只要填到 '/index.html' 即可。
app.use(history({
  htmlAcceptHeaders: ['text/html'],
  index: '/index.html'
}));
app.use(koaStatic('vue-dist'));
app.use(koaStatic('public'));
```

一个完整的 `index.prod.js` 可以查看[这里](https://github.com/easonyq/vue-nodejs-template/blob/master/index.prod.js)。

## 这么多配置新手看了还是很懵怎么办？

虽然我们讨论了这么多，但不要害怕，实际上的重点只有三个，我们来总结一下：

1. 我们需要自己编写 Vue 的 webpack.config.js，处理 loader, plugins 等等

2. 我们需要合并前后端的两个 package.json，把两方的依赖合并，并编写三条脚本 (`dev`, `build`, `start`)

3. 我们需要改动 `index.js`，处理路由顺序，并在开发环境调用 webpack-dev-middleware

为了简单上手，我把项目中的业务代码抽离，留下了一个骨架，可以作为 Vue + koa 项目的启动模板，放在 [easonyq/vue-nodejs-template](https://github.com/easonyq/vue-nodejs-template)。不过我觉得我们还是应当掌握配置方法和原理，这样以后如果技术栈的某一块发生了变化（例如 webpack 出了 5），我们也能够自己研究修改，而不是每次都以解决任务为最优先，能跑起来就不管了。

愿我们大家在前端道路上都能越走越顺！
