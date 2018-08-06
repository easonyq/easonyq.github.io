# pointer-events & 穿透

这是一个 CSS3 的属性，当值为 `'none'` 时可以用来控制元素强制“穿透”。

比如我们想在页面上显示一个浮层，但又不想这个浮层干扰其他正常元素的点击或者滚动，就可以设置这个属性来强制穿透。

一个注意点是，如果一个元素设置了这个属性之后，监听这个元素的 `click`, `touchmove` 等事件均不会触发，因为它被穿透了。

使用 js 设置时，key 为 `pointerEvents`。其他可用值（多用在 svg ）可以参考 [MDN](https://developer.mozilla.org/zh-CN/docs/Web/CSS/pointer-events)
