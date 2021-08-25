import { isNil } from 'lodash';

class LocalStorage {
  constructor() {
    this.storage = {};

    try {
      Object.entries(window.localStorage).forEach(([key, value]) => this.storage[key] = value);
    } catch (err) {
      console.warn(err);
    }
  }

  get(key, json = false) {
    const value = this.storage[key];

    return json ? JSON.parse(value) : this.storage[key];
  }

  has(key) {
    return !isNil(this.storage[key]);
  }

  set(key, value) {
    this.storage[key] = value;

    try {
      window.localStorage.setItem(key, value);
    } catch (err) {
      console.warn(err);
    }
  }

  remove(key) {
    delete this.storage[key];

    try {
      window.localStorage.removeItem(key);
    } catch (err) {
      console.warn(err);
    }
  }

  clear() {
    this.storage = {};

    try {
      localStorage.clear();
    } catch (err) {
      console.warn(err);
    }
  }
}

export default new LocalStorage();
