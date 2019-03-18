# iframe BUG 集锦和解决方案

最近几个月在开发 MIP 的过程中专注于使用 iframe 将多页融合为单页，因此碰上了 iframe 众多奇奇怪怪的问题，通过网上查阅，同事帮忙，自己乱改，了解/解决了许多问题，特此记录。

参考于同事的博客: [在 iOS 下使用 iframe 的种种问题](https://xiaoiver.github.io/coding/2018/05/20/%E5%9C%A8-iOS-%E4%B8%8B%E4%BD%BF%E7%94%A8-iframe-%E7%9A%84%E7%A7%8D%E7%A7%8D%E9%97%AE%E9%A2%98.html)

## 滚动

当页面嵌入 iframe 之后，就有两种滚动方案，分别是：

1. 在 body 上滚动，默认情况。可以把 iframe 看成一个 div，把高度设到足够防止它自身出滚动条。浏览器会对 body 滚动进行优化，例如滚动时把上面的标题栏，下面的菜单栏给隐藏掉，让可视区域变大。

  在这个方案下，iframe 一定是和 body 一起滚动的，因此一定是设置 `absolute, top:0, left: 0`，这时切换动画就会变得有点困难，例如 A 页面滚动到 100px 然后需要切换 B，为了切换效果， B 一定要设置为 `top: 100px`，那 B 往上的 100px 的空白就难了。

  切换动画总能想办法解决，另外还有个问题就无法解决了。如果 iframe 平铺，那么它永远不会存在滚动，因此如 `scrollTop` 就永远是 `0`，那么被 iframe 容纳的页面中要利用这个变量来做一些事情就会难以实现，如下拉刷新。要实现下拉刷新，就必须在第一个页面中做，但这也导致了两种代码不一致，需要横加判断，不利于代码可读性和精简。

2. 在 iframe 内部滚动。把 iframe 进行绝对定位(fixed, absolute 都可以)，然后高度设置为视窗高度，内部滚动，如下：

  ```html
  <html>
  <head>
    <title>I’m a Web App and I show AMP documents</title>
    <style>
      iframe {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
      }
    </style>
  </head>
  <body>
    <iframe width="100%" height="100%"
      scrolling="yes"
      src="https://cdn.ampproject.org/c/pub1.com/doc1"></iframe>
  </body>
  </html>

  ```

注意我们没有应用 CSS overflow 属性，而是直接使用 iframe 的 scrolling 属性。这个属性已经被 HTML5 规范废弃，但是由于历史原因，很多浏览器还是支持。

但是 iOS 不支持这个属性，详见 [Bug Report](https://bugs.webkit.org/show_bug.cgi?id=149264) 以及 [Online Demo](http://output.jsbin.com/dedega)。

抛开 iOS 的问题，还有一个子方案是在 iframe 滚动时通知外层 body，用脚本进行滚动。但经验证明，这样太卡，不适合做方案，所以只是提一下。

既然我们确定了 iframe 内部滚动，我们就要解决 iOS 的问题。

### iOS iframe的滚动

这里我们主要参考的是 Google AMP 的做法。他们有过2个方案，因为 iOS7 的兼容性问题 （不允许 `Object.defineProperty(document)` ），我们采用了所谓 “原始方案”。 如果抛弃 iOS7，可以参考 AMP 的进阶方案（2个 `<html>` 标签），[传送门](https://xiaoiver.github.io/coding/2018/05/20/%E5%9C%A8-iOS-%E4%B8%8B%E4%BD%BF%E7%94%A8-iframe-%E7%9A%84%E7%A7%8D%E7%A7%8D%E9%97%AE%E9%A2%98.html#%E6%94%B9%E8%BF%9B%E6%96%B9%E6%A1%88)

AMP 最初使用了这么一种解决方法：虽然 iframe 不能滚动，但是可以把 HTML 和 BODY 作为滚动容器，让其中的内容滚动。

```html
<html style="overflow-y: auto; -webkit-overflow-scrolling: touch;">
<head></head>
<body
  style="
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
  ">
</body>
</html>
```

虽然 iframe 可以滚动了，但是这种方法存在以下问题：

1. 在 AMP 中，用户定义在 body 上的部分 CSS 规则会失效，例如 margin
2. 由于在容器内滚动，body.scrollTop 会始终为 0，body.scrollHeight 也等于视口高度而非实际全部内容高度

第二个问题影响很大，例如要实现“回到顶部”这样的组件，就无法通过 window.scrollTo() 完成了。 针对这个问题，AMP 给出了这样的 HACK 方案：

* 向 <body> 中插入一个不可见的定位元素 <div>，使用绝对定位 position:absolute;top:0;left:0;width:0;height:0;visibility:hidden;
* 这样在滚动时，通过 -topElement.getBoundingClientRect().top 得到顶部的滚动距离
* 类似的，插入另一个底部定位元素，通过 endElement.offsetTop 获取滚动高度
* 创建一个新的顶部定位元素，在执行滚动到某个位置时，改变其 top，然后调用 scrollIntoView()

所以总共要插入3个 div，可以在 MIP2 的 HTML 结构中看到它们。

最终成品效果可以参考任何一个 MIP 页面，或者[文档](https://mip-project.github.io/api/viewport.html)。在 MIP 中，监听滚动事件不能再单纯使用 `$el.addEventListener('scroll')` 而是要使用 `viewport` 的相关方法。

## iOS iframe + fixed + 滚动

这个 BUG 在于当 iOS 的 iframe 内部有 `position: fixed` 的元素时，上下滚动页面会导致这个元素上下跳动，非常难看。[Demo](https://drive.google.com/file/d/0B_v8thsbiGyDMXZMZkRFZGFRbjA/view)。这个 BUG 存在于 2016 年。[Bug Report](https://bugs.webkit.org/show_bug.cgi?id=154399)

解决方案依然参考于 AMP，在普通的 `<body>` 之外再创建一个容器，把所有 fixed 元素都移动进去。之后可以理解为在第一个 `<body>` 上滚动，而不影响第二个容器，因此也不会再跳动了。

MIP2 为了考虑到例如 `body div` 这样的选择器依然生效，且 `div` 直接挂在 `html` 下和 `body` 平级不太好看，因此把这所谓“第二个容器”依然设定为 `body`，如下：（当然两个 `<body>` 也略有奇怪）

```html
<html>
<body>
    <!-- 原始内容 -->
</body>
<body class="mip-fixedlayer"
  style="position: absolute; top: 0px; left: 0px; height: 0px; width: 0px; pointer-events: none; overflow: hidden; animation: none; border: none; box-sizing: border-box; box-shadow: none; display: block; float: none; margin: 0px; opacity: 1; outline: none; transform: none; transition: none; visibility: visible; background: none;">
  <!-- 所有 fixed 元素 -->
</body>
</html>
```

这个方案的缺点有几个：

1. 虽然考虑了 `body div`，但依然有些选择器会匹配不上，例如 `.wrapper .fixed`，因为 `.fixed` 被移动走了，所以这条肯定也无法生效了。
2. 当 fixed 元素是 CustomElement 时，移动会导致触发 disconnectedCallback 和 connectedCallback，相当于多走了一遍生命周期，可能会有意料之外的效果。
3. （MIP2特有问题）MIP2 的组件允许使用 Vue 的语法，因此这个移动操作相当于强行使用 DOM API 把模板中的一部分拿走，导致了 DOM 和 VDOM 不一致，且绑定的事件，数据等全部失效。

即便如此，比起上下跳动，这还是个值得应用的好方法。

## iframe 宽度 (依然是 iOS)

当我们给 iframe 设置宽度 100% 时，例如 `<iframe width="100%"></iframe>` ，我们不希望 iframe 出现滚动条。

但是 iOS 下存在问题，`width: 100%` 似乎被浏览器的默认设置覆盖了，无法得到应用。

有一种 HACK 方式如下，在 [AMP ISSUE](https://github.com/ampproject/amphtml/issues/11133) 中也采用了类似思路，使用 `min-width` 覆盖掉 iOS Safari 对于 `width` 的默认设置：

```css
iframe {
  width: 1px;
  min-width: 100%;
}
```

## 滚动穿透 （又 TM 是 iOS）

点击穿透我们都听说过，滚动穿透类似：如对话框浮层场景下，在 `fixed` 定位的浮层上滚动时，当滚动位置到底了之后，继续滚动可能会让下面的 body 也一起滚动，造成穿透。演示效果来自 xiaop 。

![滚动穿透](https://xiaoiver.github.io/assets/img/1_default_bottom_overscroll.gif)

首先想到的是让下面的页面不要滚动，因此加上

```css
html, body {
  overflow: hidden;
}
```

这个方案有2个问题：

1. 安卓有效，iOS 无效（依然可以滚动）。
2. 关闭浮层时应该要恢复底部的滚动距离。但这样最会丢失滚动距离，因此没法恢复。

问题 2 比较容易解决，在隐藏之前记录滚动位置，在下次显示之后再设置滚动到刚才那个位置即可。麻烦的是 iOS。

在 iOS 下，只有浮层滚动超过顶部或者底部才会带来问题 （也就是 overscroll）。所以如果我们不让滚动超越顶部和底部，不就没问题了？

```javascript
el.addEventListener('touchstart', function() {
    let top = el.scrollTop;
    let totalScroll = el.scrollHeight;
    let currentScroll = top + el.offsetHeight;

    if (top === 0) {
        el.scrollTop = 1;
    } else if (currentScroll === totalScroll) {
        el.scrollTop = top - 1;
    }
});
```

但是有两点需要注意：

1. 最好不要监听 `touchstart` 事件，否则在一些第三方浏览器（UC）中，会出现点击输入框弹起软键盘时出现不必要的滚动。应该监听 scroll 事件。
2. 如果监听的是 `scroll` 事件，页面在初始状态就应该触发一次，或者直接调用滚动 1px。`touchstart` 则不需要。

补充一句：MIP 注册的是 `scroll` 事件，在 `viewer.js` 中。

## 页面未响应 （还是 iOS）

在 MIP 中会有隐藏 iframe 的情况出现，这是 Google AMP 所不具备的场景，也因此这个 BUG 也是他们没有的。

简单来说，当加载页面时会新建一个 iframe，而后退时则会隐藏这个 iframe （通过设置 `display: none`）。但是在 iOS 的 UC/手机百度 这两个情况下，后退过后的页面完全无法点击，就跟程序未响应一样（但是浏览器的按钮是可以响应的，例如刷新，前进后退等）。

一个精简过的例子在[https://github.com/mipengine/mip2/blob/master/packages/mip/examples/page/iframe/uc.html](https://github.com/mipengine/mip2/blob/master/packages/mip/examples/page/iframe/uc.html)，操作步骤是：

1. 创建 `<iframe src="scroll.html">`
2. 隐藏这个 iframe
3. 页面不响应了

此外还有两个观察结果：

1. 如果第二个页面换成 m.baidu.com 或者其他 google AMP 也能复现，但如果是 PC 百度首页就一切正常。
2. 把 `display:none` 换成 `visibility: hidden` 或者 `opacity:0;height:0;width:0` 依然能够复现问题

基本可以断定这是一个浏览器的 BUG，准确地说是 UIWebView 的 BUG，而使用了较新的 WKWebview 的例如微信就不存在这个问题。 UIWebView 除了这个奇怪的问题，还有诸如 scroll 事件延迟 等其他滚动相关的问题，可以搜索一下。

最终发现，只要是使用了 `-webkit-overflow-scrolling: touch` 的页面都会有这个问题。而不使用这个属性，或者在隐藏 iframe 的同时去掉这个属性，就能正常运行。

最终在 MIP 中使用的方案是每次隐藏 iframe 时添加一条规则 `<style>* {-webkit-overflow-scrolling: auto!important;}</style>`，而显示时再去掉这条规则，可以解决这个问题。

## pushState 引发页面刷新 (iOS QQ浏览器)

这个问题只在 iOS 的 __QQ浏览器__ 下能够复现，其他如 Safari, UC 等均正常。

当我们使用 `pushState` 和 创建 iframe 这两个操作时，需要注意先后顺序，如果：

* 先 `pushState` 再创建 iframe，会导致页面直接 __跳转__ 到 `pushState` 的目标页面。
* 先创建 iframe 再 `pushState` 则没有任何问题。

示例代码如下：

```html
<body>
  <button id="btn1">IFrame1</button>

  <script>
    document.querySelector('#btn1').addEventListener('click', () => {
      // 放这里不OK
      // history.pushState({key: 'key'}, '', './page2.html')
      let iframe = document.createElement('iframe')
      iframe.src = './page2.html'
      iframe.name = 'A normal name'
      document.body.appendChild(iframe)
      // 放这里就OK
      history.pushState({key: 'key'}, '', './page2.html')
    })
  </script>
</body>
```
