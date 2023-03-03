import Faker from 'faker';

import { createVueInstance, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';

const localVue = createVueInstance();

const stubs = {
  'pbehavior-general-form': true,
  'pbehavior-comments-field': true,
  'pbehavior-filter-field': true,
  'pbehavior-recurrence-rule-field': true,
  'c-enabled-color-picker-field': true,
};

const selectPbehaviorGeneralForm = wrapper => wrapper.find('pbehavior-general-form-stub');
const selectPbehaviorCommentsField = wrapper => wrapper.find('pbehavior-comments-field-stub');
const selectPbehaviorFilterField = wrapper => wrapper.find('pbehavior-filter-field-stub');
const selectPbehaviorRruleField = wrapper => wrapper.find('pbehavior-recurrence-rule-field-stub');
const selectEnabledColorPickerField = wrapper => wrapper.find('c-enabled-color-picker-field-stub');

describe('pbehavior-form', () => {
  const factory = generateShallowRenderer(PbehaviorForm, { localVue, stubs });
  const snapshotFactory = generateRenderer(PbehaviorForm, { localVue, stubs });

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

  test('Comments updated after trigger pbehavior comments field', () => {
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

    const newComments = [
      {
        key: Faker.datatype.string(),
        message: Faker.datatype.string(),
      },
    ];

    selectPbehaviorCommentsField(wrapper).vm.$emit('input', newComments);

    expect(wrapper).toEmit('input', {
      ...form,
      comments: newComments,
    });
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

    selectPbehaviorFilterField(wrapper).vm.$emit('input', newPatterns);

    expect(wrapper).toEmit('input', {
      ...form,
      patterns: newPatterns,
    });
  });

  test('Recurrence rule updated after trigger general rrule field', () => {
    const form = {
      name: Faker.datatype.string(),
      rrule: '',
      patterns: {},
      comments: [],
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newForm = {
      ...form,
      rrule: Faker.datatype.string(),
    };

    selectPbehaviorRruleField(wrapper).vm.$emit('input', newForm);

    expect(wrapper).toEmit('input', newForm);
  });

  test('Color changed after trigger color field', () => {
    const form = {
      name: Faker.datatype.string(),
      rrule: '',
      patterns: {},
      comments: [],
    };

    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newColor = Faker.internet.color();

    selectEnabledColorPickerField(wrapper).vm.$emit('input', newColor);

    expect(wrapper).toEmit('input', {
      ...form,

      color: newColor,
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
