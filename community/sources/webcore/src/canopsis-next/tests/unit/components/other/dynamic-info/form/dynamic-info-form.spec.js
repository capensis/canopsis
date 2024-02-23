import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';

const stubs = {
  'dynamic-info-general-form': true,
  'dynamic-info-infos-form': true,
  'dynamic-info-patterns-form': true,
};

const selectDynamicInfoGeneralForm = wrapper => wrapper.find('dynamic-info-general-form-stub');
const selectDynamicInfoInfosForm = wrapper => wrapper.find('dynamic-info-infos-form-stub');
const selectDynamicInfoPatternsForm = wrapper => wrapper.find('dynamic-info-patterns-form-stub');

describe('dynamic-info-form', () => {
  const factory = generateShallowRenderer(DynamicInfoForm, { stubs });
  const snapshotFactory = generateRenderer(DynamicInfoForm, { stubs });

  test('Dynamic info general fields changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: {
          infos: [],
        },
      },
    });

    const dynamicInfoGeneralForm = selectDynamicInfoGeneralForm(wrapper);

    const newForm = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
      enabled: Faker.datatype.boolean(),
      disable_during_periods: [],
      description: Faker.datatype.string(),
    };

    dynamicInfoGeneralForm.triggerCustomEvent('input', newForm);

    expect(wrapper).toEmitInput(newForm);
  });

  test('Dynamic info infos changed after trigger infos form', () => {
    const wrapper = factory({
      propsData: {
        form: {
          infos: [],
        },
      },
    });

    const dynamicInfoInfosForm = selectDynamicInfoInfosForm(wrapper);

    const newInfos = [
      {
        name: Faker.datatype.string(),
        value: Faker.datatype.string(),
        actions: [],
      },
    ];

    dynamicInfoInfosForm.triggerCustomEvent('input', newInfos);

    expect(wrapper).toEmitInput({ infos: newInfos });
  });

  test('Dynamic info patterns changed after trigger patterns form', () => {
    const wrapper = factory({
      propsData: {
        form: {
          infos: [],
          patterns: {},
        },
      },
    });

    const patternsForm = selectDynamicInfoPatternsForm(wrapper);

    const newPatterns = {
      alarm_pattern: {},
      entity_pattern: {},
    };

    patternsForm.triggerCustomEvent('input', newPatterns);

    expect(wrapper).toEmitInput({
      infos: [],
      patterns: newPatterns,
    });
  });

  test('Renders `dynamic-info-form` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          infos: [{}],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `dynamic-info-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          enabled: true,
          infos: [{}],
          patterns: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `dynamic-info-form` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          infos: [{}],
        },
      },
    });

    await wrapper.setData({
      hasGeneralFormAnyError: true,
      hasInfosFormAnyError: true,
      hasPatternsFormAnyError: true,
    });

    expect(wrapper).toMatchSnapshot();
  });
});
