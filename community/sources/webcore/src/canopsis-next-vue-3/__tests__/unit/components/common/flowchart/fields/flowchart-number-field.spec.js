import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import FlowchartNumberField from '@/components/common/flowchart/fields/flowchart-number-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('flowchart-number-field', () => {
  const factory = generateShallowRenderer(FlowchartNumberField, { stubs });

  test('Value changed after trigger select field', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    const newValue = Faker.datatype.number();

    selectField.triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
