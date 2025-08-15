import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import {
  ALARM_PATTERN_FIELDS,
  EVENT_FILTER_PATTERN_FIELDS,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_RULE_TYPES,
  QUICK_RANGES,
  TIME_UNITS,
} from '@/constants';

import PatternRuleField from '@/components/forms/fields/pattern/pattern-rule-field.vue';

const stubs = {
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
};

const selectPatternAttributeField = wrapper => wrapper.find('pattern-attribute-field-stub');
const selectPatternOperatorField = wrapper => wrapper.find('pattern-operator-field-stub');
const selectPatternDateValueField = wrapper => wrapper.find('pattern-rule-field-date-value-stub');
const selectMixedInputField = wrapper => wrapper.find('c-mixed-input-field-stub');
const selectInfosAttributeField = wrapper => wrapper.find('c-infos-attribute-field-stub');
const selectInputTypeField = wrapper => wrapper.find('c-input-type-field-stub');
const selectDurationField = wrapper => wrapper.find('c-duration-field-stub');

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
  const snapshotFactory = generateRenderer(PatternRuleField, { stubs });

  test('Attribute changed after trigger input event on the attribute field', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    const patternAttributeField = selectPatternAttributeField(wrapper);

    patternAttributeField.triggerCustomEvent('input', { value: ALARM_PATTERN_FIELDS.displayName });

    expect(wrapper).toEmitInput({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.displayName,
      alias: undefined,
      definedType: 'string',
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

    const patternAttributeField = selectPatternAttributeField(wrapper);

    const attributeData = { value: ALARM_PATTERN_FIELDS.displayName };
    patternAttributeField.triggerCustomEvent('input', attributeData);

    expect(wrapper).toEmitInput({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.displayName,
      alias: undefined,
      definedType: PATTERN_FIELD_TYPES.string,
    });
  });

  test('Attribute with alias changed after trigger input event on the attribute field', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    const patternAttributeField = selectPatternAttributeField(wrapper);

    patternAttributeField.triggerCustomEvent('input', {
      value: ALARM_PATTERN_FIELDS.displayName,
      options: {
        alias: true,
        definedType: PATTERN_FIELD_TYPES.number,
      },
    });

    expect(wrapper).toEmitInput({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.displayName,
      alias: true,
      definedType: PATTERN_FIELD_TYPES.number,
      fieldType: PATTERN_FIELD_TYPES.number,
    });
  });

  test('Attribute with defined type but no alias after trigger input event on the attribute field', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    const patternAttributeField = selectPatternAttributeField(wrapper);

    patternAttributeField.triggerCustomEvent('input', {
      value: ALARM_PATTERN_FIELDS.output,
      options: {
        definedType: PATTERN_FIELD_TYPES.boolean,
      },
    });

    expect(wrapper).toEmitInput({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.output,
      alias: undefined,
      definedType: PATTERN_FIELD_TYPES.boolean,
      fieldType: PATTERN_FIELD_TYPES.boolean,
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

    const patternOperatorField = selectPatternOperatorField(wrapper);

    patternOperatorField.triggerCustomEvent('input', PATTERN_OPERATORS.beginsWith);

    expect(wrapper).toEmitInput({
      ...rule,
      operator: PATTERN_OPERATORS.beginsWith,
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

    const mixedInputField = selectMixedInputField(wrapper);

    mixedInputField.triggerCustomEvent('input', value);

    expect(wrapper).toEmitInput({
      ...rule,
      value,
    });
  });

  test('Value to string changed after trigger input event on the value type field', () => {
    const rule = {
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.output,
      operator: PATTERN_OPERATORS.equal,
      value: Faker.datatype.number(),
      field: PATTERN_RULE_INFOS_FIELDS.value,
      fieldType: PATTERN_FIELD_TYPES.string,
    };
    const wrapper = factory({
      propsData: {
        rule,
        type: PATTERN_RULE_TYPES.infos,
      },
    });

    const inputTypeField = selectInputTypeField(wrapper);

    inputTypeField.triggerCustomEvent('input', PATTERN_FIELD_TYPES.string);

    expect(wrapper).toEmitInput({
      ...rule,
      value: `${rule.value}`,
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

    const operatorField = selectPatternOperatorField(wrapper);

    operatorField.triggerCustomEvent('input', PATTERN_OPERATORS.inRangePeriod);

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

    const dateValueField = selectPatternDateValueField(wrapper);

    dateValueField.triggerCustomEvent('input', {
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

    const durationField = selectDurationField(wrapper);

    const duration = {
      unit: TIME_UNITS.hour,
      value: Faker.datatype.number(),
    };

    durationField.triggerCustomEvent('input', duration);

    expect(wrapper).toEmitInput({
      ...rule,
      duration,
    });
  });

  test('IsAnyInfosRuleOrAlias computed returns true when alias prop is true', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
        alias: true,
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

  test('IsAnyInfosRuleOrAlias computed returns false when alias is false and type is not infos', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
        alias: false,
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

  test('isInfosValueField computed returns true when alias prop is true', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
        alias: true,
      },
    });

    expect(wrapper.vm.isInfosValueField).toBe(true);
  });

  test('isInfosValueField computed returns false when alias is false and field is not value', () => {
    const rule = {
      ...emptyRule,
      field: PATTERN_RULE_INFOS_FIELDS.name,
    };
    const wrapper = factory({
      propsData: {
        rule,
        alias: false,
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

  test('notDefinedType computed returns true when definedType differs from fieldType', () => {
    const rule = {
      ...emptyRule,
      fieldType: PATTERN_FIELD_TYPES.string,
      definedType: PATTERN_FIELD_TYPES.number,
    };
    const wrapper = factory({
      propsData: {
        rule,
      },
    });

    expect(wrapper.vm.notDefinedType).toBe(true);
  });

  test('notDefinedType computed returns false when definedType equals fieldType', () => {
    const rule = {
      ...emptyRule,
      fieldType: PATTERN_FIELD_TYPES.string,
      definedType: PATTERN_FIELD_TYPES.string,
    };
    const wrapper = factory({
      propsData: {
        rule,
      },
    });

    expect(wrapper.vm.notDefinedType).toBe(false);
  });

  test('inputTypesWithDefinedType computed adds defined flag to matching type', () => {
    const rule = {
      ...emptyRule,
      definedType: PATTERN_FIELD_TYPES.number,
    };
    const inputTypes = [
      { value: PATTERN_FIELD_TYPES.string },
      { value: PATTERN_FIELD_TYPES.number },
      { value: PATTERN_FIELD_TYPES.boolean },
    ];
    const wrapper = factory({
      propsData: {
        rule,
        inputTypes,
      },
    });

    const result = wrapper.vm.inputTypesWithDefinedType;

    expect(result[0]).toEqual({
      value: PATTERN_FIELD_TYPES.string,
      defined: false,
    });
    expect(result[1]).toEqual({
      value: PATTERN_FIELD_TYPES.number,
      defined: true,
    });
    expect(result[2]).toEqual({
      value: PATTERN_FIELD_TYPES.boolean,
      defined: false,
    });
  });

  test('inputTypesWithDefinedType computed sets all defined flags to false when no matching type', () => {
    const rule = {
      ...emptyRule,
      definedType: PATTERN_FIELD_TYPES.stringArray,
    };
    const inputTypes = [
      { value: PATTERN_FIELD_TYPES.string },
      { value: PATTERN_FIELD_TYPES.number },
    ];
    const wrapper = factory({
      propsData: {
        rule,
        inputTypes,
      },
    });

    const result = wrapper.vm.inputTypesWithDefinedType;

    expect(result[0]).toEqual({
      value: PATTERN_FIELD_TYPES.string,
      defined: false,
    });
    expect(result[1]).toEqual({
      value: PATTERN_FIELD_TYPES.number,
      defined: false,
    });
  });

  test('updateAttribute method sets alias and definedType from options', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    const updateModelSpy = jest.spyOn(wrapper.vm, 'updateModel');

    wrapper.vm.updateAttribute({
      value: ALARM_PATTERN_FIELDS.displayName,
      options: {
        alias: true,
        definedType: PATTERN_FIELD_TYPES.number,
      },
    });

    expect(updateModelSpy).toHaveBeenCalledWith({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.displayName,
      alias: true,
      definedType: PATTERN_FIELD_TYPES.number,
      fieldType: PATTERN_FIELD_TYPES.number,
    });
  });

  test('updateAttribute method sets fieldType when definedType is provided', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    const updateModelSpy = jest.spyOn(wrapper.vm, 'updateModel');

    wrapper.vm.updateAttribute({
      value: ALARM_PATTERN_FIELDS.output,
      options: {
        definedType: PATTERN_FIELD_TYPES.boolean,
      },
    });

    expect(updateModelSpy).toHaveBeenCalledWith({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.output,
      alias: undefined,
      definedType: PATTERN_FIELD_TYPES.boolean,
      fieldType: PATTERN_FIELD_TYPES.boolean,
    });
  });

  test('updateAttribute method uses string as default definedType when options not provided', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    const updateModelSpy = jest.spyOn(wrapper.vm, 'updateModel');

    wrapper.vm.updateAttribute({
      value: ALARM_PATTERN_FIELDS.output,
    });

    expect(updateModelSpy).toHaveBeenCalledWith({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.output,
      alias: undefined,
      definedType: PATTERN_FIELD_TYPES.string,
    });
  });

  test('updateAttribute method does not set fieldType when definedType is not provided', () => {
    const wrapper = factory({
      propsData: {
        rule: emptyRule,
      },
    });

    const updateModelSpy = jest.spyOn(wrapper.vm, 'updateModel');

    wrapper.vm.updateAttribute({
      value: ALARM_PATTERN_FIELDS.output,
      options: {
        alias: true,
      },
    });

    expect(updateModelSpy).toHaveBeenCalledWith({
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.output,
      alias: true,
      definedType: PATTERN_FIELD_TYPES.string,
    });
  });

  test('Renders `pattern-rule-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: ALARM_PATTERN_FIELDS.displayName,
          operator: PATTERN_OPERATORS.equal,
          value: 'ruleValue',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-rule-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: ALARM_PATTERN_FIELDS.displayName,
          operator: PATTERN_OPERATORS.equal,
          value: 'ruleValue',
          field: PATTERN_RULE_INFOS_FIELDS.value,
        },
        attributes: [
          { value: 'attribute-1', text: 'Attribute text 1' },
        ],
        infos: ['Infos 1', 'Infos 2'],
        operators: [PATTERN_OPERATORS.notEqual, PATTERN_OPERATORS.equal],
        inputTypes: [
          { value: PATTERN_FIELD_TYPES.string },
          { value: PATTERN_FIELD_TYPES.stringArray },
        ],
        valueField: {
          is: 'custom-component',
          props: {
            name: 'test name',
          },
        },
        type: PATTERN_RULE_TYPES.infos,
        disabled: true,
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-rule-field` with infos type and field is name', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: ALARM_PATTERN_FIELDS.infos,
          operator: PATTERN_OPERATORS.equal,
          value: 'infos',
          field: PATTERN_RULE_INFOS_FIELDS.name,
        },
        type: PATTERN_RULE_TYPES.infos,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-rule-field` with extra infos type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: EVENT_FILTER_PATTERN_FIELDS.extraInfos,
          operator: PATTERN_OPERATORS.equal,
          value: 22,
          field: 'extra_field.name',
        },
        type: PATTERN_RULE_TYPES.extraInfos,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-rule-field` with duration type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: ALARM_PATTERN_FIELDS.duration,
          operator: PATTERN_OPERATORS.notEqual,
          duration: {
            unit: TIME_UNITS.year,
            value: 1,
          },
        },
        type: PATTERN_RULE_TYPES.duration,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-rule-field` with date type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: ALARM_PATTERN_FIELDS.creationDate,
          operator: PATTERN_OPERATORS.inRangeDates,
          range: {
            type: QUICK_RANGES.last1Hour.value,
            from: new Date(),
            to: new Date(),
          },
        },
        type: PATTERN_RULE_TYPES.date,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-rule-field` with alias prop', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: ALARM_PATTERN_FIELDS.displayName,
          operator: PATTERN_OPERATORS.equal,
          value: 'alias value',
          fieldType: PATTERN_FIELD_TYPES.string,
          definedType: PATTERN_FIELD_TYPES.number,
        },
        alias: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-rule-field` with type mismatch warning', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: ALARM_PATTERN_FIELDS.displayName,
          operator: PATTERN_OPERATORS.equal,
          value: 'test value',
          fieldType: PATTERN_FIELD_TYPES.string,
          definedType: PATTERN_FIELD_TYPES.number,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-rule-field` with timestamp field type as date rule', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: ALARM_PATTERN_FIELDS.displayName,
          operator: PATTERN_OPERATORS.inRangeDates,
          fieldType: PATTERN_FIELD_TYPES.timestamp,
          range: {
            type: QUICK_RANGES.last1Hour.value,
            from: new Date(),
            to: new Date(),
          },
        },
        type: PATTERN_RULE_TYPES.string,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
