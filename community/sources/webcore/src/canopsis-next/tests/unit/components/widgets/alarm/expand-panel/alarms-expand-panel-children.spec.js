import { mount, createVueInstance } from '@unit/utils/vue';
import { fakeStaticAlarms, fakeAlarm } from '@unit/data/alarm';
import { generateDefaultAlarmListWidget } from '@/helpers/entities';

import AlarmsExpandPanelChildren from '@/components/widgets/alarm/expand-panel/alarms-expand-panel-children.vue';

jest.mock('file-saver', () => ({
  saveAs: jest.fn(),
}));

const localVue = createVueInstance();

const stubs = {
  'alarms-list-table-with-pagination': true,
};

const snapshotFactory = (options = {}) => mount(AlarmsExpandPanelChildren, {
  localVue,
  stubs,

  ...options,
});

describe('alarms-expand-panel-children', () => {
  const nowTimestamp = 1386435600000;
  const totalCount = 15;
  const alarm = fakeAlarm();
  const childrenAlarms = fakeStaticAlarms({
    totalItems: totalCount,
    timestamp: nowTimestamp,
  });

  const children = {
    data: childrenAlarms,
    meta: {
      page: 1,
      total_count: totalCount,
    },
  };

  const query = {
    page: 1,
    limit: 10,
  };

  const widget = generateDefaultAlarmListWidget();

  it('Renders `alarms-expand-panel-children` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm,
        children,
        widget,
        query,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel-children` with editing', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm,
        children,
        widget,
        query,
        editing: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel-children` with causes alarms', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm,
        children,
        widget,
        query,
        pending: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
