# JSBridge

最早广为人知的是微信的 WeiXinJSBridge，糯米也提供过 BNJS，目的都是给在 WebView 里面的 JS (Hybrid APP) 能够通过这个“桥”来调用到一些本地 APP 的功能。目前，包括微信小程序，React Native 都离不开 JSBridge，下面介绍一下它的原理。

![](https://user-gold-cdn.xitu.io/2018/3/29/16270f34f02109eb?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

## 技术方案

目前 JSBridge 主要有三种技术方案：

1. 基于 Web 的 Hybrid 解决方案：例如微信浏览器、各公司的 Hybrid 方案
2. 非基于 Web UI 但业务逻辑基于 JavaScript 的解决方案：例如 React-Native
3. 微信小程序基于 Web UI，但是为了追求运行效率，对 UI 展现逻辑和业务逻辑的 JavaScript 进行了隔离。因此小程序的技术方案介于上面描述的两种方式之间。

## 用途

JSBridge 的用途，普遍被认为是使得前端 JS 拥有调用本地 APP 的能力。但实际上，JSBridge 的作用核心是**构建 Native 和非 Native 间消息通信的通道**，而且是**双向通信的通道**。

![](https://user-gold-cdn.xitu.io/2018/3/29/16270f744a3e61f2?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

所谓双向通信的通道，分别是：

* JS 向 Native 发送消息: 调用相关功能、通知 Native 当前 JS 的相关状态等。
* Native 向 JS 发送消息: 回溯调用结果、消息推送、通知 JS 当前 Native 的状态等。

## 原理

JavaScript 是运行在一个单独的 JS Context 中（例如，WebView 的 Webkit 引擎、JSCore）。由于这些 Context 与原生运行环境的天然隔离，我们可以将这种情况与 RPC（Remote Procedure Call，远程过程调用）通信进行类比，将 Native 与 JavaScript 的每次互相调用看做一次 RPC 调用。

在 JSBridge 的设计中，可以把前端看做 RPC 的客户端，把 Native 端看做 RPC 的服务器端，从而 JSBridge 要实现的主要逻辑就出现了：**通信调用（Native 与 JS 通信）** 和 **句柄(handler)解析调用**。（如果你是个前端，而且并不熟悉 RPC 的话，你也可以把这个流程类比成 JSONP 的流程）

### JS 调用 Native

1. 注入 API

    注入 API 方式的主要原理是，通过 WebView 提供的接口，向 JavaScript 的 Context（window）中注入对象或者方法，让 JavaScript 调用时，直接执行相应的 Native 代码逻辑，达到 JavaScript 调用 Native 的目的。

    通常在 JS 端，就使用 `window.xxx` 来调用方法，直接执行本地代码。

2. 拦截 URL Scheme

    URL Scheme 是一种类似 URL 的链接，但协议不是 http 或者 https，而是自定义的。例如 baiduboxapp://, thunder:// 等等。这里也有 host 的概念，对应 URL 的 host，就是紧接在协议后面的部分。

    拦截 URL Scheme 的主要流程是：Web 端通过某种方式（例如 iframe.src）发送 URL Scheme 请求，之后 Native 拦截到请求并根据 URL Scheme（包括所带的参数）进行相关操作。

    这种方式的主要问题有：

    1. URL 有长度限制

    2. 这是一种异步的方式，相比注入 API，耗时会更久

    但这种方式支持 ios6，**可以说是一个古老的方法。在现在的 APP 中已经基本淘汰了**。

3. prompt

    在古老的安卓（4.2以下）还有一种叫做 prompt 的方式，但这也是一种兼容方式。主流还是注入 API。

### Native 调用 JS

这要求 JS 的方法必须在全局的 window 上。本质是把 JS 方法变成字符串，拼接到本地代码中。因为都是本地代码，我就不列了。

## 实现

具体实现时，和 JSONP 很类似。JSONP 通过一个**唯一的** callback 方法来实现双向通讯（浏览器发给客户端，客户端返回调用代码，由浏览器执行）。而在这里，使用的是一个唯一的自增 id。JS 端把 ID 传递到 Native，之后 Native 以这个 id 作为标记执行 callback。

反向时（从 Native 到 JS）也是一个类似的过程，只是这个 id 变成了由 Native 生成的 responseId 取代。

此外，在引用 JSBridge 的方法上也存在两种：

1. 由 Native 注入。相当于 JS 运行起来时默认就有 `window.xxx`，因此可以直接调用。这个方法的优点是 JSBridge 完全由 Native 控制，因此版本管理很容易，可以随着 Native 一起升级，没有兼容问题。但缺点是注入有成功率，因此需要有重试的机制，此外在 JS 调用时需要事先判断是否注入成功(`window.xxx` 是否存在)。

2. 由 JS 引用（普通 `<script>` 的形式）。这个方案的优缺点和上面相反，优点在于能确保引用的成功率，因此确认一定有 `window.xxx` 的存在。缺点是如果 Native 升级，必须更新内部网页引用的桥地址。绝大部分情况需要做版本兼容。

## 参考链接

[JSBridge的原理 - 掘金](https://juejin.im/post/5abca877f265da238155b6bc)