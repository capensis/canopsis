import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { PATTERN_FIELD_TYPES } from '@/constants';

import CInputTypeField from '@/components/forms/fields/c-input-type-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectTextField = wrapper => wrapper.find('select.v-select');

describe('c-input-type-field', () => {
  const factory = generateShallowRenderer(CInputTypeField, { stubs });
  const snapshotFactory = generateRenderer(CInputTypeField);

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

  it('Renders `c-input-type-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-input-type-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: PATTERN_FIELD_TYPES.number,
        label: 'Custom label',
        name: 'name',
        disabled: true,
        flat: true,
        errorMessages: ['Message'],
        types: [
          { value: PATTERN_FIELD_TYPES.number },
          { value: PATTERN_FIELD_TYPES.string },
          { value: PATTERN_FIELD_TYPES.stringArray },
          { value: PATTERN_FIELD_TYPES.null },
          { text: 'Custom boolean', value: PATTERN_FIELD_TYPES.boolean },
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
