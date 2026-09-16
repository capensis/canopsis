import {
  createQueryModule,
  createActiveViewModule,
  createMockedStoreModules,
  createWidgetModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { omit } from 'lodash';

import { VIEW_SCREEN_MODES, WIDGET_GRID_SIZES_KEYS, WIDGET_TYPES } from '@/constants';

import { setField } from '@/helpers/immutable';

import ViewTabWidgets from '@/components/other/view/view-tab-widgets.vue';

const stubs = {
  'grid-layout': true,
  'widget-edit-drag-handler': true,
  'window-size-field': true,
  portal: true,
};

const selectGridLayout = wrapper => wrapper.find('grid-layout-stub');

describe('view-tab-widgets', () => {
  const { queryModule, removeQuery } = createQueryModule();
  const { widgetModule, updateGridPositions } = createWidgetModule();
  const { editing, activeViewModule, registerEditingOffHandler, unregisterEditingOffHandler } =
    createActiveViewModule();

  const store = createMockedStoreModules([queryModule, activeViewModule, widgetModule]);

  const widgets = [
    {
      _id: 'widget_Context_505742f9-faf5-445e-a537-2288a84fc58e',
      type: WIDGET_TYPES.text,
      grid_parameters: {
        desktop: { autoHeight: true, h: 14, w: 12, x: 0, y: 0 },
        mobile: { autoHeight: true, h: 12, w: 3, x: 0, y: 0 },
        tablet: { autoHeight: true, h: 1, w: 12, x: 0, y: 0 },
      },
    },
    {
      _id: 'widget_ServiceWeather_43a12599-5800-4a86-b6f4-50bf186c4840',
      type: WIDGET_TYPES.text,
      grid_parameters: {
        desktop: { autoHeight: true, h: 24, w: 12, x: 0, y: 14 },
        mobile: { autoHeight: true, h: 12, w: 3, x: 0, y: 12 },
        tablet: { autoHeight: true, h: 1, w: 12, x: 0, y: 1 },
      },
    },
    {
      _id: 'widget_ServiceWeather_58e5c9a5-aa04-4dc6-a59d-6fa847bc62e0',
      type: WIDGET_TYPES.text,
      grid_parameters: {
        desktop: { autoHeight: true, h: 21, w: 12, x: 0, y: 38 },
        mobile: { autoHeight: true, h: 1, w: 12, x: 0, y: 24 },
        tablet: { autoHeight: true, h: 1, w: 12, x: 0, y: 2 },
      },
    },
  ];

  const factory = generateShallowRenderer(ViewTabWidgets, {
    stubs,
    mocks: {
      $mq: 'l',
    },
  });

  afterEach(() => {
    removeQuery.mockReset();
    updateGridPositions.mockReset();
  });

  it('Each query removed after destroy', () => {
    const wrapper = factory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
      },
      store,
    });

    wrapper.destroy();

    expect(removeQuery).toHaveBeenCalledTimes(widgets.length);
    expect(removeQuery.mock.calls).toEqual(
      widgets.map(({ _id: id }) => [expect.any(Object), { id }, undefined]),
    );
  });

  it('Register and unregister editing off handler is working (visible true)', async () => {
    const wrapper = factory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
        visible: true,
      },
      store,
    });

    await flushPromises();

    expect(registerEditingOffHandler).toHaveBeenCalledTimes(1);

    wrapper.destroy();

    expect(unregisterEditingOffHandler).toHaveBeenCalledTimes(1);
  });

  it('Register and unregister editing off handler is working (visible false)', async () => {
    const wrapper = factory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
        visible: false,
      },
      store,
    });

    await flushPromises();

    expect(registerEditingOffHandler).not.toHaveBeenCalled();

    wrapper.destroy();

    expect(unregisterEditingOffHandler).toHaveBeenCalledTimes(2);
  });

  it('Update positions doesn\'t trigger updateGridPositions without changes', async () => {
    editing.mockReturnValueOnce(true);

    const wrapper = factory({
      propsData: {
        tab: {
          _id: 'tab-id',
          widgets,
        },
      },
      store,
    });

    await wrapper.vm.updatePositions();

    expect(updateGridPositions).toHaveBeenCalledTimes(0);
  });

  it('Reflects active view screen mode from store', () => {
    const wrapperDefaultScreenMode = factory({
      propsData: {
        tab: {
          _id: 'tab-id',
          widgets: [],
        },
      },
      store,
    });

    expect(wrapperDefaultScreenMode.vm.activeViewScreenMode).toBe(VIEW_SCREEN_MODES.default);
    expect(wrapperDefaultScreenMode.vm.isDefaultMode).toBe(true);

    const activeViewFullscreenModule = {
      name: 'activeView',
      getters: {
        editing: vi.fn().mockReturnValue(false),
        item: vi.fn().mockReturnValue({}),
        screenMode: vi.fn().mockReturnValue(VIEW_SCREEN_MODES.fullscreen),
        isKioskScreenMode: vi.fn().mockReturnValue(false),
        periodicRefreshPaused: vi.fn().mockReturnValue(false),
      },
      actions: {
        registerEditingOffHandler: vi.fn(),
        unregisterEditingOffHandler: vi.fn(),
        toggleEditing: vi.fn(),
        setScreenMode: vi.fn(),
        fetch: vi.fn(),
        resumePeriodicRefresh: vi.fn(),
        pausePeriodicRefresh: vi.fn(),
      },
    };

    const storeFullscreen = createMockedStoreModules([
      queryModule,
      activeViewFullscreenModule,
      widgetModule,
    ]);

    const wrapperNonDefaultScreenMode = factory({
      propsData: {
        tab: {
          _id: 'tab-id',
          widgets: [],
        },
      },
      store: storeFullscreen,
    });

    expect(wrapperNonDefaultScreenMode.vm.activeViewScreenMode).toBe(VIEW_SCREEN_MODES.fullscreen);
    expect(wrapperNonDefaultScreenMode.vm.isDefaultMode).toBe(false);

    wrapperDefaultScreenMode.destroy();
    wrapperNonDefaultScreenMode.destroy();
  });

  it('Update positions triggers updateGridPositions with changes', async () => {
    editing.mockReturnValueOnce(true);

    const wrapper = factory({
      propsData: {
        tab: {
          _id: 'tab-id',
          widgets,
        },
      },
      store,
    });

    const gridLayoutElement = selectGridLayout(wrapper);
    const newHeightForThirdWidgetOnDesktop = 25;
    const newLayouts = setField(
      gridLayoutElement.vm.layout,
      [2, 'h'],
      newHeightForThirdWidgetOnDesktop,
    );

    gridLayoutElement.triggerCustomEvent('input', newLayouts);

    await flushPromises();
    await wrapper.vm.updatePositions();

    const newPositions = widgets.map((widget, index) =>
      index !== 2
        ? omit(widget, ['type'])
        : omit(
            setField(
              widget,
              ['grid_parameters', WIDGET_GRID_SIZES_KEYS.desktop, 'h'],
              newHeightForThirdWidgetOnDesktop,
            ),
            ['type'],
          ),
    );

    expect(updateGridPositions).toHaveBeenCalledTimes(1);
    expect(updateGridPositions).toHaveBeenCalledWith(expect.any(Object), { data: newPositions });
  });
});
