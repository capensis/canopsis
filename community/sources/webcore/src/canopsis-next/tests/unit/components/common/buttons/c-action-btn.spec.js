import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import CActionBtn from '@/components/common/buttons/c-action-btn.vue';

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

const factory = (options = {}) => shallowMount(CActionBtn, {
  localVue,
  stubs,
  ...options,
});

describe('c-action-btn', () => {
  it('Action button with default type and custom icon.', () => {
    const icon = Faker.datatype.string();
    const wrapper = factory({ propsData: { icon, type: 'edit' } });

    const iconElement = wrapper.find('v-icon-stub');

    expect(iconElement.text()).toBe(icon);
  });

  it('Action button with default type and custom tooltip.', () => {
    const tooltip = Faker.datatype.string();
    const wrapper = factory({ propsData: { tooltip, type: 'duplicate' } });

    const tooltipContentElement = wrapper.find('span');

    expect(tooltipContentElement.text()).toBe(tooltip);
  });

  it('Action button with default type and custom color.', () => {
    const color = Faker.datatype.string();
    const wrapper = factory({ propsData: { color, type: 'duplicate' } });

    const iconElement = wrapper.find('v-icon-stub');

    expect(iconElement.attributes('color')).toBe(color);
  });

  it('Action button without default type and all custom props.', () => {
    const icon = Faker.datatype.string();
    const tooltip = Faker.datatype.string();
    const color = Faker.datatype.string();
    const wrapper = factory({ propsData: { color, icon, tooltip } });

    const iconElement = wrapper.find('v-icon-stub');
    const tooltipContentElement = wrapper.find('span');

    expect(iconElement.text()).toBe(icon);
    expect(tooltipContentElement.text()).toBe(tooltip);
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

    expect(buttonElement.attributes('disabled')).toBeTruthy();
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

  it('Click event working correctly.', () => {
    const onClick = jest.fn();
    const wrapper = factory({
      listeners: {
        click: onClick,
      },
    });

    const button = wrapper.find('button.v-btn');

    button.trigger('click');

    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('Renders `c-action-btn` with default edit type correctly.', () => {
    const wrapper = mount(CActionBtn, {
      localVue,
      propsData: { type: 'edit' },
    });

    const tooltipContent = wrapper.find('.v-tooltip__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `c-action-btn` with default duplicate type correctly.', () => {
    const wrapper = mount(CActionBtn, {
      localVue,
      propsData: { type: 'duplicate' },
    });

    const tooltipContent = wrapper.find('.v-tooltip__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `c-action-btn` with default delete type correctly.', () => {
    const wrapper = mount(CActionBtn, {
      localVue,
      propsData: { type: 'delete' },
    });

    const tooltipContent = wrapper.find('.v-tooltip__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `c-action-btn` with custom type correctly.', () => {
    const wrapper = mount(CActionBtn, {
      localVue,
      propsData: { icon: 'test_icon', color: 'color', tooltip: 'tooltip' },
    });

    const tooltipContent = wrapper.find('.v-tooltip__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});
