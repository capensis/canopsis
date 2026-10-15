import { recursiveRaf } from '@/helpers/render';

/**
 * Class representing asynchronous booting functionality.
 */
export class AsyncBooting {
  /**
   * Create an AsyncBooting instance.
   */
  constructor() {
    this.wasCalled = false;
    this.registered = [];
    this.cancelBooting = null;
  }

  /**
   * Register boot callback.
   *
   * @param {any} key - The key to register.
   * @param {Function} bootCallback - The callback function to execute on boot.
   * @param {Function} onFinishCallback - The callback function to execute on finish phase.
   */
  register(key, bootCallback, onFinishCallback) {
    if (this.wasCalled) {
      bootCallback();

      return;
    }

    this.registered.push({ key, bootCallback, onFinishCallback });
  }

  /**
   * Unregister a key.
   *
   * @param {any} key - The key to unregister.
   */
  unregister(key) {
    this.registered = this.registered.filter(({ key: registeredKey }) => key !== registeredKey);
  }

  /**
   * Run the asynchronous booting process.
   *
   * @param {number} [duration] - Work budget for one frame, in milliseconds.
   * @param {Function} [onFinishCallbackRoot] - Callback which will calls after all bootCallbacks
   */
  run(duration, onFinishCallbackRoot) {
    this.cancelBooting?.();

    const cancelBooting = recursiveRaf(
      this.registered,
      async ({ key, bootCallback }) => {
        if (this.registered.some(({ key: registeredKey }) => key === registeredKey)) {
          await bootCallback();
        }
      },
      {
        duration,
        onComplete: () => {
          this.cancelBooting = null;
          this.registered.forEach(({ onFinishCallback }) => onFinishCallback?.());
          onFinishCallbackRoot?.();
          this.registered = [];
          this.wasCalled = true;
        },
      },
    );

    if (this.registered.length) {
      this.cancelBooting = cancelBooting;
    }
  }

  /**
   * Clear all registered components and rAFs.
   */
  clear() {
    this.cancelBooting?.();

    this.registered = [];
    this.cancelBooting = null;
  }
}
