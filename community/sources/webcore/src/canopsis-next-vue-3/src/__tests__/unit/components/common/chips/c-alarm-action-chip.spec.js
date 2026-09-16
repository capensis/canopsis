import { generateShallowRenderer } from '@unit/utils/vue';

import CChip from '@/components/common/chips/c-chip.vue';

const stubs = {
  'v-icon': {
    template: '<span v-on="$attrs" class="v-icon"></span>',
  },
};

const selectCloseIcon = wrapper => wrapper.find('.v-icon');

describe('c-chip', () => {
  const factory = generateShallowRenderer(CChip, { stubs });

  test('Renders `c-alarm-tag-chip` with default props', () => {
    const wrapper = factory({
      propsData: {
        closable: true,
      },
    });

    const closeIcon = selectCloseIcon(wrapper);

    closeIcon.trigger('click');

    expect(wrapper).toHaveBeenEmit('close');
  });
});
