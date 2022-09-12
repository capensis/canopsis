import { mount, createVueInstance } from '@unit/utils/vue';
import { fakeStaticAlarms } from '@unit/data/alarm';

import AlarmGeneralList from '@/components/widgets/alarm/alarm-general-list.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(AlarmGeneralList, {
  localVue,

  ...options,
});

describe('alarm-general-list', () => {
  const alarms = fakeStaticAlarms({ totalItems: 4 });

  test('Renders `alarm-general-list` with empty items', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `alarm-general-list` with alarms', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: alarms,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
