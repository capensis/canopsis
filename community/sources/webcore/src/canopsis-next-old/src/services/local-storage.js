class LocalStorage {
  constructor() {
    this.storage = {};

    try {
      Object.entries(window.localStorage).forEach(([key, value]) => this.storage[key] = value);
    } catch (err) {
      console.warn(err);
    }
  }

  /**
   * Get data from storage
   *
   * @param {string | number} key
   * @returns {any}
   */
  get(key) {
    return this.storage[key];
  }

  /**
   * Check if storage has the data by key
   *
   * @param {string | number} key
   * @returns {boolean}
   */
  has(key) {
    return key in this.storage;
  }

  /**
   * Set value in storage by key
   *
   * @param {string | number} key
   * @param {any} value
   */
  set(key, value) {
    this.storage[key] = value;

    try {
      window.localStorage.setItem(key, value);
    } catch (err) {
      console.warn(err);
    }
  }

  /**
   * Remove data from storage by key
   *
   * @param {string | number} key
   */
  remove(key) {
    delete this.storage[key];

    try {
      window.localStorage.removeItem(key);
    } catch (err) {
      console.warn(err);
    }
  }

  /**
   * Get value by key and remove it from storage
   *
   * @param {string | number} key
   * @returns {*}
   */
  pop(key) {
    const value = this.get(key);

    this.remove(key);

    return value;
  }

  /**
   * Clear storage
   */
  clear() {
    this.storage = {};

    try {
      window.localStorage.clear();
    } catch (err) {
      console.warn(err);
    }
  }
}

export default new LocalStorage();
