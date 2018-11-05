# 如何让你的网页“看起来”展现地更快 —— 骨架屏二三事

让网页展现的更快，官方说法叫做首屏绘制，First Paint 或者简称 FP，直白的说法叫做白屏时间，就是从输入 URL 到真的看到内容（不必可交互，那个叫 TTI, Time to Interactive）之间经历的时间。当然这个时间越短越好。

但这里要注意，和首屏相关的除了 FP 还有两个指标，分别称为 FCP (First Contentful Paint，页面有效内容的绘制) 和 FMP (First Meaningful Paint，页面有意义的内容绘制)。虽然这几个概念可能会让我们绕晕，但我们只需要了解一点：**首屏时间 FP 并不要求内容是真实的，有效的，有意义的，可交互的**。换言之，*随便* 给用户看点啥都行。

![FP/FCP/FMP/TTI](http://boscdn.bpc.baidu.com/assets/easonyq/skeleton/RAIL.png)

这就是本文标题的玄机了：“看起来”。是的，只是看起来更快，实际上还是那样。所以本文并不讨论性能优化，讨论的是一个投机取巧的小伎俩，但的确能够实实在在的提升体验。打个比方，性能优化是修炼内功，提升你本身的各项机能；而本文接下来要讨论的是一些招式，能让你在第一时间就唬住对手。

这所谓的招式就是我接下来要谈的内容，学名骨架屏，也叫 Skeleton。你可能没听过这个名字，但你不可能没见过它。

## 骨架屏长什么样

![](http://boscdn.bpc.baidu.com/assets/pwa-book/skeleton.png)

这种应该是最常见的形式，使用各种形状的灰色矩形来模拟图片和文字。有些 APP 也会使用圆形，但重点都是和实际内容结构近似，不能差距太大。

如果追求效果，还可以在色块表面添加动画（如波纹），显示出一种动态的效果，算是致敬 Loading 了。

![](http://boscdn.bpc.baidu.com/assets/pwa-book/image-skeleton.png)

在图片居多的站点，这将会是一个很好的体验，因为图片通常加载较慢。如上图演示中的占位图片采用了低像素的图片，即大体配色和变化是和实际内容一致的。

如果无法生成这样的低像素图片，稍微降级的方案是通过算法获取图片的主体颜色，使用纯色块占位。

再退一级，还可以使用全站相同的站位图片，或者直接一个统一颜色的色块。虽说效果肯定不如上面两种，但也聊胜于无。

骨架屏完全是自定义的，想做成什么样全凭你的想象。你想做圆形的，三角形的，立体的都可以，但“占位”决定了它的特性：它不能太复杂，必须第一时间，最快展现出来。

## 骨架屏有哪些优势

大体来说，骨架屏的优势在于：

1. 在页面加载初期预先渲染内容，提升感官上的体验。

2. 一般情况骨架屏和实际内容的结构是类似的，因此之后的切换不会过于突兀。这点和传统的 Loading 动图不同，可以认为是其升级版。

3. 只需要简单的 CSS 支持 (涉及图片懒加载可能还需要 JS )，不要求 HTTPS 协议，没有额外的学习和维护成本。

4. 如果页面采用组件化开发，每个组件可以根据自身状态定义自身的骨架屏及其切换时机，同时维持了组件之间的独立性。

## 骨架屏能用在哪里

现在的 WEB 站点，大致有两种渲染模式：

### 前端渲染

由于最近几年 Angular/React/Vue 的相继推出和流行，前端渲染开始占据主导。这种模式的应用也叫单页应用（SPA, Single Page Application）。

前端渲染的模式是服务器（多为静态服务器）返回一个固定的 HTML。通常这个 HTML 包含一个空的容器节点，没有其他内容。之后内部包含的 JS 包含路由管理，页面渲染，页面切换，绑定事件等等逻辑，所以称之为前端渲染。

因为前端要管理的事情很多，所以 JS 通常很大很复杂，执行起来也要花较多的时间。**在 JS 渲染出实际内容之前，骨架屏就是一个很好的替补队员。**

### 后端渲染

在这波前端渲染流行之前，早期的传统网站采用的模式叫做后端渲染，即服务器直接返回网站的 HTML 页面，已经包含首页的全部（或绝大部分） DOM 元素。其中包含的 JS 的作用大多是绑定事件，定义用户交互后的行为等。少量会额外添加/修改一些 DOM，但无碍大局。

此外，前端渲染的模式存在 SEO 不友好的问题，因为它返回的 HTML 是一个空的容器。如果搜索引擎没有执行 JS 的能力（称为 Deep Render），那它就不知道你的站点究竟是什么内容，自然也就无法把站点排到搜索结果中去。这对于绝大部分站点来说是不可接受的，于是前端框架又相继推出了服务端渲染（简称 SSR, Server Side Rendering）模式。这个模式和传统网站很接近，在于返回的 HTML 也是包含所有的 DOM，而非前端渲染。而前端 JS 除了绑定事件之外，还会多做一个事情叫做“激活”（hydration），这里就不再赘述了。

不论是传统模式还是 SSR，只要是后端渲染，就不需要骨架屏。**因为页面的内容直接存在于 HTML，所以并没有骨架屏出场的余地。**

## 骨架屏怎么用

讨论了一波背景，我们来看如何使用。首先先无视具体的实现细节，先看思路。

### 实现思路

大体分为几个步骤：

1. 往本应为空的容器节点内部注入骨架屏的 HTML。

    骨架屏为了尽快展现，要求快速和简单，所以骨架屏多数使用静态的图片。而且把图片编译成 base64 编码格式可以节省网络请求，使得骨架屏更快展现，更加有效。

    ```html
    <html>
        <head>
            <style>
                .skeleton-wrapper {
                    // styles
                }
            </style>
            <!-- 声明 meta 或者引入其他 CSS -->
        </head>
        <body>
            <div id="app">
                <div class="skeleton-wrapper">
                    <img src="data:image/svg+xml;base64,XXXXXX">
                </div>
            </div>
            <!-- 引用 JS -->
        </body>
    </html>
    ```

2. 在执行 JS 开始真正内容的渲染之前，清空骨架屏 HTML

    以 Vue 为例，即在 `mount` 之前清空内容即可。

    ```javascript
    let app = new Vue({...})
    let container = document.querySelector('#app')
    if (container) {
        container.innerHTML = ''
    }
    app.$mount(container)
    ```

仅此两步，并不牵涉多么复杂的机制和高端的 API，因此非常容易应用，赶快用起来！

### 示例

我编写了一个示例，用于快速展现骨架屏的效果，[代码在此](https://github.com/easonyq/easonyq.github.io/blob/master/%E5%AD%A6%E4%B9%A0%E8%AE%B0%E5%BD%95/demo/skeleton/normal/)。

* `index.html`

    默认包含了骨架屏，并且内联了样式（以 `<style>` 标签添加在头部）。

* `render.js`

    它负责创建 DOM 元素并添加到 `<body>` 上，渲染页面实际的内容，用来模拟常见的前端渲染模式。

* `index.css`

    页面实际内容的样式表，不包含骨架屏的样式。

代码的三个文件各司其职，配合上面的实现思路，应该还是很好理解的。可以在 [这里](https://easonyq.github.io/%E5%AD%A6%E4%B9%A0%E8%AE%B0%E5%BD%95/demo/skeleton/normal/index.html) 查看效果。

因为这个示例的逻辑太过简单，而实际的前端渲染框架复杂得多，包含的功能也不单纯是渲染，还有状态管理，路由管理，虚拟 DOM 等等，所以文件大小和执行时间都更大更长。**我们在查看例子的时候，把网络调成 "Fast 3G" 或者 "Slow 3G" 能够稍微真实一些。**

但匪夷所思的是，对着这个地址刷新试几次，我也基本看不到骨架屏（骨架屏的内容是一个居中的蓝色方形图片，外加一条白色横线反复侧滑的高亮动画）。是我们的实现思路有问题吗？

## 浏览器的奥秘：减少重排

为了排除肉眼的遗漏和干扰，我们用 Chrome Dev Tools 的 Performance 工具来记录刚才发生了什么，截图如下：（截图时的网络设置为 "Fast 3G"）

![normal timeline](http://boscdn.bpc.baidu.com/assets/easonyq/skeleton/skeleton-timeline.png)

我们可以很明显地看到 3 个时间点：

1. HTML 加载完成了。浏览器在解析 HTML 的同时，发现了它需要引用的 2 个外部资源 `index.js` 和 `index.css`，于是发送网络请求去获取。

2. 获取成功后，执行 JS 并注册 CSS 的规则。

3. JS 一执行，很自然的渲染出了实际的内容，并应用了样式规则（随机颜色的横条）。

我们的骨架屏呢？按照预想，骨架屏应该出现在 1 和 2 之间，也就是在获取 JS 和 CSS 的同时，就应该渲染骨架屏了。这也是我们当时把骨架屏的 HTML 注入到 `index.html`， 还把 CSS 从 `index.css` 中分离出来的良苦用心，然而浏览器并不买账。

这其实和浏览器的渲染顺序有关。

相信大家都整理过行李箱。我们在整理行李箱时，会根据每个行李的大小合理安排，大的和小的配合，填满一层再放上面一层。现在突然有人跑来跟你说，你的电脑不用带了，你要多带两件衣服，你不能带那么多瓶矿泉水。除了想打他之外，为了重新整理行李箱，必然需要把整理好的行李拿出来再重新放。在浏览器中这个过程叫做重排 (reflow)，而那个馊主意就是新加载的 CSS。显而易见，重排的开销是很大的。

熟能生巧，箱子理多了，就能想出解决办法。既然每个 CSS 文件加载都可能触发重绘，那我能不能等所有 CSS 加载完了一起渲染呢？正是基于这一点，浏览器会等 HTML 中所有的 CSS 都加载完，注册完，一起应用样式，力求一次排列完成工作，不要反复重排。看起来浏览器的设计者经常出差，因为这是一个很正确的优化思路，但应用在骨架屏上就出了问题。

我们为了尽早展现骨架屏，把骨架屏的样式从 `index.css` 分离出来。但浏览器不知道，它以为骨架屏的 HTML 还依赖 `index.css`，所以必须等它加载完。而它加载完之后，`render.js` 也差不多加载完开始执行了，于是骨架屏的 HTML 又被替换了，自然就看不到了。而且在等待 JS, CSS 加载的时候依然是个白屏，骨架屏的效果大打折扣。

所以我们要做的是告诉浏览器，你放心大胆的先画骨架屏，它和后面的 `index.css` 是无关的。那怎么告诉它呢？

## 告诉浏览器先渲染骨架屏

我们在引用 CSS 时，会使用 `<link rel="stylesheet" href="xxxx>` 这样的语法。但实际上，浏览器还提供了其他一些机制确保（后续）页面的性能，我们称之为 preload，中文叫预加载。具体来说，使用 `<link rel="preload" href="xxxx">`，提前把后续要使用的资源先声明一下。在浏览器空闲的时候会提前加载并放入缓存。之后再使用就可以节省一个网络请求。

这看似无关的技术，在这里将起到很大的作用，因为 **预加载的资源是不会影响当前页面的**。

我们可以通过这种方式，告诉浏览器：先不要管 `index.css`，直接画骨架屏。之后 `index.css` 加载回来之后，再应用这个样式。具体来说代码如下：

```html
<link rel="preload" href="index.css" as="style" onload="this.rel='stylesheet'">
```

方法的核心是通过改变 `rel` 可以让浏览器重新界定 `<link>` 标签的角色，从预加载变成当页样式。（另外也有文章采用修改 `media` 的方法，但浏览器支持度较低，这里不作展开了。我把文章列在最后了）这样的话，浏览器在 CSS 尚未获取完成时，会先渲染骨架屏（因为此时的 CSS 还是 `preload`，也就是后续使用的，并不妨碍当前页面）。而当 CSS 加载完成并修改了自己的 `rel` 之后，浏览器重新应用样式，目的达成。

## 不得不考虑的注意点

事实上，并不是把 `rel="stylesheet"` 改成 `rel="preload"` 就完事儿了。在真正应用到生产环境之前，我们还有很多事情要考虑。

### 兼容性考虑

首先，在 `<link>` 内部我们使用了 `onload`，也就是使用了 JS。为了应对用户的浏览器没有开启脚本功能的情况，我们需要添加一个 fallback

```html
<noscript><link rel="stylesheet" href="index.css"></noscript>
```

其次，`rel="preload"` 并不是没有兼容性问题。对于不支持 preload 的浏览器，我们可以添加一些 [polyfill 代码](https://github.com/filamentgroup/loadCSS/blob/master/src/cssrelpreload.js)（来使所有浏览器获得一致的效果。

```html
<script>
/*! loadCSS rel=preload polyfill. [c]2017 Filament Group, Inc. MIT License */
(function(){ ... }());
</script>
```

polyfill 的压缩代码可以参见 Lavas 的 SPA 模板[第 29 行](https://github.com/lavas-project/lavas-template-vue/blob/release-basic/core/spa.html.tmpl#L29)。

### 加载顺序

不同于传统页面，我们的实际 DOM 是通过 `render.js` 生成的。所以如果 JS 先于 CSS 执行，那将会发生跳动。（因为先渲染了实际内容却没有样式，而后样式加载，页面出现很明显的变化）**所以这里我们需要严格控制 CSS 早于渲染。**

```html
<link rel="preload" href="index.css" as="style" onload="this.rel='stylesheet';window.STYLE_READY=true;window.mountApp && window.mountApp()">
```

JS 对外暴露一个 `mountApp` 方法用于渲染页面（其实是模拟 Vue 的 `mount`）

```javascript
// render.js

function mountApp() {
    // 方法内部就是把实际内容添加到 <body> 上面
}

// 本来直接调用方法完成渲染
// mountApp()

// 改成挂到 window 由 CSS 来调用
window.mountApp = mountApp()
// 如果 JS 晚于 CSS 加载完成，那直接执行渲染。
if (window.STYLE_READY) {
    mountApp()
}
```

如果 CSS 更快加载完成，那么通过设置 `window.STYLE_READY` 允许 JS 加载完成后直接执行；而如果 JS 更快，则先不自己执行，而是把机会留给 CSS 的 `onload`。

### 清空 onload

[loadCSS](https://github.com/filamentgroup/loadCSS) 的开发者提出，某些浏览器会在 `rel` 改变时重新出发 `onload`，导致后面的逻辑走了两次。为了消除这个影响，我们再在 `onload` 里面添加一句 `this.onload=null`。

### 最终的 CSS 引用方式

```html
<link rel="preload" href="index.css" as="style" onload="this.onload=null;this.rel='stylesheet';window.STYLE_READY=true;window.mountApp && window.mountApp()">

<!-- 为了方便阅读，折行重复一遍 -->
<!-- this.onload=null -->
<!-- this.rel='stylesheet' -->
<!-- window.STYLE_READY=true -->
<!-- window.mountApp && window.mountApp() -->
```

## 修改后的效果

修改后的代码在 [这里](https://github.com/easonyq/easonyq.github.io/blob/master/%E5%AD%A6%E4%B9%A0%E8%AE%B0%E5%BD%95/demo/skeleton/preload/index.html)，访问地址在 [这里](https://easonyq.github.io/%E5%AD%A6%E4%B9%A0%E8%AE%B0%E5%BD%95/demo/skeleton/preload/index.html)。（为了简便，我省去了处理兼容性的代码，即 `<noscript>` 和 preload polyfill）

Performance 截图如下：（依然采用了 "Fast 3G" 的网络设置）

![preload timeline](http://boscdn.bpc.baidu.com/assets/easonyq/skeleton/skeleton-timeline-after.png)

这次在 `render.js` 和 `index.css` 还在加载的时候页面已经呈现出骨架屏的内容，实际肉眼也可以观测到。在截图的情况下，骨架屏的展现大约持续了 300ms，占据整个网络请求的大约一半时间。

至于说为什么不是 HTML 加载完成立马展现骨架屏，而是还要等大约 300ms 才展现，从图上看是浏览器 ParseHTML 所花费的时间，可能在 Dev Tools 打开的情况下耗时明显了一些，但可优化空间已经不大。（可能简化骨架屏的结构能起作用）

## 后记

这个优化点最早由我的前同事 [xiaop 同学](https://xiaoiver.github.io/) 在开发 Lavas 的 SPA 模板中发现并完成的，ISSUE 记录[在此](https://github.com/lavas-project/lavas/issues/73)。我在他的基础上，做了一个分离 Lavas 和 Vue 环境并且更直白的例子，让截图也尽可能简单，方便阅读。在此非常感谢他的工作！

## 参考文章

* [让骨架屏更快渲染](https://zhuanlan.zhihu.com/p/34550387) - xiaop 同学原作

* [Loading CSS without blocking render](https://keithclark.co.uk/articles/loading-css-without-blocking-render/) - 使用修改 `media` 的方式达成目的。

* [filamentgroup/loadCSS](https://github.com/filamentgroup/loadCSS) - 同样使用修改 `rel` 的方式，并提供了 preload polyfill
