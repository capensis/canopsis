import { AsyncBooting } from '@/services/async-booting';

describe('AsyncBooting', () => {
  let animationFrames;

  beforeEach(() => {
    animationFrames = [];

    vi.spyOn(window, 'requestAnimationFrame').mockImplementation((callback) => {
      animationFrames.push(callback);

      return animationFrames.length;
    });
    vi.spyOn(window, 'cancelAnimationFrame').mockImplementation(() => {});
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('boots registered callbacks within a duration measured in milliseconds', async () => {
    const asyncBooting = new AsyncBooting();
    const bootCallbacks = [vi.fn(), vi.fn(), vi.fn()];
    const finishCallbacks = [vi.fn(), vi.fn(), vi.fn()];
    const rootFinishCallback = vi.fn();

    vi.spyOn(performance, 'now')
      .mockReturnValueOnce(0)
      .mockReturnValueOnce(4)
      .mockReturnValueOnce(9)
      .mockReturnValueOnce(16);

    bootCallbacks.forEach((bootCallback, index) => {
      asyncBooting.register(index, bootCallback, finishCallbacks[index]);
    });
    asyncBooting.run(8, rootFinishCallback);

    await animationFrames.shift()();

    expect(bootCallbacks[0]).toHaveBeenCalledOnce();
    expect(bootCallbacks[1]).toHaveBeenCalledOnce();
    expect(bootCallbacks[2]).not.toHaveBeenCalled();
    expect(rootFinishCallback).not.toHaveBeenCalled();
    expect(asyncBooting.wasCalled).toBe(false);

    await animationFrames.shift()();

    expect(bootCallbacks[2]).toHaveBeenCalledOnce();
    finishCallbacks.forEach(callback => expect(callback).toHaveBeenCalledOnce());
    expect(rootFinishCallback).toHaveBeenCalledOnce();
    expect(asyncBooting.wasCalled).toBe(true);
  });

  it('skips a callback unregistered while booting', async () => {
    const asyncBooting = new AsyncBooting();
    const firstBootCallback = vi.fn();
    const removedBootCallback = vi.fn();
    const firstFinishCallback = vi.fn();
    const removedFinishCallback = vi.fn();

    asyncBooting.register('first', firstBootCallback, firstFinishCallback);
    asyncBooting.register('removed', removedBootCallback, removedFinishCallback);
    asyncBooting.run(0);
    asyncBooting.unregister('removed');

    await animationFrames.shift()();
    await animationFrames.shift()();

    expect(firstBootCallback).toHaveBeenCalledOnce();
    expect(removedBootCallback).not.toHaveBeenCalled();
    expect(firstFinishCallback).toHaveBeenCalledOnce();
    expect(removedFinishCallback).not.toHaveBeenCalled();
  });

  it('cancels scheduled booting when cleared', () => {
    const asyncBooting = new AsyncBooting();
    const bootCallback = vi.fn();

    asyncBooting.register('item', bootCallback);
    asyncBooting.run();
    asyncBooting.clear();

    expect(window.cancelAnimationFrame).toHaveBeenCalledWith(1);
    expect(bootCallback).not.toHaveBeenCalled();
    expect(asyncBooting.registered).toEqual([]);
  });
});
