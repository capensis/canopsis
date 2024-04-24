import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { PATTERN_FIELD_TYPES, PATTERN_OPERATORS, QUICK_RANGES, TIME_UNITS } from '@/constants';

import PatternRulesField from '@/components/forms/fields/pattern/pattern-rules-field.vue';

const stubs = {
  'pattern-rule-field': true,
  'c-action-btn': true,
  'c-btn-with-error': true,
};

const selectAddButton = wrapper => wrapper.find('c-btn-with-error-stub');
const selectPatternRulesField = wrapper => wrapper.findAll('pattern-rule-field-stub');
const selectPatternRuleFieldByIndex = (wrapper, index) => selectPatternRulesField(wrapper)
  .at(index);
const selectPatternRemoveRuleButtonByIndex = (wrapper, index) => wrapper.findAll('c-action-btn-stub')
  .at(index);

describe('pattern-rules-field', () => {
  const rules = [
    {
      attribute: 'attribute 1',
      operator: PATTERN_OPERATORS.equal,
      value: 'attribute value',
      fieldType: PATTERN_FIELD_TYPES.string,
      key: 'key 1',
    },
    {
      attribute: 'attribute 2',
      operator: PATTERN_OPERATORS.notEqual,
      value: 'attribute value 2',
      fieldType: PATTERN_FIELD_TYPES.string,
      key: 'key 2',
    },
    {
      attribute: 'attribute 3',
      operator: PATTERN_OPERATORS.contains,
      value: 'attribute contains',
      fieldType: PATTERN_FIELD_TYPES.string,
      key: 'key 3',
    },
  ];

  const factory = generateShallowRenderer(PatternRulesField, { stubs });
  const snapshotFactory = generateRenderer(PatternRulesField, { stubs });

  test('Rule removed after trigger remove event on the pattern Rule field', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    const removeSecondRuleButton = selectPatternRemoveRuleButtonByIndex(wrapper, 1);

    removeSecondRuleButton.triggerCustomEvent('click');

    expect(wrapper).toEmitInput([
      rules[0],
      rules[2],
    ]);
  });

  test('Rule updated after trigger update event on the pattern Rule field', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    const lastRule = selectPatternRuleFieldByIndex(wrapper, 2);

    const updatedRule = {
      attribute: 'new attribute',
      operator: '',
      value: '',
      key: 'new key',
    };

    lastRule.triggerCustomEvent('input', updatedRule);

    expect(wrapper).toEmitInput([
      rules[0],
      rules[1],
      {
        ...updatedRule,
        dictionary: '',
        field: '',
        value: undefined,
      },
    ]);
  });

  test('Operator and value cleared after update rule with new attribute', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    const lastRule = selectPatternRuleFieldByIndex(wrapper, 1);

    const updatedRule = {
      ...rules[1],
      attribute: 'new attribute',
    };

    lastRule.triggerCustomEvent('input', updatedRule);

    expect(wrapper).toEmitInput([
      rules[0],
      {
        ...updatedRule,
        dictionary: '',
        field: '',
        operator: '',
        value: undefined,
      },
      rules[2],
    ]);
  });

  test('Value changed to array after update rule with array operator', () => {
    const attribute = {
      text: 'Attribute text',
      value: rules[1].attribute,
      options: {
        operators: [PATTERN_OPERATORS.notEqual, PATTERN_OPERATORS.hasNot],
      },
    };

    const wrapper = factory({
      propsData: {
        rules,
        attributes: [attribute],
      },
    });

    const lastRule = selectPatternRuleFieldByIndex(wrapper, 1);

    const updatedRule = {
      ...rules[1],
      operator: PATTERN_OPERATORS.hasNot,
    };

    lastRule.triggerCustomEvent('input', updatedRule);

    expect(wrapper).toEmitInput([
      rules[0],
      {
        ...updatedRule,
        value: [updatedRule.value],
      },
      rules[2],
    ]);
  });

  test('Rule added after click on the add button', () => {
    const attribute = {
      value: 'test',
    };
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [
          attribute,
        ],
      },
    });

    selectAddButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toEmitInput([
      ...rules,
      {
        attribute: attribute.value,
        dictionary: '',
        field: '',
        fieldType: PATTERN_FIELD_TYPES.string,
        operator: '',
        value: '',
        range: {
          type: QUICK_RANGES.last1Hour.value,
          from: 0,
          to: 0,
        },
        key: expect.any(String),
        duration: {
          unit: TIME_UNITS.second,
          value: 1,
        },
      },
    ]);
  });

  test('Renders `pattern-rules-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rules: [],
        attributes: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-rules-field` with custom props', () => {
    const attribute = {
      text: 'Attribute text',
      value: 'attribute value',
      options: {
        operators: [PATTERN_OPERATORS.notEqual],
        customProp: 'customPropValue',
      },
    };
    const wrapper = snapshotFactory({
      propsData: {
        rules: [
          ...rules,
          {
            attribute: attribute.value,
            fieldType: PATTERN_FIELD_TYPES.string,
          },
        ],
        attributes: [
          attribute,
        ],
        required: true,
        disabled: true,
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
