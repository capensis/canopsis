import { generateShallowRenderer } from '@unit/utils/vue';
import { cloneDeep } from 'lodash';

import { TIME_UNITS } from '@/constants';

import { toSeconds } from '@/helpers/date/duration';

import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

describe('widget-periodic-refresh', () => {
  const mockWidget = {
    parameters: {
      periodic_refresh: {
        enabled: true,
        value: 10,
        unit: TIME_UNITS.second,
      },
    },
  };

  let mockContext;

  beforeEach(() => {
    vi.useFakeTimers();

    mockContext = {
      widget: cloneDeep(mockWidget),
      activeViewPeriodicRefreshPaused: false,
      isPeriodicRefreshVisible: true,
      periodicRefreshEnabled: true,
      periodicRefreshInterval: null,
      intervalStartedAt: null,
      intervalDelay: null,
      periodicRefreshSeconds: 10,
      fetchList: vi.fn(),
      resetPeriodicRefresh: widgetPeriodicRefreshMixin.methods.resetPeriodicRefresh,
    };
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.clearAllMocks();
  });

  describe('computed', () => {
    test('Should return periodicRefreshEnabled correctly', () => {
      const { periodicRefreshEnabled } = widgetPeriodicRefreshMixin.computed;

      expect(periodicRefreshEnabled.call(mockContext)).toBe(true);

      mockContext.widget.parameters.periodic_refresh.enabled = false;
      expect(periodicRefreshEnabled.call(mockContext)).toBe(false);

      mockContext.widget.parameters.periodic_refresh = undefined;
      expect(periodicRefreshEnabled.call(mockContext)).toBeUndefined();
    });

    test('Should return periodicRefreshSeconds correctly', () => {
      const { periodicRefreshSeconds } = widgetPeriodicRefreshMixin.computed;

      expect(periodicRefreshSeconds.call(mockContext)).toBe(10);

      mockContext.widget.parameters.periodic_refresh = {
        value: 5,
        unit: TIME_UNITS.minute,
      };
      expect(periodicRefreshSeconds.call(mockContext)).toBe(toSeconds(5, TIME_UNITS.minute));

      mockContext.widget.parameters.periodic_refresh = undefined;
      expect(periodicRefreshSeconds.call(mockContext)).toBe(toSeconds(undefined, undefined));
    });
  });

  describe('methods', () => {
    test('Should start periodic refresh correctly', () => {
      const context = {
        ...mockContext,
        periodicRefreshSeconds: 10,
        startPeriodicRefresh: widgetPeriodicRefreshMixin.methods.startPeriodicRefresh,
        stopPeriodicRefresh: widgetPeriodicRefreshMixin.methods.stopPeriodicRefresh,
      };

      context.startPeriodicRefresh();

      expect(context.periodicRefreshInterval).toBeDefined();
      expect(context.intervalStartedAt).toBeDefined();
      expect(Date.now() - context.intervalStartedAt).toBeLessThan(10);
    });

    test('Should stop existing interval before starting new one', () => {
      const context = {
        ...mockContext,
        periodicRefreshSeconds: 10,
        periodicRefreshInterval: setInterval(() => {}, 1000),
        startPeriodicRefresh: widgetPeriodicRefreshMixin.methods.startPeriodicRefresh,
        stopPeriodicRefresh: widgetPeriodicRefreshMixin.methods.stopPeriodicRefresh,
      };

      const oldInterval = context.periodicRefreshInterval;
      vi.spyOn(global, 'clearInterval');

      context.startPeriodicRefresh();

      expect(clearInterval).toHaveBeenCalledWith(oldInterval);
      expect(context.periodicRefreshInterval).not.toBe(oldInterval);
    });

    test('Should call fetchList when interval triggers', () => {
      const context = {
        ...mockContext,
        periodicRefreshSeconds: 1,
        startPeriodicRefresh: widgetPeriodicRefreshMixin.methods.startPeriodicRefresh,
        stopPeriodicRefresh: widgetPeriodicRefreshMixin.methods.stopPeriodicRefresh,
      };

      context.startPeriodicRefresh();

      expect(context.fetchList).not.toHaveBeenCalled();

      vi.advanceTimersByTime(1000);

      expect(context.fetchList).toHaveBeenCalledTimes(1);

      vi.advanceTimersByTime(1000);

      expect(context.fetchList).toHaveBeenCalledTimes(2);
    });

    test('Should pause periodic refresh correctly', () => {
      const context = {
        ...mockContext,
        periodicRefreshSeconds: 10,
        intervalStartedAt: Date.now() - 3000,
        periodicRefreshInterval: setInterval(() => {}, 10000),
        pausePriodicRefresh: widgetPeriodicRefreshMixin.methods.pausePriodicRefresh,
        stopPeriodicRefresh: widgetPeriodicRefreshMixin.methods.stopPeriodicRefresh,
      };

      vi.spyOn(global, 'clearInterval');

      context.pausePriodicRefresh();

      expect(clearInterval).toHaveBeenCalledWith(context.periodicRefreshInterval);
      expect(context.intervalDelay).toBe(7000); // 10000ms - 3000ms
    });

    test('Should resume periodic refresh correctly', () => {
      const context = {
        ...mockContext,
        intervalDelay: 2000,
        resumePriodicRefresh: widgetPeriodicRefreshMixin.methods.resumePriodicRefresh,
        startPeriodicRefresh: vi.fn(),
      };

      vi.spyOn(global, 'setTimeout');

      context.resumePriodicRefresh();

      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), 2000);

      vi.advanceTimersByTime(2000);

      expect(context.fetchList).toHaveBeenCalled();
      expect(context.startPeriodicRefresh).toHaveBeenCalled();
    });

    test('Should resume periodic refresh with zero delay when no delay set', () => {
      const context = {
        ...mockContext,
        intervalDelay: null,
        resumePriodicRefresh: widgetPeriodicRefreshMixin.methods.resumePriodicRefresh,
        startPeriodicRefresh: vi.fn(),
      };

      vi.spyOn(global, 'setTimeout');

      context.resumePriodicRefresh();

      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), 0);
    });

    test('Should stop periodic refresh correctly', () => {
      const testInterval = setInterval(() => {}, 1000);
      const context = {
        ...mockContext,
        periodicRefreshInterval: testInterval,
        intervalStartedAt: Date.now(),
        intervalDelay: 5000,
        stopPeriodicRefresh: widgetPeriodicRefreshMixin.methods.stopPeriodicRefresh,
      };

      vi.spyOn(global, 'clearInterval');

      context.stopPeriodicRefresh();

      expect(clearInterval).toHaveBeenCalledWith(testInterval);
      expect(context.periodicRefreshInterval).toBeNull();
      expect(context.intervalStartedAt).toBeNull();
      expect(context.intervalDelay).toBeNull();
    });
  });

  describe('watchers', () => {
    test('Should start periodic refresh when enabled and has value', () => {
      const watcher =
        widgetPeriodicRefreshMixin.watch.periodicRefreshSettings.handler;
      const context = {
        ...mockContext,
        periodicRefreshSeconds: 10,
        startPeriodicRefresh: vi.fn(),
        stopPeriodicRefresh: vi.fn(),
      };

      const newValue = {
        enabled: true,
        value: 10,
        unit: TIME_UNITS.second,
      };

      watcher.call(context, newValue, {});

      expect(context.startPeriodicRefresh).toHaveBeenCalled();
    });

    test('Should reset the delay without starting periodic refresh when paused', () => {
      const watcher =
        widgetPeriodicRefreshMixin.watch.periodicRefreshSettings.handler;
      const context = {
        ...mockContext,
        activeViewPeriodicRefreshPaused: true,
        startPeriodicRefresh: vi.fn(),
        stopPeriodicRefresh: vi.fn(),
      };

      const newValue = {
        enabled: true,
        value: 10,
        unit: TIME_UNITS.second,
      };

      watcher.call(context, newValue, {});

      expect(context.startPeriodicRefresh).not.toHaveBeenCalled();
      expect(context.stopPeriodicRefresh).not.toHaveBeenCalled();
      expect(context.intervalDelay).toBe(10000);
    });

    test('Should stop periodic refresh when disabled', () => {
      const watcher =
        widgetPeriodicRefreshMixin.watch.periodicRefreshSettings.handler;
      const context = {
        ...mockContext,
        stopPeriodicRefresh: vi.fn(),
      };

      const newValue = {
        enabled: false,
        value: 10,
        unit: TIME_UNITS.second,
      };

      watcher.call(context, newValue, {});

      expect(context.stopPeriodicRefresh).toHaveBeenCalled();
    });

    test('Should restart periodic refresh when value changes', () => {
      const watcher =
        widgetPeriodicRefreshMixin.watch.periodicRefreshSettings.handler;
      const context = {
        ...mockContext,
        periodicRefreshInterval: setInterval(() => {}, 1000),
        periodicRefreshSeconds: 15,
        startPeriodicRefresh: vi.fn(),
        stopPeriodicRefresh: vi.fn(),
      };

      const newValue = {
        enabled: true,
        value: 15,
        unit: TIME_UNITS.second,
      };

      const oldValue = {
        enabled: true,
        value: 10,
        unit: TIME_UNITS.second,
      };

      watcher.call(context, newValue, oldValue);

      expect(context.startPeriodicRefresh).toHaveBeenCalled();
      expect(context.intervalDelay).toBe(15000);
    });

    test('Should not restart periodic refresh when periodicRefreshSeconds is zero', () => {
      const watcher =
        widgetPeriodicRefreshMixin.watch.periodicRefreshSettings.handler;
      const context = {
        ...mockContext,
        periodicRefreshSeconds: 0,
        startPeriodicRefresh: vi.fn(),
        stopPeriodicRefresh: vi.fn(),
      };

      const newValue = {
        enabled: true,
        value: 0,
        unit: TIME_UNITS.second,
      };

      watcher.call(context, newValue, {});

      expect(context.startPeriodicRefresh).not.toHaveBeenCalled();
    });

    test('Should pause periodic refresh when activeViewPeriodicRefreshPaused becomes true', () => {
      const watcher = widgetPeriodicRefreshMixin.watch.activeViewPeriodicRefreshPaused;
      const context = {
        ...mockContext,
        pausePriodicRefresh: vi.fn(),
        resumePriodicRefresh: vi.fn(),
      };

      watcher.call(context, true);

      expect(context.pausePriodicRefresh).toHaveBeenCalled();
      expect(context.resumePriodicRefresh).not.toHaveBeenCalled();
    });

    test('Should resume periodic refresh when activeViewPeriodicRefreshPaused becomes false', () => {
      const watcher = widgetPeriodicRefreshMixin.watch.activeViewPeriodicRefreshPaused;
      const context = {
        ...mockContext,
        pausePriodicRefresh: vi.fn(),
        resumePriodicRefresh: vi.fn(),
      };

      watcher.call(context, false);

      expect(context.resumePriodicRefresh).toHaveBeenCalled();
      expect(context.pausePriodicRefresh).not.toHaveBeenCalled();
    });
  });

  describe('settings changes', () => {
    const fetchList = vi.fn();
    const factory = generateShallowRenderer({
      mixins: [widgetPeriodicRefreshMixin],
      props: {
        widget: { type: Object, required: true },
        paused: { type: Boolean, default: false },
      },
      computed: {
        activeViewPeriodicRefreshPaused() {
          return this.paused;
        },
      },
      methods: { fetchList },
      template: '<div />',
    });

    test.each([
      [{ value: 20 }, 20000],
      [{ unit: TIME_UNITS.minute }, 600000],
    ])('Should restart the timer after changing settings in place: %j', async (settings, delay) => {
      const wrapper = factory({ props: { widget: cloneDeep(mockWidget) } });

      vi.advanceTimersByTime(3000);
      Object.assign(wrapper.vm.widget.parameters.periodic_refresh, settings);
      await wrapper.vm.$nextTick();

      vi.advanceTimersByTime(delay - 1);

      expect(fetchList).not.toHaveBeenCalled();

      vi.advanceTimersByTime(1);

      expect(fetchList).toHaveBeenCalledTimes(1);
    });

    test('Should restart the timer when saved settings replace the widget', async () => {
      const wrapper = factory({ props: { widget: cloneDeep(mockWidget) } });

      vi.advanceTimersByTime(3000);
      await wrapper.setProps({
        widget: {
          parameters: {
            periodic_refresh: { enabled: true, value: 20, unit: TIME_UNITS.second },
          },
        },
      });
      vi.advanceTimersByTime(19999);

      expect(fetchList).not.toHaveBeenCalled();

      vi.advanceTimersByTime(1);

      expect(fetchList).toHaveBeenCalledTimes(1);
    });

    test('Should reset a paused timer and wait for an explicit resume', async () => {
      const wrapper = factory({ props: { widget: cloneDeep(mockWidget) } });

      vi.advanceTimersByTime(3000);
      await wrapper.setProps({ paused: true });

      const interval = wrapper.vm.periodicRefreshInterval;

      wrapper.vm.widget.parameters.periodic_refresh.value = 20;
      await wrapper.vm.$nextTick();

      expect(wrapper.vm.intervalDelay).toBe(20000);
      expect(wrapper.vm.periodicRefreshInterval).toBe(interval);
      expect(vi.getTimerCount()).toBe(0);

      vi.advanceTimersByTime(30000);

      expect(fetchList).not.toHaveBeenCalled();

      await wrapper.setProps({ paused: false });
      vi.advanceTimersByTime(19999);

      expect(fetchList).not.toHaveBeenCalled();

      vi.advanceTimersByTime(1);

      expect(fetchList).toHaveBeenCalledTimes(1);
    });

    test('Should reset the delay without starting a stopped timer', async () => {
      const wrapper = factory({ props: { widget: cloneDeep(mockWidget) } });

      wrapper.vm.stopPeriodicRefresh();
      wrapper.vm.widget.parameters.periodic_refresh.value = 20;
      await wrapper.vm.$nextTick();
      vi.advanceTimersByTime(30000);

      expect(wrapper.vm.intervalDelay).toBe(20000);
      expect(wrapper.vm.periodicRefreshInterval).toBeNull();
      expect(fetchList).not.toHaveBeenCalled();
      expect(vi.getTimerCount()).toBe(0);
    });
  });

  describe('beforeUnmount', () => {
    test('Should stop periodic refresh on component destroy', () => {
      const context = {
        ...mockContext,
        stopPeriodicRefresh: vi.fn(),
      };

      widgetPeriodicRefreshMixin.beforeUnmount.call(context);

      expect(context.stopPeriodicRefresh).toHaveBeenCalled();
    });
  });
});
