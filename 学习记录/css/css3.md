# matrix

## box-shadow

box-shadow: h-shadow v-shadow blur spread color inset

blur: 模糊距离
spread: 阴影尺寸
inset: 默认 outset 为外阴影。可以改成 inset 为内阴影

## transition

transition: property duration timing-function delay

## animation

animation: name duration timing-function delay iteration-count direction

iteration-count: 播放次数，可以 infinite
direction: 默认 normal 正常播放，可以改成 alternate，表示在 infinite 的情况下，一次正常播放，一次反向播放，轮流进行。不设置为 infinite时这个属性无效
animation-fill-mode: **不在 animation 简写中，需要单独列出，但也十分有用**。用来确定播放完毕后应该停止在什么状态，默认 none 表示不改变状态（即动画播放完恢复播放开始之前的状态）。backwards 表示动画开始之前应用第一帧的状态，forwards 表示动画结束后固定在最后一帧的状态，both 就是两者都保持。

有一个比较特别的 timing-function 叫做 steps。他可以让动画运行成阶梯状，在使用雪碧图通过 background-position 变化来实现动画的时候尤其有效。
steps 中的 start 表示直接从第 1 个 step 开始，而 end 表示从第 0 个 step （也就是初始状态）开始。
举例来说，动画总计是移动 100px，持续时间 10 秒。steps(9, start) 表示直接从 10px 开始，到 1 秒后移动到 20px，以此类推。而 steps(9, end) 表示从 0 开始， 1 秒后移动到 10px。
所以 end 如果配合 forwards，就可以实现一个完整的移动效果了。

```css
/* 通常 steps end 和 forwards 配合来实现组合效果 */
animation-timing-function: steps(20, end);
animation-fill-mode: forwards
```

steps 还可以详见[这里](http://www.divcss5.com/css3-style/c50603.shtml)

animation 还要配合 @keyframes 一起使用

## matrix

transform 中的 translate, rotate, skew, scale 背后都是由 matrix 来实现的。

transform-origin 可以移动坐标系的中心点

假设 transform: matrix(a,b,c,d,e,f)，对于图形的每一个坐标点 (x,y) 都会变化为

a c e   x   ax + cy + e
b d f * y = bx + dy + f
0 0 1   1        1

ax+cy+e 是变化后的横坐标，bx+dy+f 是变化后的纵坐标。根据 a,b,c,d,e,f 不同的取值，就可以实现上述的四种变化，更可以实现其他混合的变化。这四种变化只是 matrix 的四种特殊情况，用来方便程序员使用而已。

举例来说，如果 a=1, b=0, c=0, d=1, e=30, f=30，那么 x => ax + cy + e = x + 30; y => bx + dy + f = y + 30。所以实际上是向右向下各平移30px，等价于 translate(30px, 30px)

### matrix vs translate

如上述例子，matrix(1, 0, 0, 1, x, y) 等价于 translate(x, y)

### matrix vs scale

matrix(x, 0, 0, y, 0, 0) 等价于 scale(x, y)

### matrix vs rotate

设旋转角度为a

matrix(cos(a), sin(a), -sin(a), cos(a), 0, 0) 等价于 rotate(a) (CS-SC)

x => cos(a)\*x - sin(a)\*y
y => sin(a)\*x + cos(a)\*y 

### matrix vs skew

设倾斜角度为a

matrix(1, tan(b), tan(a), 1, 0, 0) 等价于 skew(a, b)

这个计算有点复杂，就不列了。

### matrix 的作用

如果是不在这四种变化之中的变化，就需要使用 matrix 了。例如 __镜像对称__。 左右翻转为 matrix(-1, 0, 0, 1, 0, 0)，上下翻转为 matrix(1, 0, 0, -1, 0, 0)
