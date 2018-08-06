# Object.create(null) vs {}

```javascript
let o = Object.create(null);
let o2 = {}
```

两者的区别是，`o` 不继承任何东西，是一个**真的**空对象；而 `o2` 继承了 `Object` 的所有方法，因此 `o2.constructor.prototype === Object.prototype`。因此 `o2` 等价于 `Object.create(Object.prototype)`。

说到应用场景，如果一个对象被用作 Map 来存储映射关系的，尽量使用前者。因为使用后者的话，在使用 `for (var key in o)` 时，还需要额外使用 `if (o.hasOwnProperty(key))` 来过滤那些我们并没有手动添加，而是继承自 `Object` 的 key。当然如果使用 `Object.keys(o).forEach` 也能绕开这个问题。

# NodeList vs Array

我们常用的方法 `document.querySelectorAll` 的返回值并不是 Array，而是 NodeList。这个 NodeList 只拥有 __部分__ Array 的方法，并不是全部。

常用的可用方法有：

* length
* forEach
* [i] (用下标访问)

常用的不可用方法有：

* map
* reduce
* find
* filter

所以在使用时必须要注意。有一个比较简便的转化方法是借用 es6 的 `...` 操作符，如 `[...document.querySelectorAll()]` 就是一个常规的数组了。不用 es6 的话，也可以使用 `Array.prototype.slice.call(document.querySelectorAll())`。

与 NodeList 类似的还有 HTMLCollection （由 `dom.children` 返回）

# location.href 可能不添加历史记录

按照常理来说，`location.href = xxx` 和 `location.assign(xxx)` 都能够进行跳转，并且把历史记录添加到 history 中；而 `location.replace(xxx)` 则在跳转后不添加历史记录，类似于 `history.pushState` 和 `history.replaceState` 的情况。

但如果当页面 load 事件 __尚未触发时__ 就调用 `location.href = xxx` 或者 `location.assign(xxx)`，浏览器也不会将当前页面添加到历史记录中（猜测可能浏览器是在 load 事件之后才添加历史记录的）。

实际上从业务角度来说，如果某个页面 A 一进入立马跳转到页面 B，那么 A 的历史记录不添加到 history 也并非不合理。否则 B 后退到 A 又立马跳转到 B，相当于在 B 原地循环，也很奇怪。

如果非要解决，方法有两种：

1. 调用之前添加 `setTimeout`。（比较简单）
2. 检查 `load` 事件是否触发了。如果触发了则直接调用，还没触发就把跳转包在 `window.onload` 里面。（比较保险）
