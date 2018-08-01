# Object.create(null) vs {}

```javascript
let o = Object.create(null);
let o2 = {}
```

两者的区别是，`o` 不继承任何东西，是一个**真的**空对象；而 `o2` 继承了 `Object` 的所有方法，因此 `o2.constructor.prototype === Object.prototype`。因此 `o2` 等价于 `Object.create(Object.prototype)`。

说到应用场景，如果一个对象被用作 Map 来存储映射关系的，尽量使用前者。因为使用后者的话，在使用 `for (var key in o)` 时，还需要额外使用 `if (o.hasOwnProperty(key))` 来过滤那些我们并没有手动添加，而是继承自 `Object` 的 key。当然如果使用 `Object.keys(o).forEach` 也能绕开这个问题。
