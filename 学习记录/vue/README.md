基于 Vue 版本 2.5.0。2017-10-13 发布。2018-4-16日编写时 master 分支所在版本。

Vue 实例生命周期

![lifecycle](https://cn.vuejs.org/images/lifecycle.png)

声明周期钩子：
* beforeCreate (均有)
* created (均有)
* beforeMount
* mounted
* beforeUpdate
* updated
* beforeDestroy
* destroyed

除了 create 的两个，其他钩子在 SSR 中是没有的。
