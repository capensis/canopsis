import { generateShallowRenderer } from '@unit/utils/vue';

import CChip from '@/components/common/chips/c-chip.vue';

const stubs = {
  'v-chip': { template: '<div><slot /></div>' },
  'v-icon': {
    template: '<span class="v-icon" @click="$emit(\'click\')"></span>',
  },
};

const selectCloseIcon = wrapper => wrapper.find('v-icon-stub');

describe('c-chip', () => {
  const factory = generateShallowRenderer(CChip, { stubs });

  test('Renders `c-alarm-tag-chip` with default props', () => {
    const wrapper = factory({
      propsData: {
        closable: true,
      },
    });

    const closeIcon = selectCloseIcon(wrapper);

    wrapper.vm.$emit('close');

    expect(wrapper.emitted('close')).toHaveLength(1);
  });
});
