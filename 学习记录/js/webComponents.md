# web component

[参考文章](https://developers.google.com/web/fundamentals/web-components/)

* 简单来说就是在 HTML 中使用自定义标签
* **组件提供了HTML、CSS、JavaScript封装的方法，实现了与同一页面上其他代码的隔离**
* 包含4种技术： Template, Custom Element, Shadow DOM, HTML Import。最重要的是 Custom Element 和 Shadow DOM。

## Custom Element

使用 `customElements.define()` 来定义一个新的标签。

示例 - 定义一个移动抽屉面板 `<app-drawer>`：

```javascript
class AppDrawer extends HTMLElement {

  // A getter/setter for an open property.
  get open() {
    return this.hasAttribute('open');
  }

  set open(val) {
    // Reflect the value of the open property as an HTML attribute.
    if (val) {
      this.setAttribute('open', '');
    } else {
      this.removeAttribute('open');
    }
    this.toggleDrawer();
  }

  // A getter/setter for a disabled property.
  get disabled() {
    return this.hasAttribute('disabled');
  }

  set disabled(val) {
    // Reflect the value of the disabled property as an HTML attribute.
    if (val) {
      this.setAttribute('disabled', '');
    } else {
      this.removeAttribute('disabled');
    }
  }

  // Can define constructor arguments if you wish.
  constructor() {
    // If you define a ctor, always call super() first!
    // This is specific to CE and required by the spec.
    super();

    // Setup a click listener on <app-drawer> itself.
    this.addEventListener('click', e => {
      // Don't toggle the drawer if it's disabled.
      if (this.disabled) {
        return;
      }
      this.toggleDrawer();
    });
  }

  toggleDrawer() {
    // ...
  }
}

customElements.define('app-drawer', AppDrawer);

```

上述代码定义了这个 `<app-drawer>` 标签，并且添加了两个属性，`open` 和 `disabled`。`click` 事件定义在构造函数中，`this` 指向标签本身，可以使用例如 `this.children` 或者 `this.querySelectorAll` 之类的方法。

标签的示例用法：

```html
<app-drawer open></app-drawer>
```

另外自定义元素也可以继承其他自定义元素，如 `class FancyDrawer extends AppDrawer` 这样。

### 扩展原生 HTML

Custom Elements API 也可以用于对已有的内置 HTML 进行扩展。但不同的标签的类并不是统一的 `HTMLElement`。例如 `<button>` 要继承 `HTMLButtonElement`，而 `<img>` 要继承 `HTMLImageElement`。[完整列表](https://html.spec.whatwg.org/multipage/indices.html#element-interfaces)。下面的例子扩展了原生的按钮，增加了按钮点击时的波纹效果。

```javascript
// See https://html.spec.whatwg.org/multipage/indices.html#element-interfaces
// for the list of other DOM interfaces.
class FancyButton extends HTMLButtonElement {
  constructor() {
    super(); // always call super() first in the ctor.
    this.addEventListener('click', e => this.drawRipple(e.offsetX, e.offsetY));
  }

  // Material design ripple animation.
  drawRipple(x, y) {
    let div = document.createElement('div');
    div.classList.add('ripple');
    this.appendChild(div);
    div.style.top = `${y - div.clientHeight/2}px`;
    div.style.left = `${x - div.clientWidth/2}px`;
    div.style.backgroundColor = 'currentColor';
    div.classList.add('run');
    div.addEventListener('transitionend', e => div.remove());
  }
}

customElements.define('fancy-button', FancyButton, {extends: 'button'});
```

使用方法有:

1. 添加 `is` 属性

    ```html
    <!-- This <button> is a fancy button. -->
    <button is="fancy-button" disabled>Fancy button!</button>
    ```

2. 使用 JS 添加

    ```javascript
    // Custom elements overload createElement() to support the is="" attribute.
    let button = document.createElement('button', {is: 'fancy-button'});
    button.textContent = 'Fancy button!';
    button.disabled = true;
    document.body.appendChild(button);

    // Or use new
    let button = new FancyButton();
    button.textContent = 'Fancy button!';
    button.disabled = true;
    ```

TODO 未完
