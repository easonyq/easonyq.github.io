# babel

简单来说把 ES6/7/8 解析成 ES5 的工具。当前版本是 6.x

参考文档：https://segmentfault.com/a/1190000011155061#articleHeader8

## 使用方法

可以使用单体文件 (standalone script)，命令行 (cli)，和构建工具的插件 (webpack 的 babel-loader, rollup 的 rollup-plugin-babel)。

常规情况是使用最后一种。

## 运行方式

babel 总共分为三个阶段：解析，转换，生成。

初始情况下，babel 不对代码做任何事情，相当于 `const babel = code => code`。

当我们添加语法插件之后，在解析这一步就使得 babel 能够解析更多的语法。（顺带一提，babel 内部试用的解析类库叫做 babylon，并非 babel 自行开发）

当我们添加转译插件之后，在转换这一步就增加了一系列操作。

## 配置文件

使用插件和使用 preset 都需要配置来告诉 babel。

配置可以通过命令参数的方式传入 (package.json 中 script 那部分)，但更多是写在 `.babelrc` 里面。

### preset

一组插件的集合，因为常用，所以不必重复定义 & 安装。（单点和套餐的差别，只不过这里的套餐只是方便，并没有打折）

preset 分为两种：

1. 官方内容，目前包括 env, react, flow 三者。

2. stage-x，这里面包含的都是 ES7 的各个阶段的草案。

    这里面还细分为
    * Stage 0 - 稻草人: 只是一个想法，可能是 babel 插件。
    * Stage 1 - 提案: 初步尝试。
    * Stage 2 - 初稿: 完成初步规范。
    * Stage 3 - 候选: 完成规范和浏览器初步实现。
    * Stage 4 - 完成: 将被添加到下一年度发布。

    例如 `syntax-dynamic-import` 就是 stage-2 的内容，`transform-object-rest-spread` 就是 stage-3 的内容。
    此外，低一级的 stage 会包含所有高级 stage 的内容。如 stage-1 会包含 stage-2, stage-3 的所有内容。
    stage-4 在下一年更新会直接放到 env 中，所以没有单独的 stage-4 可供使用。`transform-async-to-generate` 就属于 stage-4。

### plugin

syntax-xxxx 开头的为语法插件，它帮助 babel 认识某种语法。
transform-xxxx 开头的为转译插件，它能够让 babel 转译某种语法为其他可运行的语法。__在使用转译插件时不必在重复引用语法插件__。

所有官方的插件可以查看[babeljs的官网列表](https://babeljs.cn/docs/plugins)。

### 执行顺序

* Plugin 会运行在 Preset 之前。
* Plugin 会从第一个开始顺序执行。ordering is first to last.
* Preset 的顺序则 __刚好相反__(从最后一个逆序执行)。

preset 的逆向顺序主要是为了保证向后兼容，因为大多数用户会在 “stage-0” 之前列出 “es2015” 。因此，我们编排 preset 的时候，应当把 es2015 写在 stage-x 的前面。

### 配置属性

preset 和 plugin 都可以通过把名字和配置项放在一个数组中来实现每个单项的配置。例如我们最熟悉的 env 就是这样。

```
presets: [
    ['env', {
        module: false
    }],
    'stage-2'
]
```

### env

因为 env 最为常用，因此单独拿出来讲一下。

如果不写任何配置项，env 等价于 latest，也等价于 es2015 + es2016 + es2017 三个写在一起（不包含 stage-x 中的插件）。env 包含的插件列表维护在[这里](https://github.com/babel/babel-preset-env/blob/master/data/plugin-features.js)

env 的核心目的是通过配置得知目标环境的特点，然后只做必要的转换。例如目标浏览器支持 ES6，那么 es2015 这个 preset 其实是不需要的，于是代码就可以小一点。（一般转化后的代码总是更长）

下面列出几种比较常用的配置方法：

```json
{
  "presets": [
    ["env", {
      "targets": {
        "browsers": ["last 2 versions", "safari >= 7"]
      }
    }]
  ]
}
```

如上配置将考虑所有浏览器的最新2个版本（safari大于等于7.0的版本）的特性，将必要的代码进行转换。而这些版本已有的功能就不进行转化了。这里的语法可以参考 [browserslist](https://github.com/browserslist/browserslist)

```json
{
  "presets": [
    ["env", {
      "targets": {
        "node": "6.10"
      }
    }]
  ]
}
```

如上配置将目标设置为 nodejs，并且支持 6.10 及以上的版本。也可以使用 `'current'` 来支持最新稳定版本。例如箭头函数在 nodejs 6及以上将不被转化，但如果是 nodejs 0.12 就会被转化了。

另外一个有用的配置项是 `modules`。它的取值可以是 `amd`, `umd`, `systemjs`, `commonjs` 和 `false`。这可以让 babel 以特定的模块化格式来输出代码。如果选择 `false` 就不进行模块化处理。

## 其他配套工具

除了 babel 本身，我们接触 babel 还会碰到很多 babel 开头的名词。其实他们各不相同，但又广泛使用。我们需要逐个了解一下。

### babel-node

`babel-node` 是 `babel-cli` 的一部分，它不需要单独安装。

它的作用是在 node 环境中，直接运行 ES6 的代码，而不需要额外进行转码。例如我们有一个 js 文件以 ES6 的语法进行编写（如使用了箭头函数）。我们可以直接使用 `babel-node es6.js` 进行执行，而不用再进行转码了。

`babel-node` = `babel-polyfill` + `babel-register`。那这两位又是谁呢？

### babel-register

babel-register 模块改写 `require` 命令，为它加上一个钩子。此后，每当使用 `require` 加载 `.js`、`.jsx`、`.es` 和 `.es6` 后缀名的文件，就会先用 Babel 进行转码。

使用时，必须首先加载 `require('babel-register')`。

需要注意的是，babel-register 只会对 `require` 命令加载的文件转码，而 __不会对当前文件转码__。

另外，由于它是实时转码，所以 __只适合在开发环境使用__。

### babel-polyfill

Babel 默认只转换 js 语法，而不转换新的 API，比如 Iterator、Generator、Set、Maps、Proxy、Reflect、Symbol、Promise等全局对象，以及一些定义在全局对象上的方法（比如 `Object.assign`）都不会转码。

举例来说，ES6 在 Array 对象上新增了 `Array.from` 方法。Babel 就不会转码这个方法。如果想让这个方法运行，必须使用 `babel-polyfill`。(内部集成了 `core-js` 和 `regenerator`)

使用时，在所有代码运行之前增加 `require('babel-polyfill')`。或者在 `webpack.config.js` 中将 `babel-polyfill` 作为第一个 entry。因此必须把 `babel-polyfill` 作为 `dependencies` 而不是 `devDependencies`

`babel-polyfill` 主要有两个缺点：

1. 使用 `babel-polyfill` 会导致打出来的包非常大，因为 `babel-polyfill` 是一个整体，把所有方法都加到原型链上。比如我们只使用了 `Array.from`，但它把 `Object.defineProperty` 也给加上了，这就是一种浪费了。这个问题可以通过单独使用 core-js 的某个类库来解决，core-js 都是分开的。

2. `babel-polyfill` 会污染全局变量，给很多类的原型链上都作了修改，如果我们开发的也是一个类库供其他开发者使用，这种情况就会变得非常不可控。

因此在实际使用中，如果我们无法忍受这两个缺点（尤其是第二个），通常我们会倾向于使用 `babel-plugin-transform-runtime`。

但如果代码中包含高版本 js 中类型的实例方法 (例如 `[1,2,3].includes(1)`)，这还是要使用 polyfill

### babel-runtime 和 babel-plugin-transform-runtime (重点)

我们时常在项目中看到 .babelrc 中使用 `babel-plugin-transform-runtime`，而 `package.json` 中的 `dependencies` (注意不是 `devDependencies`) 又包含了 `babel-runtime`，那这两个是不是成套使用的呢？他们又起什么作用呢？

先说 `babel-plugin-transform-runtime`。

babel 会转换 js 语法，之前已经提过了。以 `async/await` 举例，如果不使用这个 plugin （即默认情况），转换后的代码大概是：

```javascript
// babel 添加一个方法，把 async 转化为 generator
function _asyncToGenerator(fn) { return function () { ....}} // 很长很长一段

// 具体使用处
var _ref = _asyncToGenerator(function* (arg1, arg2) {
  yield (0, something)(arg1, arg2);
});
```

不用过于纠结具体的语法，只需看到，这个 `_asyncToGenerator` 在当前文件被定义，然后被使用了，以替换源代码的 `await`。但每个被转化的文件都会插入一段 `_asyncToGenerator` 这就导致重复和浪费了。

在使用了 `babel-plugin-transform-runtime` 了之后，转化后的代码会变成

```javascript
// 从直接定义改为引用，这样就不会重复定义了。
var _asyncToGenerator2 = require('babel-runtime/helpers/asyncToGenerator');
var _asyncToGenerator3 = _interopRequireDefault(_asyncToGenerator2);

// 具体使用处是一样的
var _ref = _asyncToGenerator3(function* (arg1, arg2) {
  yield (0, something)(arg1, arg2);
});
```

从定义方法改成引用，那重复定义就变成了重复引用，就不存在代码重复的问题了。

但在这里，我们也发现 `babel-runtime` 出场了，它就是这些方法的集合处，也因此，__在使用 `babel-plugin-transform-runtime` 的时候必须把 `babel-runtime` 当做依赖。__

再说 `babel-runtime`，它内部集成了

1. `core-js `: 转换一些内置类 (`Promise`, `Symbols`等等) 和静态方法 (`Array.from` 等)。绝大部分转换是这里做的。自动引入。

2. `regenerator`: 作为 `core-js` 的拾遗补漏，主要是 `generator/yield` 和 `async/await` 两组的支持。当代码中有使用 `generators/async` 时自动引入。

3. helpers, 如上面的 `asyncToGenerator` 就是其中之一，其他还有如 `jsx`, `classCallCheck` 等等，可以查看 [babel-helpers](https://github.com/babel/babel/blob/6.x/packages/babel-helpers/src/helpers.js)。在代码中有内置的 helpers 使用时（如上面的第一段代码）移除定义，并插入引用（于是就变成了第二段代码）。

`babel-plugin-transform-runtime` __不支持__ 实例方法 (例如 `[1,2,3].includes(1)`)

此外补充一点，把 helpers 抽离并统一起来，避免重复代码的工作还有一个 plugin 也能做，叫做 `babel-plugin-external-helpers`。但因为我们使用的 `transform-runtime` 已经包含了这个功能，因此不必重复使用。而且 babel 的作者们也已经开始讨论这两个插件过于类似，正在讨论在 babel 7 中把 `external-helpers` 删除，讨论在 [issue#5699](https://github.com/babel/babel/issues/5699) 中。

### babel-loader

TODO
