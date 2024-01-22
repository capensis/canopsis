import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';

const stubs = {
  'pbehavior-general-form': true,
  'pbehavior-comments-field': true,
  'recurrence-rule-form': true,
  'pbehavior-recurrence-rule-exceptions-field': true,
  'c-enabled-color-picker-field': true,
  'c-collapse-panel': true,
  'pbehavior-patterns-form': true,
};

const selectPbehaviorGeneralForm = wrapper => wrapper.find('pbehavior-general-form-stub');
const selectPbehaviorCommentsField = wrapper => wrapper.find('pbehavior-comments-field-stub');
const selectEnabledColorPickerField = wrapper => wrapper.find('c-enabled-color-picker-field-stub');
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

    selectPbehaviorGeneralForm(wrapper).triggerCustomEvent('input', newForm);

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

    selectPbehaviorCommentsField(wrapper).triggerCustomEvent('input', newComments);

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

    selectPatternsField(wrapper).triggerCustomEvent('input', newPatterns);

    expect(wrapper).toEmit('input', {
      ...form,
      patterns: newPatterns,
    });
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

    selectEnabledColorPickerField(wrapper).triggerCustomEvent('input', newColor);

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

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
