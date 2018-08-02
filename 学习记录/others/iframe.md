# iframe BUG 集锦和解决方案

最近几个月专注于使用 iframe 将多页融合为单页，因此碰上了 iframe 众多奇奇怪怪的问题，通过网上查阅，同事帮忙，自己乱改，了解/解决了许多问题，特此记录。

参考于同事的博客: [在 iOS 下使用 iframe 的种种问题](https://xiaoiver.github.io/coding/2018/05/20/%E5%9C%A8-iOS-%E4%B8%8B%E4%BD%BF%E7%94%A8-iframe-%E7%9A%84%E7%A7%8D%E7%A7%8D%E9%97%AE%E9%A2%98.html)

## 滚动

当页面嵌入 iframe 之后，就有两种滚动方案，分别是：

1. 在 body 上滚动，默认情况。可以把 iframe 看成一个 div，把高度设到足够防止它自身出滚动条。浏览器会对 body 滚动进行优化，例如滚动式把上面的标题栏，下面的菜单栏给隐藏掉，让可视区域变大。

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
