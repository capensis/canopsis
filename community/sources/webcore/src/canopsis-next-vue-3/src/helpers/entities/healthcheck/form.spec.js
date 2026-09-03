import { describe, expect, it } from 'vitest';

import { HEALTHCHECK_ENGINES_NAMES } from '@/constants';

import { healthcheckParametersToForm } from './form';

describe('healthcheckParametersToForm', () => {
  it('uses default minimal and optimal instance values when engine parameters are missing', () => {
    const form = healthcheckParametersToForm();

    Object.values(HEALTHCHECK_ENGINES_NAMES).forEach((engineName) => {
      expect(form[engineName]).toEqual({
        enabled: false,
        minimal: 0,
        optimal: 0,
      });
    });
  });
});
