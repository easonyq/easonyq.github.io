# call & bind & apply

一般如果面试 JavaScript，新手总会倒在大杀器 `this` 的面前。凡是面到 `this`，大概率会提到 `call`, `bind`, `apply` 这三个方法。此三者都是 Function 原型链上定义的方法，也就是每个方法或者类（类也是由方法实现的）都有的方法，所以我们应该详细了解一下，即便自己不使用，也不至于看不懂或者犹豫不决。

## this - 好用又难用的前端分界线

作为面向对象语言，`this` 是必不可少的。但能否正确使用 `this` 是区分一个前端新手和中手的差别，因为它其实还是有很多文章可做的（所以可以常年用作面试题）。

我们需要先理解一些大体的原则：

1. **`this` 永远指向一个函数的调用者**，如果没有，那就是 `window`。此外在最外层直接使用 `this` 也等价于 `window`。
2. 基本上每定义一个方法，就会新建一个 `this`。所以在编写 `function` 关键词的时候，心里就要提醒自己一下：`this` 变了。
3. W3C 也意识到第二点带来的不便，为了减少大家在定义函数前常规的 `let me = this` 这样的操作，现在推出了不会修改 `this` 的箭头函数，用起来更放心。不过偶尔可能也会有反作用。

第三点是跟随最近几年 ES2015 等新标准的推行以及浏览器的支持而逐渐普及的。但在早先没有箭头函数的时期，我们通常是使用第二点提到的方法来保存 `this`，供后续使用。那么除此之外，是否有更加优雅（或者高级）的方法来解决呢？

这就牵涉到了 `call`, `bind` 和 `apply` 这三个方法。

## bind

官方文档对 `bind` 方法的解释如下：它创造一个新的函数，且这个函数中的 `this` 指向一个特定的值（The `bind()` method creates a new function that, when called, has its this keyword set to the provided value）。所以首先一点，它**返回**一个新的方法，**并不直接调用**。这是它和 `call` 及 `apply` 最大的区别。

### 固定 this

按照定义，`bind` 就是用来固定 `this` 的，代码如下：

```javascript
class MyForm {
    constructor() {
        this.button = xxx;
        this.formName = 'MyForm';
    }
    bindEvents() {
        // 错误示范
        this.button.addEventListener('click', function () {
            console.log(this.formName) // undefined
            // 因为 function 改变了 this，这里的 this 不再指向 MyForm 实例，而是 this.button 这个 DOM 元素了。
        })

        // 正确示范
        let clickHandler = function () {
            console.log(this.formName) // MyForm
        }
        // 固定 this 为 MyForm 实例。
        this.button.addEventListener('click', clickHandler.bind(this))
    }
}
```

经过 `bind` 的固定，`clickHandler.bind(this)` 的返回值无论如何调用，都可以保证 `this` 就指向实例，也就确保了 `this.formName` 是我们预期的值。这种写法在 react 中绑定事件比较常见，也是 `bind` 最本质的用法。

### 借用方法

在实际使用中，大家还开发出了 `bind` 的第二个用法，叫做“借用方法”。

```javascript
class A {
    showName() {
        console.log(this.name);
    }
}
class B {
    constructor() {
        this.name = 'b\'s name';
    }
}

let a = new A();
let b = new B();
let lendShowName = a.showName.bind(b);
lendShowName(); // b's name
```

`b` 实例并没有 `showName` 方法，但它有 `name` 属性。正巧 `a` 实例有这个方法，于是就“借用”过来使用。因为 `bind` 返回之后还需要再次调用，所以更常规的做法是使用 `call` 来直接借用并调用，这在后续会提到。

### 固定参数 & 科里化 (curry)

`bind` 还有第三种用法，叫做“固定参数”。通过例子看会比较直观，如下：

```javascript
function add(a, b) {
    return a + b;
}

console.log(add(1, 2)) // 3

let increase = add.bind(null, 1);
console.log(increase(2)) // 3
console.log(increase(3)) // 4
```

通过 `bind` 的后续参数，可以将被绑定方法的前几个参数固定。例如例子中的 `add` 的第一个参数被固定为 `1` 之后，它就变成了 `increase` 方法，只需要传入一个参数，就可以实现累加。

这个思路和函数式编程中的科里化 (curry) 是一致的。科里化要求每个函数只接一个参数，因此接多个参数的函数会拆分成多个函数，`bind` 就是实现的方法之一。

额外提一句，经过 `bind` 之后的函数，如果给的参数超过限制，则会被忽略。因此如果调用 `increase(3, 4, 5)` 效果等价于 `increase(3)`。

## call & apply

`call` 和 `apply` 比较类似，因此放在一起说。这两个方法在调用之后，都是**立即执行**的，这也是和 `bind` 最大的差别。

`call` 和 `apply` 的定义都是执行一个函数，并且它的 `this` 指向特定的值。这两者的区别在于在传递后续参数是，`call` 需要单个单个传入（以逗号分隔），而 `apply` 需要放到一个数组里面，如下：

```javascript
function add(a, b) {
    return a + b;
}

console.log(add.call(null, 1, 2)); // 3
console.log(add.apply(null, [1, 2])); // 3
```

因为这两个方法和 `bind` 本意都是用来调整 `this` 的，因此 `bind` 的前面两种用法（固定 `this` 和借用方法）都是可以用这两个方法来实现的。

下面这段代码是早期没有引入 `class` 关键词时候的继承写法（不过其实 `class` 也只是语法糖）。通过 `Parent.call(this)`，把父类的方法和属性全部借给子类，所以从面向对象的角度来看，父子关系就成立了。

```javascript
var Parent = function () {
    this.sayHello = function () {
        console.log('hello');
    }
};

var Child = function () {
    Parent.call(this);
};

let c = new Child();
c.sayHello(); // hello
```

之前在 `bind` 列出过”借用方法“的例子。不过实际上就“借用方法”来说，用 `call` 或者 `apply` 更加常见，比如我们最最常用的**转换成数组**：

```javascript
let contentArr = document.querySelector('.content');

// 错误示范：contentArr 是一个类数组，而不是真的数组，因此它其实没有 map 方法
contentArr.map();

// 正确示范：借用 Array 原型链的 slice 方法，不传参数来转化为数组。这里 call 和 apply 都可以。
Array.prototype.slice.call(contentArr).map();
```

此外，基于 `apply` 的后续参数是数组的特性，就可以实现动态个数的参数。如下：

```javascript
let numberArr = [1, 2, 3, 4, 5];

let getSum = function () {
    let sum = 0;
    for (let i = 0; i < arguments.length; i++) {
        sum += arguments[i];
    }
    return sum;
}

// 在不知道 numberArr 长度的情况下，用 apply 会更加合理。
console.log(getSum.apply(null, numberArr));
```

最后，因为 `call` 和 `apply` 都是直接调用的，因此不存在 `bind` 的第三种用法。如果参数给多了，超过部分就会被忽略；如果参数给的不够，就被当做 `undefined` 处理。