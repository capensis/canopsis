import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { PATTERN_FIELD_TYPES } from '@/constants';

import CInputTypeField from '@/components/forms/fields/compat/c-input-type-field.vue';

const stubs = {
  'c-select': createSelectInputStub('v-select'),
};

const selectTextField = wrapper => wrapper.find('select.v-select');

describe('c-input-type-field', () => {
  const factory = generateShallowRenderer(CInputTypeField, { stubs });

  it('Input type changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: PATTERN_FIELD_TYPES.string,
      },
    });

    const textField = selectTextField(wrapper);

    textField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.number);

    expect(wrapper).toEmitInput(PATTERN_FIELD_TYPES.number);
  });
});
