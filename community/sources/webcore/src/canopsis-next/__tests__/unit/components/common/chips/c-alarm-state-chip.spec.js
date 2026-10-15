import { generateShallowRenderer } from '@unit/utils/vue';

import CAlarmStateChip from '@/components/common/chips/c-alarm-state-chip.vue';

const stubs = {
  'c-alarm-chip': true,
};

const selectAlarmChip = wrapper => wrapper.find('c-alarm-chip-stub');

describe('c-alarm-state-chip', () => {
  const factory = generateShallowRenderer(CAlarmStateChip, { stubs });

  test('Click emitted after trigger click on chip', () => {
    const wrapper = factory();

    selectAlarmChip(wrapper).triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('click');
  });
});
