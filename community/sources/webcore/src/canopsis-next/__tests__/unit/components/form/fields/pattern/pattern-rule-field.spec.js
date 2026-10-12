import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import {
  ALARM_PATTERN_FIELDS,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_RULE_TYPES,
  QUICK_RANGES,
  TIME_UNITS,
} from '@/constants';

import PatternRuleField from '@/components/forms/fields/pattern/pattern-rule-field.vue';

const stubs = {
  PatternAttributeField: {
    name: 'PatternAttributeField',
    props: {
      returnObject: {
        type: Boolean,
        default: false,
      },
    },
    template: '<div class="pattern-attribute-field-stub" />',
  },
  'pattern-attribute-field': true,
  'c-infos-attribute-field': true,
  'c-quick-date-interval-type-field': true,
  'c-date-time-interval-field': true,
  'c-input-type-field': true,
  'pattern-operator-field': true,
  'pattern-rule-field-date-value': true,
  'c-mixed-input-field': true,
  'c-duration-field': true,
  'custom-component': true,
  'c-alert': true,
  'v-row': {
    template: '<div><slot /></div>',
  },
  'v-col': {
    template: '<div><slot /></div>',
  },
};

const selectPatternAttributeField = wrapper =>
  wrapper.findComponent({ name: 'PatternAttributeField' });
const selectInfosAttributeField = wrapper => wrapper.find('c-infos-attribute-field-stub');
const selectInputTypeField = wrapper => wrapper.find('c-input-type-field-stub');

describe('pattern-rule-field', () => {
  const emptyRule = {
    attribute: '',
    operator: '',
    value: '',
    field: '',
    fieldType: PATTERN_FIELD_TYPES.string,
    dictionary: '',
    range: {
      type: '',
      from: 0,
      to: 0,
    },
    duration: {
      unit: TIME_UNITS.second,
      value: 1,
    },
  };

  const factory = generateShallowRenderer(PatternRuleField, { stubs });

  test('Attribute changed after trigger input event on the attribute field', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    wrapper.vm.updateAttribute({ value: ALARM_PATTERN_FIELDS.displayName });

    expect(wrapper).toEmitInput({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.displayName,
      alias: false,
    });
  });

  test('Pattern attribute field has return-object prop set to true', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    const patternAttributeField = selectPatternAttributeField(wrapper);

    expect(patternAttributeField.props('returnObject')).toBe(true);
  });

  test('Pattern attribute field emits input when updateAttribute is triggered', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    const attributeData = { value: ALARM_PATTERN_FIELDS.displayName };
    wrapper.vm.updateAttribute(attributeData);

    expect(wrapper).toEmitInput({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.displayName,
      alias: false,
    });
  });

  test('Attribute with alias changed after trigger input event on the attribute field', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    wrapper.vm.updateAttribute({
      value: ALARM_PATTERN_FIELDS.displayName,
      alias: true,
    });

    expect(wrapper).toEmitInput({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.displayName,
      alias: true,
      fieldType: PATTERN_FIELD_TYPES.string,
    });
  });

  test('Attribute with definedType sets fieldType on updateAttribute', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    wrapper.vm.updateAttribute({
      value: ALARM_PATTERN_FIELDS.output,
      definedType: PATTERN_FIELD_TYPES.boolean,
    });

    expect(wrapper).toEmitInput({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.output,
      alias: false,
      fieldType: PATTERN_FIELD_TYPES.boolean,
    });
  });

  test('Attribute without definedType resets fieldType to string on updateAttribute', () => {
    const rule = {
      ...emptyRule,
      fieldType: PATTERN_FIELD_TYPES.number,
    };
    const wrapper = factory({
      propsData: {
        rule,
      },
    });

    wrapper.vm.updateAttribute({
      value: ALARM_PATTERN_FIELDS.displayName,
    });

    expect(wrapper).toEmitInput({
      ...rule,
      attribute: ALARM_PATTERN_FIELDS.displayName,
      alias: false,
      fieldType: PATTERN_FIELD_TYPES.string,
    });
  });

  test('Operator changed after trigger input event on the operator field', () => {
    const rule = {
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.output,
    };
    const wrapper = factory({
      propsData: {
        rule,
      },
    });

    wrapper.vm.updateField('operator', PATTERN_OPERATORS.beginWith);

    expect(wrapper).toEmitInput({
      ...rule,
      operator: PATTERN_OPERATORS.beginWith,
    });
  });

  test('Value changed after trigger input event on the value field', () => {
    const rule = {
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.output,
      operator: PATTERN_OPERATORS.equal,
    };
    const wrapper = factory({
      propsData: {
        rule,
      },
    });

    const value = Faker.datatype.string();

    wrapper.vm.updateValue(value);

    expect(wrapper).toEmitInput({
      ...rule,
      value,
    });
  });

  test.each([
    [true, 'common.true'],
    [false, 'common.false'],
  ])('Boolean value field displays the %s label', (value, label) => {
    const wrapper = factory({
      propsData: {
        rule: {
          ...emptyRule,
          attribute: ALARM_PATTERN_FIELDS.output,
          operator: PATTERN_OPERATORS.equal,
          fieldType: PATTERN_FIELD_TYPES.boolean,
          value,
        },
      },
    });

    expect(wrapper.vm.valueComponent).toEqual({
      is: 'c-enabled-field',
      props: expect.objectContaining({
        modelValue: value,
        label,
      }),
    });
  });

  test('Field and dictionary changed after trigger input event on the infos attribute field', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
        type: PATTERN_RULE_TYPES.infos,
      },
    });

    const field = Faker.datatype.string();
    const dictionary = Faker.datatype.string();

    const patternInfosAttributeField = selectInfosAttributeField(wrapper);

    patternInfosAttributeField.triggerCustomEvent('input', {
      ...emptyRule,
      field,
      dictionary,
    });

    expect(wrapper).toEmitInput({
      ...emptyRule,
      field,
      dictionary,
    });
  });

  test('Field changed after trigger input event on the extra infos attribute field', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
        type: PATTERN_RULE_TYPES.extraInfos,
      },
    });

    const field = Faker.datatype.string();

    const patternExtraInfosAttributeField = selectInfosAttributeField(wrapper);

    patternExtraInfosAttributeField.triggerCustomEvent('input', {
      ...emptyRule,
      field,
    });

    expect(wrapper).toEmitInput({
      ...emptyRule,
      field,
    });
  });

  test('Infos value type field has a dedicated validation name', () => {
    const name = 'rule-key';
    const wrapper = factory({
      propsData: {
        name,
        rule: {
          ...emptyRule,
          attribute: ALARM_PATTERN_FIELDS.infos,
          dictionary: 'dictionary',
          field: PATTERN_RULE_INFOS_FIELDS.value,
        },
        type: PATTERN_RULE_TYPES.infos,
      },
    });

    expect(selectInputTypeField(wrapper).attributes('name')).toBe(`${name}.fieldType`);
  });

  test('Range operator changed after trigger input event on the operator field', () => {
    const rule = {
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.creationDate,
    };
    const wrapper = factory({
      propsData: {
        rule,
        type: PATTERN_RULE_TYPES.date,
      },
    });

    wrapper.vm.updateField('operator', PATTERN_OPERATORS.inRangePeriod);

    expect(wrapper).toEmitInput({
      ...rule,

      operator: PATTERN_OPERATORS.inRangePeriod,
    });
  });

  test('Range type changed after trigger input event on the date value field', async () => {
    const rule = {
      ...emptyRule,
      operator: PATTERN_OPERATORS.within,
      attribute: ALARM_PATTERN_FIELDS.creationDate,
    };
    const wrapper = factory({
      propsData: {
        rule,
        type: PATTERN_RULE_TYPES.date,
      },
    });

    wrapper.vm.updateRange({
      ...rule.range,
      type: QUICK_RANGES.last15Minutes.value,
    });

    expect(wrapper).toEmitInput({
      ...rule,

      operator: PATTERN_OPERATORS.within,
      range: {
        ...rule.range,
        type: QUICK_RANGES.last15Minutes.value,
      },
    });
  });

  test('Duration changed after trigger input event on the date duration field', () => {
    const rule = {
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.duration,
      operator: PATTERN_OPERATORS.higher,
    };
    const wrapper = factory({
      propsData: {
        rule,
        type: PATTERN_RULE_TYPES.duration,
      },
    });

    const duration = {
      unit: TIME_UNITS.hour,
      value: Faker.datatype.number(),
    };

    wrapper.vm.updateDuration(duration);

    expect(wrapper).toEmitInput({
      ...rule,
      duration,
    });
  });

  test('IsAnyInfosRuleOrAlias computed returns true when rule alias is true', () => {
    const wrapper = factory({
      propsData: {
        rule: {
          ...emptyRule,
          alias: true,
        },
      },
    });

    expect(wrapper.vm.isAnyInfosRuleOrAlias).toBe(true);
  });

  test('IsAnyInfosRuleOrAlias computed returns true when type is infos', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
        type: PATTERN_RULE_TYPES.infos,
      },
    });

    expect(wrapper.vm.isAnyInfosRuleOrAlias).toBe(true);
  });

  test('IsAnyInfosRuleOrAlias computed returns true when type is extra infos', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
        type: PATTERN_RULE_TYPES.extraInfos,
      },
    });

    expect(wrapper.vm.isAnyInfosRuleOrAlias).toBe(true);
  });

  test('IsAnyInfosRuleOrAlias computed returns false when rule alias is false and type is not infos', () => {
    const wrapper = factory({
      propsData: {
        rule: {
          ...emptyRule,
          alias: false,
        },
        type: PATTERN_RULE_TYPES.string,
      },
    });

    expect(wrapper.vm.isAnyInfosRuleOrAlias).toBe(false);
  });

  test('isInfosValueField computed returns true when rule field is value', () => {
    const rule = {
      ...emptyRule,
      field: PATTERN_RULE_INFOS_FIELDS.value,
    };
    const wrapper = factory({
      propsData: {
        rule,
      },
    });

    expect(wrapper.vm.isInfosValueField).toBe(true);
  });

  test('isInfosValueField computed returns true when rule alias is true', () => {
    const wrapper = factory({
      propsData: {
        rule: {
          ...emptyRule,
          alias: true,
        },
      },
    });

    expect(wrapper.vm.isInfosValueField).toBe(true);
  });

  test('isInfosValueField computed returns false when rule alias is false and field is not value', () => {
    const rule = {
      ...emptyRule,
      field: PATTERN_RULE_INFOS_FIELDS.name,
    };
    const wrapper = factory({
      propsData: {
        rule: {
          ...rule,
          alias: false,
        },
      },
    });

    expect(wrapper.vm.isInfosValueField).toBe(false);
  });

  test('isDateRule computed returns true when rule fieldType is timestamp', () => {
    const rule = {
      ...emptyRule,
      fieldType: PATTERN_FIELD_TYPES.timestamp,
    };
    const wrapper = factory({
      propsData: {
        rule,
        type: PATTERN_RULE_TYPES.string,
      },
    });

    expect(wrapper.vm.isDateRule).toBe(true);
  });

  test('isDateRule computed returns true when type is date', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
        type: PATTERN_RULE_TYPES.date,
      },
    });

    expect(wrapper.vm.isDateRule).toBe(true);
  });

  test('isDateRule computed returns false when type is not date and fieldType is not timestamp', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
        type: PATTERN_RULE_TYPES.string,
      },
    });

    expect(wrapper.vm.isDateRule).toBe(false);
  });

  test('notDefinedType computed returns true when item definedType differs from rule fieldType', () => {
    const wrapper = factory({
      propsData: {
        rule: {
          ...emptyRule,
          attribute: 'attribute-1',
          field: PATTERN_RULE_INFOS_FIELDS.value,
          fieldType: PATTERN_FIELD_TYPES.number,
        },
        attributes: [
          {
            value: 'attribute-1',
            text: 'Attribute text 1',
            definedType: PATTERN_FIELD_TYPES.string,
          },
        ],
        inputTypes: [
          { value: PATTERN_FIELD_TYPES.string },
          { value: PATTERN_FIELD_TYPES.number },
          { value: PATTERN_FIELD_TYPES.boolean },
        ],
      },
    });

    expect(wrapper.vm.notDefinedType).toBe(true);
  });

  test('notDefinedType computed returns false when item definedType equals rule fieldType', () => {
    const wrapper = factory({
      propsData: {
        rule: {
          ...emptyRule,
          attribute: 'attribute-1',
          field: PATTERN_RULE_INFOS_FIELDS.value,
          fieldType: PATTERN_FIELD_TYPES.string,
        },
        attributes: [
          {
            value: 'attribute-1',
            text: 'Attribute text 1',
            definedType: PATTERN_FIELD_TYPES.string,
          },
        ],
        inputTypes: [
          { value: PATTERN_FIELD_TYPES.string },
          { value: PATTERN_FIELD_TYPES.number },
          { value: PATTERN_FIELD_TYPES.boolean },
        ],
      },
    });

    expect(wrapper.vm.notDefinedType).toBe(false);
  });
});
