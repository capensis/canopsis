import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import FlowchartNumberField from '@/components/common/flowchart/fields/flowchart-number-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('flowchart-number-field', () => {
  const factory = generateShallowRenderer(FlowchartNumberField, { stubs });
  const snapshotFactory = generateRenderer(FlowchartNumberField);

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
