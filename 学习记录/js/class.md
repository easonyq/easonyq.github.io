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

补充说明：ES5 的继承，是先创造子类的实例对象 `this`，然后把父类的方法属性都加上去；ES6 反之，先创造父类的实例对象 `this`，再用子类的构造函数去修改它。所以子类构造函数要先调用 `super`。

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

最后注意

1. 两种 `super` 在执行过程中，`this` 是指向子类而非父类。因此如果在父类中使用了类似 `this.xxx = xxx`，这里其实是调用了子类的 `xxx` 而非父类。
2. `super` 在子类中使用时，如果写入 `super.xxx = xxx`，是写入子类，等价于 `this.xxx = xxx`；如果读取 `super.xxx`，是读取父类，等价于 `Parent.prototype.xxx`。

    ```javascript
    class A {
      constructor() {
        this.x = 1;
      }
    }

    class B extends A {
      constructor() {
        super();
        this.x = 2;
        super.x = 3; // 等价于 this.x = 3;
        console.log(super.x); // undefined
        console.log(this.x); // 3
      }
    }

    let b = new B();
    ```

TODO
