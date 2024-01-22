import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CExpandBtn from '@/components/common/buttons/c-expand-btn.vue';

const stubs = {
  'v-btn': {
    template: `
      <button
        class="v-btn"
        @click="$listeners.click"
      >
        <slot />
      </button>
    `,
  },
};

describe('c-expand-btn', () => {
  const factory = generateShallowRenderer(CExpandBtn, { stubs });
  const snapshotFactory = generateRenderer(CExpandBtn);

  it('Expand button expand is worked.', () => {
    const wrapper = factory({ propsData: { expanded: false } });

    wrapper.find('button.v-btn').trigger('click');

    expect(wrapper).toEmit('expand', [true]);
  });

  it('Expand button collapse is worked.', () => {
    const wrapper = factory({ propsData: { expanded: true } });

    wrapper.find('button.v-btn').trigger('click');

    expect(wrapper).toEmit('expand', [false]);
  });

  it('Renders `c-expand-btn` correctly.', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-expand-btn` expanded correctly.', () => {
    const wrapper = snapshotFactory({
      propsData: { expanded: true },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-expand-btn` with custom color correctly.', () => {
    const wrapper = snapshotFactory({
      propsData: { color: 'custom-color' },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
