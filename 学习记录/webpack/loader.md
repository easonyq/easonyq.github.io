# webpack loader

webpack1 loader 的写法：（用 `!` 连接的字符串）

```javascript
module: {
    loaders: {
        test: /\.css$/,
        loader: 'style-loader!css-loader'
    }
}
```

webpack2 loader 的写法：（数组表示）

```javascript
module: {
    rules: {
        test: /\.css$/,
        use: [
            'style-loader',    
            { loader: 'css-loader', options: { importLoaders: 1 } },
            'postcss-loader'
        ]
    }
}

```

以上两种写法都是全局的规则配置。还有一种针对个别文件的单独配置方法，即在引用时添加 loader，如：

```javascript
import css from 'style-loader!css-loader!./file.css';
```

__这三种写法的执行顺序都是从右到左（针对1,3），从下到上（针对2）。__

这主要是 loader 内部采用函数式写法有关。如果我们把方法考虑成函数，则函数调用是类似 `f(g(x))`，计算时先算 `g` 再算 `f`，也就是从右到左的。而 webpack2 内核依然是函数式，只是写法变成了数组而已。

## 几个样式相关的 loader 的作用

* less-loader, stylus-loader, sass-loader 等

    用以让 webpack 识别对应后缀的样式语言，并处理成 css

* css-loader

    用于处理 css 中的 `@import` 和 `url()`，把引用的资源一起引入到 webpack 中。

    ```javascript
    url('image.png') => require('./image.png')
    url('~module/image.png') => require('module/image.png')
    ```

    配置项 `importLoaders` 值为数字，表示在 css-loader 之前有多少个 loader 作用。

* style-loader

    把所有 css 变成一个 `<style>` 标签并插入到 `<head>` 中。可以通过配置项 `insertAt: 'top'` 放到最开头。（放在越后面的样式在规则权重相同时优先级越高）

* postcss-loader

    主要用于给 CSS 自动添加各类浏览器前缀。需要使用 autoprefixer 并读取 browserslist。

    postcss-loader 和 less-loader 一起使用时，必须把 postcss 写在前面（即先执行 less-loader，因为它只针对 CSS 有效）