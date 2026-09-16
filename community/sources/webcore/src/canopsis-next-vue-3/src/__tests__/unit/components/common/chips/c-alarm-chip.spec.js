import { generateShallowRenderer } from '@unit/utils/vue';

import CAlarmChip from '@/components/common/chips/c-alarm-chip.vue';

const selectChipContainer = wrapper => wrapper.find('.chip-container');

describe('c-alarm-chip', () => {
  const factory = generateShallowRenderer(CAlarmChip);

  test('Click event emitted', () => {
    const wrapper = factory();

    selectChipContainer(wrapper).trigger('click');

    expect(wrapper).toHaveBeenEmit('click');
  });
});
