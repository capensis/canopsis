import Faker from 'faker';

import { createVueInstance, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import PbehaviorPatternsForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-patterns-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-patterns-field': true,
};

const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');

describe('pbehavior-patterns-form', () => {
  const factory = generateShallowRenderer(PbehaviorPatternsForm, { localVue, stubs });
  const snapshotFactory = generateRenderer(PbehaviorPatternsForm, { localVue, stubs });

  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const newPatterns = {
      title: Faker.datatype.string(),
      entity_pattern: {
        id: Faker.datatype.string(),
      },
    };

    selectPatternsField(wrapper).vm.$emit('input', newPatterns);

    expect(wrapper).toEmit('input', newPatterns);
  });

  test('Renders `pbehavior-patterns-form` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          entity_pattern: {},
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehavior-patterns-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          entity_pattern: {},
        },
        readonly: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
