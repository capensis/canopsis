import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';

import GridEditWidgets from '@/components/widgets/grid-edit-widgets.vue';

const localVue = createVueInstance();

const stubs = {
  portal: true,
  'widget-wrapper-menu': true,
  'grid-layout': true,
  'grid-item': true,
};

const snapshotStubs = {
  'widget-wrapper-menu': true,
  portal: {
    props: ['to'],
    template: `
      <div class="portal">
        Portal to: {{ to }}
        <slot />
      </div>
    `,
  },
};

const factory = (options = {}) => shallowMount(GridEditWidgets, {
  localVue,
  stubs,
  mocks: {
    $mq: 'l',
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(GridEditWidgets, {
  localVue,
  stubs: snapshotStubs,
  mocks: {
    $mq: 'l',
  },

  ...options,
});

const selectGridItems = wrapper => wrapper.findAll('grid-item-stub');

describe('grid-edit-widgets', () => {
  const widgets = [
    {
      _id: 'widget_Context_505742f9-faf5-445e-a537-2288a84fc58e',
      grid_parameters: {
        desktop: { autoHeight: true, h: 14, w: 12, x: 0, y: 0 },
        mobile: { autoHeight: true, h: 12, w: 3, x: 0, y: 0 },
        tablet: { autoHeight: true, h: 1, w: 12, x: 0, y: 0 },
      },
    },
    {
      _id: 'widget_ServiceWeather_43a12599-5800-4a86-b6f4-50bf186c4840',
      grid_parameters: {
        desktop: { autoHeight: true, h: 24, w: 12, x: 0, y: 14 },
        mobile: { autoHeight: true, h: 12, w: 3, x: 0, y: 0 },
        tablet: { autoHeight: true, h: 1, w: 12, x: 0, y: 0 },
      },
    },
    {
      _id: 'widget_ServiceWeather_58e5c9a5-aa04-4dc6-a59d-6fa847bc62e0',
      grid_parameters: {
        desktop: { autoHeight: true, h: 21, w: 12, x: 0, y: 38 },
        mobile: { autoHeight: true, h: 1, w: 12, x: 0, y: 12 },
        tablet: { autoHeight: true, h: 1, w: 12, x: 0, y: 1 },
      },
    },
  ];

  it('Widgets watchers updating layouts', async () => {
    const wrapper = factory({
      propsData: {
        tab: {
          widgets: [],
        },
        updateTabMethod: jest.fn(),
      },
    });

    expect(selectGridItems(wrapper)).toHaveLength(0);

    await wrapper.setProps({
      tab: {
        widgets,
      },
    });

    expect(selectGridItems(wrapper)).toHaveLength(widgets.length);
  });

  it('Auto height updated after clock on the button', async () => {
    const widget = widgets[0];
    const wrapper = factory({
      propsData: {
        tab: {
          widgets: [widget],
        },
        updateTabMethod: jest.fn(),
      },
    });

    const updateAutoHeightButton = selectGridItems(wrapper)
      .at(0)
      .find('v-btn-stub');

    updateAutoHeightButton.vm.$emit('click');

    const gridLayout = wrapper.find('grid-layout-stub');

    gridLayout.vm.$emit('layout-updated');

    const updateWidgetsEvents = wrapper.emitted('update:widgets-fields');

    expect(updateWidgetsEvents).toHaveLength(1);

    const [updateWidgets] = updateWidgetsEvents[0];

    expect(updateWidgets).toEqual({
      [widget._id]: {
        'grid_parameters.mobile': expect.any(Function),
        'grid_parameters.tablet': expect.any(Function),
        'grid_parameters.desktop': expect.any(Function),
      },
    });
    const updateWidgetParameterFunction = updateWidgets[widget._id]['grid_parameters.desktop'];
    const { desktop } = widget.grid_parameters;

    expect(updateWidgetParameterFunction()).toEqual({
      ...desktop,
      autoHeight: !desktop.autoHeight,
    });
  });

  it('Update tab method called', () => {
    const updateTabMethod = jest.fn();
    const wrapper = factory({
      propsData: {
        tab: {
          widgets,
        },
        updateTabMethod,
      },
    });

    const widget = { id: Faker.datatype.string() };

    const widgetWrapperMenu = selectGridItems(wrapper)
      .at(0)
      .find('widget-wrapper-menu-stub');

    widgetWrapperMenu.vm.updateTabMethod(widget);

    expect(updateTabMethod).toHaveBeenCalledTimes(1);
    expect(updateTabMethod).toHaveBeenCalledWith(widget);
  });

  it('Renders `grid-edit-widgets` with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          widgets: [],
        },
        updateTabMethod: jest.fn(),
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `grid-edit-widgets` with widgets', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          widgets,
        },
        updateTabMethod: jest.fn(),
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it.each(['m', 't', 'l', 'xl'])('Renders `grid-edit-widgets` with widgets on the %s window size', async (size) => {
    const wrapper = snapshotFactory({
      propsData: {
        tab: {
          widgets,
        },
        updateTabMethod: jest.fn(),
      },
      mocks: {
        $mq: size,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
