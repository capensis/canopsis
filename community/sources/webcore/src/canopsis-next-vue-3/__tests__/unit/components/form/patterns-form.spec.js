import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import FilterForm from '@/components/other/filter/form/filter-form.vue';

const stubs = {
  'c-name-field': true,
  'c-patterns-field': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');

describe('patterns-form', () => {
  const factory = generateShallowRenderer(FilterForm, { stubs });

  test('Title changed after trigger text field', () => {
    const form = {
      title: '',
      alarm_pattern: {},
    };
    const wrapper = factory({
      propsData: {
        form,
        withTitle: true,
        withAlarm: true,
      },
    });

    const title = Faker.datatype.string();

    selectNameField(wrapper).triggerCustomEvent('update:modelValue', title);

    expect(globalThis.__TEST_FORM__.values.value).toEqual({ ...form, title });
  });

  test('Patterns changed after trigger patterns field', () => {
    const form = {
      title: Faker.datatype.string(),
      entity_pattern: {},
    };
    const wrapper = factory({
      propsData: {
        form,
        withTitle: true,
        withEntity: true,
      },
    });

    const patternsField = selectPatternsField(wrapper);

    const newForm = {
      ...form,
      entity_pattern: {
        id: Faker.datatype.string(),
      },
    };

    patternsField.triggerCustomEvent('update:value', newForm);

    expect(globalThis.__TEST_FORM__.values.value).toEqual(newForm);
  });
});
