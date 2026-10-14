import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_LIST_STEPS } from '@/constants';

import CAlarmStateChip from '@/components/common/chips/c-alarm-state-chip.vue';

const stubs = {
  'c-alarm-chip': true,
};

const selectAlarmChip = wrapper => wrapper.find('c-alarm-chip-stub');

describe('c-alarm-state-chip', () => {
  const snapshotFactory = generateRenderer(CAlarmStateChip, { stubs });
  const factory = generateShallowRenderer(CAlarmStateChip, { stubs });

  test('Click emitted after trigger click on chip', () => {
    const wrapper = factory();

    selectAlarmChip(wrapper).triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('click');
  });

  test('Renders `c-alarm-state-chip` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-alarm-state-chip` with alarm state', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'custom-state-val',
        type: ALARM_LIST_STEPS.changeState,
        badgeValue: 'Events count',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-alarm-state-chip` with alarm state, small and appendIconName', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'custom-state-val',
        type: ALARM_LIST_STEPS.changeState,
        badgeValue: 'Events count',
        small: true,
        appendIconName: 'lock',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
