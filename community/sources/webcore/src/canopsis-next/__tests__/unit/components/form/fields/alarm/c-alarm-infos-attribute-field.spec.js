import { generateShallowRenderer } from '@unit/utils/vue';

import CAlarmInfosAttributeField from '@/components/forms/fields/alarm/c-alarm-infos-attribute-field.vue';

const comboboxStub = {
  name: 'CComboboxStub',
  props: {
    modelValue: {
      type: [Object, String],
      default: undefined,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  template: '<div />',
};

const stubs = {
  'c-combobox': comboboxStub,
};

const selectComboboxes = wrapper => wrapper.findAllComponents({ name: 'CComboboxStub' });
const selectRuleCombobox = wrapper => selectComboboxes(wrapper)[0];
const selectDictionaryCombobox = wrapper => selectComboboxes(wrapper)[1];

describe('c-alarm-infos-attribute-field', () => {
  const factory = generateShallowRenderer(CAlarmInfosAttributeField, { stubs });

  test('Rule and dictionary are required when the field is required', () => {
    const wrapper = factory({
      props: {
        value: {},
        required: true,
      },
    });

    expect(selectRuleCombobox(wrapper).props('required')).toBe(true);
    expect(selectDictionaryCombobox(wrapper).props('required')).toBe(true);
  });

  test('Rule is required when an optional dictionary is filled', () => {
    const wrapper = factory({
      props: {
        value: { dictionary: 'info-name' },
      },
    });

    expect(selectRuleCombobox(wrapper).props('required')).toBe(true);
    expect(selectDictionaryCombobox(wrapper).props('required')).toBe(false);
  });

  test('Dictionary is required when an optional rule is selected', () => {
    const wrapper = factory({
      props: {
        value: { rule: 'rule-id' },
        rules: [{ _id: 'rule-id', name: 'Rule' }],
      },
    });

    expect(selectRuleCombobox(wrapper).props('required')).toBe(false);
    expect(selectDictionaryCombobox(wrapper).props('required')).toBe(true);
  });

  test('Rule update is emitted through the value model', async () => {
    const wrapper = factory({
      props: {
        value: { dictionary: 'info-name' },
        rules: [{ _id: 'rule-id', name: 'Rule' }],
      },
    });

    await selectRuleCombobox(wrapper).triggerCustomEvent('update:modelValue', {
      _id: 'rule-id',
      name: 'Rule',
    });

    expect(wrapper).toEmit('update:value', {
      dictionary: 'info-name',
      rule: 'rule-id',
    });
  });

  test('Dictionary update is emitted through the value model', async () => {
    const wrapper = factory({
      props: {
        value: { rule: 'rule-id' },
      },
    });

    await selectDictionaryCombobox(wrapper).triggerCustomEvent('update:modelValue', 'info-name');

    expect(wrapper).toEmit('update:value', {
      rule: 'rule-id',
      dictionary: 'info-name',
    });
  });
});
