import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import IdleRulePatternsForm from '@/components/other/idle-rule/form/partials/idle-rule-patterns-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-patterns-field': true,
};

const factory = (options = {}) => shallowMount(IdleRulePatternsForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(IdleRulePatternsForm, {
  localVue,
  stubs,

  ...options,
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
