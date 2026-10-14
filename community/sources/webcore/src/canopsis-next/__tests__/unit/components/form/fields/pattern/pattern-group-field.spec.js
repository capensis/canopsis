import { generateShallowRenderer } from '@unit/utils/vue';

import PatternGroupField from '@/components/forms/fields/pattern/pattern-group-field.vue';

const stubs = {
  'pattern-operator-information': true,
  'pattern-rules-field': true,
};

const selectPatternRulesField = wrapper => wrapper.find('pattern-rules-field-stub');

describe('pattern-group-field', () => {
  const factory = generateShallowRenderer(PatternGroupField, { stubs });

  test('Rules updated after trigger input event on pattern rules field', async () => {
    const wrapper = factory({
      propsData: {
        group: {
          rules: [],
        },
        attributes: [],
      },
    });

    const newRules = [{}, {}];

    selectPatternRulesField(wrapper).triggerCustomEvent('update:rules', newRules);

    expect(wrapper).toEmitInput({
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

    selectPatternRulesField(wrapper).triggerCustomEvent('update:rules', []);

    expect(wrapper).toHaveBeenEmit('remove');
  });
});
