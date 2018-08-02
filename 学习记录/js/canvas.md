# canvas

canvas 可以用来绘制图形和处理图片。这里重点介绍图片相关的功能

## 图片相关

总计有4个方法：

* drawImage() 用于将图片绘制在画布中
* getImageData() 用于获取画布中的图片（像素）信息
* putImageData() 用于将像素信息写会到画布中成为图片
* toDataURL() 用于把画图中的内容转成一个 URI 对象。使用 `<img src=xxx>` 可以加载这个图片

每个方法的参数就省略了，详情可以查询 API。每个方法都有多种参数。

### imageData

其中 `getImageData()` 的返回值是一个对象，包含 `{width, height, data}`。这个 `data` 是个数组，是保存像素信息的关键，之后 `putImageData()` 的参数也是使用这个 `data` 格式。

`data` 中每4个元素记录一个点的像素信息，分别是 R,G,B和透明值。

### 高清屏幕

在高清屏幕（手机屏幕）中，直接读取图片会导致图片模糊，代码如下：

```html
<!-- 通过 img 标签引入图片，以便绘制到 canvas 中 -->
<img src="html5rocks.png" alt="" width="300" height="90">

<!-- canvas -->
<canvas width="300" height="90"></canvas>

<script>
    function init() {
        var canvas = document.querySelector('canvas');
        var ctx = canvas.getContext('2d');
        ctx.drawImage(document.querySelector('img'), 0, 0, 300, 90);
    }
    window.onload = init;
</script>
```

### 模糊原因

原因和浏览器处理 canvas 的方式有关，[链接](https://www.html5rocks.com/en/tutorials/canvas/hidpi/)

在 devicePixelRatio = 2 的屏幕 (iPhone4S) 上，一个 100px 的内容实际占用了 200px 的大小。也因此，如果我们在界面上显示 100x100 的图片，必须使用 200x200 的原图，否则就会模糊。

而 canvas 中的内容没有这样 \*2 尺寸的处理，因此实际上 canvas 内的全部内容（包括图片，绘制的文字，线条等）都会模糊

### 解决方案

1. 使用 [hidpi-canvas-polyfill](https://github.com/jondavidjohn/hidpi-canvas-polyfill) 这个类库会处理 canvas 中的图形，但没有处理图片。
2. 单独处理图片，修改绘制代码。

    ```javascript
    function init() {
        var canvas = document.querySelector('canvas');
        var ctx = canvas.getContext('2d');

        // polyfill 提供了这个方法用来获取设备的 pixel ratio
        var getPixelRatio = function(context) {
            var backingStore = context.backingStorePixelRatio ||
                context.webkitBackingStorePixelRatio ||
                context.mozBackingStorePixelRatio ||
                context.msBackingStorePixelRatio ||
                context.oBackingStorePixelRatio ||
                context.backingStorePixelRatio || 1;

            return (window.devicePixelRatio || 1) / backingStore;
        };

        var ratio = getPixelRatio(ctx);
        
        // 注意，这里的 width 和 height 变成了 width * ratio 和 height * ratio
        ctx.drawImage(document.querySelector('img'), 0, 0, 300 * ratio, 90 * ratio);
    }
    ```

3. 原理：和普通图片的处理方案相同，依然是放大 canvas，再使用 CSS 将其缩小来达到目的。
