const PRELOAD_ERROR_RELOAD_KEY = 'canopsis:preload-error-reload';
const PRELOAD_ERROR_RELOAD_COOLDOWN = 10_000;

/**
 * Reload the application once when a lazy-loaded chunk from an older deployment is unavailable.
 *
 * @param {Object} [options] - Handler dependencies.
 * @param {EventTarget} [options.eventTarget] - Target emitting Vite preload errors.
 * @param {Storage} [options.storage] - Session storage used to prevent reload loops.
 * @param {Function} [options.reload] - Page reload callback.
 * @param {Function} [options.now] - Current timestamp provider.
 * @returns {Function} Function unregistering the handler.
 */
export const registerVitePreloadErrorHandler = ({
  eventTarget = window,
  storage = window.sessionStorage,
  reload = () => window.location.reload(),
  now = Date.now,
} = {}) => {
  let reloadRequested = false;

  const handler = (event) => {
    if (reloadRequested) {
      return;
    }

    const currentTimestamp = now();

    try {
      const previousTimestamp = Number(storage.getItem(PRELOAD_ERROR_RELOAD_KEY));

      if (
        previousTimestamp &&
        currentTimestamp - previousTimestamp < PRELOAD_ERROR_RELOAD_COOLDOWN
      ) {
        return;
      }

      storage.setItem(PRELOAD_ERROR_RELOAD_KEY, String(currentTimestamp));
    } catch {
      return;
    }

    reloadRequested = true;
    event.preventDefault();
    reload();
  };

  eventTarget.addEventListener('vite:preloadError', handler);

  return () => eventTarget.removeEventListener('vite:preloadError', handler);
};
