# class

参考 [ES6 入门](http://es6.ruanyifeng.com/#docs/class)，但进行了一定的精简。

## prototype & __proto__

单类情况：

```javascript
class Point {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }

    toString() {
        console.log('x=', x, 'y=', y);
    }
};

let p1 = new Point(1, 2);

p1.constructor === Point; // true。constructor 指向构造类
p1 instanceof Point; // true

Object.getOwnPropertyNames(Point.prototype); // ['constructor', 'toString']。 说明方法全部定义在原型链上
Object.keys(Point.prototype); // []。类原型链上的方法不可枚举

p1.hasOwnProperty('toString'); // false。类的方法定义在类的原型链上，并不在实例上。
p1.hasOwnProperty('x'); // true。通过构造函数中的 this 定义的就在实例上。
p1.constructor.prototype.hasOwnProperty('toString'); // true

p1.constructor.prototype === p1.__proto__; // true

let p2 = new Point(2, 3);

// 也可写作 Point.prototype.sayHello = xxx
p1.__proto__.sayHello = () => console.log('hello!');

p1.sayHello(); // hello!
p2.sayHello(); // hello! 定义在原型链上的方法，实例全部可以使用。也因此要慎用，因为会影响所有实例。

```

继承情况：

1. 子类的 `__proto__` 属性，表示构造函数的继承，总是指向父类。

2. 子类 `prototype` 属性的 `__proto__` 属性，表示方法的继承，总是指向父类的prototype属性。

这里的两种情况说的都是类，并不是上面的实例。

```javascript
class A {
}

class B extends A {
}

B.__proto__ === A // true
B.prototype.__proto__ === A.prototype // true
```

## this

类的内部 `this` 指向实例，如上的 `this.x = x` 就是一个例子。但使用解构会导致 `this` 指向 `window` 从而报错。

```javascript
class Logger {
  printName(name = 'there') {
    this.print(`Hello ${name}`);
  }

  print(text) {
    console.log(text);
  }
}

const logger = new Logger();
const { printName } = logger;

printName(); // TypeError: Cannot read property 'print' of undefined, window 并没有 print 方法
logger.printName(); // ok
```

解决方案有三种。（实际上这种问题很少见，大多数情况还是使用 `logger.printName()` 才是正统）

使用 `bind`：

```javascript
class Logger {
  constructor() {
    this.printName = this.printName.bind(this);
  }

  // ...
}
```

使用箭头函数：

```javascript
class Logger {
  constructor() {
    this.printName = (name = 'there') => {
      this.print(`Hello ${name}`);
    };
  }

  // ...
}
```

使用 Proxy 过于复杂，这里不列了。

## 静态方法

静态方法是指定义在类上而非原型链上的方法。调用时使用类来调用；而非实例。如下：

```javascript
class Foo {
  static classMethod() {
    return 'hello';
  }
}

Foo.classMethod() // 'hello'

var foo = new Foo();
foo.classMethod(); // TypeError: foo.classMethod is not a function
```

静态方法中，`this` 指向类而不是实例。静态方法会被子类继承。静态方法甚至可以和普通方法重名。

## 静态属性

静态属性是指类本身的属性，即 `Class.propName`，并不是定义在实例 (`this`) 上的属性。如：

```javascript
class Foo {
}

Foo.prop = 1;
Foo.prop // 1
```

但是注意，在类的内部是**不能**定义静态属性的：

```javascript
// 以下两种写法都无效
class Foo {
  // 写法一
  prop: 2

  // 写法二
  static prop: 2
}

Foo.prop // undefined
```

类内部的静态属性还在提案阶段，这里就不展开了。

## 私有变量/方法

方法一：

```javascript
class Widget {
  foo (baz) {
    bar.call(this, baz);
  }

  // ...
}

function bar(baz) {
  return this.snaf = baz;
}
```

`bar` 实际上是私有的，在类的外部无法使用。解决方案就是把它放到类的声明之外，因为类的内部都是可以访问的。

方法二：

```javascript
const bar = Symbol('bar');
const snaf = Symbol('snaf');

export default class myClass{

  // 公有方法
  foo(baz) {
    this[bar](baz);
  }

  // 私有方法
  [bar](baz) {
    return this[snaf] = baz;
  }

  // ...
};
```

利用 `Symbol` 值的唯一性，在外部无法得知 `bar` 的值，因此也就无法调用了。

## 继承和 super

super 有两种用法。第一，在子类的构造函数中，**必须**调用以建立子类的 `this`。

补充说明：ES5 的继承，是先创造子类的实例对象 `this`，然后把父类的方法属性都加上去；ES6 反之，先创造父类的实例对象 `this`，再用子类的构造函数去修改它。所以子类构造函数要先调用 `super`，作用等价于把父类 `constructor` 的代码拿过来执行一遍。

```javascript
class ColorPoint extends Point {
  constructor(x, y, color) {
    super(x, y); // 调用父类的constructor(x, y)
    this.color = color;
  }

  toString() {
    return this.color + ' ' + super.toString(); // 调用父类的toString()
  }
}
```

第二，当成一个对象，使用 `super.xx()` 来调用父类实例的方法。（或者父类的静态方法）

```javascript
class A {
  p() {
    return 2;
  }
}

class B extends A {
  constructor() {
    super();
    console.log(super.p()); // 2
  }
}

let b = new B();
```

`super.p()` 指向父类的原型对象，等价于 `A.prototype.p()`，因此返回 `2`。

**注意**，`super` 指向的是父类的**原型对象**而不是实例，因此如果父类定义的是实例的属性

```javascript
class A {
  constructor() {
    this.p = 2;
  }
}
```

那么 `super.p()` 是访问不到的。

```javascript
class A {
  constructor() {
    this.x = 1;
  }
  print() {
    console.log(this.x);
  }
}

class B extends A {
  constructor() {
    super();
    this.x = 2;
  }
  m() {
    super.print();
  }
}


let b = new B();
b.m() // 2
```

在把 `super` 当做对象，调用 `super.xxx` 方法或者属性时：
1. 执行的是父类的 __原型链__ 上面的方法或者属性。因此如上面写法，`super.x` 获取到的是 `undefined`，因为父类的 `x` 是实例的属性而非原型链的属性。

    因此，如果写成 `A.prototype.x = 1` 就可以获取到了。

2. 执行时，`this` 是指向子类的。因此执行 `super.print()` 时，打印的 `this.x` 是子类的 `x`，所以打印 `2` 而不是 `1`。
3. 如果使用 `super` 进行赋值（估计只有面试题这么做）等价于 `this`。如 `super.x = 3` 等价于 `this.x = 3`。
4. 如果在静态方法中调用 `super.xxx`，指向的是父类。因此如果调用方法，那就会去找父类的静态方法。在这个过程中如果使用了 `this`，它指向子类。所以 `this.xxx` 就会访问子类的静态变量或者方法。

## __proto__

### 类

```javascript
class A {
}

class B extends A {
}

B.__proto__ === A // true
B.prototype.__proto__ === A.prototype // true
```

1. 子类的 `__proto__` 属性，表示构造函数的继承，总是指向父类。

2. 子类 `prototype` 属性的 `__proto__` 属性，表示方法的继承，总是指向父类的 `prototype` 属性。

实际上继承的内部也是通过设置 `__proto__` 属性来实现的。

### 实例

最最基础和根本的原则：**实例的 `__proto__` 等于实例构造类的原型链。**

```javascript
p1.__proto__ === p1.constructor.prototype === Point.prototype; // true
```

1. 子类实例的 `__proto__` 等于子类构造类的原型链，也就是 `Child.prototype`。
2. `Child.prototype` 的 `__proto__` 等于父类的原型链，也就是 `Parent.prototype`。
3. 父类实例的 `__proto__` 等于父类构造类的原型链，也是 `Parent.prototype`。

所以 `child.__proto__.__proto__ === parent.__proto__;` 成立。


```javascript
let parent = new Parent();
let child = new Child();

child.__proto__.__proto__ === parent.__proto__;
```

## 特殊的继承关系

### 继承 `Object`

```javascript
class A extends Object {
}

A.__proto__ === Object // true
A.prototype.__proto__ === Object.prototype // true
```

### 不继承

```javascript
class A {
}

A.__proto__ === Function.prototype // true
A.prototype.__proto__ === Object.prototype // true， 特殊点
```

A 不继承任何父类，他就是个普通函数，所以继承 `Function.prototype`。

但调用 `new A()` 之后返回的是一个空对象，因此 A 的原型链继承对象的原型链。

### 继承 `null`

```javascript
class A extends null {
}

A.__proto__ === Function.prototype // true
A.prototype.__proto__ === undefined // true
```

A 的情况和第二种一样，是个普通函数，继承 `Function.prototype`。

而 A 的实例实质上执行了 `return Object.create(null)`， 因此也没有父类。

### 继承源生类

源生类就是如 Object, Array, Boolean, Number, String 之类的内置的类。ES6 也可以继承这些类，获得他们相同的行为和方法。

```javascript
class MyArray extends Array {
  constructor(...args) {
    super(...args);
  }
}

var arr = new MyArray();
arr[0] = 12;
arr.length // 1

arr.length = 0;
arr[0] // undefined
```
