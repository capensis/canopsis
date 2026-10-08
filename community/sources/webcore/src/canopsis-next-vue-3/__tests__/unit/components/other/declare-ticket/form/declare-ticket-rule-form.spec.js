import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';
import {
  createMockedStoreModules,
  createAuthModule,
  createTemplateVarsModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';

import DeclareTicketRuleForm from '@/components/other/declare-ticket/form/declare-ticket-rule-form.vue';

const declareTicketRuleTestQueryStub = {
  name: 'DeclareTicketRuleTestQuery',
  props: ['form'],
  template: '<div />',
};

const stubs = {
  'c-enabled-field': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
  'declare-ticket-rule-general-form': true,
  'declare-ticket-rule-patterns-form': true,
  'declare-ticket-rule-test-query': declareTicketRuleTestQueryStub,
  'template-testing-test-variables-tab': true,
  'template-testing-test-variables': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectDeclareTicketRuleGeneralForm = wrapper =>
  wrapper.find('declare-ticket-rule-general-form-stub');
const selectDeclareTicketRulePatternsForm = wrapper =>
  wrapper.find('declare-ticket-rule-patterns-form-stub');
const selectDeclareTicketRuleTestQuery = wrapper =>
  wrapper.findComponent(declareTicketRuleTestQueryStub);

describe('declare-ticket-rule-form', () => {
  const { authModule } = createAuthModule();
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

  test('General form uses the shared form context', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    expect(selectDeclareTicketRuleGeneralForm(wrapper).exists()).toBe(true);
    expect(globalThis.__TEST_FORM__.values.value).toEqual(form);
  });

  test('Patterns form uses the shared form context', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    expect(selectDeclareTicketRulePatternsForm(wrapper).exists()).toBe(true);
    expect(globalThis.__TEST_FORM__.values.value.patterns).toEqual(form.patterns);
  });

  test('Test query form receives form prop', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    expect(selectDeclareTicketRuleTestQuery(wrapper).props('form')).toEqual(form);
  });
});
