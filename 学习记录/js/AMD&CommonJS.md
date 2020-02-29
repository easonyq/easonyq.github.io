# CommonJS

nodejs 是 CommonJS 规范的实现。

简单来说，一个文件就是一个模块，文件中定义的变量方法都不对外可见，除非显式的指定，如 `global.xxx` 或者 `module.exports` 等等。使用 `require` 进行模块加载。

CommonJS 是以在浏览器环境之外构建 JavaScript 生态系统为目标而产生的项目，比如在服务器和桌面环境中。

# AMD

requireJS 是 AMD 规范的实现。（此外 curl.js 也是，但不如 requireJS 那么有名而已）

AMD（异步模块定义）是为浏览器环境设计的，因为 CommonJS 模块系统是同步加载的，当前浏览器环境还没有准备好同步加载模块的条件。举例来说

```javascript
var math = require('math');
math.add(2, 3);
```

在服务端执行时，所有模块都在本地磁盘，因此同步等待 `math` 的载入不会花费太久；但浏览器端，所有模块都是远程资源。如果这类加载操作全都等待，浏览器就会假死。因此浏览器只能采取 __异步加载__ 的模式。

AMD 也是用 `reuqire`，但它有两个参数

```javascript
require([module], callback)
```

如上述例子，会改写成

```javascript
require(['math'], function (math) {
    math.add(2, 3)
});
```

AMD 还有 `define` 用以定义模块。总之参考 requireJS 的语法就行了。

# UMD

UMD = 通用模块规范。它兼容了 AMD 和 CommonJS，为的就是创造一种通用模式。使用这种模式编写的代码，在两种环境都能够运行。但是代码量更大，冗余代码多，可读性不强。

简单来说，UMD 的步骤是
1. 先判断是否支持 Node.js 模块格式（`exports` 是否存在），存在则使用Node.js模块格式。
2. 再判断是否支持 AMD（`define` 是否存在），存在则使用AMD方式加载模块。
3. 前两个都不存在，则将模块公开到全局（`window` 或 `global`）。

# ES

es2015 提出的模块化编程方式，使用 `import`, `export` 等进行引入和导出。引入导出的单位是方法，变量或者类，而不是文件。因此一个文件可以有多个导出。但只能有一个默认导出。

```javascript
// a.js
export var a = 'a';
export function sayA() {
    console.log('sayA');
}
export default class A {
    greet() {
        console.log('This is A');
    }
}

// b.js
import {a, sayA} from './a';
console.log(a) // a
sayA() // sayA

// c.js
import A from './a';
let a = new A();
a.greet(); // This is A
```

注意 `import` 后面如果跟 `{ xxx }` 并不是对象解构的意思，是获取普通导出的内容。如果不跟大括号，则使用默认导出。一般多用普通导出。

`import` 还有两种写法： 使用 `as` 更换导入内容的命名，以及使用 `* as XXX` 直接导入全部，之后使用 `XXX.yy` 来使用。

# rollup 的 format

rollup 在定义配置项时，有一个 format 可以指定输出格式。因为和模块类型有关，也记录到一起。

* amd – 异步模块定义，用于像RequireJS这样的模块加载器
* cjs – CommonJS，适用于 Node 和 Browserify/Webpack
* es – 将软件包保存为ES模块文件
* iife – 一个自动执行的功能，适合作为 &lt;script&gt; 标签。（如果要为应用程序创建一个捆绑包，您可能想要使用它，因为它会使文件大小变小。）
* umd – 通用模块定义，以amd，cjs 和 iife 为一体
