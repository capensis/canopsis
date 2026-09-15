import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createEntityInfoPropertyModule } from '@unit/utils/store';
import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';

import { metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';

const stubs = {
  'c-enabled-field': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
  'meta-alarm-rule-general-form': true,
  'meta-alarm-rule-parameters-form': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectMetaAlarmRuleGeneralForm = wrapper => wrapper.find('meta-alarm-rule-general-form-stub');
const selectMetaAlarmRuleParametersForm = wrapper => wrapper.find('meta-alarm-rule-parameters-form-stub');

describe('meta-alarm-rule-form', () => {
  const form = metaAlarmRuleToForm();

  const { entityInfoPropertyModule } = createEntityInfoPropertyModule();

  const store = createMockedStoreModules([entityInfoPropertyModule]);

  const factory = generateShallowRenderer(MetaAlarmRuleForm, { stubs });
  const snapshotFactory = generateRenderer(MetaAlarmRuleForm, { stubs });

  test('General form is rendered in general tab', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
      store,
    });

    expect(selectMetaAlarmRuleGeneralForm(wrapper).exists()).toBe(true);
  });

  test('Enabled changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
      store,
    });

    selectEnabledField(wrapper).triggerCustomEvent('input', false);

    expect(wrapper).toEmitInput({
      ...form,
      enabled: false,
    });
  });

  test('General fields updated after trigger general form', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
      store,
    });

    const newFields = {
      ...form,
      name: Faker.datatype.string(),
      output_template: Faker.datatype.string(),
    };

    await selectMetaAlarmRuleGeneralForm(wrapper).triggerCustomEvent('input', newFields);

    expect(wrapper).toEmitInput(newFields);
  });

  test('Parameters changed after trigger parameters form', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
      store,
    });

    const newFields = {
      ...form,
      patterns: {
        alarm_pattern: {},
      },
    };

    selectMetaAlarmRuleParametersForm(wrapper).triggerCustomEvent('input', newFields);

    expect(wrapper).toEmitInput(newFields);
  });

  test('Renders `meta-alarm-rule-form` with default props', () => {
    const wrapper = snapshotFactory({
      store,
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `meta-alarm-rule-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
        disabledIdField: true,
      },
      store,
    });

    expect(wrapper).toMatchSnapshot();
  });
});
