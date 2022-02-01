import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import ViewTabWidgets from '@/components/other/view/view-tab-widgets.vue';

const localVue = createVueInstance();

const createWidgetsGridStub = className => ({
  props: ['update-tab-method'],
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
  const widgets = [
    { title: 'Widget 1', _id: 'id' },
    { title: 'Widget 2', _id: 'id2' },
  ];

  it('Each query removed after destroy', () => {
    const removeQuery = jest.fn();
    const wrapper = factory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
      },
      store: createMockedStoreModules([{
        name: 'query',
        actions: {
          remove: removeQuery,
        },
      }]),
    });

    wrapper.destroy();

    expect(removeQuery).toHaveBeenCalledTimes(widgets.length);
    expect(
      removeQuery.mock.calls,
    ).toEqual(
      widgets.map(({ _id: id }) => [expect.any(Object), { id }, undefined]),
    );
  });

  it('Event emitted after trigger edition grid', () => {
    const wrapper = factory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
        isEditingMode: true,
      },
      store: createMockedStoreModules([{
        name: 'query',
        actions: {
          remove: jest.fn(),
        },
      }]),
    });

    const gridEditWidgetsElement = wrapper.find('.grid-edit-widgets');

    const newWidgets = [{ _id: Faker.datatype.string() }];

    gridEditWidgetsElement.vm.$emit('update:widgets-fields', newWidgets);

    const updateWidgetsFields = wrapper.emitted('update:widgets-fields');
    expect(updateWidgetsFields).toHaveLength(1);

    const [eventData] = updateWidgetsFields[0];

    expect(eventData).toEqual(newWidgets);
  });

  it('Update tab method called', () => {
    const updateTabMethod = jest.fn();
    const wrapper = factory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
        isEditingMode: true,
        updateTabMethod,
      },
      store: createMockedStoreModules([{
        name: 'query',
        actions: {
          remove: jest.fn(),
        },
      }]),
    });

    const widget = { id: Faker.datatype.string() };

    const gridEditWidgetsElement = wrapper.find('.grid-edit-widgets');

    gridEditWidgetsElement.vm.updateTabMethod(widget);

    expect(updateTabMethod).toHaveBeenCalledTimes(1);
    expect(updateTabMethod).toHaveBeenCalledWith(widget);
  });

  it('Default update tab method called', () => {
    const wrapper = factory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets,
        },
        isEditingMode: true,
      },
      store: createMockedStoreModules([{
        name: 'query',
        actions: {
          remove: jest.fn(),
        },
      }]),
    });

    const widget = { id: Faker.datatype.string() };

    const gridEditWidgetsElement = wrapper.find('.grid-edit-widgets');

    gridEditWidgetsElement.vm.updateTabMethod(widget);
  });

  it('Renders `view-tab-widgets` with default and required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          id: 'tab-id',
          widgets: [],
        },
      },
      store: createMockedStoreModules([{ name: 'query' }]),
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
        isEditingMode: true,
      },
      store: createMockedStoreModules([{ name: 'query' }]),
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
      store: createMockedStoreModules([{ name: 'query' }]),
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
