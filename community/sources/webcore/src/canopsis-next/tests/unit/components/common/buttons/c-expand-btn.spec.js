import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import CExpandBtn from '@/components/common/buttons/c-expand-btn.vue';

const localVue = createVueInstance();

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

const factory = (options = {}) => shallowMount(CExpandBtn, {
  localVue,
  stubs,
  ...options,
});

describe('c-expand-btn', () => {
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
    const wrapper = mount(CExpandBtn, {
      localVue,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-expand-btn` expanded correctly.', () => {
    const wrapper = mount(CExpandBtn, {
      localVue,
      propsData: { expanded: true },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-expand-btn` with custom color correctly.', () => {
    const wrapper = mount(CExpandBtn, {
      localVue,
      propsData: { color: 'custom-color' },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
