# 概述

使用 C# 来编写 webUI，取代 JS。因此是 C# + HTML + CSS。

优点
1. 前端在 webAssembly 上运行，后端在 .NET framework 上运行。前后端同构(C#)，因此可以重用代码和lib。包括自己写的和所有nuget上的包。
2. 后端运行并回传：Blazor 可以选择把前端代码放到后端运行（使用 SignalR 传递前端事件），执行完成后再把UI变化回传到前端并和当前的 DOM 合并，表现出变化。
3. 基于 open web standards，因此可以在所有浏览器工作，包括移动浏览器。
4. C# 前端代码可以调用 JS 的代码，因此所有JS的框架，社区lib都能被使用。 (叫做 Javascript interop) （JS代码能放到后端运行吗？https://docs.microsoft.com/en-us/aspnet/core/blazor/call-javascript-from-dotnet?view=aspnetcore-5.0 ）
5. 自带UI组件库， Telerik, DevExpress, Syncfusion, Radzen, Infragistics, GrapeCity, JQWidgets
6. 开源，免费

创建时有两种模式：
1. Blazor Server App: 在后端运行。使用 SignalR 将前端事件传递到后端。对应上面第2点。
2. Blazor WebAssembly App: 在前端运行。

注意：
razor 本身是一个C#的模板语言，出现比较早，当时是用来增强 V 的一种新写法。
之后出现的 razor views/pages 是一个在ASP .NET Core MVC中取代 V 和 C 的技术。通过它可以直接从 M 获取数据并展现。这两者都遵循传统的网页应用逻辑（每次请求返回全部HTML，跳转页面都会刷新）

Blazor 的主要目的是前后端同构，并实现 SPA（首次请求返回所有组件，后续请求数据更改页面，并不刷新页面）。
有两种模式：Blazor WebAssembly 通过 web assembly 在浏览器里运行 C# (dll)；Blazor Server 使用 SignalR (基于websocket) 完成前后端通信。

所以两者目的不同。但通常情况 blazor 会使用 razor components 进行开发。可以不必把两者区分到这么细节。

# Blazor WebAssembly

![](https://docs.microsoft.com/en-us/aspnet/core/blazor/index/_static/blazor-webassembly.png?view=aspnetcore-5.0)

![](https://docs.microsoft.com/en-us/aspnet/core/blazor/hosting-models/_static/blazor-webassembly.png?view=aspnetcore-5.0)

简单来说就是使用C#, razor components写前端。这种模式并不关心后端实现，代码也都跑在前端（借助 WebAssembly 跑 dll），从角色上说和 Angular, Vue, React 是等价的。

入口 wwwroot/index.html

项目启动时需要下载 .NET Runtime, app code 以及 app 的依赖，因此启动时间较长。另外必须在支持 WebAssembly 的浏览器上运行。

这里又细分为两种：

1. standalone Blazor WebAssembly app: 背后的后端系统不是由 ASP.NET Core 来提供服务的
2. hosted Blazor WebAssembly app: 背后的后端由 ASP.NET Core 提供服务，中间使用 SignalR 进行通讯。

# Blazor Server

![](https://docs.microsoft.com/en-us/aspnet/core/blazor/index/_static/blazor-server.png?view=aspnetcore-5.0)

![](https://docs.microsoft.com/en-us/aspnet/core/blazor/hosting-models/_static/blazor-server.png?view=aspnetcore-5.0)

这个模式在上一个模式的基础上又增加了后端，且前后端同构，是新项目的推荐使用方式。这种模式下前端代码(razor components)实际是在后端运行的（运行完再回传到前端更新HTML），这是两种模式的本质差异。从角色上说他可以跟 VueSSR 等带SSR的方式等价。

和前一种相比（尤其是跟 hosted Blazor WebAssembly app 相比）
1. 它的启动速度会快很多，因为它下载的东西更少。
2. 它能够使用后端的能力（因为代码实际是在后端运行的）
3. 可以在不支持 WebAssembly 的浏览器上运行（也是因为代码是在后端运行的，浏览器只要支持 websocket 即可）

缺点：
1. 每次交互需要到后端运行代码并返回混合HTML，因此需要一次网络请求，存在延时。
2. 无法离线
3. 对服务器压力大（相当于回到了老式的 web service，在 client 端几乎什么都不做）

如果有多个后端提供服务，需要使用 sticky sessions 来确保同一个 client 在断开连接并重连时每次都连接到相同的后端，因为状态是保存在后端的。推荐使用 Azure SignalR Service，有默认支持。

状态可以被永久存储，这样可以防止一些关键数据的丢失（例如购物车或者多步的表单）。可以选择存入数据库，blob等，也可以放在URL或者browser storage里面。

和 ASP.NET Core MVC 类似，有 Program.cs 和 Startup.cs。

在 Startup.cs 的 `ConfigureServices` 中注意

```cs
services.AddRazorPages(); // 引入 razor
services.AddServerSideBlazor(); // 使用 Blazor Server
```

另外在 `Configure` 中

```cs
app.UseEndpoints(endpoints =>
{
    endpoints.MapBlazorHub();
    endpoints.MapFallbackToPage("/_Host"); // 表示使用 Pages/_Host.cshtml 来做最外层的页面
});
```

其他文件还有：
1. App.razor - 前端页面的根组件（但不是根页面），在这里引用所有其他组件构成页面。比如 MainLayout。
    注意：整体的包含关系是：Pages/_Host.cshtml -> App.razor -> MainLayout.razor -> NavMenu.razor & Router body
2. Pages/ - 所有主体页面，有 .razor 也有 .cshtml。其中三个 .razor 就是三个页面（外侧有框架包裹，这里只是页面的主体内容，框架内容在其他地方）
3. Shared/ - layout 相关的组件。其中 MainLayout 就是整个页面的外侧框架，包括左侧sidebar, 上侧navbar和当中的content。这里左侧又单独写成了一个组件 NavMenu.razor，而主体内容则使用 `@Body` 标记留空。

# 框架特性

## 路由和事件

MainLayout 中 @Body 的内容根据当前路由来决定，类似于 react-router 中的 <Switch> 和 <Router> 以及 vue-router 中的 <RouterView>。
每个页面组件的开头会通过使用 `@page "/"` 来声明自己的路由。

页面跳转是通过把事件传递到后端，后端计算后给出新的HTML，并和当前HTML **混合** 而成，因此没有实质跳转，是通过 pushAPI 来更改路由的，并不刷新页面。
页面点击事件同理。（使用C#编写逻辑，而非传统的JS）
此外，SignalR 是通过 websocket 实现的，因此任何点击或者切换页面不会再发送新的请求，而是在页面开始就建立好的path为 `/_blazor` 的 websocket 请求中持续收发 message 来实现的。


通过 `@onclick="Handler"` 来绑定处理函数，可以接受异步函数，返回 Task。处理函数可以通过参数接受事件（和JS相同），完整列表：https://docs.microsoft.com/en-us/aspnet/core/blazor/components/event-handling?view=aspnetcore-5.0#event-argument-types
也可以通过 `@onclick="(e => Console.WriteLine("Hello"))"` 来定义匿名方法。


## 数据驱动

Blazor 的数据流模式和 react 是一样的，都遵循“数据驱动”的模式，即通过创建状态（在 `@code` 中声明变量），绑定状态（在 HTML 标签中使用 `@variable`）并修改状态（给 `variable` 赋值）来完成数据和页面同步变化的。

展示数据：直接使用 `@variable` 进行输出
读取数据：使用 `@bind` 绑定到某个变量

```html
<input placeholder="input your name" @bind="name">
<span>Your name is @name</span>

@code {
    private string name = string.Empty;
}
```

有些组件（例如 checkbox）通过使用 `@bind` 可以同时绑定读写

```html
<input type="checkbox" @bind="flag">

@code {
    private bool flag = true;
}
```

可以使用 `@bind:event="oninput"` 来指定只处理特殊的事件。

针对 DateTime 类型的变量绑定，还可以通过 `@bind:format="yyyy-MM-dd"` 来规定格式。

如果绑定的对象类型和实际输入的类型并不一致，Blazor 会拒绝这次改动，保留原来的值。

## 父子组件的通讯

和 react 等前端框架的思路一致，Blazor 也遵循单向数据流的规则：当子组件需要使用并修改父组件的状态时，应当把父组件的状态更改函数递到子组件。在子组件中调用这些方法来改变状态，而不是直接改变。

one-way flow of data:
Change notifications flow up the hierarchy.
New parameter values flow down the hierarchy.

子组件

```html
<div>Parent state is @Number</div>
<button onclick="OnClickHandler">Trigger a parent handler</button>

@code {
    [Parameter]
    public int Number { get; set; }

    [Parameter]
    public EventCallback<int> ChangeNumber { get; set; }

    private async Task OnClickHandler () {
        await ChangeNumber.Invoke(Number + 1);
        // Number++; is wrong
    }
}
```

父组件

```html
<ChildComponent Number="number" ChangeNumber="ChangeNumber"></ChildComponent>

@code {
    private int number = 0;

    private void ChangeNumber(int newNumber) {
        number = newNumber;
    }
}
```

父子组件之间的绑定和传递：https://docs.microsoft.com/en-us/aspnet/core/blazor/components/data-binding?view=aspnetcore-5.0#binding-with-component-parameters
（相当于子组件定义了两个参数：值和改变值的回调。父组件在调用子组件时先用变量绑定，再传入回调，来达成父子组件的通信）

## DOM Element & ref

Blazor 也有 `@ref` 用于获取 DOM 元素：

```html
<div @ref="wrapper"></div>

@code {
    private ElementReference wrapper;
}
```

这个 `wrapper` 可以调用 C# 方法，也可以传给 JS 操作，效果等同于 HTMLElement。
ref 在 prerendering 时不可使用，因为这时 C# 还没有与浏览器建立连接。

ref 尽量保持只读。如果要写不要在内嵌 C# 代码的元素写，否则会导致 Blazor 内部的状态管理混乱，或者VDOM应用时 DOM 结构匹配不上。

## 组件和引用

所有页面都可以直接使用标签来引用其他组件，所以页面和组件是等价的概念（和react是一样的）。

被引用的组件本身可以是个页面（可以包含 `@page`）。路由不影响它被引用时的效果。

## 组件参数

和 react 一样，组件支持参数（react中称为property）。在组件中声明参数，并在引用组件时通过 HTML attribute 的方式传入。

1. 声明参数
    在组件的 `@code` 部分，声明一个变量并且增加 `[Parameter]` attributre。必须是 public。

    ```html
    <div class="component-wrapper">
        <span>@IncrementAmount</span>
    </div>

    @code {
        [Parameter]
        public int IncrementAmount {get; set; }
    }
    ```

2. 传入参数
    在使用组件时 `<Counter IncrementAmount="10" />` 以 HTML 属性的形式传入。
    这里的 "" 是表达式而不是字符串，因此 `10` 可以直接转成 int 类型。如果输入其他类型，也会被强转。

3. slot
    和 react 类似的 slot 逻辑。通过约定的变量名 `ChildContent` 传入

    ```html
    <div class="component-wrapper">
        <h2>Hello</h2>
        <div>@ChildContent</div>
    </div>

    @code {
        // 类型和变量名都固定
        private RenderFragment ChildContent { get; set; }
    }
    ```

    ```html
    <Component>
        <div>This div will be insert</div>
    </Component>
    ```

4. Cascading values and parameters
    类似于 react 中的 `<provider>` 和 `<consumer>`，可以实现跨层级的参数传递，而不用逐层传递。

    使用 `<CascadingVaule Value="">` 标签来提供值（类似 `<provider>`）。例如在最外层的 Layout.razor 中

    ```html
    <div class="col-sm-9">
        <CascadingValue Value="theme">
            <div class="content">@Body</div>
        </CascadingVaule>
    </div>

    @code {
        private string theme = "some-theme";
    }
    ```

    凡是套在 CascadingValue 中的组件都可以使用这个值。**值是通过类型而非变量名绑定的，类似 .NET CORE 的注入**。

    ```html
    <!-- 这段必须套在刚才的 CascadingValue 里面 -->
    <div>@ThemeClassName</div>

    @code {
        [CascadingParameter]
        protected string ThemeClassName {get; set;}
    }
    ```

    当存在多个相同类型的值时会发生冲突，冲突时取内层的 CascadingVaule，相当于覆盖。可以通过指定 CascadingValue 的 Name 属性来给予不同的命名，然后在使用时指定名字

    ```html
    @code {
        [CascadingParameter(Name = "xxx")]
        protected string AnotherClassName {get; set;}
    }
    ```

## 组件的控制逻辑

组件可以直接在HTML中使用 `@foreach`, `@if` 等语法进行控制逻辑，基本和 C# 是一样的，例如

```html
<tbody>
    @foreach (var forecast in forecasts)
    {
        <tr>
            <td>@forecast.Date.ToShortDateString()</td>
            <td>@forecast.TemperatureC</td>
            <td>@forecast.TemperatureF</td>
            <td>@forecast.Summary</td>
        </tr>
    }
</tbody>
```

## 表单自动验证

blazor 内置了一些通用组件，例如 <EditForm> 组件支持传入一个 model，然后来自动验证输入是否合法：

```html
<EditForm Model="@personModel" OnValidSubmit="HandleSubmit">
    <DataAnnotationsValidator/> <!-- 告诉 blazor 需要根据 model 上的 annotation 来验证表单 -->
    <InputText id="firstName" @bind-Value="personModel.FirstName" />
    <button type="submit">Submit</button>
</EditForm>

@code {
    private PersonModel personModel = new PersonModel();

    private void HandleSubmit() {
        // take action here, the form is valid!
    }
}
```

Model:

```cs
public class PersonModel {
    [Required]
    public string FirstName {get; set;}
}
```

## 组件调用后端代码

在组件开头编写

```cs
@using MyApp.Data // using namespace，获取model类型
@inject WeatherForecastService service // 类似注入的方式获取service实例
```

之后在 `@code` 中调用 `service` 的方法

```cs
private WeatherForecast[] forecasts;

protected override async Task OnInitializedAsync()
{
    forecasts = await ForecastService.GetForecastAsync(DateTime.Now);
}
```

来获取数据。（service的方法返回的是 `Task<WeatherForecase[]>`）

另外注意到 `onInitializedAsync` 是一个类似 react 的生命周期钩子函数，会在页面初始化时候执行。

## 生命周期钩子函数 (lifecycle hook)

主要分两类：同步和异步。同步直接执行，异步返回的task等待后继续执行。

```cs
// 通用签名
protected override void Hooks() {

}

protected override async Task HooksAsync() {

}
```

![SetParameter, OnInitialized and OnParametersSet](https://docs.microsoft.com/en-us/aspnet/core/blazor/components/lifecycle/_static/lifecycle1.png?view=aspnetcore-5.0)

* SetParametersAsync(ParameterView parameters) - 父组件设置当前组件参数之后（**只在组件初次渲染时触发**，参数就在 ParameterView 里面）。可以在这里决定如何设置当前组件的参数，默认行为是将 ParameterView 中的值给 [Parameter] 中的同名变量。
* OnInitializedAsync() - 组件初始化完成后（**只在组件初次渲染时触发**。一般用来发请求）
    通常判断一下变量是否为 null 来决定是否展示 loading 界面。
    当 Blazor Server 使用 prerender 模式时，这个钩子会被调用**两次**，第一次是 prerender，第二次是和浏览器建立连接后的 HTML 混合。为了防止两次重绘，可以进行一些判断，或者缓存结果。
* OnParametersSetAsync() - 组件的参数变化时（不止在初次渲染时触发）。对于对象类型的参数，只要赋值就算变化，并不做比较。
* OnAfterRenderAsync(bool firstRender) - 组件渲染完成后（此时开始有 DOM 有 JS 了）
    这个生命周期函数比较特别，修改状态不会自动触发重绘（防止死循环），需要手动调用 `StateHasChanged` 通知 Blazor 重绘页面。注意避免死循环（最好套在 if 里面）

## C# 和 JS 互相调用 (JS interop)

### C# 调用 JS 方法

使用场景：使用一些 JS 类库，使用系统方法。

C# 可以调用注册在 window 上的全局方法：

```cs
@page "/"
@inject IJSRuntime JS; // 注入 JS 对象，调用方法要用

<button @onclick="OnClick">Click Me</button>

@code {
    private async Task OnClick() {
        // 调用时必须 await，也因此方法签名需要改成 async Task
        // 第一个参数是方法名，第二个参数往后是方法本身的参数。
        // 泛型是返回值的类型
        var something = await JS.InvokeAsync<string>("prompt", "Input something");
        // do something
        await JS.InvokeVoidAsync("alert", "return void or do not need return value");
    }
}

```

如果是要调用自己写/引入的方法，在wwwroot/index.html 或者 Pages/_Host.cshtml 中定义/引入，必须要挂在 window 上才行。不能在 razor 中定义/引入。

因为要挂在 window 上，所以尽量建立自己的 namespace，这样在调用时第一个参数就变成 `namespace.functionName`。

如果不想挂在 window 上，Blazor 还支持使用模块化的 JS 来定义，参见 https://docs.microsoft.com/en-us/aspnet/core/blazor/call-javascript-from-dotnet?view=aspnetcore-5.0#blazor-javascript-isolation-and-object-references

注意：
当Blazor 进行预渲染 (prerendering) 时因为是在后端进行的，此时与浏览器的连接尚未建立，因此是无法调用到 JS 方法的。(同理，任何 JS 相关的东西，比如 document 等，也都没有)
为了解决这个问题，可以使用生命周期函数 `OnAfterRenderAsync(bool firstRender)`，在这个周期里面的函数可以确保连接已经建立。（当然在 click 事件里写也能确保）
在钩子最后调用 `StateHasChanged` 可以通知 Blazor 获取最新状态并更新 DOM。这些代码最好在 `firstRender=true`时调用，否则会导致死循环。

### JS 调用 C# 方法

使用场景：

JS 可以调用在 razor component 中 `@codes` 段落中的方法。

**静态**方法被标记为 `[JSInvokable]` 时，可以通过 `DotNet.invokeMethodAsync('App Assembly Name', 'Function Name', arguments...).then(returnValue => {...})` 在 JS 中调用。

其他实例方法的调用参见 https://docs.microsoft.com/en-us/aspnet/core/blazor/call-dotnet-from-javascript?view=aspnetcore-5.0

## 开发

如果使用 vs 开发，则和 webpack 一样，开发状态时会有 websocket 连接服务器，因此修改后端代码后前端页面直接自动刷新，不需要重启服务器也不需要手动刷新页面。

对于 razor 中 @code 的断点也在后端，而不再前端，因此对于 C# 程序员来说，所有 DEBUG 都在 VS 中进行。 （Blazor WebAssembly App 是否也是这样尚不确定）
