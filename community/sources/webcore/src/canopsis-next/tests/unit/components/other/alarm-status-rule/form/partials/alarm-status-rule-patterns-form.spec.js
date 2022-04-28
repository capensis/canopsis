import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import AlarmStatusRulePatternsForm from '@/components/other/alarm-status-rule/form/partials/alarm-status-rule-patterns-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-patterns-field': true,
};

const factory = (options = {}) => shallowMount(AlarmStatusRulePatternsForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmStatusRulePatternsForm, {
  localVue,
  stubs,

  ...options,
});

const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');

describe('alarm-status-rule-patterns-form', () => {
  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory();

    const patternsField = selectPatternsField(wrapper);

    const newPatterns = {
      alarm_pattern: {},
      entity_pattern: {},
    };

    patternsField.vm.$emit('input', newPatterns);

    expect(wrapper).toEmit('input', newPatterns);
  });

  test('Renders `alarm-status-rule-patterns-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `alarm-status-rule-patterns-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          alarm_pattern: {},
          entity_pattern: {},
        },
        flapping: true,
        disabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
