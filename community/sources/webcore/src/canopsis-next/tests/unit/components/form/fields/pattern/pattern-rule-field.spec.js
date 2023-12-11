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
  'c-mixed-input-field': true,
  'c-duration-field': true,
  'custom-component': true,
};

const selectPatternAttributeField = wrapper => wrapper.find('pattern-attribute-field-stub');
const selectPatternOperatorField = wrapper => wrapper.find('pattern-operator-field-stub');
const selectMixedInputField = wrapper => wrapper.find('c-mixed-input-field-stub');
const selectInfosAttributeField = wrapper => wrapper.find('c-infos-attribute-field-stub');
const selectQuickDateIntervalTypeField = wrapper => wrapper.find('c-quick-date-interval-type-field-stub');
const selectDateTimeIntervalField = wrapper => wrapper.find('c-date-time-interval-field-stub');
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

    patternAttributeField.vm.$emit('input', ALARM_PATTERN_FIELDS.displayName);

    expect(wrapper).toEmit('input', {
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.displayName,
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

    patternOperatorField.vm.$emit('input', PATTERN_OPERATORS.beginsWith);

    expect(wrapper).toEmit('input', {
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

    mixedInputField.vm.$emit('input', value);

    expect(wrapper).toEmit('input', {
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

    inputTypeField.vm.$emit('input', PATTERN_FIELD_TYPES.string);

    expect(wrapper).toEmit('input', {
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

    patternInfosAttributeField.vm.$emit('input', {
      ...emptyRule,
      field,
      dictionary,
    });

    expect(wrapper).toEmit('input', {
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

    patternExtraInfosAttributeField.vm.$emit('input', {
      ...emptyRule,
      field,
    });

    expect(wrapper).toEmit('input', {
      ...emptyRule,
      field,
    });
  });

  test('Range type changed after trigger input event on the date interval type field', () => {
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

    const quickDateIntervalTypeField = selectQuickDateIntervalTypeField(wrapper);

    quickDateIntervalTypeField.vm.$emit('input', QUICK_RANGES.last15Minutes.value);

    expect(wrapper).toEmit('input', {
      ...rule,
      range: {
        ...rule.range,
        type: QUICK_RANGES.last15Minutes.value,
      },
    });
  });

  test('Interval changed after trigger input event on the date interval type field', () => {
    const rule = {
      ...emptyRule,
      attribute: ALARM_PATTERN_FIELDS.creationDate,
      range: {
        ...emptyRule.range,
        type: QUICK_RANGES.custom.value,
      },
    };
    const wrapper = factory({
      propsData: {
        rule,
        type: PATTERN_RULE_TYPES.date,
      },
    });

    const dateTimeIntervalField = selectDateTimeIntervalField(wrapper);

    const newRange = {
      type: rule.range.type,
      from: Faker.datatype.number(),
      to: Faker.datatype.number(),
    };

    dateTimeIntervalField.vm.$emit('input', newRange);

    expect(wrapper).toEmit('input', {
      ...rule,
      range: newRange,
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

    durationField.vm.$emit('input', duration);

    expect(wrapper).toEmit('input', {
      ...rule,
      duration,
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

    expect(wrapper.element).toMatchSnapshot();
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
        intervalRanges: [QUICK_RANGES.last15Minutes, QUICK_RANGES.custom],
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pattern-rule-field` with date type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rule: {
          attribute: ALARM_PATTERN_FIELDS.creationDate,
          operator: PATTERN_OPERATORS.higher,
          range: {
            type: QUICK_RANGES.custom.value,
            from: 1000000,
            to: 2000000,
          },
        },
        type: PATTERN_RULE_TYPES.date,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
