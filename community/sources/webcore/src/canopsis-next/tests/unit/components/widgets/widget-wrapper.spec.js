import { generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import { CANOPSIS_EDITION, WIDGET_TYPES } from '@/constants';

import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

const stubs = {
  'c-alert-overlay': true,

  'alarms-list-widget': true,
  'entities-list-widget': true,
  'service-weather-widget': true,
  'stats-calendar-widget': true,
  'text-widget': true,
  'counter-widget': true,
  'testing-weather-widget': true,
  'map-widget': true,
  'bar-chart-widget': true,
  'line-chart-widget': true,
  'pie-chart-widget': true,
  'numbers-widget': true,
  'user-statistics-widget': true,
  'alarm-statistics-widget': true,
  'availability-widget': true,
};

describe('widget-wrapper', () => {
  const types = Object.values(WIDGET_TYPES);
  const tabId = 'tab-id';

  const snapshotFactory = generateRenderer(WidgetWrapper, {

    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  it('Renders `widget-wrapper` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {
          type: WIDGET_TYPES.alarmList,
        },
        tab: {
          _id: tabId,
        },
        editing: false,
      },
      store: createMockedStoreModules([{
        name: 'info',
        getters: {
          edition: CANOPSIS_EDITION.pro,
        },
      }]),
    });

    expect(wrapper).toMatchSnapshot();
  });

  it.each(types)('Renders `widget-wrapper` with type %s and core edition', (type) => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {
          type,
        },
        tab: {
          _id: tabId,
        },
        editing: false,
      },
      store: createMockedStoreModules([{
        name: 'info',
        getters: {
          edition: CANOPSIS_EDITION.community,
        },
      }]),
    });

    expect(wrapper).toMatchSnapshot();
  });

  it.each(types)('Renders `widget-wrapper` with type %s and cat edition', (type) => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {
          type,
        },
        tab: {
          _id: tabId,
        },
        editing: false,
      },
      store: createMockedStoreModules([{
        name: 'info',
        getters: {
          edition: CANOPSIS_EDITION.pro,
        },
      }]),
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `widget-wrapper` with widget title', () => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {
          type: WIDGET_TYPES.alarmList,
          title: 'Title',
        },
        tab: {
          _id: tabId,
        },
        editing: false,
      },
      store: createMockedStoreModules([{
        name: 'info',
        getters: {
          edition: CANOPSIS_EDITION.pro,
        },
      }]),
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `widget-wrapper` with editing mode', () => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {
          type: WIDGET_TYPES.alarmList,
        },
        tab: {
          _id: tabId,
        },
        editing: true,
      },
      store: createMockedStoreModules([{
        name: 'info',
        getters: {
          edition: CANOPSIS_EDITION.pro,
        },
      }]),
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `widget-wrapper` with title and editing mode', () => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {
          type: WIDGET_TYPES.alarmList,
          title: 'Title',
        },
        tab: {
          _id: tabId,
        },
        editing: true,
      },
      store: createMockedStoreModules([{
        name: 'info',
        getters: {
          edition: CANOPSIS_EDITION.pro,
        },
      }]),
    });

    expect(wrapper).toMatchSnapshot();
  });
});
