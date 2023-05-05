import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import IdleRulePatternsForm from '@/components/other/idle-rule/form/partials/idle-rule-patterns-form.vue';

const stubs = {
  'c-patterns-field': true,
};

const factory = generateShallowRenderer(IdleRulePatternsForm, { stubs,
});

const snapshotFactory = generateRenderer(IdleRulePatternsForm, { stubs,
});

const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');

describe('idle-rule-patterns-form', () => {
  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory();

    const patternsField = selectPatternsField(wrapper);

    const newPatterns = {
      alarm_pattern: {},
    };

    patternsField.vm.$emit('input', newPatterns);

    expect(wrapper).toEmit('input', newPatterns);
  });

  test('Renders `idle-rule-patterns-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `idle-rule-patterns-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          alarm_pattern: {},
        },
        isEntityType: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
