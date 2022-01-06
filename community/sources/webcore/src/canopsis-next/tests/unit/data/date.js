import Faker from 'faker';

export const fakeTimestamp = () => Faker.datatype.number({
  min: 1000000000,
  max: Date.now(),
});
