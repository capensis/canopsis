import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import MetaAlarmRulePatternsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-patterns-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-patterns-field': true,
};

const factory = (options = {}) => shallowMount(MetaAlarmRulePatternsForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(MetaAlarmRulePatternsForm, {
  localVue,
  stubs,

  ...options,
});

const selectPatternsFieldNode = wrapper => wrapper.vm.$children[0];

describe('meta-alarm-rule-patterns-form', () => {
  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory();

    const patternsFieldNode = selectPatternsFieldNode(wrapper);

    const patterns = {
      alarm_pattern: {},
      entity_pattern: {},
      event_pattern: {},
    };

    patternsFieldNode.$emit('input', patterns);

    expect(wrapper).toEmit('input', patterns);
  });

  test('Renders `meta-alarm-rule-patterns-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `meta-alarm-rule-patterns-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          alarm_pattern: {},
          entity_pattern: {},
          event_pattern: {},
        },
        disabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
