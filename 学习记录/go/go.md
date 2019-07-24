# GO

https://golang.org/

使用安装包后 go 会被自动安装到 /usr/local/go 目录中。**这个目录不能被指定为 workspace**

这个目录中的 src 目录中列出了所有内置包的源码，例如 fmt, net/http 等等。

安装完成后，还有两个路径配置步骤：

1. export GOBIN=$(go env GOPATH)/bin。这是为了解决在使用 go install xxx/xxx.go 时出现 GOBIN not set 的错误
2. export PATH=$PATH:$GOBIN，将 bin 添加到 PATH 中，这样可以直接使用 bin 中 install 的命令。

https://golang.org/doc/code.html


## 文件结构

1. 所有 GO 的代码都放在一个 workspace 里面（区别于常规的语言）。
2. workspace 的默认路径是 $HOME/go，也就是 ~/go。使用命令 go env GOPATH 可以查看当前workspace 路径所在。这个路径可以通过设置 GOPATH 来修改 (export GOPATH = xxx)
3. 在 workspace 里面通常有两个目录，为 bin 和 src
4. src 里面再分项目，例如 src/myFirstGo, src/example 等。bin 内部存放的是可执行的文件。
5. 版本管理以 src 内部的项目为单位。例如 src/myFirstGo/.git, src/example/.git 等。

## 编译执行

1. 执行 go build xxx.go 会进行编译。如果文件中包含 main 函数则会在平级生成可执行文件；如果不包含则不生成实际文件，而是放入到 pkg 目录中。
2. 如果是可执行文件执行 go install xxx.go，会执行编译并把编译产生的可执行文件移动到 bin 目录下。这要求这个文件处在 workspace 下，否则会报错。

## 导出和引入

1. 使用 `package xxx` 来声明包的名字，必须在 GO 的第一行。所有同目录下的文件必须使用同一个名字。这个名字在 import 时作为最后一个部分。**可执行文件必须使用 package main**。包名可以重复，但和路径拼接后不能重复，因为 import 的路径必须是唯一的。
2. 使用 `package main` 并包含 `main` 函数来指定这是一个可执行的程序，也就是入口。不包含 main 的则作为 lib（pkg）使用。
3. 作为 pkg 时，导出的函数**必须大写开头！！**，如果小写开头的方法则被认为是内部方法，只能在同一个包内部使用。要使用导出方法之前，需要先 go build。
4. 引入时使用路径名字。例如 `import "myLib/util"` 指向 src/myLib/util/ 目录。之后使用例如 util.Foo() 来调用 util 目录下定义了 Foo 函数的那个 go 文件。**GO 中 import 的是包，而不是某个 go 文件**，调用的具体方法才是文件导出的内容，GO 会自动搜索，这点和 Java 是一致的。
5. 同一个包内部的方法可以直接使用，不需要 import，也不受限于开头是大写还是小写。

## 测试

GO 内置了测试模块，叫做 `testing`。测试以文件为单位，需要命名为 `*_test.go`，例如 hello.go 和 hello_test.go。在测试文件中，需要包含 TestXXX 这样的函数，签名为 `func (t *testing.T)`。

如果在测试代码中调用了 `t.Error` 或者 `t.Fail` 则视为测试失败。

使用命令 `go test myLib/util` 来进行测试。测试的单位也是包，而不是 go 文件。

## 使用远程包

使用 `go get github.com/golang/example/hello` 来获取远程包。这个命令后会自动下载，编译和安装，所以之后可以直接用命令来执行。

GO 包没有一个特定的发布流程，只要放在一些公有服务器上可以被访问即可，例如 github 就可以。

https://go-search.org/  这个网站可以查询到所有的 GO 包。另外顶部的 AddPackages 可以添加自己的包。

https://github.com/golang/go/wiki/PackagePublishing 这是官方的发布WIKI

## 语法简述

1. 方法的参数和返回值需要指明类型，而且类型放在方法名或者参数**之后**。这点和 typescript 相同（不过不使用 `:`），而和 Java, C 等不同。

    ```go
    func add(x int, y int) int {
        return x + y
    }
    ```

    另外，如果多个参数或者返回值类型相同，则之前的可以省略，最后一个不能省略。例如上述方法可以等价为：

    ```go
    func add(x, y int) int {
        return x + y
    }
    ```

2. 使用 `,` 可以轻松地操作多个值，类似 JS 中的数组或者对象解构。例如

    ```go
    func swap(x, y string) (string, string) {
        return y, x
    }

    var a, b = "hello", "world" // 有初始值，不必再写 a, b 的类型，因为可以推倒得到。
    a, b := swap(a, b)
    ```

3. 方法的返回值也可以命名，之后使用赋值的方式，最后使用没有返回值的 `return` 就可以返回刚才命名的返回值。例如

    ```go
    func divide2(num int) (x, y int) {
        x = num / 2
        y = num - x * 2
        return
    }
    ```

4. `:=` 表示简短赋值，等价于声明+赋值，因此可以代替 `var`。**只能用在方法内部**。这个操作符不用带类型，因为可以从右侧的值中推导。

5. GO 中只有 for 循环，没有 while 或者 do。另外和其他语言不同，GO 的 for 后面不写括号。另外如果省略了初始化部分和累加部分，只剩下判断条件部分，for 此时就和 while 等价了，例如：

    ```go
    sum := 1
    // 等价于常规语言的 while
    for sum < 1000 {
        sum += sum
    }
    // 等价于常规语言的 while(true)
    for {
        sum++
    }
    ```

6. if 和 for 一样不需要写括号。另外 if 可以在条件之前添加一个简短的语句来执行（通常用 `:=`）。这个语句声明变量的作用域只在 if 和与之配对的 else 中。**else 和 if 最后的 `}` 必须在同一行**

    ```go
    if tmp := num / 2; tmp > 100 {
        // do something
    }
    ```

7. GO 中的 switch 不需要在每个 case 最后写 `break`，它只会执行一个分支。如果确实想执行以后的所有分支，需要在这个 case 的最后一句写 `fallthrough`。switch 也支持再判断变量之前加一句临时赋值语句。

    switch 可以省略条件，等价于 `switch(true)`，这样下面的 case 中哪个为真就可以被运行，是一种 if-elseif-else 的变体。

8. 方法中可以使用 `defer` 关键词，在它后面跟随的语句会在这个方法执行完成后再执行，例如

    ```go
    func main() {
        // 最后执行
        defer fmt.Println("world")
        fmt.Println("hello")
    }
    // hello
    // world
    ```

    如果方法中有多个 `defer`，会采用**栈**的方式先进后出，例如

    ```go
    func main() {
        fmt.Println("counting")

        for i := 0; i < 10; i++ {
            defer fmt.Println(i)
        }

        fmt.Println("done")
    }
    // counting
    // done
    // 9 8 7 6 5 4 3 2 1 0
    ```

    defer 一般用来处理方法执行后的回收工作，例如某个方法打开文件但是之后报错，那文件就不会被关闭。比较推荐的写法是在打开文件后马上使用 defer 并关闭文件，这样能确保文件被关闭。

    另外 defer 的语句虽然延迟执行，但这个语句中的参数的值是在运行过程中就确定的，并不等到最后。可以理解当时就为把参数记录下来了。

    最后，defer 语句可以用来修改返回值。例如

    ```go
    func c() (i int) {
        defer func() { i++ }()
        return 1
    }
    // c() returns 2
    ```

9. 星号加类型表示指针类型，如 `*int`，这样可以用来声明变量。另外可以使用 `&` 加具体的值来获得这个值的地址，使得左边成为一个指针。

    在指针前使用 `*` 表示获取这个地址的内容，可以使用（获取值）也可以赋值。赋值后原来变量的值也跟着一起变化（因为本质上内存中就存在这里）

    ```go
    i := 42

    var p *int // 没有初始化，值为 nil
    q := &i // 获取 i 的地址，这样 q 也是一个指针，并且已经被赋值

    fmt.Println(*q) // 42
    *q = 21 // 通过修改指针
    fmt.Println(i) // 21
    ```

    概括来说，`&` + 变量 返回一个指针，相当于取地址的操作；`*` + 指针 返回一个变量，相当于寻址的操作。

10. 结构体 `struct` 相当于 JS 中的对象，是自定义类型。使用 `type` 配合可以定义一个新的类型。

    ```go
    type Vertex struct {
        x, y int
    }

    v := Vertex{1, 2}
    v.x = 4

    p := &v // p 是一个指向结构体的指针
    (*p).y = 8 // 等价于 v.y = 8
    p.y = 8 // 为了书写简便， GO 允许直接使用 p.y，等价于 (*p).y

    v2 := Vertex{y: 10} // 没有指定 x, x 默认为 0
    ```

    也可以单独使用结构体，相当于一次性的类型，之后不再复用了。

    ```go
    object := struct {
        num int
        isEven bool
    }{2, true}
    ```


11. 在类型**前**使用 `[n]` 来声明数组，这点和 Java 相反。例如 `[10]string` 表示长度为10的字符串数组类型。GO 的数组长度不可改变。获取值时依然还在右边，如 `arr[0]`。在初始化时可以直接使用大括号，如 `a := [6]int{1,2,3,4,5,6}`。

    因为数组长度固定的限制，实际使用中我们大多使用“切片”(Slices)，写法是 `[]string`。它是数组中的一段，使用 `a[low:high]` 定义的切片包含 low 但是不包含 high。如 `s := a[1: 4]` 包括了数组 a 的第2，3，4个元素。切片没有新建数据，他是描述数组中的一段，所以可以理解为**数组的指针**，更改切片中的元素也会更改数组中的元素，所有针对同一个数组的切片也都会同步被修改。

    使用 `a := []int{1,2,3}` 相当于先构建一个数组，再构建包含它所有元素的切片。

    `len(a)` 表示切片的长度，是指切片中包含的元素个数。`cap(s)` 表示切片的容量，是指切片第一个元素映射到原始数组中的位置往后总共的元素数量（包括第一个）。

    对切片重新切片可以重新定义它的上下标，返回一个新的切片。例如

    ```go
    s := []int{1,2,3,4,5,6}
    s = s[:0] // 等价于 s[0:0]，则 s 中没有元素，起始位置对应到原数组中的第一个元素，cap(s) = 6
    s = s[:4] // 虽然看上去是对 s 这个空切片重新切，但实际上是针对原始数组重新切，这下 s = [1,2,3,4], cap(s) = 6
    s = s[2:] // 和上面一样，只是这次是减少元素，丢掉前面2个，所以 s = [3,4], cap(s) = 4
    ```

    可以使用内置的函数 `make` 来创建切片，用于自定义它的 len 和 cap。语法 `s := make([]T, len, cap)`，例如

    ```go
    s := make([]int, 0, 5) // len(s) = 0, cap(s) = 5
    // make 的第三个参数也可以省略，这样 cap = len

    s = s[:cap(s)] // len(s) = 5, cap(s) = 5
    s = s[1:] // len(s) = 4, cap(s) = 4
    ```

    `append` 方法可以向切片中追加元素。如果追加了元素之后容量超过原来切片的容量，则创建一个新的底层数组，并创建一个新的切片，把值都复制过去，再追加。所以 cap 会增加。`s = append(s, 2, 3, 4)`

    `range` 可以用来对切片进行遍历，通常用在 for 循环中，例如

    ```go
    number := []int{1,2,3,4,5}
    sum := 0
    // range 返回的第一个值是 index，但这里没有使用，所以需要设置为 _，否则编译器会报错。
    for _, value := range number {
        sum += value // 这个 value 只是副本，改变它不影响切片或者数组的值
    }
    // sum = 15
    ```

12. `map` 关键词用以创建映射关系，通常和 `make` 配合使用，例如

    ```go
    m := make(map[string]int) // 初始化
    m["Hello"] = 1 // 赋值

    n := map[string]int{
        "Hello": 1,
    } // 初始化带赋值

    delete(n, "Hello") // 删除映射
    value, ok := n["Hello"] // value = 0, ok = false

    ```

13. GO 中没有类，但是可以使用结构体和类型来定义方法(method，区别于 function)，模拟类或者对象的行为。

    ```go
    type Vertex struct {
        x, y: float64
    }

    // 可以理解为定义在 Vertex 类上的实例方法 Abs
    // (v Vertex) 可以理解为一种特殊的参数，称为 receiver
    func (v Vertex) Abs() float64 {
        return math.Sqrt(v.x * v.x + v.y * v.y)
    }

    func main() {
        v := Vertex{3, 4}
        v.Abs() // 5
    }
    ```

    **注意**：只能为同一个包内定义的类型创建接受者方法，不能为其他包的类型（包括基本类型）创建接受者方法。