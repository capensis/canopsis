import { mount } from '@vue/test-utils';
import { createPinia } from 'pinia';
import { describe, expect, it, vi } from 'vitest';
import { defineComponent } from 'vue';

import { i18n } from '@/i18n';

import {
  getTemplateFieldsValidationErrors,
  getTemplateRuleValidationErrors,
  useTemplateTestActions,
  useTemplateTestingValidationForm,
} from '@/components/other/template-testing/test-variables/hooks/template-testing-test-variables';

describe('getTemplateRuleValidationErrors', () => {
  it('extracts rule errors from a template-validate response', () => {
    const errors = getTemplateRuleValidationErrors(
      {
        errors: {
          'rule.config.actions.0.value': 'Invalid value',
          'testdata.event': 'The field must match the rule',
        },
      },
      {
        config: {
          actions: [{ key: 'action-key' }],
        },
      },
    );

    expect(errors).toEqual({
      'config.actions.action-key.value': 'Invalid value',
    });
  });
});

describe('getTemplateFieldsValidationErrors', () => {
  it('maps result errors to configured fields using their result keys', () => {
    const errors = getTemplateFieldsValidationErrors(
      {
        'actions.0.parameters.url': {
          err: { line: 2, message: 'Invalid template' },
        },
        unknown: {
          err: { line: 1, message: 'Unknown error' },
        },
      },
      [
        {
          key: 'actions.0.parameters.webhook.request.url',
          resultKey: 'actions.0.parameters.url',
        },
      ],
    );

    expect(errors).toEqual({
      'templateTesting.actions.0.parameters.url': '2|Invalid template',
    });
  });

  it('uses the form key when a result key is not needed', () => {
    const errors = getTemplateFieldsValidationErrors(
      {
        output_template: {
          err: { line: 4, message: 'Unexpected closing brace' },
        },
      },
      [{ key: 'output_template' }],
    );

    expect(errors).toEqual({
      'templateTesting.output_template': '4|Unexpected closing brace',
    });
  });
});

describe('useTemplateTestingValidationForm', () => {
  it('resets last run values without passing symbols to template field props', () => {
    const TestComponent = defineComponent({
      template: '<div />',
      setup() {
        return useTemplateTestingValidationForm({
          active: false,
          generalForm: {},
        });
      },
    });
    const wrapper = mount(TestComponent, {
      global: {
        plugins: [createPinia(), i18n],
      },
    });

    wrapper.vm.setLastRunVariables({ template: 'previous value' });
    wrapper.vm.resetLastRunVariables();

    expect(wrapper.vm.lastRunVariables).toEqual({ template: null });
  });
});

describe('useTemplateTestActions', () => {
  it('validates test data before opening the save-as-new modal', async () => {
    const validateForSubmit = vi.fn().mockResolvedValue(true);
    const TestComponent = defineComponent({
      template: '<div />',
      setup() {
        return useTemplateTestActions(
          { ruleId: 'rule-id', type: 4 },
          {
            testFormContext: { validateForSubmit: vi.fn() },
            testDataFormContext: { validateForSubmit },
          },
        );
      },
    });
    const wrapper = mount(TestComponent, {
      global: {
        plugins: [createPinia(), i18n],
      },
    });

    await expect(wrapper.vm.saveAsNew(() => ({}))).resolves.toBeUndefined();

    expect(validateForSubmit).toHaveBeenCalledOnce();
  });
});
