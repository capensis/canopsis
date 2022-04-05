import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { PATTERN_OPERATORS, QUICK_RANGES, TIME_UNITS } from '@/constants';

import CPatternRulesField from '@/components/forms/fields/pattern/c-pattern-rules-field.vue';

const localVue = createVueInstance();

const stubs = {
  'c-pattern-rule-field': true,
  'c-action-btn': true,
};

const factory = (options = {}) => shallowMount(CPatternRulesField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CPatternRulesField, {
  localVue,
  stubs,

  ...options,
});

const selectAddButton = wrapper => wrapper.find('v-btn-stub');
const selectPatternRulesField = wrapper => wrapper.findAll('c-pattern-rule-field-stub');
const selectPatternRuleFieldByIndex = (wrapper, index) => selectPatternRulesField(wrapper)
  .at(index);
const selectPatternRemoveRuleButtonByIndex = (wrapper, index) => wrapper.findAll('c-action-btn-stub')
  .at(index);

describe('c-pattern-rules-field', () => {
  const rules = [
    {
      attribute: 'attribute 1',
      operator: PATTERN_OPERATORS.equal,
      value: 'attribute value',
      key: 'key 1',
    },
    {
      attribute: 'attribute 2',
      operator: PATTERN_OPERATORS.notEqual,
      value: 'attribute value 2',
      key: 'key 2',
    },
    {
      attribute: 'attribute 3',
      operator: PATTERN_OPERATORS.contains,
      value: 'attribute contains',
      key: 'key 3',
    },
  ];

  test('Rule removed after trigger remove event on the pattern Rule field', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    const removeSecondRuleButton = selectPatternRemoveRuleButtonByIndex(wrapper, 1);

    removeSecondRuleButton.vm.$emit('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
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

    lastRule.vm.$emit('input', updatedRule);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      rules[0],
      rules[1],
      updatedRule,
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

    lastRule.vm.$emit('input', updatedRule);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      rules[0],
      {
        ...updatedRule,
        operator: '',
        value: '',
      },
      rules[2],
    ]);
  });

  test('Value changed to array after update rule with array operator', () => {
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [],
      },
    });

    const lastRule = selectPatternRuleFieldByIndex(wrapper, 1);

    const updatedRule = {
      ...rules[1],
      operator: PATTERN_OPERATORS.hasNot,
    };

    lastRule.vm.$emit('input', updatedRule);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      rules[0],
      {
        ...updatedRule,
        value: [updatedRule.value],
      },
      rules[2],
    ]);
  });

  test('Rule added after click on the add button', () => {
    const attributeItem = {
      value: 'test',
    };
    const wrapper = factory({
      propsData: {
        rules,
        attributes: [
          attributeItem,
        ],
      },
    });

    const addButton = selectAddButton(wrapper);

    addButton.vm.$emit('click');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      ...rules,
      {
        attribute: attributeItem.value,
        dictionary: '',
        field: '',
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

  test('Renders `c-pattern-rules-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rules: [],
        attributes: [],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-pattern-rules-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        rules,
        attributes: [
          {
            text: 'Attribute text',
            value: 'attribute value',
            options: {
              operators: [PATTERN_OPERATORS.notEqual],
              customProp: 'customPropValue',
            },
          },
        ],
        required: true,
        disabled: true,
        name: 'customName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
