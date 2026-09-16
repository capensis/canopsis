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
    const wrapper = factory({ propsData: { expanded: false } });

    wrapper.find('button.v-btn').trigger('click');

    expect(wrapper).toEmit('expand', true);
  });

  it('Expand button collapse is worked.', () => {
    const wrapper = factory({ propsData: { expanded: true } });

    wrapper.find('button.v-btn').trigger('click');

    expect(wrapper).toEmit('expand', false);
  });
});
