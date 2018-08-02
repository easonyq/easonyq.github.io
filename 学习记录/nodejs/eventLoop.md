# 事件循环

nodejs的事件驱动模型一般要注意下面几个点：

* 因为是单线程的，所以当顺序执行js文件中的代码的时候，事件循环是被暂停的。
* 当js文件执行完以后，事件循环开始运行，并从消息队列中取出消息，开始执行回调函数
* 因为是单线程的，所以当回调函数被执行的时候，事件循环是被暂停的
* 当涉及到I/O操作的时候，nodejs会开一个独立的线程来进行异步I/O操作，操作结束以后将消息压入消息队列。

因此，只有在同步代码执行完成后才会进入异步回调的进行。而多个异步回调的先后顺序如何确定就是事件循环要解决的问题。从这里可以看出，如果主进程代码（同步代码）卡住，异步代码也是不会执行的。时此外事件循环中列出的所有顺序都只针对异步回调。

事件循环机制是使单线程的 JavaScript 支持高性能非阻塞 I/O 操作的原因。当 Node.js 启动的时候就会初始化一个事件循环，并开始执行 js 主代码，这中间可能会产生一些定时器（schedule timers），异步 I/O API调用，或者process.nextTick调用等等， 然后进入事件循环。

   ┌───────────────────────┐
┌─>│        timers         │ 这个阶段执行 `setTimeout()` 和 `setInterval()` 中的回调函数
│  └──────────┬────────────┘
│  ┌──────────┴────────────┐
│  │     I/O callbacks     │ 这个阶段执行除了 `close` 回调函数以外的几乎所有的 I/0 回调函数
│  └──────────┬────────────┘
│  ┌──────────┴────────────┐
│  │     idle, prepare     │ 这个阶段仅仅 Node.js 内部使用
│  └──────────┬────────────┘      ┌───────────────┐
│  ┌──────────┴────────────┐      │   incoming:   │
│  │         poll          │<─────┤  connections, │ 执行队列中的回调函数、检索新的回调函数
│  └──────────┬────────────┘      │   data, etc.  │
│  ┌──────────┴────────────┐      └───────────────┘
│  │        check          │ `setImmediate()` 将在这里被调用
│  └──────────┬────────────┘
│  ┌──────────┴────────────┐
└──┤    close callbacks    │ `close` 回调函数被调用如：socket.on('close', ...)
   └───────────────────────┘

## timers

setTimeout() 和 setInterval() 都要指定一个运行时间，这个运行时间其实不是确切的运行时间，而是一个期望时间，Event Loop 会在 timers 阶段执行超过期望时间的定时器回调函数，但由于你不确定在其他阶段甚至主进程中的事件执行时间，所以定时器不一定会按时执行。

```javascript
var asyncApi = function (callback) {
  setTimeout(callback, 90)
}

const timeoutScheduled = Date.now();
setTimeout(() => {
  const delay = Date.now() - timeoutScheduled;
  console.log(`${delay}ms setTimeout 被执行`); // 140ms 之后被执行
}, 100);

asyncApi(() => {
  const startCallback = Date.now();
  while (Date.now() - startCallback < 50) {
    // do nothing
  }
})
```

## I/O callbacks

这个阶段主要执行一些系统操作带来的回调函数，如 TCP 错误，如果 TCP 尝试链接时出现 ECONNREFUSED 错误 ，一些 *nix 会把这个错误报告给 Node.js。而这个错误报告会先进入队列中，然后在 I/O callbacks 阶段执行。

## poll

poll 阶段有两个主要功能：

1. 也会执行时间定时器到达期望时间的回调函数
2. 执行事件循环列表（poll queue）里的函数

当 Event Loop 进入 poll 阶段并且没有其余的定时器，那么：

1. 如果事件循环列表不为空，则迭代同步的执行队列中的函数。
2. 如果事件循环列表为空，则判断是否有 setImmediate() 函数待执行。如果有结束 poll 阶段，直接到 check 阶段。如果没有，则等待回调函数进入队列并立即执行。

## check

在 poll 阶段结束之后，执行 setImmediate()。

## close

突然结束的事件的回调函数会在这里触发，如果 socket.destroy()，那么 close 会被触发在这个阶段，也有可能通过 process.nextTick() 来触发。

## setImmediate()、setTimeout()、process.nextTick()

这里要说明一下 process.nextTick() 是在下次事件循环之前运行，如果把 process.nextTick() 和 setImmediate() 写在一起，那么是 process.nextTick() 先执行。next 比 immediate 快，官方也说这个函数命名有问题，但是因为历史存留没办法解决。

```javascript
process.nextTick(() => {
  console.log('nextTick');
});
setImmediate(() => {
  console.log('setImmediate');
});
setTimeout(() => {
  console.log('setTimeout'); 
}, 0)

// 执行结果，nextTick, setTimeout, setImmediate
// 查看 Node.js 源码，setTimeout(fun, 0) 会转化成 setTimeout(fun, 1)，所以在这种简单的情况下，对于不同设备，setImmediate 有可能早于 setTimeout 执行。
```

## 从event loop机制的角度上区分setImmediate()与setTimeout()

从poll和check阶段的逻辑，我们可以看出setImmediate和setTimeout、setInterval都是在poll 阶段执行完当前的I/O队列中相应的回调函数后触发的。但是这两个函数却是由不同的路径触发的。

* setImmediate函数，是在当前的poll queue对列执行后为空或是执行的数目达到上限后，event loop直接调入check阶段执行setImmediate函数。
* setTimeout、setInterval则是在当前的poll queue对列执行后为空或是执行的数目达到上限后，event loop去timers检查是否存在已经到期的定时器，如果存在直接执行相应的回调函数。

如果程序中既有setTimeout和setImmediate，两者的执行顺序是什么？

```javascript
// timeout_vs_immediate.js
setTimeout(function timeout() {
  console.log('timeout');
}, 0);

setImmediate(function immediate() {
  console.log('immediate');
});
```

上面的程序执行的结果并不是唯一的，有时immediate在前，有时timeout在qian。主要是由于他们运行的当前上下文环境中存在其他的程序影响了他们执行顺序。

```javascript
// timeout_vs_immediate.js
const fs = require('fs');

fs.readFile(__filename, () => {
  setTimeout(() => {
    console.log('timeout');
  }, 0);
  setImmediate(() => {
    console.log('immediate');
  });
});
```

上面的程序把setImmediate和setTimeout放到了I/O循环中，此时他们的执行顺序永远都是immediate在前，timeout在后。

## 从event loop机制的角度上区分process.nextTick()与setImmediate()

尽管process.nextTick()也是一个异步的函数，但是它并没有出现在上面event loop的结构图中。不管当前正在event loop的哪个阶段，在当前阶段执行完毕后，跳入下个阶段前的瞬间执行process.nextTick()函数。
由于process.nextTick()函数的特性，很可能出现一种恶劣的情形：在event loop进入poll前调用该函数，就会阻止程序进入poll阶段
但是也正是nodejs的一个设计哲学：每个函数都可以是异步的，即使它不必这样做。例如下面程序片段，如果不对内部函数作异步处理就可能出现异常。

```javascript
let bar;

// this has an asynchronous signature, but calls callback synchronously
function someAsyncApiCall(callback) { callback(); }

// the callback is called before `someAsyncApiCall` completes.
someAsyncApiCall(() => {

  // since someAsyncApiCall has completed, bar hasn't been assigned any value
  console.log('bar', bar); // undefined

});

bar = 1;
```

由于someAsyncApiCall函数在执行时，内部函数是同步的，这是变量bar还没有被赋值。如果改为下面的就会正常。

```javascript
let bar;

function someAsyncApiCall(callback) {
  process.nextTick(callback);
}

someAsyncApiCall(() => {
  console.log('bar', bar); // 1
});

bar = 1;
```

process.nextTick() 函数是在任何阶段执行结束的时刻
setImmediate() 函数是在poll阶段后进去check阶段事执行
