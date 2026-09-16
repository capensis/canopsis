import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { PATTERN_OPERATORS } from '@/constants';

import PatternOperatorField from '@/components/forms/fields/pattern/pattern-operator-field.vue';

const stubs = {
  'c-select': createSelectInputStub('c-select'),
};

const selectSelectField = wrapper => wrapper.find('select.c-select');

describe('pattern-operator-field', () => {
  const factory = generateShallowRenderer(PatternOperatorField, { stubs });

  test('Value changed after trigger the select', () => {
    const wrapper = factory({
      propsData: {
        value: PATTERN_OPERATORS.notExist,
      },
    });

    selectSelectField(wrapper).triggerCustomEvent('update:modelValue', PATTERN_OPERATORS.equal);

    expect(wrapper).toEmitInput(PATTERN_OPERATORS.equal);
  });
});
