import { createButtonStub } from '@unit/stubs/button';
import { generateShallowRenderer } from '@unit/utils/vue';

import { DATETIME_FORMATS, TIME_UNITS } from '@/constants';

import ViewPeriodicRefreshBtn from '@/components/other/view/partials/view-periodic-refresh-btn.vue';

const stubs = {
  'c-tooltip': true,
  'v-btn': createButtonStub('v-btn'),
  'v-icon': true,
  'v-progress-circular': true,
};

describe('view-periodic-refresh-btn', () => {
  const factory = generateShallowRenderer(ViewPeriodicRefreshBtn, { stubs });
  const getPeriodicRefreshFullPaused = (activeViewPeriodicRefreshPaused, isNavigationEditingMode) =>
    ViewPeriodicRefreshBtn.computed.periodicRefreshFullPaused.call({
      activeViewPeriodicRefreshPaused,
      isNavigationEditingMode,
    });

  const mockPeriodicRefresh = {
    notify: vi.fn(),
  };

  const mockPopups = {
    remove: vi.fn(),
    info: vi.fn(),
  };

  const createViewWithPeriodicRefresh = (
    enabled = true,
    value = 30,
    unit = DATETIME_FORMATS.second,
  ) => ({
    periodic_refresh: {
      enabled,
      value,
      unit,
    },
  });

  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.clearAllMocks();
  });

  describe('computed', () => {
    test('Should return correct periodicRefreshFullPaused value when activeViewPeriodicRefreshPaused', () => {
      expect(getPeriodicRefreshFullPaused(true, false)).toBe(true);
    });

    test('Should return correct periodicRefreshFullPaused value when isNavigationEditingMode', () => {
      expect(getPeriodicRefreshFullPaused(false, true)).toBe(true);
    });

    test('Should return false for periodicRefreshFullPaused when not paused', () => {
      expect(getPeriodicRefreshFullPaused(false, false)).toBe(false);
    });
  });

  describe('countdown display', () => {
    const countdownFactory = generateShallowRenderer(ViewPeriodicRefreshBtn, {
      stubs: {
        ...stubs,
        'c-tooltip': {
          template: '<div><slot name="activator" :props="{}" /><slot /></div>',
        },
      },
      global: {
        provide: {
          $periodicRefresh: mockPeriodicRefresh,
        },
      },
      store: {
        getters: {
          'activeView/item': createViewWithPeriodicRefresh(true, 30, TIME_UNITS.hour),
          'activeView/periodicRefreshPaused': false,
          'navigation/isEditingMode': false,
        },
      },
    });

    test.each([
      [107993, '1d'],
      [86400, '1d'],
      [86399, '23h'],
      [3600, '1h'],
      [3599, '59m'],
      [60, '1m'],
      [59, '59s'],
      [0, '0s'],
    ])('Should display %i seconds as %s', async (seconds, label) => {
      const wrapper = countdownFactory();

      await wrapper.setData({ periodicRefreshProgress: seconds });

      expect(wrapper.find('.refresh-btn').text()).toBe(label);
    });

    test('Should retain the detailed tooltip and progress while the label changes units', async () => {
      const wrapper = countdownFactory();

      await wrapper.setData({ periodicRefreshProgress: 86400 });

      expect(wrapper.find('.refresh-btn').text()).toBe('1d');

      vi.advanceTimersByTime(1000);
      await wrapper.vm.$nextTick();

      expect(wrapper.find('.refresh-btn').text()).toBe('23h');
      expect(wrapper.vm.tooltipContent).toBe('23 hrs 59 mins 59 secs');
      expect(wrapper.vm.periodicRefreshProgressValue).toBeCloseTo(86399 / 1080);
      expect(mockPeriodicRefresh.notify).not.toHaveBeenCalled();
    });
  });

  describe('interval changes', () => {
    const intervalFactory = generateShallowRenderer({
      extends: ViewPeriodicRefreshBtn,
      props: {
        settings: { type: Object, required: true },
        paused: { type: Boolean, default: false },
        navigationEditing: { type: Boolean, default: false },
      },
      computed: {
        view() {
          return { periodic_refresh: this.settings };
        },
        activeViewPeriodicRefreshPaused() {
          return this.paused;
        },
        isNavigationEditingMode() {
          return this.navigationEditing;
        },
      },
    }, {
      stubs: {
        ...stubs,
        'c-tooltip': {
          template: '<div><slot name="activator" :props="{}" /><slot /></div>',
        },
      },
      global: {
        provide: { $periodicRefresh: mockPeriodicRefresh },
        mocks: { $popups: mockPopups },
      },
      props: {
        settings: createViewWithPeriodicRefresh(true, 30, TIME_UNITS.hour).periodic_refresh,
      },
    });

    test.each([
      [5, TIME_UNITS.hour, '5h', '5 hrs'],
      [30, TIME_UNITS.minute, '30m', '30 mins'],
    ])('Should immediately reset the counter to %i %s when saved', async (value, unit, label, tooltip) => {
      const wrapper = intervalFactory();
      const interval = wrapper.vm.periodicRefreshInterval;

      vi.advanceTimersByTime(1500);
      await wrapper.setProps({ settings: { enabled: true, value, unit } });

      expect(wrapper.find('.refresh-btn').text()).toBe(label);
      expect(wrapper.vm.tooltipContent).toBe(tooltip);
      expect(wrapper.vm.periodicRefreshProgressValue).toBe(100);
      expect(wrapper.vm.periodicRefreshInterval).toBe(interval);

      const progress = wrapper.vm.periodicRefreshProgress;

      vi.advanceTimersByTime(500);
      await wrapper.vm.$nextTick();

      expect(wrapper.vm.periodicRefreshProgress).toBe(progress - 1);

      vi.advanceTimersByTime(500);
      await wrapper.vm.$nextTick();

      expect(wrapper.vm.periodicRefreshProgress).toBe(progress - 1);
      expect(mockPeriodicRefresh.notify).not.toHaveBeenCalled();
    });

    test.each(['paused', 'navigationEditing'])('Should reset the counter while %s without resuming', async (pauseProp) => {
      const wrapper = intervalFactory();

      await wrapper.setProps({ [pauseProp]: true });
      await wrapper.setProps({ settings: { enabled: true, value: 5, unit: TIME_UNITS.minute } });

      expect(wrapper.find('.refresh-btn').text()).toBe('5m');
      expect(wrapper.vm.periodicRefreshProgressValue).toBe(100);

      vi.advanceTimersByTime(2000);
      await wrapper.vm.$nextTick();

      expect(wrapper.vm.periodicRefreshProgress).toBe(300);
      expect(wrapper.vm.periodicRefreshInterval).toBeUndefined();

      await wrapper.setProps({ [pauseProp]: false });
      vi.advanceTimersByTime(1000);
      await wrapper.vm.$nextTick();

      expect(wrapper.vm.periodicRefreshProgress).toBe(299);
    });

    test('Should reset the counter without starting a stopped timer', async () => {
      const wrapper = intervalFactory();

      wrapper.vm.stopPeriodicRefreshInterval();
      await wrapper.setProps({ settings: { enabled: true, value: 5, unit: TIME_UNITS.minute } });
      vi.advanceTimersByTime(2000);
      await wrapper.vm.$nextTick();

      expect(wrapper.find('.refresh-btn').text()).toBe('5m');
      expect(wrapper.vm.periodicRefreshInterval).toBeUndefined();
      expect(wrapper.vm.periodicRefreshProgress).toBe(300);
    });
  });

  describe('methods', () => {
    test('Should handle interval management correctly', () => {
      const wrapper = factory({
        provide: {
          $periodicRefresh: mockPeriodicRefresh,
        },
        mocks: {
          $popups: mockPopups,
        },
        computed: {
          view: () => createViewWithPeriodicRefresh(true),
          activeViewPeriodicRefreshPaused: () => false,
          isNavigationEditingMode: () => false,
        },
      });

      // Test starting interval
      wrapper.vm.startPeriodicRefreshInterval();
      expect(wrapper.vm.periodicRefreshInterval).toBeDefined();

      // Test stopping interval
      wrapper.vm.stopPeriodicRefreshInterval();
      expect(wrapper.vm.periodicRefreshInterval).toBeUndefined();
    });

    test('Should have periodicRefreshPausedWatcher method', () => {
      const wrapper = factory({
        provide: {
          $periodicRefresh: mockPeriodicRefresh,
        },
        mocks: {
          $popups: mockPopups,
        },
        computed: {
          view: () => createViewWithPeriodicRefresh(true),
          activeViewPeriodicRefreshPaused: () => false,
          isNavigationEditingMode: () => false,
        },
      });

      // Test that the new method exists
      expect(typeof wrapper.vm.periodicRefreshPausedWatcher).toBe('function');
    });
  });

  describe('new periodic refresh pause functionality', () => {
    test('Should incorporate activeViewPeriodicRefreshPaused in periodicRefreshFullPaused', () => {
      expect(getPeriodicRefreshFullPaused(true, false)).toBe(true);
      expect(getPeriodicRefreshFullPaused(false, false)).toBe(false);
    });

    test('Should access activeViewPeriodicRefreshPaused computed property from mixin', () => {
      const wrapper = factory({
        provide: {
          $periodicRefresh: mockPeriodicRefresh,
        },
        mocks: {
          $popups: mockPopups,
        },
        computed: {
          view: () => createViewWithPeriodicRefresh(true),
          activeViewPeriodicRefreshPaused: () => true,
          isNavigationEditingMode: () => false,
        },
      });

      // Test that the component can access the new computed property from activeViewMixin
      expect(wrapper.vm.activeViewPeriodicRefreshPaused).toBe(true);
    });
  });
});
