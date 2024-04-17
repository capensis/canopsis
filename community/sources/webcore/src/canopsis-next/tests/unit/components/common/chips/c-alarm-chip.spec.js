import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_STATES, ALARM_STATUSES, ENTITY_INFOS_TYPE } from '@/constants';

import CAlarmChip from '@/components/common/chips/c-alarm-chip.vue';

const selectChipContainer = wrapper => wrapper.find('.chip-container');

describe('c-alarm-chip', () => {
  const snapshotFactory = generateRenderer(CAlarmChip);
  const factory = generateShallowRenderer(CAlarmChip);

  test('Click event emitted', () => {
    const wrapper = factory();

    selectChipContainer(wrapper).trigger('click');

    expect(wrapper).toHaveBeenEmit('click');
  });

  test('Renders `c-alarm-chip` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-alarm-chip` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 2,
        type: ENTITY_INFOS_TYPE.status,
        badgeValue: 3,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test.each(Object.entries(ALARM_STATUSES))('Renders `c-alarm-chip` with status: %s', (_, value) => {
    const wrapper = snapshotFactory({
      propsData: {
        value,
        type: ENTITY_INFOS_TYPE.status,
        badgeValue: value + 10,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test.each(Object.entries(ALARM_STATES))('Renders `c-alarm-chip` with state: %s', (_, value) => {
    const wrapper = snapshotFactory({
      propsData: {
        value,
        type: ENTITY_INFOS_TYPE.state,
        badgeValue: value + 20,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
