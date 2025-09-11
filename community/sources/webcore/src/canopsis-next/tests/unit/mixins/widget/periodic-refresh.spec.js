import { cloneDeep } from 'lodash';

import { DATETIME_FORMATS } from '@/constants';

import { toSeconds } from '@/helpers/date/duration';

import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

describe('widget-periodic-refresh', () => {
  const mockWidget = {
    parameters: {
      periodic_refresh: {
        enabled: true,
        value: 10,
        unit: DATETIME_FORMATS.second,
      },
    },
  };

  let mockContext;

  beforeEach(() => {
    jest.useFakeTimers();

    mockContext = {
      widget: cloneDeep(mockWidget),
      activeViewPeriodicRefreshPaused: false,
      periodicRefreshInterval: null,
      intervalStartedAt: null,
      intervalDelay: null,
      fetchList: jest.fn(),
    };
  });

  afterEach(() => {
    jest.useRealTimers();
    jest.clearAllMocks();
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
        unit: DATETIME_FORMATS.minute,
      };
      expect(periodicRefreshSeconds.call(mockContext)).toBe(toSeconds(5, DATETIME_FORMATS.minute));

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
      jest.spyOn(global, 'clearInterval');

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

      jest.advanceTimersByTime(1000);

      expect(context.fetchList).toHaveBeenCalledTimes(1);

      jest.advanceTimersByTime(1000);

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

      jest.spyOn(global, 'clearInterval');

      context.pausePriodicRefresh();

      expect(clearInterval).toHaveBeenCalledWith(context.periodicRefreshInterval);
      expect(context.intervalDelay).toBe(7000); // 10000ms - 3000ms
    });

    test('Should resume periodic refresh correctly', () => {
      const context = {
        ...mockContext,
        intervalDelay: 2000,
        resumePriodicRefresh: widgetPeriodicRefreshMixin.methods.resumePriodicRefresh,
        startPeriodicRefresh: jest.fn(),
      };

      jest.spyOn(global, 'setTimeout');

      context.resumePriodicRefresh();

      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), 2000);

      jest.advanceTimersByTime(2000);

      expect(context.fetchList).toHaveBeenCalled();
      expect(context.startPeriodicRefresh).toHaveBeenCalled();
    });

    test('Should resume periodic refresh with zero delay when no delay set', () => {
      const context = {
        ...mockContext,
        intervalDelay: null,
        resumePriodicRefresh: widgetPeriodicRefreshMixin.methods.resumePriodicRefresh,
        startPeriodicRefresh: jest.fn(),
      };

      jest.spyOn(global, 'setTimeout');

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

      jest.spyOn(global, 'clearInterval');

      context.stopPeriodicRefresh();

      expect(clearInterval).toHaveBeenCalledWith(testInterval);
      expect(context.periodicRefreshInterval).toBeNull();
      expect(context.intervalStartedAt).toBeNull();
      expect(context.intervalDelay).toBeNull();
    });
  });

  describe('watchers', () => {
    test('Should start periodic refresh when enabled and has value', () => {
      const watcher = widgetPeriodicRefreshMixin.watch['widget.parameters.periodic_refresh'].handler;
      const context = {
        ...mockContext,
        periodicRefreshSeconds: 10,
        startPeriodicRefresh: jest.fn(),
        stopPeriodicRefresh: jest.fn(),
      };

      const newValue = {
        enabled: true,
        value: 10,
        unit: DATETIME_FORMATS.second,
      };

      watcher.call(context, newValue, {});

      expect(context.startPeriodicRefresh).toHaveBeenCalled();
    });

    test('Should not start periodic refresh when paused', () => {
      const watcher = widgetPeriodicRefreshMixin.watch['widget.parameters.periodic_refresh'].handler;
      const context = {
        ...mockContext,
        activeViewPeriodicRefreshPaused: true,
        startPeriodicRefresh: jest.fn(),
        stopPeriodicRefresh: jest.fn(),
      };

      const newValue = {
        enabled: true,
        value: 10,
        unit: DATETIME_FORMATS.second,
      };

      watcher.call(context, newValue, {});

      expect(context.startPeriodicRefresh).not.toHaveBeenCalled();
      expect(context.stopPeriodicRefresh).not.toHaveBeenCalled();
    });

    test('Should stop periodic refresh when disabled', () => {
      const watcher = widgetPeriodicRefreshMixin.watch['widget.parameters.periodic_refresh'].handler;
      const context = {
        ...mockContext,
        stopPeriodicRefresh: jest.fn(),
      };

      const newValue = {
        enabled: false,
        value: 10,
        unit: DATETIME_FORMATS.second,
      };

      watcher.call(context, newValue, {});

      expect(context.stopPeriodicRefresh).toHaveBeenCalled();
    });

    test('Should restart periodic refresh when value changes', () => {
      const watcher = widgetPeriodicRefreshMixin.watch['widget.parameters.periodic_refresh'].handler;
      const context = {
        ...mockContext,
        periodicRefreshInterval: setInterval(() => {}, 1000),
        periodicRefreshSeconds: 15,
        startPeriodicRefresh: jest.fn(),
        stopPeriodicRefresh: jest.fn(),
      };

      const newValue = {
        enabled: true,
        value: 15,
        unit: DATETIME_FORMATS.second,
      };

      const oldValue = {
        enabled: true,
        value: 10,
        unit: DATETIME_FORMATS.second,
      };

      watcher.call(context, newValue, oldValue);

      expect(context.stopPeriodicRefresh).toHaveBeenCalled();
      expect(context.startPeriodicRefresh).toHaveBeenCalled();
    });

    test('Should not restart periodic refresh when periodicRefreshSeconds is zero', () => {
      const watcher = widgetPeriodicRefreshMixin.watch['widget.parameters.periodic_refresh'].handler;
      const context = {
        ...mockContext,
        periodicRefreshSeconds: 0,
        startPeriodicRefresh: jest.fn(),
        stopPeriodicRefresh: jest.fn(),
      };

      const newValue = {
        enabled: true,
        value: 0,
        unit: DATETIME_FORMATS.second,
      };

      watcher.call(context, newValue, {});

      expect(context.startPeriodicRefresh).not.toHaveBeenCalled();
    });

    test('Should pause periodic refresh when activeViewPeriodicRefreshPaused becomes true', () => {
      const watcher = widgetPeriodicRefreshMixin.watch.activeViewPeriodicRefreshPaused;
      const context = {
        ...mockContext,
        pausePriodicRefresh: jest.fn(),
        resumePriodicRefresh: jest.fn(),
      };

      watcher.call(context, true);

      expect(context.pausePriodicRefresh).toHaveBeenCalled();
      expect(context.resumePriodicRefresh).not.toHaveBeenCalled();
    });

    test('Should resume periodic refresh when activeViewPeriodicRefreshPaused becomes false', () => {
      const watcher = widgetPeriodicRefreshMixin.watch.activeViewPeriodicRefreshPaused;
      const context = {
        ...mockContext,
        pausePriodicRefresh: jest.fn(),
        resumePriodicRefresh: jest.fn(),
      };

      watcher.call(context, false);

      expect(context.resumePriodicRefresh).toHaveBeenCalled();
      expect(context.pausePriodicRefresh).not.toHaveBeenCalled();
    });
  });

  describe('beforeDestroy', () => {
    test('Should stop periodic refresh on component destroy', () => {
      const context = {
        ...mockContext,
        stopPeriodicRefresh: jest.fn(),
      };

      widgetPeriodicRefreshMixin.beforeDestroy.call(context);

      expect(context.stopPeriodicRefresh).toHaveBeenCalled();
    });
  });
});
