class SocketRoom {
  constructor(name, authNeeded) {
    this.name = name;
    this.authNeeded = authNeeded;
    this.count = 1;
    this.listeners = [];
  }

  /**
   * Add listener to message listeners array
   *
   * @param {Function} listener
   * @return {SocketRoom}
   */
  addListener(listener) {
    this.listeners.push(listener);

    return this;
  }

  /**
   * Remove listener from message listeners array
   *
   * @param {Function} listener
   * @return {SocketRoom}
   */
  removeListener(listener) {
    this.listeners = this.listeners.filter(item => item !== listener);

    return this;
  }

  /**
   * Increment count of joins to room
   *
   * @return {SocketRoom}
   */
  increment() {
    this.count += 1;

    return this;
  }

  /**
   * Decrement count of joins to room
   *
   * @return {SocketRoom}
   */
  decrement() {
    this.count -= 1;

    return this;
  }

  /**
   * Check if no joins to room
   *
   * @returns {boolean}
   */
  isEmpty() {
    return this.count <= 0;
  }

  /**
   * Call all message listeners with context
   *
   * @param {any} context
   * @param {Array} restArgs
   * @return {SocketRoom}
   */
  call(context, ...restArgs) {
    this.listeners.forEach(listener => listener.call(context, ...restArgs));

    return this;
  }
}

export default SocketRoom;
