import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import CActionBtn from '@/components/common/buttons/c-action-btn.vue';

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
  'c-simple-tooltip': true,
};

describe('c-action-btn', () => {
  const factory = generateShallowRenderer(CActionBtn, { stubs });
  const snapshotFactory = generateRenderer(CActionBtn, {
    stubs,
    attachTo: document.body,
  });

  it('Action button with default type and custom icon.', () => {
    const icon = Faker.datatype.string();
    const wrapper = factory({ propsData: { icon, type: 'edit' } });

    const iconElement = wrapper.find('v-icon-stub');

    expect(iconElement.text()).toBe(icon);
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
    const wrapper = snapshotFactory({
      propsData: { type: 'edit' },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-action-btn` with default duplicate type correctly.', () => {
    const wrapper = snapshotFactory({
      propsData: { type: 'duplicate' },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-action-btn` with default delete type correctly.', () => {
    const wrapper = snapshotFactory({
      propsData: { type: 'delete' },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-action-btn` with custom type correctly.', () => {
    const wrapper = snapshotFactory({
      propsData: { icon: 'test_icon', color: 'color', tooltip: 'tooltip' },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-action-btn` with badge.', () => {
    const wrapper = snapshotFactory({
      propsData: {
        type: 'edit',
        tooltip: 'TOOLTIP',
        badgeValue: true,
        badgeTooltip: 'BADGE TOOLTIP',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
