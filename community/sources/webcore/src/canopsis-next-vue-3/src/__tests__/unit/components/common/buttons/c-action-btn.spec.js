import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import CActionBtn from '@/components/common/buttons/c-action-btn.vue';

const stubs = {
  'v-tooltip': {
    template: '<div><slot name="activator" :props="{}" /><slot /></div>',
  },
  'v-btn': {
    template: `
      <button
        class="v-btn"
        v-bind="$attrs"
        @click="$emit('click', $event)"
      >
        <slot />
      </button>
    `,
  },
  'c-icon': {
    inheritAttrs: false,
    template: '<span class="c-icon" v-bind="$attrs"><slot /></span>',
  },
};

describe('c-action-btn', () => {
  const factory = generateShallowRenderer(CActionBtn, { stubs });

  it('Action button with default type and custom icon.', () => {
    const icon = Faker.datatype.string();
    const wrapper = factory({ propsData: { icon, type: 'edit' } });

    const iconElement = wrapper.find('.c-icon');

    expect(iconElement.text()).toBe(icon);
  });

  it('Action button with default type and custom color.', () => {
    const color = Faker.datatype.string();
    const wrapper = factory({ propsData: { color, type: 'duplicate' } });

    const iconElement = wrapper.find('.c-icon');

    expect(iconElement.attributes('color')).toBe(color);
  });

  it('Action button without default type and all custom props.', () => {
    const icon = Faker.datatype.string();
    const tooltip = Faker.datatype.string();
    const color = Faker.datatype.string();
    const wrapper = factory({ propsData: { color, icon, tooltip } });

    const iconElement = wrapper.find('.c-icon');

    expect(iconElement.text()).toBe(icon);
    expect(iconElement.attributes('color')).toBe(color);
  });

  it('Check loading property.', () => {
    const wrapper = factory({ propsData: { loading: true } });

    const buttonElement = wrapper.find('button.v-btn');

    expect(buttonElement.attributes('loading')).toBeTruthy();
  });

  it('Check disabled property.', () => {
    const wrapper = factory({ propsData: { disabled: true } });

    const buttonElement = wrapper.find('button.v-btn');

    expect(buttonElement.attributes()).toHaveProperty('disabled');
  });

  it('Check button slot.', () => {
    const wrapper = factory({
      slots: {
        button: '<div class="name-slot" />',
      },
    });

    const slotElement = wrapper.find('div.name-slot');

    expect(slotElement.exists()).toBeTruthy();
  });

  it('Click event working correctly.', async () => {
    const onClick = vi.fn();
    const wrapper = factory({
      attrs: { onClick },
    });

    const button = wrapper.find('button.v-btn');

    wrapper.vm.$emit('click');

    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
