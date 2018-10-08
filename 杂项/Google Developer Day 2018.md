# 内容分级

限于无法分身，我只参加了部分的演讲。作为一个前端开发者，我把内容分为 3 个等级，分别是：
1. 核心：和 WEB 相关的核心内容，以后的工作会使用到的内容。
2. 扩展：和 WEB 有一些关系，或者至少和前端，大前端有一些关系的内容。就算平时工作不直接用到，也值得了解的内容。
3. 乱入：和 WEB 甚至技术都没啥关系，纯粹当了解一下，不看也无所谓的内容。（我是空档期闲着没事或者纯粹兴趣爱好乱入的）

# 长求总

因为两天参加了大概十几个讲座，不管我听得懂还是听不懂，我都尽量记录，因此每个讲座的内容都不少。我先在这里做一个高度的概括，下面的“详细内容”再逐个记录每个讲座的实际内容。

## WEB 相关

* 对比了 WEB 和 APP，强推 PWA。并且今年加入了桌面 PWA （支持 ChromeOS，之后会陆续支持 Windows, Linux, MacOS）
* 浏览器事件循环 [整理后的内容](https://zhuanlan.zhihu.com/p/45111890)
* 介绍了众多浏览器的最新 API，例如人脸识别，蓝牙，本地分享，多媒体等等
* 介绍了 Google 推出的 2 个好用的工具：lighthouse（提供网站评分 & 改进意见）和 puppeteer（轻量小型浏览器内核）。
* 介绍了 AMP 的大体结构和效果

## Flutter

虽然和 WEB 无关，但安卓和 IOS 这些“端”也被称为大前端，算是小有关系吧。

* 介绍了 Flutter：一套代码，两处运行。由 Flutter 来屏蔽安卓和 IOS 的区别。
* Flutter 的特点：美丽，快速(渲染和更新速度)，高效(stateful hotreload)，开放(open source)。
* 以闲鱼作为线上示例进行安利
* Flutter 的内部绘制流程 (build -> render -> paint)
* Flutter 的详细绘制结构：(widget tree -> element tree -> render tree -> layer)
* 内部使用 dart 作为编程语言。和 RN 最大的区别在于：RN 最终调用本地的接口进行绘制，而 Flutter 是自行绘制每个像素点，因此自由度更高。另外 dart:html 支持 WEB，因此和 React 相比在功能上并没有缺失。

## 其他内容

其他内容和 WEB，或者说和技术关系不大，权当了解了解。

* 机器学习的 7 个步骤
* Google 统一广告系统 UAC
* 电子商务网站的设计和性能要求
* Android Things：可以装在设备上，让它变成一个智能设备，从而实现 IoT。

# 详细内容

## 主演讲（扩展）

*主要是介绍 GOOGLE 的各个技术和产品，为接下来 2 天这些方向的详细讲座做一个整体的了解。*

### tenserflow

机器学习和 AI 的开发平台，很多公司有使用，例如美团。
开设了机器学习入门课程，免费视频。
（演讲者是一个金发外国萌妹子，全程中文。虽然发音有点奇怪，但儿化音非常有特点，例如“画面上的【蓝点儿】就是 tensorflow 的使用区域。可以看到它已经覆盖了下至【南极圈儿】……）

### android

推出新版本 android 9 pie，比较新的功能包括：

* 自适应电池
  使用机器学习，根据用户使用习惯把 APP 分为四类，如常用的，不常用的，后台运行的等等。给他们分配不同的运行资源，最终使得 CPU 唤醒时间下降30%，增加电池使用时间

* android jetpack

* kotlin
  一种全新的安卓开发语言，有更好的错误处理，增加代码的稳定性。目前已有40%的开发者在使用 kotlin 开发安卓，还有配套的 extensions (KTX)

* android bundle
  为了解决 APK 越来越大的问题。一般的 APK 会把所有的内容打包，例如支持的各种语言，支持的各类手机底层环境等等。而通过这个技术，可以把 APK 中分解为不同的小单位。在 google play 下载时根据使用者手机的特性，只下载需要的包，大大缩小需要下载的 APK 的体积。

  用前端的思路来理解，就是 bable 的 env preset 了。

### wearOS

主要聚焦于手表，内置了非常多的便捷功能，例如运动数据，天气预报，微信，邮箱，支付宝的支付码等等。有点类似于 apple watch。也给 APP 的开发人员提供了新的开发设备和入口。

### firebase

是一个开发 APP 的应用平台，为了减少开发 APP 之前需要搭建环境的麻烦。它把本地服务器移动到云，即开发时直接请求云端的数据。另外提供了 WEB 界面去修改或者监控云上的数据。最核心的点在于这个数据的更新非常实时，所以才能成功取代本地服务器。

也可以使用 API 访问到云端的数据，和访问本地文件的复杂度和时间是类似的（据说）。

## WEB 的形式 现在 & 未来 （核心）

*算是对接下来 2 天 WEB 相关的课题的一个预告。阐述当前的问题，提出解决方案（主要是 PWA）*

WEB 是可以通过链接访问且可以在大多数设备上使用的内容集，网页的内容可以自行更新，不需要安装补丁或者更新软件版本，因此是更优秀的形式。

### 体验和用户方面

对比全球前 1000 的手机 APP 和前 1000 的手机网页，发现：

* 每个用户的月均花费时间：APP 大幅领先网页
* 每个月不同的用户数量（UV）：网页 1140 万，APP 400 万。网页大幅领先。
* 78% 的用户时间放在了前五名的 APP

结论：
1. 沉浸式的 APP 体验优于传统网页，因此更能吸引用户留存
2. APP 吸引用户的能力不如普通网页，因为 APP 有需要安装的先天问题。
3. 非 TOP APP 能分到的蛋糕其实不大

因此我们如果能把网页的体验做到和 APP 类似的话，PWA 应用就应运而生了。
PWA：快速，集成，可靠，有吸引力。

Google 的两个比较有代表性的例子：
1. Google Maps Go 谷歌最大的 PWA，已经被预装在 Android Go 设备上，但本质只是一个常规的网站，非 APP。 [https://google.com/maps/](https://google.com/maps/)

2. Offline GMail: 使用 Service Worker 来实现保存和阅读离线邮件。

总体来看：
1. PWA 应用的各项指标（广告点击率，用户花费时间等）均不输给 APP
2. PWA 在更新时需要的流量大小也远远小于 APP。（因为仅仅需要更新网页，而 APP 需要打补丁或者重新安装）
3. 平均一天有 10 个 PWA APP 被安装（浏览 OR 添加到桌面）。

未来的方向上，谷歌正在考虑桌面上的 PWA 应用（即不运行在浏览器内，而是直接运行在设备本身，和浏览器平级）。三星，微软的EDGE 分别引入了桌面 PWA APP。CHROME OS 也引入了 PWA APP。打开网页后会弹出提示是否添加 APP，添加后即可以快捷访问，和APP基本是一致的。WINDOWS 7,8,10 CHROME 70 以后会包含这些功能。MAC & LINUX 在明年早些时间。

### 运行速度方面

高端和低端手机在运行相同的JS上时间相差非常非常大。
平均每个页面运用 360K 的 JS 文件，和8年前相比基本是6倍大小。

解决方案：
1. 推荐两个工具：puppeteer & lighthouse （具体在后面有转场演讲）

2. 从 AMP 学到的 web packaging。把网页内容存到离用户更近的位置。把网页内容类似镜像的概念，从远处的原服务器复制到镜像服务器，由游览器识别出来。（国内可能没戏？）

3. web assembly 用来解决 JS 的性能问题，创造一个执行环境，提升处理性能，额外使用 30% - 40% 的CPU能力。例如网页版的绘图软件，CAD。

网页包含了越来越多的新功能，会解决网页的性能就是不如 APP 这个问题。

## 浏览器的 event loop （核心）

已经整理到[知乎](https://zhuanlan.zhihu.com/p/45111890)

## 打造跨平台的 WEB 站点 （核心）

*先强推 PWA，再强推桌面 PWA*

PWA：能够在浏览器和在桌面同时使用的站点。

### 一些优秀体验的移动站的建议

* __在需要时才请求权限，而不是在用户打开应用程序的时候就请求。__（这点和之前不一样了，Chrome 已经支持了）
* 自动登录，流畅的使用流程
* request payment API （W3C标准）能够集成 GOOGLE PAY。（国内悬）
* 有 42% 的站点没有为输入框 (input) 指定类型 (type)，因此体验欠佳
* 谷歌列出了一些最佳实践，来指导用户如何去把自己的移动站点做得更好。

Service Worker 已经可以安装在几乎全部的浏览器上。
腾讯新闻接入了 Service Worker 之后，性能提升，浏览次数，转化率均有提升

演示了 Starbucks PWA，通过 PWA 下达的订单增长超过 12%，每日和每月的活跃用户几乎翻倍。桌面用户无需使用移动设备即可下单。

__53% 的用户会放弃加载时间超过3秒的网站。__

### PWA 应用的特点

四个特点：速度快，可安装，可依赖，体验好。

* 速度快
    * 使用 placeholder content 控件（类似于 skeleton，也可以是低精度的站位图片）
    * 预缓存内容

* 可安装
    * 外观和行为与其他本地 APP 类似。（添加到桌面并从桌面打开，没有浏览器样式）
    * Web APK：PWA 可以像普通 APP 一样出现在引用程序中。（例如使用某程序打开的列表，分享的列表等均可以出现。目前安卓已经实现）
    * CHROME 显示添加到主屏的条件，还包括 __“必须包含一个监听 fetch 事件的 Service Worker”__
    * 避免一进入 APP 就弹出添加到首屏的提示。
        * 监听 `beforeinstallprompt` 保存 `event`
        * 之后调用 `event.prompt()` 弹出添加到主屏的提示
        * 安装成功后有 `appinstalled` 事件发射出来

* 可信赖
    * 使用 workbox (一个快速生成 Service Worker 的工具，也由 Google 开发)
    * 预缓存内容 (precache)
    * 运行时缓存 (runningCache)
    * 使用 indexDB 缓存内容

* 体验好
    * 恰当的后退导航按钮（不要一下子退到最外面，要一步一步）
    * 使用 toast 最小化影响主体内容。

### GOOGLE 的 PWA

* GOOGLE 搜索：
    * 外部的 JS 请求减少 50%
    * 由加载 JS 引起的用户延迟减少 6%

* Bulletin
    * 体积比 APP 更小
    * 支持包括照片和视频在内的多媒体捕获（拍照，拍视频）

* GOOGLE 地图
    * 从根本上改善低端设备或有限网络环境中的体验
    * 核心用户应用场景：
        1. 找到自己的位置
        2. 寻找一个位置
        3. 寻找附近的位置
        4. 寻找路线&导航
    * 页面加载成功率提升 20%

    * 缓存策略：（多种缓存配合）
        * 浏览器缓存 maps tiles
        * indexDB 记录用户搜索和 map files 版本等。

### 桌面 PWA

根据统计：白天10点到7点，desktop 的使用时间超过 phone 或者 tablet。

在桌面应用上，常规的做法是自定义构建一个简易的浏览器内核，并使用它容纳网页。但实际上用户的 PC 上可能已经有不止一个浏览器。因此我们实际上应该聚焦在应用的内容本身，而不需要那个浏览器外壳。

因此我们需要跨浏览器，跨操作系统的 PWA APP。（WINDOWS 和 MacOS 都可以运行的）

实现方面，同样使用 manifest.json，重要的是 scope 属性（和 Service Worker 的 scope 类似）

在桌面应用也需要使用响应式设计，根据宽度和大小的不同，显示不同的内容。（例如天气预报，可以分7天，5天，3天，小图标等等）

更多详情可以参阅

* [Chromium Blog](https://blog.chromium.org)
* [网络应用安装横幅](https://developers.google.com/web/fundamentals/app-install-banners/)
* [Progressive Web Apps on the Desktop](https://developers.google.com/web/updates/2018/05/dpwa)
* [Scope in Manifest.json](https://developer.mozilla.org/en-US/docs/Web/Manifest#scope)

## Google AMP （核心）

*对 AMP 进行了大概的介绍*

### 原因

需要使用 AMP 的原因主要是因为传统网页加载太慢，loading 时间太长。以下是一些统计数据：

* 53% 的用户放弃加载用时超过3秒的网页
* 3g 下的平均加载页面用时 19 秒
* 60% 的全球移动网络用的是 2g

### 组成部分

AMP 是由几个部分组成的：

* html (普通 HTML + AMP 组件)
* js (内联的脚本或者绑定属性)
* cache (Google AMP cache 自动抓取)

### 做法

AMP 的做法包括

* 阻止加载耗时的内容
* 待加载完后才显示内容
* 严禁投放令用户分神的广告
* 直到用户需要才加载相关内容 (lazy loading)

### amp-bind

根据用户的交互，使用数据绑定 (data binding) 和表达式 (expressions) 来动态变化页面的显示内容。

通过 3 个步骤实现这个过程：

1. state （设定初始状态，例子如下）

    ```html
    <amp-state id="team">
        <script type="application/json">
            {"star": "Yao Ming"}
        </script>
    </amp-state>
    ```

2. bind （在页面某个位置将显示和状态关联起来，例如 `<p [text]="team.star + ' is tall!'">`）
3. mutation (通过 `AMP.setState` 来更新状态，例如 `<button on="tap:AMP.setState({team: {star: '姚明'}})">`)

### AliExpress 迁移到 AMP 的优缺点分析

AliExpress （海外的阿里） 使用了 AMP，因此 Google 以他们为范例，阐述了 AMP 的得失

* 劣势
    1. 只能使用 AMP 组件
    2. 不能使用 cookies 和 localStorage
    3. 无法直接支持 touch 事件
* 优势
    1. 重点关注在业务逻辑上，花更少的精力在性能方面，开发效率更高
    2. 能够给开源项目输出代码 （例如 `<amp-date-countdown>` 组件）
    3. 依靠 AMP 获得了更好的性能
    4. SEO 效果更好

### 更多信息

最后是 3 个参考网站

* [ampproject.org](http://ampproject.org)
* [ampbyexample.com](http://ampbyexample.com)
* [ampstart.com](http://ampstart.com)

## Google 的两款工具推荐：lighthouse & puppeteer (核心)

*lighthouse & puppeteer，值得使用和学习*

### lighthouse

网站评分工具，目前是 3.0 版本。它能够衡量一个移动战的各类指标，并指出网站提升的方向。它的审查内容包括：

* PWA 功能
* 最佳范例
* 可访问性
* SEO
* 性能

使用方法：(任选其一)

* Chrome Dev Tools (F12)
* Chrome Extensions
* npm (使用 nodejs 线下跑分)
* web (直接去 lighthouse 官网输入网址在线测试)
* github + travis 可以作为 PR 的 task （类似于自动化测试那样，每次发起 PR 都运行一下，给出分数变化趋势）。具体使用方法可以参考 [lighthouse-ci](https://github.com/ebide/lighthouse-ci/)，可以设置及格线，站点 URL 等等。

### 衡量页面的一些指标

* FP - 第一次页面显示
* FCP - 第一次有内容的页面显示
* FMP - 第一次有意义的页面显示
* TTI - 可以开始用户交互

![指标](http://boscdn.bpc.baidu.com/assets/easonyq/gdd/benchmark.jpg)

![指标](http://boscdn.bpc.baidu.com/assets/easonyq/gdd/benchmark-2.jpg)

![指标](http://boscdn.bpc.baidu.com/assets/easonyq/gdd/benchmark-4.jpg)

![指标](http://boscdn.bpc.baidu.com/assets/easonyq/gdd/benchmark-3.jpg)

这部分还可以参阅一篇微博：[以用户为中心的性能指标](https://xiaoiver.github.io/coding/2017/06/09/%E4%BB%A5%E7%94%A8%E6%88%B7%E4%B8%BA%E4%B8%AD%E5%BF%83%E7%9A%84%E6%80%A7%E8%83%BD%E6%8C%87%E6%A0%87.html)

### puppeteer

首先讲一个概念，叫做 headless chrome。简单来说就是没有头尾的浏览器，也就是浏览器内核。它的特点是：

* 内核与最新的 Chrome 保持一致
* 可以使用最新的功能和接口，例如 streams, cssgrid, service worker 等等
* 能够通过代码操作 Dev Tools 的功能和数据，例如模拟网络延时，模拟设备等等

puppeteer 基于这样一个内核，给开发者提供一套接口，来做一些事情，简单结构如下：

![puppeteer 结构](http://boscdn.bpc.baidu.com/assets/easonyq/gdd/puppeteer-structure.jpg)

利用这个可以做哪些事情呢？下面举几个例子：

1. 截屏

    ```javascript
    puppeteer.launch().then(async browser => {
        const page = await browser.newPage()
        await page.goto('https://example.com')
        await page.screenshot({path: 'example.png'})

        await browser.close()
    })
    ```

2. 获取页面数据

    ```javascript
    const metrics = await page.metrics()

    // metrics.ScriptDuration
    // metrics.LayoutDuration
    // metrics.RecalcStyleDuration
    // metrics.JSHeapUsedSize
    // metrics.NodeCount
    ```

    ![页面数据](http://boscdn.bpc.baidu.com/assets/easonyq/gdd/page-data.jpg)

3. 拦截网络请求

    ```javascript
    await page.setRequestInterception(true)

    page.on('request', req => {
        if (req.resourceType === 'image') {
            req.abort()
        }
        req.continue()
    })

    await page.goto('https://www.youtube.com/')
    ```

4. 生成 PDF

    ```javascript
    const page = await browser.newPage()

    await page.setContent(`
        <!doctype html>
        <h1>Some Report in PDF</h1>
        ...
    `)

    await page.setViewPort({
        width: 1280,
        height: 1024,
        deviceScaleFactor: 2
    })

    await page.pdf({
        path: 'report.pdf',
        margin: {top: '16px', ...}
    })
    ```

5. 验证每个请求是否可以离线访问

    ```javascript
    const resp = request.response().fromServiceWorker()
    console.log(url, resp ? '√' : '×')
    ```

更多的例子可以参考 [puppeteer-examples](https://github.com/GoogleChromeLabs/puppeteer-examples) 和 [pptraas.com](http://pptraas.com)

最后，还额外推荐了两个开发工具：
* [ndb](https://github.com/GoogleChromeLabs/ndb) - an improved debugging experience for Node.js, enabled by Chrome DevTools
* [page speed insights](https://developers.google.com/speed/pagespeed/insights/)

## 深入探讨 WEB 上的新功能 (核心)

*海量最新 API 正在袭来！*

先介绍了 PWA 和 Service Worker。 OFO 把 PWA 应用于共享单车，在美国上线了。

随后是一大堆最新的 API。这些 API 有些刚刚加入标准，有些尚未加入标准。但均已经在 Google Chrome 上实现了。

因为 API 和代码太多，我也没能全部记住 & 查阅，因此这里仅列出大致内容和关键词。

### 操作系统整合

1. 其实就是添加到桌面，manifest.json，但是增加了安装成功的事件。

    ```javascript
    window.addEventListener('appinstalled', e => app.logEvent('a2hs', 'installed'))
    ```

2. `<input type="file" accept="image/*">` 像 APP 那样选择图片。其中 `accept = image` 是新增的选项

3. `navigator.share` 分享功能

    ```javascript
    let result = await navigator.share({
      title: 'Paul Rocks',
      text: 'He really does!',
      url: 'https://paul.kinlan.me/'
    })
    ```

4. Share Receiver。 能够像本地 APP 一样，在其他网页分享时，显示在分享程序的列表中。通过在 manifest.json 中增加 `share_target` 对象来实现这个功能。

    ```javascript
    // manifest.json
    "share_target": {
      "action": "compose/tweet",
      "params": {
        title,
        text,
        url
      }
    }
    ```

5. download manager 后台下载，断点续传，完成后的通知等等。

6. `navigator.mediaSession` 控制媒体（视频，音频等）能够控制播放的标题，图片，进度，控制前进后退等等。

7. `document.pictureInPictureElement` 允许浏览器退到后台时，画面依然在设备的桌面上（类似悬浮窗口）。

### 高级多媒体

统计数据：__有 70% 的网络流量来自视频__

1. `navigator.mediaDevices.enumerateDevices`：获取系统上可用的多媒体输入和输出设备的信息，如麦克风，摄像头等

2. `new ImageCapture`：截图

3. `navigator.mediaDevices.getUserMedia`：向用户请求权限获取音频或者视频流等。

    ```javascript
    let stream = await navigator.mediaDevices.getUserMedia({video: true})
    let video = document.querySelector('video');
    video.srcObject = stream;
    video.onloadedmetadata = function(e) {
        video.play();
    };
    ```

4. `canvas.captureStream(25)`：实时捕获 canvas 画布上的内容，输出为流，参数为帧率。

### 识别相关

1. 识别二维码

    ```javascript
    let detector = new BarcodeDetector()
    let codes = await detctor.detect(image)
    ```

2. `new FaceDetector()` 识别人脸
3. `new TextDetector()` 识别文字

### 硬件

1. Web BlueTooth

    ```javascript
    const device = await navigator.bluetooth.requestDevice(...)
    ```

2. Web USB

    ```javascript
    let device = await navigator.usb.requestDevice(...)
    ```

3. Ambient Light Sensor (环境光传感器)

    ```javascript
    let als = new AmbientLightSensor({frequency: 10})
    ```
4. Presentation API

    ```javascript
    const pr = new PresentationRequest('https://airhorner.com/')
    ```

developers.google.cn/web 有列出更多的信息

## WEB 电子商务 （扩展）

*列出一些优秀的电子商务网站，分析他们的共性和优化点*

在全球在线购物中有 66% 的用户通过 WEB 购物 （34% 通过 APP）。但是转化率上 WEB 要比 APP 低得多。

现存 WEB 的劣势：手动，繁琐，速度慢，多次点击

完善电子商务网站的几种方式

### 性能目标

性能目标必须要和商业价值挂钩，否则毫无意义。

性能目标可以包括：
* 网页加载时间
* 首次有效呈现时间
* 可交互时间
* 可连续交互时间
* 页面重量

有几个简单的原则：
1. 不要为了性能优化而删减核心功能（正常的公司我觉得不太可能会这样）
2. 要首先显示重要的内容。例如侧边栏看可以放后面，确保主体内容先渲染。

### 产品展示

热烈推荐 Google AMP

[DHgate](m.dhgate.com) 一个来自中国的B2B电商网站，采用了 AMP。

### 图片

图片需要优化
1. 格式：png -> jpg。
2. 根据屏幕尺寸选择正确的图片尺寸
3. 尽量避免用户等待图片加载的时间（例如提前加载）
4. 使用低像素图片占位（和skeleton类似）

lighthouse 可以测试和提供优化建议

### 查找产品 - 浏览和搜索

多使用 prefetch 进行预加载

网站内部的搜索需要处理：

* 拼写错误
* 同义词
* 自动填充
* 缩写
* 分面搜索。

### 购物车和结算

这是电子商务中最重要的环节。56%的美国消费者因为结算问题放弃在移动端的购物。可能因为速度慢，或者要填太多的信息等等。

chrome 自动填充功能可以帮助用户登录。每年帮助80亿用户登录。

跨平台：例如手机上浏览商品，到PC上进行付款。
顺畅的跨平台付款 payment request API
顺畅的跨平台身份验证 credential management API （这两个都是浏览器包含的 API 功能）

可以利用 PWA 增加更好的体验。

好的例子 [ecer.com](ecer.com),  [m.jd.id](m.jd.id) ( JD 的印度尼西亚版本)

## Flutter (扩展)

*开发一套代码，同时在安卓和 IOS 上运行，构建 APP 的简便方式。*

### APP 的现状

新的 APP 要尽快上市的重要性不言而喻，即便是已有的 APP 的更新速度其实也至关重要。数据显示，仅有 __3%__ 的移动  APP 新用户在 30 天后仍然保持活跃状态。因此，如何提升用户粘性，如何通过持续的更新来吸引 & 留存用户，是移动 APP 成功与否的关键。

移动 APP 目前面临的最主要挑战是：

* 业务挑战方面
    * 上市期 （越短越好）
    * 碎片化
    * 灵活性
    * 招聘
    * 互动
    * 测试新想法/原型开发
    * 高昂的开发费用

* 技术原因方面
    * 同时面向 IOS 和安卓平台构建和发布应用
    * 招聘人员拆分代码库
    * 设计/开发协调方面的挑战
    * 同步版本
    * 更新和刷新应用
    * 快速发布和迭代

说到底，核心就是一句话：APP 需要开发 IOS 和安卓两个版本，因此需要两个代码库，需要两套开发人员，需要两批设计（适应不同的屏幕尺寸 OR 屏幕特点，例如 IOS 的刘海）。而这两套又不能完全隔离，因为它们是同一个产品，因此还需要在保持特性的前提下拥有相当的一致性，例如体验，样式等等。

现有的技术也有致力于解决这方面问题的，比较著名(且早期的)的是 Phonegap, Adobe Air, Appcelerator 等 (React Native 后面再说)，但都存在一个共同的缺点：比较简陋，无论从样式上，体验上还是从功能上，全方位的不如原生开发。因此为了造就理想中的城堡，获得最优的效果，公司和开发者依然不惜耗费巨大的成本走迷宫。

![城堡](http://boscdn.bpc.baidu.com/assets/easonyq/gdd2018/flutter1.jpg)

### Flutter 的特点

主要有 4 个特点：美丽，快速，高效，开放。

* 美丽

    可以画页面上的每个像素（这也是和 RN 最大的区别），以精美界面获奖。可以理解为高性能的渲染引擎。

* 快速

    能够保证 60fps 的帧率，可以调用 GPU 加速。因为高层代码被编译为机器代码，因此在低端手机上也能取得相当好的效果。

* 高效

    高效指的更多的是开发体验。其中最具特色的是 stateful hotreload （修改代码热加载时组件状态能够保持住，而不是从初始状态重新开始）。

* 开放

    Flutter 是开源 & 免费的，可以在 Github 上看到代码并参与。另外还有中文官网和中国镜像，访问容易。

渲染机制方面，它总共有 3 个阶段：

build -> render -> paint

而实际渲染，它通过构造以下几个内容来进行：

widget tree -> element tree -> render tree -> layer

### 使用 Flutter 的四种方法

1. 从无到有构建全新 APP
    在 Flutter 中实现新偶像，并同时落实到 IOS 和安卓 APP 中。

2. 针对新的 APP 构想进行原型设计
    使用 Flutter 在前所未有的短时间内测试 APP 构思或者想法

3. 针对另一个平台构建 APP
    如果已经有了 IOS 和安卓 APP 的其中之一，使用 Flutter 针对另外一个平台构建 APP。验证无误后再扩展并合并代码库。

4. 将 Flutter 用于 APP 的某一部分 （线上演示的闲鱼就是这种）
    在生产环境中使用 Flutter 测试现有 APP 中的一个或者几个页面

## APP 营销最佳实践 (乱入)

*主要聚焦在如何利用 google 的各条产品线和合作品牌来提升自身 APP 的广告效果。但我觉得多数是针对出海的 APP，因为搜索，youtube，google play 等国内使用并不十分广泛。*

Google Play 是用户获取游戏/应用的最主要渠道。广告吸引的安装次数超过60亿次。

本次介绍主题：UAC 通用应用广告系列

### 解决的问题

* 如果投放搜索广告，要设置关键词，并且为每个关键词设置点击成本。
* 如果投放到 YOUTUBE，要投放到特定的频道。
* 等等。

UAC 可以免去这些设置工作，在搜索，play store, emails. web, youtube等，使用机器学习来投放广告。
例如如果用户看完游戏视频后下载游戏的数量非常多，那 UAC 就会增加这条路径的权重。

这里又学到了两个概念：
* CPI 单次安装成本  安装，推广。主要目标是下载/激活。
* CPA 单次事件成本  留存，活跃，营收。主要目标是用户质量。

### 优秀的视频广告特点

1. 15-30S
2. 加入音乐和字幕
3. 尽早的抓住眼球，显示出品牌。最好在1/4时长之前更容易获得更高的转化率
4. 结尾处（最后一帧）加入号召文字或者按钮图案。（马上下载，马上游玩等等）

H5 广告：

1. 有清晰的号召按钮
2. 增加激励性，如折扣优惠等。

谷歌拥有丰富多样的用户行为数据

* 浏览记录
* 应用内的购买频率
* 玩过的游戏类型
* 搜索词条
* 位置 （考虑时差）
* 时段 （例如工作日推游戏就不是很有效）
* 设备类型 （例如有没有买最新的 iphone）
* 从 youtube 看过什么视频（有73%的玩家喜欢看别人打游戏，48%的玩家喜欢看别人打游戏超过自己打游戏，61%玩家会在购买游戏之前看 youtube 视频）

根据这些行为数据，能够猜测用户画像，从而有针对性地去对这类用户制作视频，精准投放。举了网易的荒野行动进军日本的例子，有针对性的转化率为4倍。
