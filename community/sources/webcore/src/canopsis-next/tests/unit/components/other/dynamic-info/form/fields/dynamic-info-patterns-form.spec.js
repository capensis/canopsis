import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import DynamicInfoPatternsForm from '@/components/other/dynamic-info/form/fields/dynamic-info-patterns-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-patterns-field': true,
};

const factory = (options = {}) => shallowMount(DynamicInfoPatternsForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(DynamicInfoPatternsForm, {
  localVue,
  stubs,

  ...options,
});

const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');

describe('dynamic-info-patterns-form', () => {
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

  test('Renders `dynamic-info-patterns-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `dynamic-info-patterns-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          alarm_pattern: {},
          entity_pattern: {},
        },
        isEntityType: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
