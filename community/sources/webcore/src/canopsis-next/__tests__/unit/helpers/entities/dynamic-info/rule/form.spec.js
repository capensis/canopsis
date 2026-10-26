import { describe, expect, it } from 'vitest';

import {
  dynamicInfoToForm,
  formToDynamicInfo,
} from '@/helpers/entities/dynamic-info/rule/form';

describe('dynamic info form helpers', () => {
  it('does not expose the system-managed ID in the form', () => {
    const form = dynamicInfoToForm({ _id: 'dynamic-info-id' });

    expect(form).not.toHaveProperty('_id');
  });

  it('omits the system-managed ID from the request', () => {
    const dynamicInfo = formToDynamicInfo({
      _id: 'dynamic-info-id',
      patterns: {},
      infos: [],
    });

    expect(dynamicInfo).not.toHaveProperty('_id');
  });
});
