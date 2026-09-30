import { registerVitePreloadErrorHandler } from '@/services/vite-preload-error';

const createStorage = (initialValue) => {
  const values = new Map();

  if (initialValue !== undefined) {
    values.set('canopsis:preload-error-reload', String(initialValue));
  }

  return {
    getItem: vi.fn(key => values.get(key) ?? null),
    setItem: vi.fn((key, value) => values.set(key, value)),
  };
};

const dispatchPreloadError = (eventTarget) => {
  const event = new Event('vite:preloadError', { cancelable: true });

  eventTarget.dispatchEvent(event);

  return event;
};

describe('registerVitePreloadErrorHandler', () => {
  it('reloads the page when a dynamically imported chunk cannot be loaded', () => {
    const eventTarget = new EventTarget();
    const storage = createStorage();
    const reload = vi.fn();

    registerVitePreloadErrorHandler({ eventTarget, storage, reload, now: () => 15_000 });

    const event = dispatchPreloadError(eventTarget);

    expect(event.defaultPrevented).toBe(true);
    expect(storage.setItem).toHaveBeenCalledWith('canopsis:preload-error-reload', '15000');
    expect(reload).toHaveBeenCalledOnce();
  });

  it('reloads only once for multiple preload errors on the same page', () => {
    const eventTarget = new EventTarget();
    const reload = vi.fn();

    registerVitePreloadErrorHandler({
      eventTarget,
      storage: createStorage(),
      reload,
      now: () => 15_000,
    });

    dispatchPreloadError(eventTarget);
    dispatchPreloadError(eventTarget);

    expect(reload).toHaveBeenCalledOnce();
  });

  it('does not cause a reload loop when the refreshed page still cannot load a chunk', () => {
    const eventTarget = new EventTarget();
    const reload = vi.fn();

    registerVitePreloadErrorHandler({
      eventTarget,
      storage: createStorage(14_000),
      reload,
      now: () => 15_000,
    });

    const event = dispatchPreloadError(eventTarget);

    expect(event.defaultPrevented).toBe(false);
    expect(reload).not.toHaveBeenCalled();
  });

  it('allows another recovery reload after the cooldown', () => {
    const eventTarget = new EventTarget();
    const reload = vi.fn();

    registerVitePreloadErrorHandler({
      eventTarget,
      storage: createStorage(1_000),
      reload,
      now: () => 15_000,
    });

    dispatchPreloadError(eventTarget);

    expect(reload).toHaveBeenCalledOnce();
  });

  it('leaves the original error untouched when session storage is unavailable', () => {
    const eventTarget = new EventTarget();
    const reload = vi.fn();
    const storage = {
      getItem: vi.fn(() => {
        throw new Error('Storage unavailable');
      }),
      setItem: vi.fn(),
    };

    registerVitePreloadErrorHandler({ eventTarget, storage, reload });

    const event = dispatchPreloadError(eventTarget);

    expect(event.defaultPrevented).toBe(false);
    expect(reload).not.toHaveBeenCalled();
  });
});
