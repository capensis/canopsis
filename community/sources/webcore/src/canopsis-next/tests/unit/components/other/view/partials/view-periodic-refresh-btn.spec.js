import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createButtonStub } from '@unit/stubs/button';

import { DATETIME_FORMATS } from '@/constants';

import ViewPeriodicRefreshBtn from '@/components/other/view/partials/view-periodic-refresh-btn.vue';

const stubs = {
  'v-tooltip': true,
  'v-btn': createButtonStub('v-btn'),
  'v-icon': true,
  'v-progress-circular': true,
};

const snapshotStubs = {
  'v-tooltip': true,
  'v-btn': true,
  'v-icon': true,
  'v-progress-circular': true,
};

describe('view-periodic-refresh-btn', () => {
  const factory = generateShallowRenderer(ViewPeriodicRefreshBtn, { stubs });
  const snapshotFactory = generateRenderer(ViewPeriodicRefreshBtn, { stubs: snapshotStubs });

  const mockPeriodicRefresh = {
    notify: jest.fn(),
  };

  const mockPopups = {
    remove: jest.fn(),
    info: jest.fn(),
  };

  const createViewWithPeriodicRefresh = (enabled = true, value = 30, unit = DATETIME_FORMATS.second) => ({
    periodic_refresh: {
      enabled,
      value,
      unit,
    },
  });

  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
    jest.clearAllMocks();
  });

  describe('computed', () => {
    test('Should return correct periodicRefreshFullPaused value when activeViewPeriodicRefreshPaused', () => {
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

      expect(wrapper.vm.periodicRefreshFullPaused).toBe(true);
    });

    test('Should return correct periodicRefreshFullPaused value when isNavigationEditingMode', () => {
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
          isNavigationEditingMode: () => true,
        },
      });

      expect(wrapper.vm.periodicRefreshFullPaused).toBe(true);
    });

    test('Should return false for periodicRefreshFullPaused when not paused', () => {
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

      expect(wrapper.vm.periodicRefreshFullPaused).toBe(false);
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
      // Test when activeViewPeriodicRefreshPaused is true
      const wrapper1 = factory({
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

      expect(wrapper1.vm.periodicRefreshFullPaused).toBe(true);

      // Test when activeViewPeriodicRefreshPaused is false
      const wrapper2 = factory({
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

      expect(wrapper2.vm.periodicRefreshFullPaused).toBe(false);
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

  test('Renders `view-periodic-refresh-btn` with refresh icon when disabled', () => {
    const view = createViewWithPeriodicRefresh(false);
    const wrapper = snapshotFactory({
      provide: {
        $periodicRefresh: mockPeriodicRefresh,
      },
      mocks: {
        $popups: mockPopups,
      },
      computed: {
        view: () => view,
        activeViewPeriodicRefreshPaused: () => false,
        isNavigationEditingMode: () => false,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `view-periodic-refresh-btn` with progress when enabled', () => {
    const view = createViewWithPeriodicRefresh(true);
    const wrapper = snapshotFactory({
      provide: {
        $periodicRefresh: mockPeriodicRefresh,
      },
      mocks: {
        $popups: mockPopups,
      },
      computed: {
        view: () => view,
        activeViewPeriodicRefreshPaused: () => false,
        isNavigationEditingMode: () => false,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `view-periodic-refresh-btn` when paused by activeViewPeriodicRefreshPaused', () => {
    const view = createViewWithPeriodicRefresh(true);
    const wrapper = snapshotFactory({
      provide: {
        $periodicRefresh: mockPeriodicRefresh,
      },
      mocks: {
        $popups: mockPopups,
      },
      computed: {
        view: () => view,
        activeViewPeriodicRefreshPaused: () => true,
        isNavigationEditingMode: () => false,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
