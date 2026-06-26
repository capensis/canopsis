class Observer {
  constructor() {
    this.handlers = [];
    this.childrenHandlers = [];
  }

  register(handler) {
    this.handlers.push(handler);
  }

  registerChild(handler) {
    this.childrenHandlers.push(handler);
  }

  unregister(handler) {
    this.handlers = this.handlers.filter(h => h !== handler);
  }

  unregisterChild(handler) {
    this.childrenHandlers = this.childrenHandlers.filter(h => h !== handler);
  }

  unregisterAll() {
    this.handlers = [];
    this.childrenHandlers = [];
  }

  async notify(data) {
    await Promise.all(this.handlers.map(subscriber => subscriber(data)));
    await Promise.all(this.childrenHandlers.map(subscriber => subscriber(data)));
  }

  async notifyInSeries(data) {
    for (const handler of this.handlers) {
       
      await handler(data);
    }

    for (const handler of this.childrenHandlers) {
       
      await handler(data);
    }
  }
}

export default Observer;
