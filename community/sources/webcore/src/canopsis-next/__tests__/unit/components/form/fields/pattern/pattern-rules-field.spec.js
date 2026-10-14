import { generateShallowRenderer } from '@unit/utils/vue';

import { PATTERN_FIELD_TYPES, PATTERN_OPERATORS, QUICK_RANGES, TIME_UNITS } from '@/constants';

import PatternRulesField from '@/components/forms/fields/pattern/pattern-rules-field.vue';

const PatternRuleFieldStub = {
  name: 'pattern-rule-field',
  props: {
    rule: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: '',
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledField: {
      type: Boolean,
      default: false,
    },
    disabledOperator: {
      type: Boolean,
      default: false,
    },
    operators: {
      type: Array,
      default: () => [],
    },
    type: {
      type: String,
      default: '',
    },
  },
  emits: ['update:rule'],
  template: '<div class="pattern-rule-field-stub" />',
};

const ActionBtnStub = {
  name: 'c-action-btn',
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
    tooltip: {
      type: String,
      default: '',
    },
    type: {
      type: String,
      default: '',
    },
  },
  emits: ['click'],
  template: '<button class="c-action-btn-stub" :disabled="disabled" @click="$emit(\'click\')" />',
};

const BtnWithErrorStub = {
  name: 'c-btn-with-error',
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
    error: {
      type: String,
      default: '',
    },
  },
  emits: ['click'],
  template: `
    <button
      class="c-btn-with-error-stub"
      :data-error="error"
      :disabled="disabled"
      @click="$emit('click')"
    >
      <slot />
    </button>
  `,
};

const LayoutStub = {
  template: '<div><slot /></div>',
};

const stubs = {
  PatternRuleField: PatternRuleFieldStub,
  CActionBtn: ActionBtnStub,
  CBtnWithError: BtnWithErrorStub,
  VRow: LayoutStub,
  VCol: LayoutStub,
  'pattern-rule-field': PatternRuleFieldStub,
  'c-action-btn': ActionBtnStub,
  'c-btn-with-error': BtnWithErrorStub,
  'v-row': LayoutStub,
  'v-col': LayoutStub,
};

const selectAddButton = wrapper => wrapper.findComponent({ name: 'c-btn-with-error' });
const selectPatternRulesFields = wrapper =>
  wrapper.findAllComponents({ name: 'pattern-rule-field' });
const selectPatternRuleFieldByIndex = (wrapper, index) =>
  selectPatternRulesFields(wrapper).at(index);
const selectPatternRemoveRuleButtons = wrapper =>
  wrapper.findAllComponents({ name: 'c-action-btn' });
const selectPatternRemoveRuleButtonByIndex = (wrapper, index) =>
  selectPatternRemoveRuleButtons(wrapper).at(index);
const lastRulesEmission = wrapper => wrapper.emitted('update:rules').at(-1)[0];

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

  test('renders a field and a remove button for every rule', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    expect(selectPatternRulesFields(wrapper)).toHaveLength(rules.length);
    expect(selectPatternRemoveRuleButtons(wrapper)).toHaveLength(rules.length);
  });

  test('passes rule data and attributes to pattern rule fields', () => {
    const attributes = [{ value: 'attribute 1' }];

    const wrapper = factory({
      propsData: {
        rules,
        attributes,
      },
    });

    const firstRuleField = selectPatternRuleFieldByIndex(wrapper, 0);

    expect(firstRuleField.props()).toEqual(
      expect.objectContaining({
        rule: rules[0],
        name: rules[0].key,
        attributes,
        disabledField: false,
        disabledOperator: false,
      }),
    );
  });

  test('removes selected rule after remove click', async () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    await selectPatternRemoveRuleButtonByIndex(wrapper, 1).trigger('click');

    expect(lastRulesEmission(wrapper)).toEqual([rules[0], rules[2]]);
  });

  test('updates selected rule after child update', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    const updatedRule = {
      ...rules[2],
      attribute: 'new attribute',
      operator: '',
      value: '',
      key: 'new key',
    };

    selectPatternRuleFieldByIndex(wrapper, 2).vm.$emit('update:rule', updatedRule);

    expect(lastRulesEmission(wrapper)).toEqual([
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

  test('clears operator, field, dictionary and value when attribute changes', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    const updatedRule = {
      ...rules[1],
      attribute: 'new attribute',
    };

    selectPatternRuleFieldByIndex(wrapper, 1).vm.$emit('update:rule', updatedRule);

    expect(lastRulesEmission(wrapper)).toEqual([
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

  test('uses attribute default value when attribute changes', () => {
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

    selectPatternRuleFieldByIndex(wrapper, 0).vm.$emit('update:rule', updatedRule);

    expect(lastRulesEmission(wrapper)).toEqual([
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

  test('converts value to array when operator changes to array operator', () => {
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

    const updatedRule = {
      ...rules[1],
      operator: PATTERN_OPERATORS.hasNot,
    };

    selectPatternRuleFieldByIndex(wrapper, 1).vm.$emit('update:rule', updatedRule);

    expect(lastRulesEmission(wrapper)).toEqual([
      rules[0],
      {
        ...updatedRule,
        value: [updatedRule.value],
      },
      rules[2],
    ]);
  });

  test('clears operator when updated operator is not available for rule', () => {
    const attribute = {
      value: rules[0].attribute,
      options: {
        operators: [PATTERN_OPERATORS.equal],
      },
    };
    const wrapper = factory({
      propsData: {
        rules: [rules[0]],
        attributes: [attribute],
      },
    });

    selectPatternRuleFieldByIndex(wrapper, 0).vm.$emit('update:rule', {
      ...rules[0],
      operator: PATTERN_OPERATORS.notEqual,
    });

    expect(lastRulesEmission(wrapper)).toEqual([
      {
        ...rules[0],
        operator: '',
      },
    ]);
  });

  test('automatically selects the only available operator after value changes', () => {
    const attribute = {
      value: rules[0].attribute,
      options: {
        operators: [PATTERN_OPERATORS.equal],
      },
    };
    const wrapper = factory({
      propsData: {
        rules: [rules[0]],
        attributes: [attribute],
      },
    });

    selectPatternRuleFieldByIndex(wrapper, 0).vm.$emit('update:rule', {
      ...rules[0],
      value: 'new value',
    });

    expect(lastRulesEmission(wrapper)).toEqual([
      {
        ...rules[0],
        operator: PATTERN_OPERATORS.equal,
        value: 'new value',
      },
    ]);
  });

  test('adds a rule based on first available attribute', async () => {
    const attribute = {
      value: 'test',
    };
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [attribute],
      },
    });

    await selectAddButton(wrapper).trigger('click');

    expect(lastRulesEmission(wrapper)).toEqual([
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
          from: '',
          to: '',
        },
        key: expect.any(String),
        duration: {
          unit: TIME_UNITS.second,
          value: 1,
        },
      },
    ]);
  });

  test('uses empty attribute when adding without available attributes', async () => {
    const wrapper = factory({
      propsData: {
        rules: [],
        attributes: [],
      },
    });

    await selectAddButton(wrapper).trigger('click');

    expect(lastRulesEmission(wrapper)).toEqual([
      expect.objectContaining({
        attribute: '',
        fieldType: PATTERN_FIELD_TYPES.string,
      }),
    ]);
  });

  test('hides add button when readonly prop is true', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
        readonly: true,
      },
    });

    expect(selectAddButton(wrapper).exists()).toBe(false);
  });

  test('disables add button when disabled prop is true', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
        disabled: true,
      },
    });

    expect(selectAddButton(wrapper).props('disabled')).toBe(true);
  });

  test('disables remove buttons when readonly prop is true', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
        readonly: true,
      },
    });

    selectPatternRemoveRuleButtons(wrapper).forEach((removeButton) => {
      expect(removeButton.props('disabled')).toBe(true);
    });
  });

  test('disables remove buttons when disabled prop is true', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
        disabled: true,
      },
    });

    selectPatternRemoveRuleButtons(wrapper).forEach((removeButton) => {
      expect(removeButton.props('disabled')).toBe(true);
    });
  });

  test('disables remove button when rule disables remove action', () => {
    const disabledRemoveRule = {
      ...rules[0],
      disabled: {
        remove: true,
      },
    };
    const wrapper = factory({
      propsData: {
        rules: [disabledRemoveRule],
        attributes: [],
      },
    });

    expect(selectPatternRemoveRuleButtonByIndex(wrapper, 0).props('disabled')).toBe(true);
  });

  test('passes attribute options to pattern rule field', () => {
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

    const ruleProps = selectPatternRuleFieldByIndex(wrapper, 0).props();

    expect(ruleProps.operators).toEqual(attributeWithOptions.options.operators);
    expect(ruleProps.type).toBe(attributeWithOptions.options.type);
    expect(ruleProps.disabled).toBe(true);
  });

  test('passes readonly and disabled state to pattern rule field', () => {
    const attribute = { value: 'test-attr', options: {} };

    const readonlyWrapper = factory({
      propsData: {
        rules: [{ ...rules[0], attribute: attribute.value }],
        attributes: [attribute],
        readonly: true,
      },
    });
    const disabledWrapper = factory({
      propsData: {
        rules: [{ ...rules[0], attribute: attribute.value }],
        attributes: [attribute],
        disabled: true,
      },
    });

    expect(selectPatternRuleFieldByIndex(readonlyWrapper, 0).props('disabled')).toBe(true);
    expect(selectPatternRuleFieldByIndex(disabledWrapper, 0).props('disabled')).toBe(true);
  });

  test('passes disabled field and operator flags from rule', () => {
    const rule = {
      ...rules[0],
      disabled: {
        field: true,
        operator: true,
      },
    };
    const wrapper = factory({
      propsData: {
        rules: [rule],
        attributes: [],
      },
    });

    expect(selectPatternRuleFieldByIndex(wrapper, 0).props()).toEqual(
      expect.objectContaining({
        disabledField: true,
        disabledOperator: true,
      }),
    );
  });

  test('keeps one-of string values in rule props', () => {
    const rule = {
      attribute: 'string-attribute',
      operator: PATTERN_OPERATORS.isOneOf,
      value: ['first value', 'second value'],
      fieldType: PATTERN_FIELD_TYPES.string,
      key: 'string-key',
    };
    const wrapper = factory({
      propsData: {
        rules: [rule],
        attributes: [{ value: rule.attribute }],
      },
    });

    const props = selectPatternRuleFieldByIndex(wrapper, 0).props();

    expect(props.operators).toEqual(
      expect.arrayContaining([PATTERN_OPERATORS.isOneOf, PATTERN_OPERATORS.isNotOneOf]),
    );
    expect(props.rule).toMatchObject({
      operator: PATTERN_OPERATORS.isOneOf,
      value: rule.value,
    });
  });

  test('does not show add button error by default', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    expect(selectAddButton(wrapper).props('error')).toBe('');
  });
});
