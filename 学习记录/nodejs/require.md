# require

[参考链接](http://www.ruanyifeng.com/blog/2015/05/require.html)

require 是 CommonJS 的基础，也是 nodejs 的基础。

## 处理顺序

处理顺序是指如果使用了 `reuquire('X')`，nodejs 以什么样的顺序去寻找这个 `X` 的包。具体如下：

1. 内置模块。如 `path`, `fs`, `http` 等等。

2. 如果 X 以 `./`, `/`, `../` 开头
    1. 根据 X 的相对路径确定它的绝对路径
    2. 将 X 作为文件的形式依次查找 `X`, `X.js`, `X.json`, `X.node`
    3. 将 X 作为路径的形式依次查找 `X/package.json`（此时读取其中的 main 字段），`X/index.js`, `X/index.json`, `X/index.node`

3. 如果 X 直接以普通字符开头（即不带路径又不是内置模块）
    1. 从当前路径的 `node_modules` 目录里面查找 X
    2. 如果没有找到，则往上一级目录查找。

4. 找不到，抛错

针对第 3 步，举一个例子来说。如果我们在 `/home/ry/projects/foo.js` 中执行 `require('bar')`，那么查找的顺序依次为：

```
/home/ry/projects/node_modules/bar
/home/ry/node_modules/bar
/home/node_modules/bar
/node_modules/bar
```

具体到每个目录，和第 2 步相同，即先作为文件查找，再作为目录查找。

## Module 构造函数

node 源码的 `lib/module.js` 文件中定义了所有模块的基类 Module。

```javascript
function Module(id, parent) {
  this.id = id;
  this.exports = {};
  this.parent = parent;
  this.filename = null;
  this.loaded = false;
  this.children = [];
}

module.exports = Module;

var module = new Module(filename, parent);
```

我们创建一个 `a.js`，输出这些变量：

```javascript
// a.js

console.log('module.id: ', module.id);
console.log('module.exports: ', module.exports);
console.log('module.parent: ', module.parent);
console.log('module.filename: ', module.filename);
console.log('module.loaded: ', module.loaded);
console.log('module.children: ', module.children);
console.log('module.paths: ', module.paths);
```

输出为：

```
$ node a.js

module.id:  .
module.exports:  {}
module.parent:  null
module.filename:  /home/ruanyf/tmp/a.js
module.loaded:  false
module.children:  []
module.paths:  [ '/home/ruanyf/tmp/node_modules',
  '/home/ruanyf/node_modules',
  '/home/node_modules',
  '/node_modules' ]
```

可以看到，如果没有父模块，直接调用当前模块，parent 属性就是 null，id 属性就是一个点。filename 属性是模块的绝对路径，path 属性是一个数组，包含了模块可能的位置。另外，输出这些内容时，模块还没有全部加载，所以 loaded 属性为 false 。

如果我们再建立一个 `b.js`，然后引用 `a.js`:

```javascript
// b.js

var a = require('./a');
```

再运行 b 时，输出为：

```
$ node b.js

module.id:  /home/ruanyf/tmp/a.js
module.exports:  {}
module.parent:  { object }
module.filename:  /home/ruanyf/tmp/a.js
module.loaded:  false
module.children:  []
module.paths:  [ '/home/ruanyf/tmp/node_modules',
  '/home/ruanyf/node_modules',
  '/home/node_modules',
  '/node_modules' ]
```

因为 a 被 b 调用，因此再打印时，a 的 parent 就指向了 b。

## require 方法

require 方法是定义在模块的原型链上的方法

```javascript
Module.prototype.require = function(path) {
  return Module._load(path, this);
};
```

因为 nodejs(CommonJS) 中每个文件都是模块，因此可能会误认为 require 是全局的。其实 require 是只在模块内才有的方法。换言之，如果在命令行直接输入 `node` 进入 REPL 环境(Read-eval-print-loop)，其中编写的代码就不是模块，也就没有 require 方法了。

`Module._load` 方法的源码如下：

```javascript
Module._load = function(request, parent, isMain) {

  //  计算绝对路径
  var filename = Module._resolveFilename(request, parent);

  //  第一步：如果有缓存，取出缓存
  var cachedModule = Module._cache[filename];
  if (cachedModule) {
    return cachedModule.exports;
  }

  // 第二步：是否为内置模块
  if (NativeModule.exists(filename)) {
    return NativeModule.require(filename);
  }

  // 第三步：生成模块实例，存入缓存
  var module = new Module(filename, parent);
  Module._cache[filename] = module;

  // 第四步：加载模块
  try {
    module.load(filename);
    hadException = false;
  } finally {
    if (hadException) {
      delete Module._cache[filename];
    }
  }

  // 第五步：输出模块的exports属性
  return module.exports;
};
```

## Module._resolveFilename

出现在上述所有步骤之前，用以计算模块所处的绝对路径。

```javascript
Module._resolveFilename = function(request, parent) {

  // 第一步：如果是内置模块，不含路径返回
  if (NativeModule.exists(request)) {
    return request;
  }

  // 第二步：确定所有可能的路径
  var resolvedModule = Module._resolveLookupPaths(request, parent);
  var id = resolvedModule[0];
  var paths = resolvedModule[1];

  // 第三步：确定哪一个路径为真
  var filename = Module._findPath(request, paths);
  if (!filename) {
    var err = new Error("Cannot find module '" + request + "'");
    err.code = 'MODULE_NOT_FOUND';
    throw err;
  }
  return filename;
};
```

这里又调用了 `Module._resolveLookupPaths` 方法，用来获取所有可能的路径。打印出来如下：

```
[   '/home/ruanyf/tmp/node_modules',
    '/home/ruanyf/node_modules',
    '/home/node_modules',
    '/node_modules'
    '/home/ruanyf/.node_modules',
    '/home/ruanyf/.node_libraries'，
     '$Prefix/lib/node' ]
```

基本上就是第一部分讲到的寻找 node_modules 的顺序。最后三个是出于历史原因还保留的，实际很少用到。

然后是 `Module._findPath` 方法，如下：

```javascript
Module._findPath = function(request, paths) {

  // 列出所有可能的后缀名：.js，.json, .node
  var exts = Object.keys(Module._extensions);

  // 如果是绝对路径，就不再搜索
  if (request.charAt(0) === '/') {
    paths = [''];
  }

  // 是否有后缀的目录斜杠
  var trailingSlash = (request.slice(-1) === '/');

  // 第一步：如果当前路径已在缓存中，就直接返回缓存
  var cacheKey = JSON.stringify({request: request, paths: paths});
  if (Module._pathCache[cacheKey]) {
    return Module._pathCache[cacheKey];
  }

  // 第二步：依次遍历所有路径
  for (var i = 0, PL = paths.length; i < PL; i++) {
    var basePath = path.resolve(paths[i], request);
    var filename;

    if (!trailingSlash) {
      // 第三步：是否存在该模块文件
      filename = tryFile(basePath);

      if (!filename && !trailingSlash) {
        // 第四步：该模块文件加上后缀名，是否存在
        filename = tryExtensions(basePath, exts);
      }
    }

    // 第五步：目录中是否存在 package.json
    if (!filename) {
      filename = tryPackage(basePath, exts);
    }

    if (!filename) {
      // 第六步：是否存在目录名 + index + 后缀名
      filename = tryExtensions(path.resolve(basePath, 'index'), exts);
    }

    // 第七步：将找到的文件路径存入返回缓存，然后返回
    if (filename) {
      Module._pathCache[cacheKey] = filename;
      return filename;
    }
  }

  // 第八步：没有找到文件，返回false
  return false;
};
```

这里的寻找顺序也和上面第一部分提到的相同，分别从文件和目录两个方面，依次添加扩展名进行查找。

这里额外提一句，nodejs 还对外提供了一个 `require.resolve` 方法，用以查找某个模块的绝对路径。

```javascript
require.resolve = function(request) {
  return Module._resolveFilename(request, self);
};

// 用法
require.resolve('a.js')
// 返回 /home/ruanyf/tmp/a.js
```

## Module._load

用于加载模块。

```javascript
Module.prototype.load = function(filename) {
  var extension = path.extname(filename) || '.js';
  if (!Module._extensions[extension]) extension = '.js';
  Module._extensions[extension](this, filename);
  this.loaded = true;
};
```

不同的后缀名有不同的加载方法，都记录在 `Module._extensions[extension]` 中。如下是 `.js` 和 `.json` 的方法：

```javascript
Module._extensions['.js'] = function(module, filename) {
  var content = fs.readFileSync(filename, 'utf8');
  module._compile(stripBOM(content), filename);
};

Module._extensions['.json'] = function(module, filename) {
  var content = fs.readFileSync(filename, 'utf8');
  try {
    module.exports = JSON.parse(stripBOM(content));
  } catch (err) {
    err.message = filename + ': ' + err.message;
    throw err;
  }
};
```

只看 js 的话，首先读取文件内容，然后 `stripBOM` 用于剥离 utf8 编码的 BOM 头，之后使用 `module._compile` 方法进行编译。

```javascript
Module.prototype._compile = function(content, filename) {
  var self = this;
  var args = [self.exports, require, self, filename, dirname];
  return compiledWrapper.apply(self.exports, args);
};
```

这段代码把 exports, require, module, filename 和 dirname 都插入到文件中，也就是我们最常理解的，用方法包住模块。因此等价于

```javascript
(function (exports, require, module, __filename, __dirname) {
  // 模块源码
});
```

综上，模块的加载实质上就是，注入exports、require、module三个全局变量，然后执行模块的源码，然后将模块的 exports 变量的值输出。

最后补充一个小点：nodejs 中 `exports` 和 `module.exports` 的差别：

1. 如果只设置 `exports`，则最终两者的值会相等，本质还是通过 `module.exports` 输出的。
2. 如果两者均设置，以 `module.exports` 为准。但这会降低代码可读性，所以最好不要这么写。
3. 如有多个方法或者属性输出，可以选择 `exports`；而如果是作为一个大对象整体输出，使用 `module.exports`。

例子如下：

```javascript
exports.sayHello = () => console.log('hello')
exports.sayBye = () => console.log('bye')
```

```javascript
module.exports = {
  sayHello() {
    console.log('hello')
  },
  sayBye() {
    console.log('bye')
  }
}
```