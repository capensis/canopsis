import {
  callAfterAnimationFrames,
  getNearestAndFarthestIndexes,
  recursiveRaf,
} from '@/helpers/render';

describe('render helpers', () => {
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

  it('processes queued items within a duration measured in milliseconds', async () => {
    const callback = vi.fn();
    const onComplete = vi.fn();

    vi.spyOn(performance, 'now')
      .mockReturnValueOnce(0)
      .mockReturnValueOnce(4)
      .mockReturnValueOnce(9)
      .mockReturnValueOnce(16);

    recursiveRaf([1, 2, 3], callback, { duration: 8, onComplete });

    expect(animationFrames).toHaveLength(1);

    await animationFrames.shift()();

    expect(callback).toHaveBeenNthCalledWith(1, 1, 0, [1, 2, 3]);
    expect(callback).toHaveBeenNthCalledWith(2, 2, 1, [1, 2, 3]);
    expect(onComplete).not.toHaveBeenCalled();
    expect(animationFrames).toHaveLength(1);

    await animationFrames.shift()();

    expect(callback).toHaveBeenNthCalledWith(3, 3, 2, [1, 2, 3]);
    expect(onComplete).toHaveBeenCalledOnce();
  });

  it('cancels scheduled queue processing', () => {
    const callback = vi.fn();
    const cancel = recursiveRaf([1], callback);

    cancel();

    expect(window.cancelAnimationFrame).toHaveBeenCalledWith(1);
    expect(callback).not.toHaveBeenCalled();
  });

  it('calls a callback after the requested number of animation frames', () => {
    const callback = vi.fn();

    callAfterAnimationFrames(callback, 2);

    animationFrames.shift()();
    expect(callback).not.toHaveBeenCalled();

    animationFrames.shift()();
    expect(callback).toHaveBeenCalledOnce();
  });

  it('does not add an undefined farthest item', () => {
    expect(
      getNearestAndFarthestIndexes({
        start: 1,
        end: 2,
        ids: ['first', 'second', 'third', 'fourth'],
      }),
    ).toEqual({
      nearest: ['second', 'third'],
      farthest: ['first', 'fourth'],
    });
  });
});
