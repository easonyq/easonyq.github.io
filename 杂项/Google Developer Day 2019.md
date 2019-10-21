# 主演讲

## chrome

<img loading="lazy"> 原生支持lazyload

webAuthN 支持 web 上的指纹解锁

lighthouse  performance budget 给出平均值，从而知道自己超过了多少

web.dev 许多在线框架的教程，react, angular 等。

# PWA

WEB 和 APP 各有优势（图），PWA 将两者集中起来。本地APP可以进入 share sheet（弹出的分享菜单中的选项）

PWA 比 react native ,codova 等更进一步，构建的 APP 不需要安装。

PWA 已经可以在 IOS 上安装，支持了全部的系统。

CHROME OS 运行时看起来和普通APP一样，可以在所有系统（WINDOWS, MAC, LINUX）上运行。

PWA只不过是具备了所有必要元素的网站。必要元素就是 manifest, SW 等。

## manifest

<link rel="manifest" href="/manifest.json" crossorigin="use-credentials">

manifest 中的 scope 表示当前的 manifest 作用范围是多少。

## SW

处理离线情况

使用workbox来封装SW，提升开发效率。

## 操作系统集成

安装到本地

// 自动弹出
window.addEventListener('beforeinstallprompt', e => {
    e.preventDefault();
    showInstallPrompt(e)
})

// 通过事件自定义呼出通知
button.addEventListener('click', () => {
    app.promptEent.prompt();
    app.promptEvent.userChoice.then(handlePromptResponse);
})

CHROME自动检测PWA并提示安装（地址栏右侧的加号）

## 新功能

contact picker API (https://contact-picker.glitch.me)

native file system API

注册为文件处理程序

启动事件 PWA 可指定打开方式

# WEB 新功能

针对WEB进行构建，可用的选项很有限。功能很受限，例如读取本地文件等。

针对WEB浏览器发布新的API。如果浏览器能够执行原生应用所有的操作，同事坚守安全的原则。 Project fugu(bit.ly/powerful-apis)

不止 GOOGLE，也有微软，INTEL 等参与 fugu 项目

## 最近发布的API或者原始试用API

developers.chrome.com/origintrials

goo.gle/fugu-codelab
goo.gle/capabilities

1. async clipboard API 的图片支持 (PNG only) Chrome76中推出

navigator.clipboard.writeText(location.href)
await navigator.clipboard.write([new ClipboardItem(xx)])

await navigator.clipboard.read()

2. Web Share API 的文件支持 chrome75 中支持

navigator.canShare && navigator.canShare({files: filesArray})
navigator.share()

3. Web Share Target API chrome76支持

manifest 文件中的 share_target 字段。

4. Badging API （右上角的红色数字） 目前试用状态

window.ExperimentalBadge.set(42)
window.ExperimentalBadge.clear();

window.Badge.set(42)
window.Badge.clear()

5. Contact Picker API

从通讯录中选择联系人，获取联系人的信息（电话，EMAIL等） 目前试用状态

const contacts = aqait navigator.contacts.select(['name', 'email'], {multiple: true})
console.log(contacts)

6. 定期后台同步 目前试用状态

navigator.serviceWorker.ready.then(registration => {
    registration.periodicSync.register('get-lateest-news', 24 * 60 * 1000)
})

self.addEentListener('periodicsync', e => {
    if (e.tag === 'get-latest-news') {
        // TODO
    }
})

7. Barcode Detction API

const code = await new BarcodeDetector().detect(img)

8. Native File System API

只能获取部分文件，不能获取系统文件。可以读取，写入，保存和另存为。

const handle = await window.chooseFileSystemEntries();
const file = await handle.getFile()
console.log(await file.text())

## 后续计划

1. 短信，读取和发送

2. 本地自提访问权限。可以获取本地自提

3. 本地通知 notification 根据 time/location 来触发。

goo.gle/fugu-api-tracker