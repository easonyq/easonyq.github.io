# virtual dom

参考 DDFE Blog 的[相关文章](https://github.com/DDFE/DDFE-blog/issues/18)

vdom 是 react 的一个重要特性。vdom 是在 DOM 之前做了一层映射关系，本来直接操作 DOM 的代码改为操作 vdom。 vdom 完全使用 JS 实现，和浏览器无关，因此如创建，删除节点的操作也在 JS 内完成，速度比实际 DOM 更快一些。

Vue 2.0 引入了 vdom，采用的是 [snabbdom 算法](https://github.com/snabbdom/snabbdom)并进行修改。

## vnode 的属性定义

Vue 中 vnode 类包含的属性。完整列表可以在 `src/core/vdom/vnode.js` 中找到。其中比较重要的是

* tag 属性即这个 vnode 的标签属性
* data 属性包含了最后渲染成真实 DOM 节点后，节点上的 class, attribute, style 以及绑定的事件
* children 属性是 vnode 的子节点
* text 属性是文本属性
* elm 属性为这个 vnode 对应的真实dom节点
* key 属性是 vnode 的标记，在 diff 过程中可以提高 diff 的效率，后文有讲解

例如一个 vnode 的结构是

```json
{
    tag: 'div'
    data: {
        id: 'app',
        class: 'page-box'
    },
    children: [
        {
            tag: 'p',
            text: 'this is demo'
        }
    ]
}
```

它实际渲染的 DOM 将是

```html
<div id="app" class="page-box">
   <p>this is demo</p>
</div>
```

## 如何使用 vdom 更新 DOM

Vue 中使用 vdom，是从挂载开始的，即

```javascript
Vue.prototype._init = function () {
    // ...
    vm.$mount(vm.$options.el);
    // ...
}
```

追根溯源，使用的是 `src/core/instance/lifecycle.js` 的 `mountComponent(vm, el, hydrating)` 方法。这个方法简单来说，做了以下几个事情：

1. `vm.$el = el` 记录真实 DOM
2. `callHook(vm, 'beforeMount')`
3. 定义 `updateComponent = () => {vm._update(vm._render(), hydrating)};`
    这里 `vm._render()` 调用 render 函数，返回一个VNode。这部分代码在 vue 的 compile 阶段。此外在生成VNode的过程中，会动态计算 getter, 同时推入到 dep 里面。
4. `vm._watcher = new Watcher(vm, updateComponent, noop)` 用以注册 setter/getter，并调用 `updateComponent`。
5. 如果直接传入的 el 合法，则直接挂载，那么这里直接 `callHook(vm, 'mounted')`。判断依据是 `vm.$vode == null`。如果没有，那么等待后续手动调用 `vm.$mount`时再触发钩子。

vdom 的核心在于 `vm._update(VNode, hydrating)` 中。这里会把新传入的 VNode 和老的进行 `diff`，之后完成 DOM 的更新。这个方法定义在 `src/core/instance/lifecycle.js` 中。具体来说，它找出前一个 vm (存放到 `prevActiveInstance`) 和 它的 `_vnode` (存放到 `prevVnode`)。接着调用 `vm.$el = vm.__patch__(prevVnode, vnode)` 方法对两个 vnode 进行比较，以 patch 的形式打到 `prevVnode`，完成 DOM 更新的工作，返回这个真实 DOM 并记录到 `vm.$el` 中。

这个 `vm.__patch__` 是在 `src/platforms/web/runtime/index.js` 中添加进去的，方法的定义其实在 `src/core/vdom/patch.js` 最终 return 的方法(大约 631 行, `return function patch (oldVnode, vnode, hydrating, removeOnly, parentElm, refElm) {`)。这个方法的大致流程是：

1. 如果 `oldVnode` 为空，则表示为第一次执行的初始化操作，调用 `createElm` 直接创建。
2. 如果 `oldVnode` 和 `vnode` 是 __相同类型__ 的节点，则调用 `patchVnode` 进行补丁操作，在原有 DOM 基础上更新并返回。
    这里的 __相同类型__ 节点，具体包括 `key`, `tag`, `isComment` 是否相等，`data` 是否均有定义（不求相等）和 `input type` 是否一致。
3. 如果不是相同类型的节点，那么删除旧的 DOM，重新创建新的 DOM。

### patchVnode

首先调用 `vnode.data.hook.update` 对 `vnode.data` 进行更新。

其次开始进入各种判断，并采取不同的方法更新 vnode。`oldCh` 为 `oldVnode` 的子节点，`ch` 为 `Vnode` 的子节点：

1. 如果 `vnode` 没有文本节点的情况下，进入子节点的 diff；

    1. 当 `oldCh` 和 `ch` 都存在且不相同的情况下，调用 `updateChildren` 对子节点进行 diff；后面详述。

    2. 若 `oldCh` 不存在，`ch` 存在，首先清空 `oldVnode` 的文本节点，同时调用 `addVnodes` 方法将 `ch` 添加到 `elm` 真实 DOM 节点当中；相当于清空老节点的文字（如果有），把新节点的儿子都挂上去。

    3. 若 `oldCh` 存在，`ch` 不存在，则删除 `elm` 真实节点下的 `oldCh` 子节点；相当于和上面相反，清空老节点的儿子。（新节点没有文本也没有儿子，所以也不添加什么）

    4. 若 `oldCh` 和 `ch` 都不存在，因为 `vnode` 没有文本也没有儿子，那么就清空老节点的文本即可。

2. 如果 `vnode` 有文本节点但和老的不相等，即 `oldVnode.text !== vnode.text`，那么设置把新的文本设置上去即可。

其他几种情况是比较直观的，比较复杂的是 `updateChildren` 方法。这里 [参考文章](https://github.com/DDFE/DDFE-blog/issues/18) 写的不太清晰，我重新补充一下。

`updateChildren`，或者说整个 `patch` 方法的目的都是在计算，如何在旧的节点上打补丁（做修改）让它能够变成新的节点。这里就包括了移动顺序，添加，删除，修改等等。而 `updateChildren` 就是解决在新旧两个父亲节点都有儿子的情况。如何将旧的儿子经过一套操作变成新的儿子，是这个方法的关键。

这里举了个例子来演示最终的操作流程：

老的节点，儿子分别为

```
A B C D
```

新的节点，儿子分别为

```
B A D C F
```

首先是不带 key 的优化，看文章的图一步步走就可以。

随后讲到了 key 对整个过程的优化，道理其实也很简单。为每个 vnode 添加一个唯一的 key 用以标识，key 相同的节点就认为是**可能**相同的节点（至少说嫌疑比较大，这就是优化假设的基础）。

当进行头尾比较均无果之后，按照原先的流程，应当是把新节点的第一个儿子添加到老节点之前，新的 startIndex 进一位，进入下一次循环。但这里增加一步：再检查一下老节点的儿子中是否有 key 和新节点的这第一个儿子相同的。

如果有，检查是否**相同类型**，如果是，那当我们把新节点的第一个儿子添加到最前面了之后，这个节点其实肯定是重复的（举例来说，A B C D 和 B A D C F，当新的 B 被添加到最前端之后，老的那个 B 肯定是重复的）。原先这种重复情况算法并不 care，它会留到最后因为老节点数量多了而删除。但既然这里可以肯定重复，那在这里就删掉的话，可以使后续每次循环的次数都少1，**这就是优化所在**。如果不是**相同类型**，那就当没有找到重复节点（和在老节点找不到相同 key 一样对待），和之前一样简单添加后执行下一次循环即可。
