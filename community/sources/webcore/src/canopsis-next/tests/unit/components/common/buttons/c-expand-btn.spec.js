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

    const buttonElement = wrapper.find('button.v-btn');

    buttonElement.trigger('click');

    const expandEvents = wrapper.emitted('expand');

    expect(expandEvents).toHaveLength(1);
    expect(expandEvents[0]).toEqual([true]);
  });

  it('Expand button collapse is worked.', () => {
    const wrapper = factory({ propsData: { expanded: true } });

    const buttonElement = wrapper.find('button.v-btn');

    buttonElement.trigger('click');

    const expandEvents = wrapper.emitted('expand');

    expect(expandEvents).toHaveLength(1);
    expect(expandEvents[0]).toEqual([false]);
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
