import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import PatternsForm from '@/components/forms/patterns-form.vue';

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
  'c-patterns-field': true,
};

const selectTextField = wrapper => wrapper.find('.v-text-field');
const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');

describe('patterns-form', () => {
  const factory = generateShallowRenderer(PatternsForm, { stubs });
  const snapshotFactory = generateRenderer(PatternsForm, { stubs });

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

    const textField = selectTextField(wrapper);

    textField.vm.$emit('input', title);

    expect(wrapper).toEmit('input', { ...form, title });
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

    patternsField.vm.$emit('input', newForm);

    expect(wrapper).toEmit('input', newForm);
  });

  test('Renders `patterns-form` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `patterns-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          title: 'Title',
          alarm_pattern: {},
          event_pattern: {},
          entity_pattern: {},
          pbehavior_pattern: {},
        },
        withTitle: true,
        withAlarm: true,
        withEvent: true,
        withEntity: true,
        withPbehavior: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
