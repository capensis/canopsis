import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import FlowchartNumberField from '@/components/common/flowchart/fields/flowchart-number-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(FlowchartNumberField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(FlowchartNumberField, {
  localVue,

  ...options,
});

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('flowchart-number-field', () => {
  test('Value changed after trigger select field', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    const newValue = Faker.datatype.number();

    selectField.vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  test('Renders `flowchart-number-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `flowchart-number-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 12,
        label: 'Custom label',
        min: 11,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
