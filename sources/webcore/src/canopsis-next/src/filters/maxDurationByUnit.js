import moment from 'moment';
import 'moment-duration-format';

import { TIME_UNITS } from '@/constants';
import { getUnitValueFromSeconds } from '@/helpers/time';

import momentDurationFrLocale from '@/i18n/moment-duration-fr';

moment.updateLocale('fr', momentDurationFrLocale);

const AVAILABLE_UNITS = [
  TIME_UNITS.year,
  TIME_UNITS.month,
  TIME_UNITS.day,
  TIME_UNITS.hour,
  TIME_UNITS.minute,
  TIME_UNITS.second,
];

export default function (value = 0, availableUnits = AVAILABLE_UNITS) {
  let unitValue;

  const maxUnit = availableUnits.find((unit) => {
    unitValue = Math.floor(getUnitValueFromSeconds(value, unit));
    return unitValue;
  });

  return `${unitValue}${maxUnit}`;
}
