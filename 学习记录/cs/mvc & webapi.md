# .NET Core

ASP .NET Core 是在 ASP .NET 4.x 的基础上重新设计的一套 WEB 工具，它可以完成很多任务。这里主要讲两个部分：

1. WebAPI - 只提供后端 API 服务（不带页面）
2. MVC - 适合搭建完整站点（前后端均有）。.NET 中还有两个工具涉及前端页面：
    1. Blazor - 允许在浏览器中编写 C# 从而前后端都用 .NET
    2. Razor Pages - 快速开发前端页面（类似组件一样的性质）

在 .NET Core 的下一个版本（6）里面这些模板都会合成同一个，称为 MVC6，到时候就都一样了。

## WebAPI

实际上现在前端多为 SPA 有专门的前端框架，因此可能还是 WebAPI 类型的后端项目更加常见。[参考](https://docs.microsoft.com/zh-cn/aspnet/core/tutorials/first-web-api?view=aspnetcore-3.1&tabs=visual-studio-code)

不过 MVC 也可以用来做 API，只是是线上有细微的差别（参数上）。

### 初始化项目

1. 创建项目 `dotnet new webapi -o webApi`
2. 添加依赖 Microsoft.EntityFrameworkCore.SqlServer 和 Microsoft.EntityFrameworkCore.InMemory (后面会用到)
3. 使用 `dotnet run` 运行，访问 `http://127.0.0.1:5000/WeatherForecast` 能看到一个 JSON 就可以。（ MAC 上有 https 问题的，修改 launchSettings.json 中 applicationUrl 的配置，把 https 去掉即可）

### 大致结构

1. 主程序在 Program.cs 中。在 `CreateHostBuilder` 方法中通过 `webBuilder.UseStartup<Startup>();` （泛型类型）的方式使用了 Startup.cs
2. 在 Startup.cs 中提供 `ConfigureServices` 和 `Configure` 两个方法。这两个方法会被 runtime 调用。依赖注入就在 `ConfigureServices` 里面。
3. 在 `ConfigureServices` 里面通过 `services.AddDbContext<TodoContext>()` 来添加上下文。这里引用了我们定义的数据库上下文 TodoContext。

### 添加 Model

以下均放在 Models 目录中。

1. 定义数据结构和一些操作。

    ```cs
    public class TodoItem
    {
        public long Id { get; set; }
        public string Name { get; set; }
        public bool IsComplete { get; set; }
    }
    ```

2. 添加数据库上下文（用来做 mapping 的）

    ```cs
    using Microsoft.EntityFrameworkCore;

    namespace TodoApi.Models
    {
        public class TodoContext: DbContext
        {
            public TodoContext(DbContextOptions<TodoContext> options)
                : base(options)
            {
            }

            public DbSet<TodoItem> TodoItems { get; set; }
        }
    }
    ```

### 依赖注入

先添加 nuget 包：Microsoft.VisualStudio.Web.CodeGeneration.Design 和 Microsoft.EntityFrameworkCore.Design。

在 Startup.cs 中添加注入 (DI) 代码：

```cs
using Microsoft.EntityFrameworkCore;
using webApi.Models; // 引用刚才的 Models 目录
```

在 ConfigureServices 方法中添加

```cs
// 添加上下文
services.AddDbContext<TodoContext>(opt =>
    // 指明使用内存中的数据库
    opt.UseInMemoryDatabase("TodoList"));
```

注：在其他情况时，也可以使用 `services.AddSingleton<SomeService>()` 方法来注册服务。这个 `SomeService` 可以是 MongoDB 的读写服务，也可以是读取 appsettings.json 的配置文件服务接口。本质上 `AddDbContext` 背后也是把 dc 注册成为了 service。

**注册 service 的过程就是依赖注入的过程。**被注册为注入的类型会由系统自动实例化，然后出现在 Controller 的构造函数中。例如这里的 TodoContext 类型由系统实例化，然后传给了 TodoItemsController 的构造函数，也就成了那边的 `context` 对象。

### 添加控制器

通过工具自动新建控制器：

```shell
dotnet tool install --global dotnet-aspnet-codegenerator
# 上一步结束可能要求添加 PATH
dotnet aspnet-codegenerator controller -name TodoItemsController -async -api -m TodoItem -dc TodoContext -outDir Controllers
```

1. 这会创建 Controllers/TodoItemsController。命令中制定了 model(-m) 和 dataContext(-dc)。model 会变成 URL 的一部分来自动生成增删查改的几条路由以及处理；dataContext 会创建一个对应类型的 private 实例 `_context`。在每一个路由的处理函数中基本都会用到这个 `_context`（因为操作DB就要用它）

2. Controller 在类声明的上方使用 `[ApiController]` 来标记，这表示这个控制器响应 Web API 请求。拥有这个标记的类也会在 Startup.cs 的 `services.AddControllers();` 方法执行时被系统收纳进去，从而使路由匹配生效。

3. 在上面还有一个 `[Route("api/[controller]")]` 指明了基础路由。下面每个方法的 `[HttpGet("{id}")]`, `[HttpPut]` 等都是基于这个路由之上。`[controller]`会被这个类的名字 TodoItemsController 代替，但会省略后缀，因此等价于 TodoItems。**路由不区分大小写**。

4. 对于具体的处理方法而言，特性中的 `"{id}"` 和参数一一对应。返回值类型比较复杂，会被包裹好多层。首先如返回特定 ID 的 handler 返回值是 model 类型，即 TodoItem；而返回列表的则是 `IEnumerable<TodoItem>`。然后使用 ActionResult 包一层，最后再用异步 Task 包一层，就成了 `Task<ActionResult<TodoItem>>` 或者 `Task<ActionResult<IEnumerable<TodoItem>>>`。对于没有返回值的类型（视为操作），则直接使用 `Task<IActionResult>` 作为返回值。不过不论如何，实际返回值都是最里层的那一个。

5. 所有的处理函数本质都在调用 `_content.TodoItems` 进行操作/返回。如果是修改，还需要调用 `await _context.SaveChangesAsync();`，这点和 mongoose 的 `save()` 类似。

6. BaseController 提供了 Ok 方法，可以使用 `Ok(res)` 来返回结果，算作一个快捷方式。

7. POST 请求中使用了 `CreatedAtAction` 方法：`return CreatedAtAction(nameof(GetTodoItem), new { id = todoItem.Id }, todoItem);`。这个方法的流程是：
    1. 如果操作成功，返回状态码 201 （POST 的标准响应）
    2. 向 header 中添加 Location，指向刚才新增的元素。
    3. 使用 GetTodoItem (GET 操作的 handler) 创建头信息。

    这样返回的响应中，Headers.Location 会是 http://localhost:5000/api/TodoItems/1 （即访问这个对象的路由，也是 GetTodoItem 对应的路由）

8. 方法被标记为 `[NonAction]` 的为非操作，这些方法不会被映射到路由上。

### DTO

DTO = Data Transfer Object，表示前后端传递的数据对象，是整个模型的子集（因为某些字段不需要/不方便传递到前端去）。

我们可以在 Model/TodoItem 中新增一个 Secret 字段（表示不想传递到前端去），然后再创建一个 TodoItemDTO，和原先 Model 一样有 3 个字段。

之后在 Controller 中，所有使用 TodoItem 的地方都替换为 TodoItemDTO （主要是返回值中的类型以及返回值，为此还需要创建一个转换函数从 TodoItem 实例转换为 TodoItemDTO 实例）

### 静态文件

在 Startup.cs 的 `Configure` 方法中，需要在靠前的地方添加 `app.UseDefaultFiles()` 和 `app.useStaticFiles()`，这样就可以让 .NET CORE 返回静态文件。静态文件必须位于 wwwroot/ 目录下。**注意两者的调用顺序不能颠倒。**

举例来说，现在有 wwwroot/images/banner.jpg，它的访问路径应该是 http://localhost:5000/images/banner.jpg。

如果要为 wwwroot 之外的（多数是平级的）目录提供静态服务，则应当在 Startup.Configure 方法中如下写：

```cs
app.UseDefaultFiles();
app.UseStaticFiles(); // 为 wwwroot 服务，不要删。

// 以下为 MyStaticFiles 服务。
app.UseStaticFiles(new StaticFileOptions
{
    FileProvider = new PhysicalFileProvider(
        Path.Combine(Directory.GetCurrentDirectory(), "MyStaticFiles")),
    RequestPath = "/StaticFiles"
});
```

这样可以通过 http://localhost:5000/StaticFiles/images/anotherBanner.jpg 来访问位于 MyStaticFiles/images/anotherBanner.jpg。（ MyStaticFiles 和 wwwroot 平级）

StaticFileOptions 还可以用来设置 header 信息：

```cs
var cachePeriod = env.IsDevelopment() ? "600" : "604800";
app.UseStaticFiles(new StaticFileOptions
{
    OnPrepareResponse = ctx =>
    {
        // Requires the following import:
        // using Microsoft.AspNetCore.Http;
        ctx.Context.Response.Headers.Append("Cache-Control", $"public, max-age={cachePeriod}");
    }
});
```

## MVC

1. 创建项目： `dotnet new mvc -o mvc`
2. 信任 HTTPS 开发证书：`dotnet dev-certs https --trust`
3. 跨平台的适用于 ASP .NET Core 的 Web 服务器 Kestrel。[详细信息和配置项](https://docs.microsoft.com/zh-cn/aspnet/core/fundamentals/servers/kestrel?view=aspnetcore-3.1)。默认的配置文件位于 Properties/launchSettings.json。（例如端口号）