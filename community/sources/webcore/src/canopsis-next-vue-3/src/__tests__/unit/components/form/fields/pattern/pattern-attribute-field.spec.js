import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_PATTERN_FIELDS } from '@/constants';

import PatternAttributeField from '@/components/forms/fields/pattern/pattern-attribute-field.vue';

const stubs = {
  'c-autocomplete': createSelectInputStub('c-autocomplete'),
  'c-simple-tooltip': true,
};

const selectSelectField = wrapper => wrapper.find('select.c-autocomplete');

describe('pattern-attribute-field', () => {
  const factory = generateShallowRenderer(PatternAttributeField, { stubs });

  it('Value changed after trigger the input', () => {
    const value = {
      value: ALARM_PATTERN_FIELDS.ack,
      text: 'Text',
    };
    const wrapper = factory({
      propsData: {
        value,
        items: [value],
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('update:modelValue', value.value);

    expect(wrapper).toEmitInput(value.value);
  });
});
