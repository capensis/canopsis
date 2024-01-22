import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import PatternGroupField from '@/components/forms/fields/pattern/pattern-group-field.vue';

const stubs = {
  'pattern-operator-information': true,
  'pattern-rules-field': true,
};

const selectPatternRulesField = wrapper => wrapper.find('pattern-rules-field-stub');

describe('pattern-group-field', () => {
  const factory = generateShallowRenderer(PatternGroupField, { stubs });
  const snapshotFactory = generateRenderer(PatternGroupField, { stubs });

  test('Rules updated after trigger input event on pattern rules field', async () => {
    const wrapper = factory({
      propsData: {
        group: {
          rules: [],
        },
        attributes: [],
      },
    });

    const patternRulesField = selectPatternRulesField(wrapper);

    const newRules = [{}, {}];

    patternRulesField.triggerCustomEvent('input', newRules);

    expect(wrapper).toEmit('input', {
      rules: newRules,
    });
  });

  test('Group removed after trigger input event on pattern rules field without rules', async () => {
    const wrapper = factory({
      propsData: {
        group: {
          rules: [{}, {}],
        },
        attributes: [],
      },
    });

    selectPatternRulesField(wrapper).triggerCustomEvent('input', []);

    expect(wrapper).toEmit('remove');
  });

  test('Renders `pattern-group-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        group: {
          rules: [],
        },
        attributes: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-group-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        group: {
          rules: [],
        },
        attributes: [
          { text: 'Attribute text', value: 'attribute value' },
        ],
        disabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
