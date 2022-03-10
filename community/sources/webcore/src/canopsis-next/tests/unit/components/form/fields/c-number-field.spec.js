import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createNumberInputStub } from '@unit/stubs/input';

import CNumberField from '@/components/forms/fields/c-number-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-text-field': createNumberInputStub('v-text-field'),
};

const factory = (options = {}) => shallowMount(CNumberField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CNumberField, {
  localVue,

  ...options,
});

const selectTextField = wrapper => wrapper.find('input.v-text-field');

describe('c-number-field', () => {
  it('Value changed after trigger the input', () => {
    const wrapper = factory();
    const textField = selectTextField(wrapper);

    const newValue = Faker.datatype.number();

    textField.setValue(newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(newValue);
  });

  it('Renders `c-number-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-number-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 21,
        label: 'Custom label',
        name: 'customName',
        disabled: true,
        required: true,
        min: 19,
        max: 21,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
