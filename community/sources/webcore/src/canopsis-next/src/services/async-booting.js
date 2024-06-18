/**
 * Class representing asynchronous booting functionality.
 */
export class AsyncBooting {
  /**
   * Create an AsyncBooting instance.
   */
  constructor() {
    this.wasCalled = false;
    this.registeredRAFs = new Map();
    this.registered = [];
  }

  /**
   * Register boot callback.
   *
   * @param {any} key - The key to register.
   * @param {Function} bootCallback - The callback function to execute on boot.
   * @param {Function} onFinishCallback - The callback function to execute on finish phase.
   */
  register(key, bootCallback, onFinishCallback) {
    this.registered.push({ key, bootCallback, onFinishCallback });
  }

  /**
   * Unregister a key.
   *
   * @param {any} key - The key to unregister.
   */
  unregister(key) {
    this.registered = this.registered
      .filter(({ key: registeredKey }) => key !== registeredKey);

    if (this.registeredRAFs.has(key)) {
      window.cancelAnimationFrame(this.registeredRAFs.get(key));
    }
  }

  /**
   * Perform a recursive requestAnimationFrame.
   *
   * @param {Function} bootCallback - The callback function to execute on each frame.
   * @param {number} depth - The depth of recursion.
   * @param {any} [key] - The key to perform the recursive rAF on.
   */
  recursiveRAF(bootCallback, depth, key) {
    if (depth <= 0) {
      bootCallback();

      if (key) {
        this.registeredRAFs.delete(key);
      }

      return;
    }

    const rAFId = window.requestAnimationFrame(() => this.recursiveRAF(bootCallback, depth - 1, key));

    if (key) {
      this.registeredRAFs.set(key, rAFId);
    }
  }

  /**
   * Run the asynchronous booting process.
   *
   * @param {number} [itemsPerRender = 10] - The number of items per one requestAnimationFrame
   * @param {Function} [onFinishCallbackRoot] - Callback which will calls after all bootCallbacks
   */
  run(itemsPerRender = 10, onFinishCallbackRoot) {
    this.registered.forEach(({ key, bootCallback }, index) => {
      this.recursiveRAF(bootCallback, Math.floor(index / itemsPerRender), key);
    });

    const afterRAFs = Math.floor(this.registered.length / itemsPerRender) + 1;

    this.recursiveRAF(() => {
      this.registered.forEach(({ onFinishCallback }) => onFinishCallback?.());
      onFinishCallbackRoot?.();
      this.clear();
      this.wasCalled = true;
    }, afterRAFs);
  }

  /**
   * Clear all registered components and rAFs.
   */
  clear() {
    this.registeredRAFs.forEach(rAFId => window.cancelAnimationFrame(rAFId));

    this.registered = [];
    this.registeredRAFs = new Map();
  }
}
