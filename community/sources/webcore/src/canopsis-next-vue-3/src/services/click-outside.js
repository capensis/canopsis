import Observer from './observer';
import { ref } from 'vue';

class ClickOutside extends Observer {
  constructor() {
    super();

    this.handlersCount = ref(0);
  }

  register(handler) {
    super.register(handler);

    this.handlersCount.value = this.handlers.length;
  }

  unregister(handler) {
    super.unregister(handler);

    this.handlersCount.value = this.handlers.length;
  }

  call(...args) {
    return this.handlers.every(handler => handler(...args));
  }
}

export default ClickOutside;
