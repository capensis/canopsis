import { describe, expect, it } from 'vitest';

import { getTemplateRuleValidationErrors } from '@/components/other/template-testing/test-variables/hooks/template-testing-test-variables';

describe('getTemplateRuleValidationErrors', () => {
  it('extracts rule errors from a template-validate response', () => {
    const errors = getTemplateRuleValidationErrors({
      errors: {
      'rule.config.actions.0.value': 'Invalid value',
      'testdata.event': 'The field must match the rule',
      },
    });

    expect(errors).toEqual({
      'config.actions.0.value': 'Invalid value',
    });
  });
});
