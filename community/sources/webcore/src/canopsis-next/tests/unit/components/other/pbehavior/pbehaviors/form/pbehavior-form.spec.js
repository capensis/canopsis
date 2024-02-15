import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';

const stubs = {
  'pbehavior-general-form': true,
  'pbehavior-patterns-form': true,
};

const selectPbehaviorGeneralForm = wrapper => wrapper.find('pbehavior-general-form-stub');
const selectPatternsField = wrapper => wrapper.find('pbehavior-patterns-form-stub');

describe('pbehavior-form', () => {
  const factory = generateShallowRenderer(PbehaviorForm, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  const snapshotFactory = generateRenderer(PbehaviorForm, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('General options updated after trigger general pbehavior form', () => {
    const form = {
      name: Faker.datatype.string(),
      enabled: Faker.datatype.boolean(),
      patterns: {},
      comments: [],
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newForm = {
      name: Faker.datatype.string(),
      enabled: Faker.datatype.boolean(),
      patterns: {},
      comments: [],
    };

    selectPbehaviorGeneralForm(wrapper).vm.$emit('input', newForm);

    expect(wrapper).toEmit('input', newForm);
  });

  test('Filter updated after trigger pbehavior filter field', () => {
    const form = {
      name: Faker.datatype.string(),
      enabled: Faker.datatype.boolean(),
      patterns: {},
      comments: [],
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newPatterns = [
    ];

    selectPatternsField(wrapper).vm.$emit('input', newPatterns);

    expect(wrapper).toEmit('input', {
      ...form,
      patterns: newPatterns,
    });
  });

  test('Renders `pbehavior-form` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          name: 'pbehavior',
          patterns: {},
          comments: [],
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehavior-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          name: 'pbehavior',
          patterns: {},
          comments: [],
        },
        noPattern: true,
        noEnabled: true,
        noComments: true,
        withStartOnTrigger: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
