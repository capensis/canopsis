import { EXTERNAL_METRIC_UNITS } from '@/constants';

/**
 * Converts a file size from bytes to the specified unit.
 *
 * @param {number} [fileSize=0] - The size of the file in bytes.
 * @param {string} [unit=EXTERNAL_METRIC_UNITS.byte] - The unit to convert the file size to.
 *        Supported units are defined in EXTERNAL_METRIC_UNITS: byte, kilobyte, megabyte, gigabyte, terabyte.
 * @returns {number} The file size converted to the specified unit.
 */
export const convertFileSizeToUnit = (fileSize = 0, unit = EXTERNAL_METRIC_UNITS.byte) => {
  const degree = {
    [EXTERNAL_METRIC_UNITS.terabyte]: 4,
    [EXTERNAL_METRIC_UNITS.gigabyte]: 3,
    [EXTERNAL_METRIC_UNITS.megabyte]: 2,
    [EXTERNAL_METRIC_UNITS.kilobyte]: 1,
    [EXTERNAL_METRIC_UNITS.byte]: 0,
  }[unit];

  return fileSize / (1024 ** degree);
};
