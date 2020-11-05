class Observer {
  constructor() {
    this.subscribers = [];
  }

  subscribe(callback) {
    this.subscribers.push(callback);
  }

  unsubscribe(callback) {
    this.subscribers = this.subscribers.filter(subscriber => callback !== subscriber);
  }

  async notify() {
    await Promise.all(this.subscribers.map(subscriber => subscriber()));
  }
}

export default Observer;
