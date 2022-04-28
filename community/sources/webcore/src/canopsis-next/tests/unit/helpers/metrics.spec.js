import { getDateLabelBySampling } from '@/helpers/metrics';
import { SAMPLINGS } from '@/constants';

describe('getDateLabelBySampling', () => {
  const timestamp = 1636523087405;

  it('Valid string with sampling by hour', () => {
    expect(getDateLabelBySampling(timestamp, SAMPLINGS.hour)).toBe('10/11/2021 06:44');
  });

  it('Valid string with sampling by day', () => {
    expect(getDateLabelBySampling(timestamp, SAMPLINGS.day)).toBe('10/11/2021');
  });

  it('Valid string with sampling by week', () => {
    expect(getDateLabelBySampling(timestamp, SAMPLINGS.week)).toBe('10/11/2021 - \n17/11/2021');
  });

  it('Valid string with sampling by month', () => {
    expect(getDateLabelBySampling(timestamp, SAMPLINGS.month)).toBe('November 2021');
  });
});
