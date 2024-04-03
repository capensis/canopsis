import { generateRenderer } from '@unit/utils/vue';
import { fakeStaticAlarms } from '@unit/data/alarm';

import AlarmGeneralList from '@/components/widgets/alarm/alarm-general-list.vue';

describe('alarm-general-list', () => {
  const alarms = fakeStaticAlarms({ totalItems: 4 });

  const snapshotFactory = generateRenderer(AlarmGeneralList);

  test('Renders `alarm-general-list` with empty items', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-general-list` with alarms', () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: alarms,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
