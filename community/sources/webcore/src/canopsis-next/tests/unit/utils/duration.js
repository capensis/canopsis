import Faker from 'faker';

import { randomArrayItem } from '@unit/utils/array';
import { TIME_UNITS } from '@/constants';

export const randomDurationValue = () => ({
  unit: randomArrayItem(Object.values(TIME_UNITS)),
  value: Faker.datatype.number(),
  enabled: Faker.datatype.boolean(),
});
