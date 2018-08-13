# treeShaking

webpack 2.0 开始引入 tree shaking 技术。在介绍技术之前，先介绍几个相关概念：

* AST
    对 JS 代码进行语法分析后得出的语法树 (Abstract Syntax Tree)。AST语法树可以把一段 JS 代码的每一个语句都转化为树中的一个节点。

* DCE
    Dead Code Elimination，在保持代码运行结果不变的前提下，去除无用的代码。这样的好处是：

    * 减少程序体积
    * 减少程序执行时间
    * 便于将来对程序架构进行优化

    而所谓 Dead Code 主要包括：

    * 程序中没有执行的代码 (如不可能进入的分支，return 之后的语句等)
    * 导致 dead variable 的代码(写入变量之后不再读取的代码)

tree shaking 是 DCE 的一种方式，它可以在打包时忽略没有用到的代码。

![tree shaking](https://user-gold-cdn.xitu.io/2018/1/4/160bfdcf2a31ce4a?imageslim)

## 机制简述

tree shaking 是 rollup 作者首先提出的。这里有一个比喻：

> 如果把代码打包比作制作蛋糕。传统的方式是把鸡蛋(带壳)全部丢进去搅拌，然后放入烤箱，最后把(没有用的)蛋壳全部挑选并剔除出去。而 treeshaking 则是一开始就把有用的蛋白蛋黄放入搅拌，最后直接作出蛋糕。

因此，相比于 __排除不使用的代码__，tree shaking 其实是 __找出使用的代码__。

基于 ES6 的静态引用，tree shaking 通过扫描所有 ES6 的 `export`，找出被 `import` 的内容并添加到最终代码中。 webpack 的实现是把所有 `import` 标记为有使用/无使用两种，在后续压缩时进行区别处理。因为就如比喻所说，在放入烤箱(压缩混淆)前先剔除蛋壳(无使用的 `import`)，只放入有用的蛋白蛋黄(有使用的 `import`)

## 使用方法

首先源码必须遵循 ES6 的模块规范 (`import` & `export`)，如果是 CommonJS 规范 (`require`) 则无法使用。

根据webpack官网的提示，webpack2 支持 tree-shaking，需要修改配置文件，指定babel处理js文件时不要将ES6模块转成CommonJS模块，具体做法就是：

在.babelrc设置babel-preset-es2015的modules为fasle，表示不对ES6模块进行处理。

```json
// .babelrc
{
    "presets": [
        ["es2015", {"modules": false}]
    ]
}
```

__经过测试，webpack 3 和 4 不增加这个 `.babelrc` 文件也可以正常 tree shaking__

## Tree shaking 两步走

webpack 负责对代码进行标记，把 `import` & `export` 标记为 3 类：

1. 所有 `import` 标记为 `/* harmony import */`
2. 被使用过的 `export` 标记为 `/* harmony export ([type]) */`，其中 `[type]` 和 webpack 内部有关，可能是 `binding`, `immutable` 等等。
3. 没被使用过的 `import` 标记为 `/* harmony export [FuncName] */`，其中 `[FuncName]` 为 `export` 的方法名称

之后在 Uglifyjs (或者其他类似的工具) 步骤进行代码精简，把没用的都删除。

## 实例分析

所有实例代码均在[demo/webpack 目录](https://github.com/easonyq/easonyq.github.io/tree/master/%E5%AD%A6%E4%B9%A0%E8%AE%B0%E5%BD%95/demo/webpack)

### 方法的处理

```javascript
// index.js
import {hello, bye} from './util'

let result1 = hello()

console.log(result1)

```

```javascript
// util.js
export function hello () {
  return 'hello'
}

export function bye () {
  return 'bye'
}
```

编译后的 bundle.js 如下：

```javascript
/******/ ([
/* 0 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__util__ = __webpack_require__(1);


let result1 = Object(__WEBPACK_IMPORTED_MODULE_0__util__["a" /* hello */])()

console.log(result1)


/***/ }),
/* 1 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (immutable) */ __webpack_exports__["a"] = hello;
/* unused harmony export bye */
function hello () {
  return 'hello'
}

function bye () {
  return 'bye'
}
```

注：省略了 `bundle.js` 上边 webpack 自定义的模块加载代码，那些都是固定的。

对于没有使用的 `bye` 方法，webpack 标记为 `unused harmony export bye`，但是代码依旧保留。而 `hello` 就是正常的 `harmony export (immutable)`。

之后使用 `UglifyJSPlugin` 就可以进行第二步，把 `bye` 彻底清除，结果如下：

![funciton](http://boscdn.bpc.baidu.com/assets/easonyq/tree-shaking/function.png)

只有 `hello` 的定义和调用。

### 类(class) 的处理

```javascript
// index.js
import Util from './util'

let util = new Util()
let result1 = util.hello()
console.log(result1)
```

```javascript
// util.js
export default class Util {
  hello () {
    return 'hello'
  }

  bye () {
    return 'bye'
  }
}
```

编译后的 bundle.js 如下：

```javascript
/******/ ([
/* 0 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__util__ = __webpack_require__(1);


let util = new __WEBPACK_IMPORTED_MODULE_0__util__["a" /* default */]()
let result1 = util.hello()
console.log(result1)


/***/ }),
/* 1 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
class Util {
  hello () {
    return 'hello'
  }

  bye () {
    return 'bye'
  }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = Util;
```

注意到 webpack 是对 `Util` 类整体进行标记的（标记为被使用），而不是分别针对两个方法。也因此，最终打包的代码依然会包含 `bye` 方法。这表明 __webpack tree shaking 只处理顶层内容__，例如类和对象内部都不会再被分别处理。

这主要也是由于 JS 的动态语言特性所致。如果把 `bye()` 删除，考虑如下代码：

```javascript
// index.js
import Util from './util'

let util = new Util()
let result1 = util[Math.random() > 0.5 ? 'hello', 'bye']()
console.log(result1)
```

编译器并不能识别一个方法名字究竟是以直接调用的形式出现 (`util.hello()`) 还是以字符串的形式 (`util['hello']()`) 或者其他更加离奇的方式。因此误删方法只会导致运行出错，得不偿失。

## 副作用

副作用的意思某个方法或者文件执行了之后，还会对全局其他内容产生影响的代码。例如 polyfill 在各类 `prototype` 加入方法，就是副作用的典型。（也可以看出，程序和吃药不同，副作用不全是贬义的）

副作用总共有两种形态，是精简代码不得不考虑的问题。__我们平时在重构代码时，也应当以相类似的思维去进行，否则总有踩坑的一天。__

### 模块引入带来的副作用

```javascript
// index.js
import Util from './util'

console.log('Util unused')
```

```javascript
// util.js
console.log('This is Util class')

export default class Util {
  hello () {
    return 'hello'
  }

  bye () {
    return 'bye'
  }
}

Array.prototype.hello = () => 'hello'
```

如上代码经过 webpack + uglify 的处理后，会变成这样：

![import-side-effects](http://boscdn.bpc.baidu.com/assets/easonyq/tree-shaking/import-side-effects.png)

虽然 `Util` 类被引入之后没有进行任何使用，但是不能当做没引用过而直接删除。在混合后的代码中，可以看到 `Util` 类的本体 (`export` 的内容) 已经没有了，但是前后的 `console.log` 和对 `Array.prototype` 的扩展依然保留。这就是编译器为了确保代码执行效果不变而做的妥协，因为它不知道这两句代码到底是干嘛的，所以他默认认定所有代码 __均有__ 副作用。

### 方法调用带来的副作用

```javascript
// index.js
import {hello, bye} from './util'

let result1 = hello()
let result2 = bye()

console.log(result1)
```

```javascript
// util.js
export function hello () {
  return 'hello'
}

export function bye () {
  return 'bye'
}
```

我们引入并调用了 `bye()`，但是却没有使用它的返回值 `result2`，这种代码可以删吗？（扪心自问，如果是你人肉重构代码，直接删掉这行代码的可能性有没有超过 90% ？）

![invoke-side-effects](http://boscdn.bpc.baidu.com/assets/easonyq/tree-shaking/invoke-side-effects.png)

webpack 并没有删除这行代码，至少没有删除全部。它确实删除了 `result2`，但保留了 `bye()` 的调用（压缩的代码表现为 `Object(r.a)()`）以及 `bye()` 的定义。

这同样是因为编译器不清楚 `bye()` 里面究竟做了什么。如果它包含了如 `Array.prototye` 的扩展，那删掉就又出问题了。

### 如何解决副作用？

我们很感谢 webpack 如此严谨，但如果某个方法就是没有副作用的，我们该怎么告诉 webpack 让他放心大胆的删除呢？

有 3 个方法，适用于不同的情况。

#### pure_funcs

```javascript
// index.js
import {hello, bye} from './util'

let result1 = hello()
let a = 1
let b = 2
let result2 = Math.floor(a / b)

console.log(result1)
```

util.js 和之前相同，不再重复。有差别的是 webpack.config.js，需要增加参数 `pure_funcs`，告诉 webpack `Math.floor` 是没有副作用的，你可以放心删除：

```javascript
plugins: [
  new UglifyJSPlugin({
    uglifyOptions: {
      compress: {
          pure_funcs: ['Math.floor']
      }
    }
  })
],
```

![pure-funcs-before](http://boscdn.bpc.baidu.com/assets/easonyq/tree-shaking/pure-funcs-before.png)

![pure-funcs-after](http://boscdn.bpc.baidu.com/assets/easonyq/tree-shaking/pure-funcs-after.png)

在添加了 `pure_funcs` 配置后，原来保留的 `Math.floor(.5)` 被删除了，达到了我们的预期效果。

但这个方法有一个很大的局限性，在于如果我们把 webpack 和 uglify 合并使用，经过 webpack 的代码的方法名已经被重命名了，那么在这里配置原始的方法名也就失去了意义。而例如 `Math.floor` 这类全局方法不会重命名，才会生效。因此适用性不算太强。

#### package.json 的 sideEffects

webpack 4 在 package.json 新增了一个配置项叫做 `sideEffects`， 值为 `false` 表示整个包都没有副作用；或者是一个数组列出有副作用的模块。详细的例子可以查看 webpack 官方提供的[例子](https://github.com/webpack/webpack/tree/next/examples/side-effects)。

从结果来看，如果 `sideEffects` 值为 `false`，当前包 `export` 了 5 个方法，而我们使用了 2 个，剩下 3 个也不会被打包，是符合预期的。但这要求包作者的自觉添加，因此在当前 webpack 4 推出不久的情况下，局限性也不算小。

#### concatenateModule

webpack 3 开始加入了 `webpack.optimize.ModuleConcatenateModulePlugin()`，到了 webpack 4 直接作为 `mode = 'production' 的默认配置。这是对 webpack bundle 的一个优化，把本来“每个模块包裹在一个闭包里”的情况，优化成“所有模块都包裹在同一个闭包里”的情况。本身对于代码缩小体积有很大的提升，这里也能侧面解决副作用的问题。

依然选取这样 2 个文件作为例子：

```javascript
// index.js
import {hello, bye} from './util'

let result1 = hello()
let result2 = bye()

console.log(result1)
```

```javascript
// util.js
export function hello () {
  return 'hello'
}

export function bye () {
  return 'bye'
}
```

在开启了 concatenateModule 功能后，打包出来的代码如下：

![concatenateModule](http://boscdn.bpc.baidu.com/assets/easonyq/tree-shaking/concatenateModule.png)

首先，`bye()` 方法的调用和本体都被消除了。

其次，`hello()` 方法的调用和定义被合成到了一起，变成直接 `console.log('hello')`

第三就是这个功能原有的目的：代码量减少了。

这个功能的本意是把所有模块最终输出到同一个方法内部，从而把调用和定义合并到一起。这样像 `bye()` 这样没有副作用的方法就可以在合并之后被轻易识别出来，并加以删除。有关这个功能更加详细的介绍可以看[这篇文章](https://zhuanlan.zhihu.com/p/27980441)

## 总结

1. 使用 ES6 模块语法编写代码
2. 工具类函数尽量以单独的形式输出，不要集中成一个对象或者类
3. 声明 sideEffects
4. 自己在重构代码时也要注意副作用
