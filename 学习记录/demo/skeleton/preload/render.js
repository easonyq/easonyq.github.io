function getRandomNum() {
  return parseInt(Math.random() * 256, 10)
}

function getRandomColor() {
  let r = getRandomNum()
  let g = getRandomNum()
  let b = getRandomNum()

  return `rgb(${r}, ${g}, ${b})`
}

function mountApp() {
  let appDOM = document.querySelector('#app')
  appDOM.innerHTML = ''
  for (let i = 0; i < 10; i++) {
    let div = document.createElement('div')
    div.classList.add('bar')
    div.style['background-color'] = getRandomColor()
    appDOM.appendChild(div)
  }
}

window.mountApp = mountApp
if (window.STYLE_READY) {
  console.log('style ready')
  mountApp()
}
