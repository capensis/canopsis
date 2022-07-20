import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import ViewTabWidgets from '@/components/other/view/view-tab-widgets.vue';

const localVue = createVueInstance();

const createWidgetsGridStub = className => ({
  template: `
    <div class="${className}">
      <slot :widget="{}" />
    </div>
  `,
});

const stubs = {
  'grid-overview-widget': createWidgetsGridStub('grid-overview-widget'),
  'grid-edit-widgets': createWidgetsGridStub('grid-edit-widgets'),
  'widget-wrapper': true,
};

const factory = (options = {}) => shallowMount(ViewTabWidgets, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(ViewTabWidgets, {
  localVue,
  stubs,

  ...options,
});

describe('view-tab-widgets', () => {
  const removeQuery = jest.fn();
  const fetchActiveView = jest.fn();
  const registerEditingOffHandler = jest.fn();
  const unregisterEditingOffHandler = jest.fn();
  const updateGridPositions = jest.fn();

  const queryModule = {
    name: 'query',
    actions: {
      remove: removeQuery,
    },
  };

  const activeViewModule = {
    name: 'activeView',
    getters: {
      editing: false,
    },
    actions: {
      fetch: fetchActiveView,
      registerEditingOffHandler,
      unregisterEditingOffHandler,
    },
  };

  const widgetModule = {
    name: 'view/widget',
    actions: {
      updateGridPositions,
    },
  };

  const store = createMockedStoreModules([
    queryModule,
    activeViewModule,
    widgetModule,
  ]);

  const widgets = [
    { title: 'Widget 1', _id: 'id' },
    { title: 'Widget 2', _id: 'id2' },
  ];

  afterEach(() => {
    removeQuery.mockReset();
    fetchActiveView.mockReset();
    registerEditingOffHandler.mockReset();
    unregisterEditingOffHandler.mockReset();
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
    expect(
      removeQuery.mock.calls,
    ).toEqual(
      widgets.map(({ _id: id }) => [expect.any(Object), { id }, undefined]),
    );
  });

  it('Register and unregister editing off handler is working', async () => {
    const wrapper = factory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
      },
      store,
    });

    await flushPromises();

    expect(registerEditingOffHandler).toHaveBeenCalledTimes(1);

    wrapper.destroy();

    expect(unregisterEditingOffHandler).toHaveBeenCalledTimes(1);
  });

  it('Event emitted after trigger edition grid', async () => {
    const wrapper = factory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
      },
      store: createMockedStoreModules([
        queryModule,
        widgetModule,
        {
          ...activeViewModule,

          getters: {
            editing: true,
          },
        },
      ]),
    });

    const gridEditWidgetsElement = wrapper.find('.grid-edit-widgets');

    const data = {
      'widget-id': {
        desktop: {
          autoHeight: true,
          x: 0,
          y: 0,
          h: 1,
          w: 12,
        },
        tablet: {
          autoHeight: true,
          x: 0,
          y: 0,
          h: 1,
          w: 12,
        },
        mobile: {
          autoHeight: true,
          x: 0,
          y: 0,
          h: 1,
          w: 12,
        },
      },
    };
    gridEditWidgetsElement.vm.$emit('update:widgets-grid', data);

    expect(wrapper.vm.widgetsGrid).toEqual(data);

    await wrapper.vm.updatePositions();

    expect(updateGridPositions).toHaveBeenCalledTimes(1);
    expect(updateGridPositions).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        data: Object.entries(data).map(([key, value]) => ({ _id: key, grid_parameters: value })),
      },
      undefined,
    );
    expect(fetchActiveView).toHaveBeenCalledTimes(1);
  });

  it('Renders `view-tab-widgets` with editing mode', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets: [],
        },
      },
      store,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `view-tab-widgets` with default and required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets: [],
        },
      },
      store,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `view-tab-widgets` with editing mode', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets: [],
        },
      },
      store: createMockedStoreModules([
        queryModule,
        widgetModule,
        {
          ...activeViewModule,

          getters: {
            editing: true,
          },
        },
      ]),
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `view-tab-widgets` with widgets', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
      },
      store,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `view-tab-widgets` with widgets with editing mode', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
      },
      store: createMockedStoreModules([
        queryModule,
        widgetModule,
        {
          ...activeViewModule,

          getters: {
            editing: true,
          },
        },
      ]),
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
