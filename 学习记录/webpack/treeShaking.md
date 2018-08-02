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

    * 程序中没有执行的代码 （如不可能进入的分支，return 之后的语句等）
    * 导致 dead variable 的代码（写入变量之后不再读取的代码）

tree shaking 是 DCE 的一种方式，它可以在打包时忽略没有用到的代码。

## 内部机制

tree shaking 是 rollup 作者首先提出的。这里有一个比喻：

> 如果把代码打包比作制作蛋糕。传统的方式是把鸡蛋（带壳）全部丢进去搅拌，然后放入烤箱，最后把（没有用的）蛋壳全部挑选并剔除出去。而 treeshaking 则是一开始就把有用的蛋白蛋黄放入搅拌，最后直接作出蛋糕。

因此，相比于 __排除不使用的代码__，tree shaking 其实是 __找出使用的代码__。

基于 ES6 的静态引用，tree shaking 通过扫描所有 ES6 的 `export`，找出被 `import` 的内容并添加到最终代码中。 (webpack 的实现其实相反，还是把没有 `import` 的代码进行标记)

## 使用方法

根据webpack官网的提示，webpack支持tree-shaking，需要修改配置文件，指定babel处理js文件时不要将ES6模块转成CommonJS模块，具体做法就是：

在.babelrc设置babel-preset-es2015的modules为fasle，表示不对ES6模块进行处理。

```json
// .babelrc
{
    "presets": [
        ["es2015", {"modules": false}]
    ]
}
```

仅仅使用 webpack 自身来进行打包时，它只对无用代码进行标记，而并不会删除。举例来说，如果定义：

```javascript
// module.js
export const sayHello = name => `Hello ${name}!`;
export const sayBye = name => `Bye ${name}!`;
export const sayHi = name => `Hi ${name}!`;
```

```javascript
// index.js
import { sayHello } from './module';
import { sayHi } from './module';
const element = document.createElement('h1');
element.innerHTML = sayHello('World') + sayHi('my friend');
document.body.appendChild(element);
```

编译后的 bundle.js 如下：

```javascript
/* 0 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {
 
"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return sayHello; });
/* unused harmony export sayBye */
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "b", function() { return sayHi; });
 
 
var sayHello = function sayHello(name) {
  return "Hello " + name + "!";
};
var sayBye = function sayBye(name) {
  return "Bye " + name + "!";
};
 
var sayHi = function sayHi(name) {
  return "Hi " + name + "!";
};
 
/***/ }),
/* 1 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {
 
"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__module__ = __webpack_require__(0);
 
 
 
var element = document.createElement('h1');
element.innerHTML = Object(__WEBPACK_IMPORTED_MODULE_0__module__["a" /* sayHello */])('World') + Object(__WEBPACK_IMPORTED_MODULE_0__module__["b" /* sayHi */])(' to meet you');
document.body.appendChild(element);
 
/***/ })
```

对于没有使用的 `sayBye` 方法，webpack 标记为 `unused harmony export`，但是代码依旧保留。而其他可用的方法都有进行 `export`。

这样再使用 `webpack -p` 或者 `webpack --optimize-minimize` 就会对无用的 `sayBye` 进行删除了。这两个命令的内部也是使用 `webpack.optimize.UglifyJsPlugin()` 来进行的。


## webpack treeshaking 的局限性

1. 只能是 ES6 模块，不能动态声明或者引入。
    因为需要在编译时分析出依赖关系，因此动态引入（如 if 分支分别引入两个包）不运行起来是不知道结果的。

    ES6 默认要求静态模块引入，因此默认支持。ES6 的这个要求其实也是为静态依赖分析创造可能，毕竟在 WEB 环境下，文件加载速度远远慢于 CommonJs 预设的读取本地文件。
    而 CommonJS 是支持动态引入的，因此 treeshaking 不会对 CommonJS 起效。这也是为什么在配置 babel 的时候需要 `modules: false`，否则因为 `es2015 presets` 里面的 `transform-es2015-modules-commonjs` 把 `import` 变成 `require` 之后，treeshaking 也就失效了。

2. 只处理模块级别，不能精确到函数级别。
    只对模块 export 的内容进行分析和精简。对于内部的无用函数并不会进行检查。

## 注意点

并不是不被 `import` 的代码都能被随意消除的。举例来说：

```javascript
function A() {

}

A.prototype.hello = function () {
    console.log('hello')
}

export default A
```

这个 A 即使没有 `import`，也不能删除，因为并不确定 A 的实例的 `hello` 方法是否被使用了。
例如想象 A 是 Array， `let arr = new Array(); arr.hello()` 并没有 `import Array` （因为是内置的），且编译器也无法区分 A 不是内置的，而 Array 是内置的。
因此最保险的做法就是保留这个 A。

