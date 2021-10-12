class Observer {
  constructor() {
    this.handlers = [];
  }

  register(handler) {
    this.handlers.push(handler);
  }

  unregister(handler) {
    this.handlers = this.handlers.filter(h => h !== handler);
  }

  async notify() {
    await Promise.all(this.handlers.map(subscriber => subscriber()));
  }
}

export default Observer;
