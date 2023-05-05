import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import ScenarioActionPatternsForm from '@/components/other/scenario/form/scenario-action-patterns-form.vue';

const stubs = {
  'c-patterns-field': true,
};

const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');

describe('scenario-action-patterns-form', () => {
  const factory = generateShallowRenderer(ScenarioActionPatternsForm, { stubs });
  const snapshotFactory = generateRenderer(ScenarioActionPatternsForm, { stubs });

  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory();

    const patternsField = selectPatternsField(wrapper);

    const newPatterns = {
      alarm_pattern: {},
    };

    patternsField.vm.$emit('input', newPatterns);

    expect(wrapper).toEmit('input', newPatterns);
  });

  test('Renders `scenario-action-patterns-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `scenario-action-patterns-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          alarm_pattern: {},
        },
        name: 'customName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
