import { generateShallowRenderer } from '@unit/utils/vue';

import { PATTERN_FIELD_TYPES } from '@/constants';

import CInfosAttributeField from '@/components/forms/fields/compat/c-infos-attribute-field.vue';

const stubs = {
  'c-combobox': true,
};

const selectDictionarySelect = wrapper => wrapper.find('c-combobox-stub');
const selectFieldSelect = wrapper => wrapper.find('select.c-select');

describe('c-infos-attribute-field', () => {
  const factory = generateShallowRenderer(CInfosAttributeField, { stubs });

  it('Dictionary changed after trigger the dictionary select', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
        combobox: true,
      },
    });

    const newDictionary = 'newDictionary';

    selectDictionarySelect(wrapper).triggerCustomEvent('input', newDictionary);

    expect(wrapper).toEmitInput({
      dictionary: newDictionary,
      field: value.field,
      fieldType: PATTERN_FIELD_TYPES.string,
    });
  });

  it('Field changed after trigger the field select', () => {
    const value = {
      dictionary: 'test',
      field: 'name',
    };
    const wrapper = factory({
      propsData: {
        value,
        combobox: true,
      },
    });

    const newField = 'newField';

    selectFieldSelect(wrapper).triggerCustomEvent('input', newField);

    expect(wrapper).toEmitInput({
      dictionary: value.dictionary,
      field: newField,
      fieldType: PATTERN_FIELD_TYPES.string,
    });
  });
});
