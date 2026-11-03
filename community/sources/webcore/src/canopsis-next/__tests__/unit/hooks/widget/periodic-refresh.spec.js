import { generateShallowRenderer } from '@unit/utils/vue';
import { computed } from 'vue';

import Observer from '@/services/observer';

import { useWidgetPeriodicRefresh } from '@/hooks/widget/periodic-refresh';

describe('useWidgetPeriodicRefresh', () => {
  const handler = vi.fn().mockResolvedValue();
  const factory = generateShallowRenderer({
    props: { widget: { type: Object, required: true } },
    setup(props) {
      useWidgetPeriodicRefresh({
        options: computed(() => props.widget.parameters.periodic_refresh),
        handler,
      });
    },
    template: '<div />',
  }, {
    global: { provide: { $periodicRefresh: new Observer() } },
  });
  const createWidget = (value = 10, unit = 's', enabled = true) => ({
    parameters: { periodic_refresh: { enabled, value, unit } },
  });

  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.clearAllMocks();
  });

  test.each([
    [20, 's', 20000],
    [10, 'm', 600000],
  ])('Should restart the timer when saved settings change to %i %s', async (value, unit, delay) => {
    const wrapper = factory({ props: { widget: createWidget() } });

    await vi.advanceTimersByTimeAsync(3000);
    await wrapper.setProps({ widget: createWidget(value, unit) });
    await vi.advanceTimersByTimeAsync(delay - 1);

    expect(handler).not.toHaveBeenCalled();

    await vi.advanceTimersByTimeAsync(1);

    expect(handler).toHaveBeenCalledTimes(1);
  });

  test('Should stop the timer when disabled or unmounted', async () => {
    const wrapper = factory({ props: { widget: createWidget() } });

    expect(vi.getTimerCount()).toBe(1);

    await wrapper.setProps({ widget: createWidget(10, 's', false) });
    await vi.advanceTimersByTimeAsync(10000);

    expect(handler).not.toHaveBeenCalled();
    expect(vi.getTimerCount()).toBe(0);

    await wrapper.setProps({ widget: createWidget() });

    expect(vi.getTimerCount()).toBe(1);

    wrapper.unmount();

    expect(vi.getTimerCount()).toBe(0);
  });

  test('Should not start a timer while a refresh is still pending', async () => {
    let finishRefresh;

    handler.mockImplementationOnce(() => new Promise((resolve) => {
      finishRefresh = resolve;
    }));

    const wrapper = factory({ props: { widget: createWidget() } });

    await vi.advanceTimersByTimeAsync(10000);

    expect(handler).toHaveBeenCalledTimes(1);
    expect(vi.getTimerCount()).toBe(0);

    await wrapper.setProps({ widget: createWidget(20) });

    expect(vi.getTimerCount()).toBe(0);

    finishRefresh();
    await vi.advanceTimersByTimeAsync(0);
    await vi.advanceTimersByTimeAsync(19999);

    expect(handler).toHaveBeenCalledTimes(1);

    await vi.advanceTimersByTimeAsync(1);

    expect(handler).toHaveBeenCalledTimes(2);
  });
});
