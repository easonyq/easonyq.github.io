# BFC

Blocking Fomatting Context 块级格式化上下文。

CSS内的所有内容都是 Box, BFC 用来定义如何渲染块级元素内部的所有节点。原则是 __内外互不影响__。

## BFC 的生成

要生成 BFC，需要满足一下 __任意一个__ 条件：

* 根元素 (`<body>`)
* float 不为 none
* overflow 不为 visible
* position 不为 static
* display 为 inline-block, table-cell, table-caption
* flex boxes (display 为 flex 或者 inline-flex)

## BFC 的布局规则

简单归纳如下：
1. 内部的元素会在 __垂直方向__ 一个接一个地排列，可以理解为是BFC中的一个常规流（因为块级元素默认占满父元素100%宽度）
2. 元素垂直方向的距离由margin决定，即属于同一个BFC的两个相邻盒子的margin可能会 __发生重叠__
3. 每个元素的左外边距与包含块的左边界相接触(从左往右，否则相反)，即使存在浮动也是如此，这说明BFC中的子元素不会超出它的包含块
4. BFC的区域不会与float元素区域重叠
5. 计算BFC的高度时，浮动子元素也参与计算
6. BFC就是页面上的一个隔离的独立容器，容器里面的子元素不会影响到外面的元素，反之亦然

## BFC 的应用

其实以往很多的问题 （清除浮动，margin折叠等）的原因及其解决方案都是 BFC 。

### margin 折叠

相邻的垂直元素同时设置了margin后，实际margin值会塌陷到其中较大的那个值。

这其实是 BFC 的布局规则，要解决这个问题，让这两个元素不处在同一个 BFC 中即可。

例如
```html
<body>
    <p>haha</p>
    <p>haha</p>
    <p>haha</p>
</body>
```

中 3个 `<p>` 处在同一个 BFC (`<body>`) 中，如果让他们包裹在一个 `<div style="overflow: hidden">` 中，就不会让他们处在同一个 BFC 中，也就不会折叠了。

### 清除浮动

overflow: hidden 能够清除浮动，原因就在于 __BFC 计算高度会考虑浮动子元素__，即撑开。

### 解决浮动脱离文档流的问题

浮动元素会脱离文档流单独渲染，剩余的非浮动元素则围绕浮动元素进行渲染（如果剩余元素的高度大于浮动元素的话）。想象一下文字混排的效果，如下：

```
xxxx yyyyyyyyyyy
xxxx yyyyyyyyyyy
xxxx yyyyyyyyyyy
yyyyyyyyyyyyyyyy
yyyyyyyyyyyyyyyy
```

其中 x 为浮动元素， y为非浮动元素（事实上文字围绕排列也就是用浮动实现的）

如果不想让 y 高出 x 的部分从最左边开始，而是留出距离，类似如下效果的话

```
xxxx yyyyyyyyyyy
xxxx yyyyyyyyyyy
xxxx yyyyyyyyyyy
     yyyyyyyyyyy
     yyyyyyyyyyy
```

给 y 加上 overflow: hidden 或者 float: left 让它变成 BFC，就可以做到这个效果。

__这个效果在两栏布局时非常有效__，因为比起传统的 position: absolute; left: xxxpx + margin-left: xxxpx 的方式来说，这种方式并不需要知道左栏的宽度（xxx），是真正的自适应。