import Observer from './observer';

class ClickOutside extends Observer {
  call(...args) {
    return this.handlers.every(handler => handler(...args));
  }
}

export default ClickOutside;
