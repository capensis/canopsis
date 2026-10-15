import { SAMPLINGS } from '@/constants';

import { getDateLabelBySampling, convertDataSizeValueToTickString } from '@/helpers/entities/metric/list';

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

describe('convertDataSizeValueToTickString', () => {
  const byte = 1;
  const bytesInKilobytes = 1024;
  const bytesInMegabytes = bytesInKilobytes * 1024;
  const bytesInGigabytes = bytesInMegabytes * 1024;
  const bytesInTerabytes = bytesInGigabytes * 1024;

  it('should return 1B', () => {
    expect(convertDataSizeValueToTickString(byte)).toBe('1B');
  });

  it('should return 1KB', () => {
    expect(convertDataSizeValueToTickString(bytesInKilobytes)).toBe('1KB');
  });

  it('should return 512KB', () => {
    expect(convertDataSizeValueToTickString(0.5 * bytesInKilobytes)).toBe('512B');
  });

  it('should return 1MB', () => {
    expect(convertDataSizeValueToTickString(2 * bytesInMegabytes)).toBe('2MB');
  });

  it('should return 512KB', () => {
    expect(convertDataSizeValueToTickString(0.5 * bytesInMegabytes)).toBe('512KB');
  });

  it('should return 3GB', () => {
    expect(convertDataSizeValueToTickString(3 * bytesInGigabytes)).toBe('3GB');
  });

  it('should return 512MB', () => {
    expect(convertDataSizeValueToTickString(0.5 * bytesInGigabytes)).toBe('512MB');
  });

  it('should return 4TB', () => {
    expect(convertDataSizeValueToTickString(4 * bytesInTerabytes)).toBe('4TB');
  });

  it('should return 512GB', () => {
    expect(convertDataSizeValueToTickString(0.5 * bytesInTerabytes)).toBe('512GB');
  });
});
