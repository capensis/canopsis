import { createVueInstance, mount, shallowMount } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';

import CStateTypeField from '@/components/forms/fields/c-state-type-field.vue';
import { ENTITIES_STATES } from '@/constants';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CStateTypeField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CStateTypeField, {
  localVue,

  ...options,
});

describe('c-state-type-field', () => {
  it('State type changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: ENTITIES_STATES.ok,
      },
    });

    const valueElement = wrapper.find('select.v-select');

    valueElement.vm.$emit('input', ENTITIES_STATES.critical);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(ENTITIES_STATES.critical);
  });

  it('Renders `c-state-type-field` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-state-type-field` with default custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ENTITIES_STATES.major,
        label: 'Custom label',
        name: 'name',
        disabled: true,
        required: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-state-type-field` with validator error', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        required: true,
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });
});
