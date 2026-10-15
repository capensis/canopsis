import { generateShallowRenderer } from '@unit/utils/vue';

import CExpandBtn from '@/components/common/buttons/c-expand-btn.vue';

const stubs = {
  'v-btn': {
    template: `
      <button
        class="v-btn"
        @click="$emit('click', $event)"
      >
        <slot />
      </button>
    `,
  },
};

describe('c-expand-btn', () => {
  const factory = generateShallowRenderer(CExpandBtn, { stubs });

  it('Expand button expand is worked.', () => {
    const wrapper = factory({ propsData: { modelValue: false } });

    wrapper.find('button.v-btn').trigger('click');

    expect(wrapper.emitted('update:modelValue').at(-1)).toEqual([true]);
  });

  it('Expand button collapse is worked.', () => {
    const wrapper = factory({ propsData: { modelValue: true } });

    wrapper.find('button.v-btn').trigger('click');

    expect(wrapper.emitted('update:modelValue').at(-1)).toEqual([false]);
  });
});
