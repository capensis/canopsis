import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { randomArrayItem } from '@unit/utils/array';

import { META_ALARMS_FORM_STEPS, META_ALARMS_RULE_TYPES } from '@/constants';

import { metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';

const stubs = {
  'meta-alarm-rule-general-form': true,
  'meta-alarm-rule-type-field': true,
  'meta-alarm-rule-parameters-form': true,
  'c-information-block': true,
};

const selectMetaAlarmRuleGeneralForm = wrapper => wrapper.find('meta-alarm-rule-general-form-stub');
const selectMetaAlarmRuleTypeField = wrapper => wrapper.find('meta-alarm-rule-type-field-stub');

describe('meta-alarm-rule-form', () => {
  const form = metaAlarmRuleToForm();

  const factory = generateShallowRenderer(MetaAlarmRuleForm, { stubs });
  const snapshotFactory = generateRenderer(MetaAlarmRuleForm, { stubs });

  test('General fields updated after trigger general form', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newFields = {
      ...form,
      id: Faker.datatype.string(),
      type: randomArrayItem(Object.values(META_ALARMS_RULE_TYPES)),
      name: Faker.datatype.string(),
      auto_resolve: Faker.datatype.boolean(),
      output_template: Faker.datatype.string(),
    };

    await selectMetaAlarmRuleGeneralForm(wrapper).triggerCustomEvent('input', newFields);

    expect(wrapper).toEmitInput(newFields);
  });

  test('Type changed after trigger type field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectMetaAlarmRuleTypeField(wrapper).triggerCustomEvent('input', META_ALARMS_RULE_TYPES.attribute);

    expect(wrapper).toEmitInput({
      ...form,
      type: META_ALARMS_RULE_TYPES.attribute,
    });
  });

  test('Renders `meta-alarm-rule-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `meta-alarm-rule-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        disabledIdField: true,
        activeStep: META_ALARMS_FORM_STEPS.parameters,
        alarmInfos: [{ value: 'alarm-infos' }],
        entityInfos: [{ value: 'alarm-infos' }],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
