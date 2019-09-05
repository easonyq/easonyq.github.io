# GO Web Server

使用 `"net/http"` 内置包来启动 WEB 服务器。

**注意**：GO 的网络代码直接监听了 TCP 端口，可以取代 nginx，所以使用 GO 时不再需要在前面假设 nginx 或者 apache 等服务器。

[Go Web 示例](https://gowebexamples.com/)

## 示例代码

[参考](https://github.com/easonyq/build-web-application-with-golang/blob/master/zh/03.3.md)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	for k, v := range r.Form {
		fmt.Println("---------------")
		fmt.Println("key: ", k)
        fmt.Println("val: ", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello World!")
}

type MyHandler struct{}

func (handler MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 只处理 /handler 路由
	fmt.Fprintf(w, "This is My Handler")
}

func main() {
    http.HandleFunc("/", sayHello)
    http.Handle("/handler", MyHandler{})

	fmt.Println("Listen to 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Listen error: ", err)
    }
}
```

一些参数说明：

1. `r.Form` 中保存的是 URL 中的参数 (query)。例如当访问 `localhost:8080/?a=1&b=2` 时，这里面是 `map[a:[1] b:[2]]`。可以通过下面的 `for range` 逐个取出来。

2. 在 `for range` 中，`r.Form` 循环中的 `v` 的类型是 `[]string`。所以通常使用 `strings.Join(v, "")` 来转为普通的字符串。

3. 如果 URL 中的 query 重复，例如 `/?a=1&a=2`，则打印出的 `v` 等于 `"12"`

4. GO 源码显示，如果存在多个路由（多次 `http.Handle` 或者 `http.HandleFunc`)且都能匹配 path，GO 会选择**len最长**的匹配。因此注册顺序和最终的匹配顺序无关。

## 处理流程

1. 调用 `http.ListenAndServe` 之后，底层用 TCP 协议启动一个服务，监控 8080 端口。

2. 使用 `for {}` 无限循环，通过 Listener 接收请求 (`Listener.Accept()`)，然后创建一个连接 (`*Server.newConn`)，最后使用多线程让这个连接进行服务 (`go c.serve()`)。这样每次请求都在一个子线程中，互不影响。

3. 连接处理的代码中，先解析请求 `c.readRequest()` 获取对应的 handler。例子中传入了 `nil`，则使用默认的 handler，名为 `DefaultServeMux`。这个默认 handler 会根据上面的 `http.HandleFunc(route, handler)` 进行匹配和调用。

4. 进入我们注册的 handler，也就是 `handler` 方法，获取了一些参数，最后给 `w` 写入了返回。

5. 总结：ListenAndServe -> TCP 服务 -> `srv.Serve(l net.Listener)` -> `rw = l.Accept()` -> `c = srv.newConn(rw)` -> `go c.serve()` -> 根据 path 找到我们的 handler 并执行

![整体流程](https://github.com/easonyq/build-web-application-with-golang/raw/master/zh/images/3.3.illustrator.png?raw=true)

## ServeMux

[参考](https://github.com/easonyq/build-web-application-with-golang/blob/master/zh/03.4.md)

例子中 `http.ListenAndServe` 第二个参数传入了 `nil`，因此使用了默认的 DefaultServeMux，它的类型是 ServeMux，是 GO 内部默认实现的。因此我们也可以自定义实现 ServeMux。

```go
type ServeMux struct {
	mu sync.RWMutex   // 锁，由于请求涉及到并发处理，因此这里需要一个锁机制
	m  map[string]muxEntry  // 路由规则，一个string对应一个mux实体，这里的string就是注册的路由表达式
    hosts bool // 是否在任意的规则中带有host信息
    es []muxEntry // 排序后的 entry 列表，从长到短，供匹配多个路由时使用
}

type muxEntry struct {
	explicit bool   // 是否精确匹配
	h        Handler // 这个路由表达式对应哪个handler
	pattern  string  //匹配字符串
}

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)  // 路由实现器
}
```

我们需要格外注意一个叫做 `Handler` 的类型，它有两种用法。在使用 DefaultServeMux 时，它可以作为 `http.Handle` 方法的参数，也可以当做 自定义的 ServeMux 而作为 `http.ListenAndServe` 方法的第二个参数。

当使用 DefaultServeMux 时，`http.HandleFunc` 把参数方法进行类型转化，使之拥有了 `ServeHTTP` 方法，才能正常工作。如果直接使用 `http.Handle`，那么参数必须拥有这个方法。一般在这个方法中处理自己负责的已经被 pattern 命中的路由，如开头的例子。

当使用自定义 ServeMux 时，这个 `ServeHTTP` 方法就得处理所有路由了。在特定的路由直接调用特定的方法，参数也自己传入，如下面的例子。

```go
package main

import (
	"fmt"
	"net/http"
)

type MyMux struct {}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 只服务精确匹配
	if r.URL.Path == "/" {
		sayhelloName(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello myroute!")
}

func main() {
	http.ListenAndServe(":9090", MyMux{})
}
```

## 处理表单，字段校验和文件上传

[参考](https://github.com/easonyq/build-web-application-with-golang/blob/master/zh/04.0.md)

GO 拥有自己的 HTML 模板，扩展名为 gtpl。在 `import "html/template"` 引入后，使用如下代码可以写入到响应中

```go
import "html/template"

func handler(w http.ResponseWriter, r *http.Request) {
    if r.Method === "GET" {
        // GET 请求返回页面
        // 1. 解析模板
        t, _ := template.ParseFiles("page.gtpl")
        // 2. t.Execute(w, nil) 把解析后的模板输出到 w 中
        log.Println(t.Execute(w, nil))
    } else {
        // POST 处理表单提交信息
        // 获取 Form 数据前必须 ParseForm，这不是默认进行的！
        r.ParseForm()
        r.Form["username"] // r.Form[key] 返回类型为切片，[]string
        r.Form["password"]

        // 下面两种写法不用事先调用 r.PraseForm()，但是当存在同名参数时只返回第一个，如不存在返回空字符串
        r.FormValue("username")
        r.Form.Get("username")
    }
}
```

注意：如果使用 axios.post 发送请求，需要额外查看 content-type。如果是 application/json 的话，使用 `r.ParseForm()` 无法获取数据。应当改用 `decoder := json.NewDecoder(req.Body)` 然后 `decoder.Decode(&data)`，最后 `data.MyField` 获取。

gtpl 文件里可以包含一些特定的模板语法，例如

```html
<!-- example.gtpl -->
<input type="hidden" name="token" value="{{.}}">
```

填充时的代码如下

```go
import (
    "time"
    "crypto/md5"
    "io"
)

current := time.Now().Unix()
h := md5.New()
io.WriteString(h, strconv.FormatInt(current, 10))
token := fmt.Sprintf("%x", h.Sum(nil)) // token 可以存到 session 中

t, _ := template.ParseFiles("example.gtpl")
t.Execute(w, token) // w http.ResponseWriter
```

### 几个有用的正则

```go
import "regexp"

// 验证汉字
if m, _ := regexp.MatchString("^\\p{Han}+$", r.Form.Get("realname")); !m {
	return false
}

// 验证英文字母
if m, _ := regexp.MatchString("^[a-zA-Z]+$", r.Form.Get("engname")); !m {
	return false
}

// 验证邮箱
if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,})\.([a-z]{2,4})$`, r.Form.Get("email")); !m {
	fmt.Println("no")
} else {
	fmt.Println("yes")
}

// 验证手机号
if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, r.Form.Get("mobile")); !m {
	return false
}
```

### 日期和时间

所有时间都需要使用常量作为单位，例如 `time.Sleep` 的参数，可以是 `100 * time.MilliSecond` 也可以是 `time.Second` 等。

```go
import "time"

t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
fmt.Printf("Go launched at %s\n", t.Local())
```

### 转义

为了防止用户输入的内容中包含攻击内容，可以使用 `template.HTMLEscapeString` 方法，它会把 `<` 变成 `&lt;`。

```go
import "html/template"

// 获取到变量
username := template.HTMLEscapeString(r.Form.Get("username"))

// 写入到响应
template.HTMLEscapeString(w, []byte(r.Form.Get("username")))
```

## gorilla/mux

以上示例全部使用 GO 原生的 net/http 包。github 上还有一个比较有名的第三方库叫做 [gorilla/mux](https://github.com/gorilla/mux) 也很好用。它的优势在于

1. URL pattern 中可以包含命名参数（例如 `/user/{id}/{operation}`），使用 `mux.Vars(r)["id"]` 来获取。

2. 另外在配置路由时就支持只针对某个方法（例如 GET），不必再到 handler 里面去判断。`r.HandleFunc(pattern, handlerFoo).Methods("GET")`

3. 限制其他内容，例如自身的域名，http/https 协议等

4. 支持子路由，例如

    ```go
    r := mux.NewRouter()
    bookRouter := r.PathPrefix("/books").Subrouter()
    bookRouter.HandleFunc("/", AllBooks)
    bookRouter.HandleFunc("/{title}", GetBook)
    ```

此外 [gorilla/session](https://github.com/gorilla/session) 和 [gorilla/websocket](https://github.com/gorilla/websocket) 也值得关注和使用。

## Session & Cookie

### Cookie

GO 中有 `net/http` 包提供对 cookie 的操作

```go
// 方法签名
http.setCookie(w ResponseWriter, cookie *Cookie)

// Cookie 结构定义
type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string

// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

示例代码如下：

```go
import (
    "time"
    "net/http"
)

// 设置 cookie
expiration := time.Now()
expiration = expiration.AddDate(1, 0, 0)
cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
http.SetCookie(w, &cookie)

// 读取 cookie
cookie, _ := r.Cookie("username") // r *http.Request
// 通过 for range 遍历
for _, cookie := range r.Cookies() {
    fmt.Fprintf(w, cookie.Name)
}
```

### Session

GO 原生没有支持对 session 的操作，因此我们使用 [gorilla/session](https://github.com/gorilla/session) 来处理。

```go
// sessions.go
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func main() {
	http.HandleFunc("/secret", secret)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.ListenAndServe(":8080", nil)
}
```

## Socket

Socket 遵从 Unix “一切皆文件” 的基本思路，和文件一样通过“打开 -> 读写 -> 关闭”的模式来处理，所以 Socket 也是一个文件描述符，是一种特殊的 I/O。

常用的 Socket 有两种：流式 Socket (SOCK_STREAM) 和数据报式 Socket (SOCK_DGRAM)。流式针对面向连接的 TCP 服务应用，数据报式则针对于无连接的 UDP 服务应用。

![TCP/IP协议结构图](https://github.com/easonyq/build-web-application-with-golang/raw/master/zh/images/8.1.socket.png?raw=true)

类似于本地通过 PID 来标识进程，在网络中使用网络层的 IP 地址来唯一标识一台主机，而传输层的“协议+端口”可以唯一标识应用程序，因此，IP+协议+端口就成了可以在网络中唯一标识进程的三元组。

以下内容使用 GO 的原生支持。比较常用的包 gorilla/websocket 也能解决问题。

### 处理 IP

GO 中的 `net` 包定义了很多网络相关的类型，其中 `type IP []byte` 用来表示 IPv4 的类型。下面是一些常用的 IP 相关的操作

```go
// 把一个字符串的IP地址(例如 "127.0.0.1") 转化为 IP 类型
addr := net.ParseIP(ipString)

// 把一个网络地址+端口转化为一个 TCPAddr 类型，包含 IP IP, Port int, Zone string（IPv6 使用） 三个字段。
// 返回类型是 *TCPAddr
var tcpAddr *net.TCPAddr = net.ResolveTCPAddr("tcp", "www.google.com:80")
```

### TCP Client

使用 `net.DialTCP` 来建立一个 TCP 连接，返回一个 TCPConn 类型的对象。当连接建立时服务器端也创建一个同类型的对象，此时客户端和服务器端通过各自拥有的 TCPConn 对象来进行数据交换。一般而言，客户端通过 TCPConn 对象将请求信息发送到服务器端，读取服务器端响应的信息。服务器端读取并解析来自客户端的请求，并返回应答信息，这个连接只有当任一端关闭了连接之后才失效，不然这连接可以一直在使用。建立连接的函数定义如下：

```go
// network: "tcp4", "tcp6", "tcp" 中的任意一个
// laddr: 本机地址，一般是 nil
// raddr: 远程地址，可以由上面的 ResolveTCPAddr 返回
func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
```

一个使用的例子是：

```go
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	// 取输入参数的第二个作为服务端地址
	service := os.Args[1]

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	// 获取连接 conn
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	// 写入RequestHeaders: HEAD / HTTP/1.0\r\n\r\n
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)

	// 读取 conn 中的所有数据
	// 实际使用还需要套一个 for ，直到读到 err = io.EOF 为止
	result := make([]byte, 256)
	_, err = conn.Read(result)
	checkError(err)

	fmt.Println(string(result))
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```

### TCP Server

与 `DialTCP` 对应，net 中还有一个 `ListenTCP` 方法可以用来当做服务端提供服务，它的签名是

```go
func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)
func (l *TCPListener) Accept() (Conn, error)
```

服务器的示例代码如下：

```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	// 监听本地7777端口，提供服务
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		// 使用 listener.Accept 获取连接 conn
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	daytime := time.Now().String()
	conn.write([]byte(daytime))
}
```

在 `for` 循环中，如果出错，并不是 `break` 而是 `continue`，应当使客户端报错退出。

下面的例子会针对不同的客户端返回不同的结果，只列出 `handleClient` 的代码

```go
func handleClient(conn net.Conn) {
	// 设置超时 2 分钟
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	// 读取长度不要设置的太大，防止 flood attack
	request := make([]byte, 128)
	defer conn.Close()
	for {
		// 服务端也使用 conn 来读取客户端的请求数据
		read_len, err := conn.Read(request)
		if err != nil {
			fmt.Println(err)
			break
		}

		if read_len == 0 {
			// 连接被客户端关闭，读不到结果
			break
		} else if strings.TrimSpace(string(request[:read_len])) == "timestamp" {
			// 如果客户端请求 "timestamp"，则返回时间戳
			daytime := strconv.FormatInt(time.Now().Unix(), 10)
			conn.Write([]byte(daytime))
		} else {
			// 返回可读的时间格式
			daytime := time.Now().String()
			conn.Write([]byte(daytime))
		}

		// 清空 request
		request = make([]byte, 128)
	}
}
```

这个例子中，服务器会不断读取客户端的数据，并且连接是不关闭的，即保持长连接。开头的超时设置表示当2分钟内客户端没有数据发送则关闭。

除了上面的 `SetReadDeadline` 之外，常用的控制函数还有

* `func DialTimeout(net, addr string, timeout time.Duration) (Conn, error)` 带超时的 Dail，客户端和服务器都适用。
* `func (c *TCPConn) SetWriteDeadline(t time.Time) error` 在规定时间内如果没有写入消息则关闭连接
* `func (c *TCPConn) SetKeepAlive(keepalive bool) os.Error` 设置 keepAlive 属性。操作系统层在tcp上没有数据和ACK的时候，会间隔性的发送keepalive包，操作系统可以通过该包来判断一个tcp连接是否已经断开，在windows上默认2个小时没有收到数据和keepalive包的时候认为 tcp 连接已经断开，这个功能和我们通常在应用层加的心跳包的功能类似。

### UDP

UDP 的方法和 TCP 基本一致，只是方法的名字都改成了 UDP。另外 UDP 中没有 Accept 函数。

```go
func ResolveUDPAddr(net, addr string) (*UDPAddr, os.Error)
func DialUDP(net string, laddr, raddr *UDPAddr) (c *UDPConn, err os.Error)
func ListenUDP(net string, laddr *UDPAddr) (c *UDPConn, err os.Error)
func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err os.Error)
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (n int, err os.Error)
```

### WebSocket

HTML5 中加入了 WebSocket，允许浏览器和服务器进行全双工通信，保持一个连接不断，持续在这个连接中收发消息。WebSocket URL 的协议是 ws:// 或者 wss:// (SSL，类似于 https)。

GO 中没有原生对 WebSocket 的支持。可以使用 `go get golang.org/x/net/websocket` 来操作。

## 连接数据库

GO 没有官方的数据库驱动，而是为开发数据库驱动定义了一些标准接口，根据这些接口开发者可以开发响应的驱动。[接口参考](https://github.com/easonyq/build-web-application-with-golang/blob/master/zh/05.1.md)

### MySQL

这里以 MySQL 的 github.com/go-sql-driver/mysql 为例，它实现了 database/sql 接口，全部使用 go 编写。

假设有一个数据库名为 test，内部有两张表，为 userinfo 和 userdetail，表结构如下：

```sql
CREATE TABLE `userinfo` (
	`uid` INT(10) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NULL DEFAULT NULL,
	`department` VARCHAR(64) NULL DEFAULT NULL,
	`created` DATE NULL DEFAULT NULL,
	PRIMARY KEY (`uid`)
);

CREATE TABLE `userdetail` (
	`uid` INT(10) NOT NULL DEFAULT '0',
	`intro` TEXT NULL,
	`profile` TEXT NULL,
	PRIMARY KEY (`uid`)
)
```

操作这张表的代码如下

```go
package main

import (
	"database/sql"
	"fmt"
    // _ 符号在包名的面前表示只引入包，但不使用这个包里的变量和方法
    // 这个包已经通过 database/sql.Resiger 接口注册了名为 "mysql" 的驱动，因此在使用 sql.Open("mysql") 的时候就算使用它了，不直接使用它内部的方法和变量。
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "astaxie:astaxie@/test?charset=utf8")
	checkErr(err)

	// 插入数据
	stmt, err := db.Prepare("INSERT INTO userinfo SET username=?,department=?,created=?")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

    fmt.Println(id)

	// 更新数据
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// 查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	// 删除数据
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
```

### Redis

redis 是一个 KV 存储系统，和 Memcached 类似。可以存储的类型包括字符串，链表，集合和有序集合等。GO 支持的 redis 驱动是 github.com/garyburd/redigo。

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	Pool *redis.Pool
)

func init() {
	redisHost := ":6379"
	Pool = newPool(redisHost)
	close()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		}
	}
}

func close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func Get(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}
	return data, err
}

func main() {
	test, err := Get("test")
	fmt.Println(test, err)
}
```

### MongoDB

![Mongo and mySQL](https://github.com/easonyq/build-web-application-with-golang/raw/master/zh/images/5.6.mongodb.png?raw=true)

GO 中最常用的 mongo 驱动是 mongo-go-driver，由 mongodb 官方发布，地址是 [https://github.com/mongodb/mongo-go-driver](https://github.com/mongodb/mongo-go-driver)

[教程](https://vkt.sh/go-mongodb-driver-cookbook/)

[简单用法示例](https://www.hwholiday.com/2018/how_use_mongo-go-driver/)

示例：

```go
// GO 要求 struct 中每个字段都大写开头
// 和 JSON 类似，strcut tag 中的是 DB 中实际存在的字段
// 如果仅仅是首字母大小写的差别，tag 可以不写，例如 DB 中的 name 可以自动对应到 GO 的 Name。但如果有其他变化，就必须要写。
type publishedInfo struct {
	VersionTS   time.Time `bson:"versionTS"`
	PublishedAt time.Time `bson:"publishedAt"`
	Spec        string    `bson:"spec"`
}

type dashboard struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	DraftSpec string             `bson:"draftSpec"`
	Published []publishedInfo    `bson:"published"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://iotviz:iotviz_qianmo_bjyz@bjyz-zhaomuwei.epc.baidu.com:8001/iotviz-local?authSource=iotviz2-prod"))
	checkError(err)

	collection := client.Database("iotviz-local").Collection("dashboard")

	findOptions := options.Find()
	findOptions.SetSkip(2)
	findOptions.SetLimit(2)
	findOptions.SetSort(bson.M{"_id": -1})
	findOptions.SetProjection(bson.M{
		"_id":       1,
		"name":      1,
		"draftSpec": 1,
		"published": 1,
	})

	// 如果使用 ID 查找，一方面可以使用 collection.FindOne
	// 另一方面，不能直接用字符串，而是需要构造一个 ObjectID 类型
	expectedId, _ := primitive.ObjectIDFromHex("5c4961714af86b2a87a51718")
	cur, err := collection.Find(ctx, bson.M{
		"_id": expectedId,
	}, findOptions)
	checkError(err)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result dashboard
		err := cur.Decode(&result)
		checkError(err)

		// fmt.Println(result.ID.Hex())
		fmt.Println(result.Name)
		// fmt.Println(result.Published[0].PublishedAt)
	}
	if err := cur.Err(); err != nil {
		fmt.Println(err)
	}
}
```

## 部署服务

GO 编译后是一个可执行文件，通常使用第三方工具进行管理，例如 Supervisord。具体操作步骤见[这里](https://github.com/easonyq/build-web-application-with-golang/blob/master/zh/12.3.md#supervisord) (例子中的 /data/blog/blogdemon 应该就是 go build 出来的产物)

[如何部署Golang程序到服务器](https://www.jianshu.com/p/bfaba9b6d46d) (`env GOOS=linux GOARCH=386 go build main.go`)