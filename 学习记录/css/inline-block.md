# inline-block

## inline, block, inline-block

* inline:
    使元素变成行内元素，拥有行内元素的特性，即可以与其他行内元素共享一行，不会独占一行. 
    不能更改元素的height，width的值，大小由内容撑开. 
    可以使用padding，margin的left和right产生边距效果，但是top和bottom就不行.

* block:
    使元素变成块级元素，独占一行，在不设置自己的宽度的情况下，块级元素会默认填满父级元素的宽度. 
    能够改变元素的height，width的值. 
    可以设置padding，margin的各个属性值，top，left，bottom，right都能够产生边距效果.

* inline-block:
    结合了inline与block的一些特点，结合了上述inline的第1个特点和block的第2,3个特点.
    用通俗的话讲，就是不独占一行的块级元素。

在默认情况下，当一个 Box 中包含两个 __有宽度的__ 块级子元素时，因为块级元素占满一行，所以即使有宽度，也不会让他们在同一行中，如下：

![block](https://images2015.cnblogs.com/blog/1144006/201705/1144006-20170513095231269-1572459142.png)

如果给这两个子元素设置 display: inline-block，就可以让他们在同一行显示

![inline-block](https://images2015.cnblogs.com/blog/1144006/201705/1144006-20170513095240254-1054271047.png)

注意两个点：
1. 实现效果和 float 类似，都能放到一行中
2. 但发现两个元素之间，以及元素和父元素的底边之间还有一条空白

## inline-block vs float

首先，float 元素会导致父元素高度不撑开，因此需要 overflow: hidden 或者 clearfix 来进行修复；inline-block 没有这个问题。

其次，在多个元素 float 并且高度不一的情况下，会有如下问题：

![float](https://images2015.cnblogs.com/blog/1144006/201705/1144006-20170513095302926-2090422648.png)

float 查看是否还能继续往右放，不能就换行放。而换行时并不一定找到最左边，才会出现上图情况。

但如果使用 inline-block，可以做到比较理想的底部对齐的效果

![inline-block](https://images2015.cnblogs.com/blog/1144006/201705/1144006-20170513095312144-2092703645.png)

### inline-block 的空隙

这个空隙是由于 HTML 中的换行符引起的。换行符会被浏览器解析为空白字符，添加到 HTML 中。但我们为了代码可读性，又不可能把所有 HTML 写到一行中。

解决方法是给父元素设置 font-size: 0，可以避免空隙。

## &lt;img&gt;

`<img>` 是一个行内元素(inline)，但却可以设置宽高。

严格来说，它是一个行内替换元素(replaced inline element)，`<input>`,`<select>`, `<textarea>` 等也是行内替换元素。
所谓替换元素，是指从源代码来看是看不到元素本身的内容的（图片就只有一个URL，而`<input>` 就要根据 type 来决定到底渲染什么）。与之相对，`<p>` 就是非替换元素，它的内容直接写在标签内部。

替换元素虽然是行内元素，但是可以使用 margin/padding/width/height，比较类似于 inline-block。