import { describe, expect, it } from 'vitest';

import { COLORS } from '@/config';
import { HEALTHCHECK_STATUSES } from '@/constants';

import { getHealthcheckNodeColor } from '@/helpers/entities/healthcheck/color';

describe('getHealthcheckNodeColor', () => {
  it.each([
    [HEALTHCHECK_STATUSES.healthy, COLORS.primary],
    [HEALTHCHECK_STATUSES.degraded, COLORS.healthcheck.warning],
    [HEALTHCHECK_STATUSES.unhealthy, COLORS.healthcheck.error],
    ['unexpected', COLORS.healthcheck.unknown],
  ])('maps the %s status to its healthcheck color', (status, color) => {
    expect(getHealthcheckNodeColor({ status })).toBe(color);
  });

  it('uses the unknown color when the node state is unavailable', () => {
    expect(
      getHealthcheckNodeColor({
        status: HEALTHCHECK_STATUSES.healthy,
        is_unknown: true,
      }),
    ).toBe(COLORS.healthcheck.unknown);
  });
});
