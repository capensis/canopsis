class SocketRoom {
  constructor(name) {
    this.name = name;
    this.count = 1;
    this.listeners = [];
  }

  addListener(listener) {
    this.listeners.push(listener);
  }

  removeListener(listener) {
    this.listeners = this.listeners.filter(item => item !== listener);
  }

  increment() {
    this.count += 1;
  }

  decrement() {
    this.count -= 1;
  }

  isEmpty() {
    return this.count <= 0;
  }

  call(context, ...restArgs) {
    this.listeners.forEach(listener => listener.call(context, ...restArgs));
  }
}

export default SocketRoom;
