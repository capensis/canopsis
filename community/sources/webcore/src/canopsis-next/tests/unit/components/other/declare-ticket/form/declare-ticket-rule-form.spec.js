import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createAuthModule, createTemplateVarsModule } from '@unit/utils/store';
import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';

import { USER_PERMISSIONS } from '@/constants';

import DeclareTicketRuleForm from '@/components/other/declare-ticket/form/declare-ticket-rule-form.vue';

const stubs = {
  'c-enabled-field': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
  'declare-ticket-rule-general-form': true,
  'declare-ticket-rule-patterns-form': true,
  'declare-ticket-rule-test-query': true,
  'template-testing-test-variables-tab': true,
  'template-testing-test-variables': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectDeclareTicketRuleGeneralForm = wrapper => wrapper.find('declare-ticket-rule-general-form-stub');
const selectDeclareTicketRulePatternsForm = wrapper => wrapper.find('declare-ticket-rule-patterns-form-stub');
const selectDeclareTicketRuleTestQuery = wrapper => wrapper.find('declare-ticket-rule-test-query-stub');

describe('declare-ticket-rule-form', () => {
  const { authModule, currentUserPermissionsById } = createAuthModule();
  const { templateVarsModule } = createTemplateVarsModule();
  const store = createMockedStoreModules([authModule, templateVarsModule]);

  const form = {
    enabled: true,
    name: 'name',
    system_name: 'System name',
    patterns: {
      alarm_patterns: {},
    },
    webhooks: [],
  };

  const factory = generateShallowRenderer(DeclareTicketRuleForm, { stubs, store });
  const snapshotFactory = generateRenderer(DeclareTicketRuleForm, { stubs, store });

  test('General form is rendered in general tab', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    expect(selectDeclareTicketRuleGeneralForm(wrapper).exists()).toBe(true);
  });

  test('Enabled changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectEnabledField(wrapper).triggerCustomEvent('input', false);

    expect(wrapper).toEmitInput({
      ...form,
      enabled: false,
    });
  });

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

  test('Test query form receives form prop', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    expect(selectDeclareTicketRuleTestQuery(wrapper).props('form')).toEqual(form);
  });

  test('Renders `declare-ticket-rule-form` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {},
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `declare-ticket-rule-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `declare-ticket-rule-form` with errors', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          webhooks: [],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `declare-ticket-rule-form` with template testing tab access', () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USER_PERMISSIONS.technical.templateTesting]: { actions: [] },
    });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([authModule, templateVarsModule]),
      propsData: {
        form: {
          webhooks: [],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
