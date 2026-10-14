/**
 * Mock for features service used in unit tests.
 * This prevents loading actual feature files from @/features folder.
 */
class Features {
  constructor() {
    this.cache = {};
    this.features = [];
  }

  get(key, defaultValue) {
    return this.cache[key] ?? defaultValue;
  }

  has(key) {
    return key in this.cache;
  }

  call(key, context, ...args) {
    const func = this.get(key);

    if (typeof func !== 'function') {
      throw new Error(`Feature in the path = ${key} is not function: '${func}'`);
    }

    return func.call(context, ...args);
  }
}

export const featuresService = new Features();

export default featuresService;
