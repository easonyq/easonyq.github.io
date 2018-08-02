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
