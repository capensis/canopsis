import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import DeclareTicketRuleForm from '@/components/other/declare-ticket/form/declare-ticket-rule-form.vue';

const stubs = {
  'declare-ticket-rule-general-form': true,
  'declare-ticket-rule-patterns-form': true,
  'declare-ticket-rule-test-query': true,
};

const selectTabItems = wrapper => wrapper.findAll('.v-tab');
const selectTestQueryTab = wrapper => selectTabItems(wrapper).at(2);
const selectDeclareTicketRuleGeneralForm = wrapper => wrapper.find('declare-ticket-rule-general-form-stub');
const selectDeclareTicketRulePatternsForm = wrapper => wrapper.find('declare-ticket-rule-patterns-form-stub');
const selectDeclareTicketRuleTestQuery = wrapper => wrapper.find('declare-ticket-rule-test-query-stub');

describe('declare-ticket-rule-form', () => {
  const form = {
    enabled: true,
    name: 'name',
    system_name: 'System name',
    patterns: {
      alarm_patterns: {},
    },
  };

  const factory = generateShallowRenderer(DeclareTicketRuleForm, { stubs });
  const snapshotFactory = generateRenderer(DeclareTicketRuleForm, { stubs });

  test('Form fields changed after trigger input event on general form', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newForm = {
      ...form,
      enabled: !form.enabled,
    };

    selectDeclareTicketRuleGeneralForm(wrapper).triggerCustomEvent('input', newForm);

    expect(wrapper).toEmitInput(newForm);
  });

  test('Patterns fields changed after trigger input event on patterns form', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newFormPatterns = {
      alarm_patterns: [{}, {}],
      entity_patterns: [{}, {}],
    };

    selectDeclareTicketRulePatternsForm(wrapper).triggerCustomEvent('input', newFormPatterns);

    expect(wrapper).toEmitInput({
      ...form,
      patterns: newFormPatterns,
    });
  });

  test('Patterns fields changed after trigger input event on patterns form', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    expect(selectDeclareTicketRuleTestQuery(wrapper).vm.form).toEqual(form);
  });

  test('Renders `declare-ticket-rule-form` with default props', async () => {
    const wrapper = snapshotFactory();

    await selectTestQueryTab(wrapper).trigger('click');

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `declare-ticket-rule-form` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    await selectTestQueryTab(wrapper).trigger('click');

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `declare-ticket-rule-form` with errors', async () => {
    const wrapper = snapshotFactory();

    await selectTestQueryTab(wrapper).trigger('click');

    await wrapper.setData({
      hasGeneralError: true,
      hasPatternsError: true,
    });

    expect(wrapper).toMatchSnapshot();
  });
});
