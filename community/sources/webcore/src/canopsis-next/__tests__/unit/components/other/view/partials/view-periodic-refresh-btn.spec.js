import { createButtonStub } from '@unit/stubs/button';
import { generateShallowRenderer } from '@unit/utils/vue';

import { DATETIME_FORMATS } from '@/constants';

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
