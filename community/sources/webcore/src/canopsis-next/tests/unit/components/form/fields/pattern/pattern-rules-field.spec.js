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

  const mockValidator = {
    attach: jest.fn(),
    detach: jest.fn(),
    validate: jest.fn(),
    errors: {
      has: jest.fn(() => false),
    },
  };

  const factory = generateShallowRenderer(PatternRulesField, {
    stubs,
    provide: {
      $validator: mockValidator,
    },
  });
  const snapshotFactory = generateRenderer(PatternRulesField, {
    stubs,
    provide: {
      $validator: mockValidator,
    },
  });

  beforeEach(() => {
    mockValidator.attach.mockClear();
    mockValidator.detach.mockClear();
    mockValidator.validate.mockClear();
    mockValidator.errors.has.mockReturnValue(false);
  });

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

    const expectedValue = [{
      value: updatedRule.value,
      key: expect.any(String),
    }];

    expect(wrapper).toEmitInput([
      rules[0],
      {
        ...updatedRule,
        value: expectedValue,
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
        alias: false,
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

  test('Add button is hidden when readonly prop is true', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
        readonly: true,
      },
    });

    expect(selectAddButton(wrapper).exists()).toBe(false);
  });

  test('Remove buttons are disabled when readonly prop is true', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
        readonly: true,
      },
    });

    const removeButtons = wrapper.findAll('c-action-btn-stub');

    removeButtons.wrappers.forEach((removeButton) => {
      expect(removeButton.attributes().disabled).toBe('true');
    });
  });

  test('Remove buttons are disabled when disabled prop is true', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
        disabled: true,
      },
    });

    const removeButtons = wrapper.findAll('c-action-btn-stub');

    removeButtons.wrappers.forEach((removeButton) => {
      expect(removeButton.attributes().disabled).toBe('true');
    });
  });

  test('Validation rule is attached on component creation', () => {
    factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    expect(mockValidator.attach).toHaveBeenCalledWith({
      name: 'rules',
      rules: {
        is_not: 0,
      },
      getter: expect.any(Function),
      context: expect.any(Function),
      vm: expect.any(Object),
    });
  });

  test('Validation rule is detached on component destruction', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    wrapper.destroy();

    expect(mockValidator.detach).toHaveBeenCalledWith('rules');
  });

  test('Validation getter filters out disabled rules', () => {
    const disabledAttribute = {
      value: 'disabled-attr',
      options: { disabled: true },
    };
    const enabledAttribute = {
      value: 'enabled-attr',
      options: {},
    };
    const rulesWithDisabled = [
      { ...rules[0], attribute: disabledAttribute.value },
      { ...rules[1], attribute: enabledAttribute.value },
    ];

    factory({
      propsData: {
        rules: rulesWithDisabled,
        attributes: [disabledAttribute, enabledAttribute],
      },
    });

    const validationCall = mockValidator.attach.mock.calls[mockValidator.attach.mock.calls.length - 1][0];
    const { getter } = validationCall;
    const count = getter();

    expect(count).toBe(1);
  });

  test('hasRulesErrors returns true when validator has errors for the field', () => {
    const { errors } = mockValidator;

    errors.has.mockReturnValue(true);

    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    expect(wrapper.vm.hasRulesErrors).toBe(true);
  });

  test('Error message is displayed on add button when hasRulesErrors is true', () => {
    mockValidator.errors.has.mockReturnValue(true);

    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    const addButton = selectAddButton(wrapper);
    expect(addButton.attributes().error).toBeTruthy();
  });

  test('Operator is automatically set when only one operator is available', () => {
    const singleOperatorAttribute = {
      value: 'single-op-attr',
      options: {
        operators: [PATTERN_OPERATORS.equal],
      },
    };
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [singleOperatorAttribute],
      },
    });

    const updatedRule = {
      ...rules[0],
      attribute: singleOperatorAttribute.value,
      value: 'new value',
    };

    const ruleField = selectPatternRuleFieldByIndex(wrapper, 0);
    ruleField.triggerCustomEvent('input', updatedRule);

    expect(wrapper).toEmitInput([
      {
        ...updatedRule,
        operator: PATTERN_OPERATORS.equal,
        dictionary: '',
        field: '',
        value: undefined,
      },
      rules[1],
      rules[2],
    ]);
  });

  test('getUpdatedRule method clears operator when not in available operators', () => {
    const limitedOperatorAttribute = {
      value: 'limited-op-attr',
      options: {
        operators: [PATTERN_OPERATORS.equal],
      },
    };
    const wrapper = factory({
      propsData: {
        rules: [{ ...rules[0], attribute: limitedOperatorAttribute.value }],
        attributes: [limitedOperatorAttribute],
      },
    });

    const originalRule = { ...rules[0], attribute: limitedOperatorAttribute.value };
    const newRule = { ...originalRule, operator: PATTERN_OPERATORS.notEqual };

    const updatedRule = wrapper.vm.getUpdatedRule(originalRule, newRule);

    expect(updatedRule.operator).toBe('');
  });

  test('Rule props are calculated correctly with attribute options', () => {
    const attributeWithOptions = {
      value: 'attr-with-options',
      options: {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        type: 'custom-type',
        disabled: true,
      },
    };
    const wrapper = factory({
      propsData: {
        rules: [{ ...rules[0], attribute: attributeWithOptions.value }],
        attributes: [attributeWithOptions],
      },
    });

    const ruleField = selectPatternRuleFieldByIndex(wrapper, 0);
    const ruleProps = ruleField.props();

    expect(ruleProps.operators).toEqual(attributeWithOptions.options.operators);
    expect(ruleProps.type).toBe(attributeWithOptions.options.type);
    expect(ruleProps.disabled).toBe(true);
  });

  test('rulesMap computed property maps attributes correctly', () => {
    const attributesWithOptions = [
      {
        value: 'attr1',
        options: { disabled: true, type: 'string' },
      },
      {
        value: 'attr2',
        options: { operators: [PATTERN_OPERATORS.equal] },
      },
      {
        value: 'attr3',
      },
    ];
    const wrapper = factory({
      propsData: {
        rules: [],
        attributes: attributesWithOptions,
      },
    });

    const expectedRulesMap = {
      attr1: { disabled: true, type: 'string' },
      attr2: { operators: [PATTERN_OPERATORS.equal] },
      attr3: {},
    };

    expect(wrapper.vm.rulesMap).toEqual(expectedRulesMap);
  });

  test('Default value is applied when changing attribute', () => {
    const attributeWithDefault = {
      value: 'attr-with-default',
      options: {
        defaultValue: 'default-value',
      },
    };
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [attributeWithDefault],
      },
    });

    const updatedRule = {
      ...rules[0],
      attribute: attributeWithDefault.value,
    };

    const ruleField = selectPatternRuleFieldByIndex(wrapper, 0);
    ruleField.triggerCustomEvent('input', updatedRule);

    expect(wrapper).toEmitInput([
      {
        ...updatedRule,
        operator: '',
        field: '',
        dictionary: '',
        value: 'default-value',
      },
      rules[1],
      rules[2],
    ]);
  });

  test('Custom name prop is used in validation', () => {
    const customName = 'customRulesName';

    factory({
      propsData: {
        rules,
        attributes: [],
        name: customName,
      },
    });

    expect(mockValidator.attach).toHaveBeenCalledWith(
      expect.objectContaining({
        name: customName,
      }),
    );
  });

  test('Pattern rule field receives correct props from getRuleProps', () => {
    const attributeWithCustomProps = {
      value: 'attr-custom',
      options: {
        type: 'custom-type',
        disabled: false,
        operators: [PATTERN_OPERATORS.equal],
      },
    };
    const wrapper = factory({
      propsData: {
        rules: [{ ...rules[0], attribute: attributeWithCustomProps.value }],
        attributes: [attributeWithCustomProps],
        readonly: false,
        disabled: false,
      },
    });

    const ruleField = selectPatternRuleFieldByIndex(wrapper, 0);
    const props = ruleField.props();

    expect(props.type).toBe('custom-type');
    expect(props.disabled).toBe(false);
    expect(props.operators).toEqual([PATTERN_OPERATORS.equal]);
  });

  test('Disabled state is propagated correctly through readonly and disabled props', () => {
    const attribute = { value: 'test-attr', options: {} };

    const readonlyWrapper = factory({
      propsData: {
        rules: [{ ...rules[0], attribute: attribute.value }],
        attributes: [attribute],
        readonly: true,
        disabled: false,
      },
    });

    const disabledWrapper = factory({
      propsData: {
        rules: [{ ...rules[0], attribute: attribute.value }],
        attributes: [attribute],
        readonly: false,
        disabled: true,
      },
    });

    const bothWrapper = factory({
      propsData: {
        rules: [{ ...rules[0], attribute: attribute.value }],
        attributes: [attribute],
        readonly: true,
        disabled: true,
      },
    });

    expect(selectPatternRuleFieldByIndex(readonlyWrapper, 0).props().disabled).toBe(true);
    expect(selectPatternRuleFieldByIndex(disabledWrapper, 0).props().disabled).toBe(true);
    expect(selectPatternRuleFieldByIndex(bothWrapper, 0).props().disabled).toBe(true);
  });

  test('Add button uses first attribute when adding new rule', () => {
    const firstAttribute = { value: 'first-attr' };
    const secondAttribute = { value: 'second-attr' };

    const wrapper = factory({
      propsData: {
        rules: [],
        attributes: [firstAttribute, secondAttribute],
      },
    });

    selectAddButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toEmitInput([
      expect.objectContaining({
        attribute: firstAttribute.value,
      }),
    ]);
  });

  test('Value is converted when operator changes', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    const updatedRule = {
      ...rules[0],
      operator: PATTERN_OPERATORS.hasNot, // Array operator
    };

    const ruleField = selectPatternRuleFieldByIndex(wrapper, 0);
    ruleField.triggerCustomEvent('input', updatedRule);

    expect(wrapper).toEmitInput([
      expect.objectContaining({
        value: [rules[0].value], // Value converted to array
      }),
      rules[1],
      rules[2],
    ]);
  });

  test('beforeDestroy hook calls detachMinValueRule', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    mockValidator.detach.mockClear();
    wrapper.destroy();

    expect(mockValidator.detach).toHaveBeenCalledWith('rules');
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
