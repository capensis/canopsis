import Faker from 'faker';

import { randomArrayItem } from '@unit/utils/array';

import { TIME_UNITS } from '@/constants';

export const randomTimeUnit = () => randomArrayItem(Object.values(TIME_UNITS));

export const randomDurationValue = () => ({
  unit: randomTimeUnit(),
  value: Faker.datatype.number(),
  enabled: Faker.datatype.boolean(),
});

export const randomPeriodicRefresh = () => ({
  value: Faker.datatype.number(),
  unit: randomTimeUnit(),
});
