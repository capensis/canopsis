import { describe, expect, it } from 'vitest';

import {
  formToMetaAlarmRule,
  metaAlarmRuleToForm,
} from '@/helpers/entities/meta-alarm/rule/form';

describe('meta alarm rule form helpers', () => {
  it('does not expose the system-managed ID in the form', () => {
    const form = metaAlarmRuleToForm({ _id: 'meta-alarm-rule-id' });

    expect(form).not.toHaveProperty('_id');
  });

  it('omits the system-managed ID from the request', () => {
    const form = metaAlarmRuleToForm();
    const rule = formToMetaAlarmRule({ ...form, _id: 'meta-alarm-rule-id' });

    expect(rule).not.toHaveProperty('_id');
  });
});
