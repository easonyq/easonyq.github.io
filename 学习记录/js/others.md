# 目录

* `Object.create(null)` vs `{}`
* NodeList vs Array
* location.href 可能不添加历史记录
* debounce & throttle

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

# debounce & throttle

据说这两个都是函数式编程中常用的方法，它们也普遍集成在很多类库中(例如 lodash)或者出现在开发者自己写的 util 中(因为这俩实现起来并不长)。

两个都是用来处理频繁密集的操作的，但使用场景和效果略有不同：

## debounce

它的作用是不要让每次操作都真的执行，只记录最终状态，等待一定时间后，再执行一遍即可。举例来说，监听用户的输入事件 `keyup`，输入完成后发送 ajax 请求查询结果。我们并不需要每次 `keyup` 都去发消息，只要等用户输入完成再发就可以了。

它内部的实现是 `setTimeout` 和 `clearTimeout`。没等 `delay` 到，再次调用就会清空上一次的 timer，而设置一个新的，于是就阻止了不必要的执行。

```javascript
/**
 *
 * @param fn {Function}   实际要执行的函数
 * @param delay {Number}  延迟时间，单位是毫秒（ms）
 *
 * @return {Function}     返回一个“防反跳”了的函数
 */
function debounce(fn, delay) {

  // 定时器，用来 setTimeout
  var timer

  // 返回一个函数，这个函数会在一个时间区间结束后的 delay 毫秒时执行 fn 函数
  return function () {

    // 保存函数调用时的上下文和参数，传递给 fn
    var context = this
    var args = arguments

    // 每次这个返回的函数被调用，就清除定时器，以保证不执行 fn
    clearTimeout(timer)

    // 当返回的函数被最后一次调用后（也就是用户停止了某个连续的操作），
    // 再过 delay 毫秒就执行 fn
    timer = setTimeout(function () {
      fn.apply(context, args)
    }, delay)
  }
}
```

lodash 的 `_.debounce` 参数更多，作用更多，但核心就是这样。

## throttle

debounce 的作用是让频繁被调用的方法只在最后执行一次。而 throttle 的作用是以一个给定的频率（通常比调用频率低）来多次执行方法。这比较适用于需要跟随事件实时执行，但又不想执行太多次消耗性能的情况，例如游戏界面的渲染。

throttle 的实现核心是要记录上一次的执行时间，然后判断当前时间是否应该执行。此外还应该和 debounce 一样维持一个 timer，确保最后额外执行一次作为最终的结果。

```javascript
/**
*
* @param fn {Function}   实际要执行的函数
* @param delay {Number}  执行间隔，单位是毫秒（ms）
*
* @return {Function}     返回一个“节流”函数
*/

function throttle(fn, threshhold) {

  // 记录上次执行的时间
  var last

  // 定时器
  var timer

  // 默认间隔为 250ms
  threshhold || (threshhold = 250)

  // 返回的函数，每过 threshhold 毫秒就执行一次 fn 函数
  return function () {

    // 保存函数调用时的上下文和参数，传递给 fn
    var context = this
    var args = arguments

    var now = +new Date()

    // 如果距离上次执行 fn 函数的时间小于 threshhold，那么就放弃
    // 执行 fn，并重新计时
    if (last && now < last + threshhold) {
      clearTimeout(timer)

      // 保证在当前时间区间结束后，再执行一次 fn
      timer = setTimeout(function () {
        last = now
        fn.apply(context, args)
      }, threshhold)

    // 在时间区间的最开始和到达指定间隔的时候执行一次 fn
    } else {
      last = now
      fn.apply(context, args)
    }
  }
}
```

附上一个图形化的[例子](http://demo.nimius.net/debounce_throttle/)，非常直观！
