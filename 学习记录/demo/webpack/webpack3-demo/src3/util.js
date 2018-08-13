console.log('This is Util class')

export default class Util {
  hello () {
    return 'hello'
  }

  bye () {
    return 'bye'
  }
}

Array.prototype.hello = () => {'hello'}
