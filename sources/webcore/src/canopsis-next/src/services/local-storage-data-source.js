/**
 * Class representing a LocalStorageDataSource
 */
class LocalStorageDataSource {
  constructor() {
    try {
      Object.keys(localStorage).map(key => this[key] = localStorage.getItem(key));
    } catch (err) {
      console.warn(err);
    }
  }

  /**
   * @param {string} key - key localStorage
   * @return {string}
   */
  getItem(key) {
    return this[key];
  }

  /**
   * @param {string} key - key local Storage
   * @param {string} value - value local Storage
   * @return {LocalStorageDataSource}
   */
  setItem(key, value) {
    this[key] = value;

    try {
      localStorage.setItem(key, value);
    } catch (err) {
      console.warn(err);
    }

    return this;
  }

  /**
   * @param {string} key - key local Storage
   * @return {LocalStorageDataSource}
   */
  removeItem(key) {
    delete this[key];

    try {
      localStorage.removeItem(key);
    } catch (err) {
      console.warn(err);
    }

    return this;
  }
}

export default new LocalStorageDataSource();
