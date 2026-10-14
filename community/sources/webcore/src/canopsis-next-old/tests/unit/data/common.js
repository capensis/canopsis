import Faker from 'faker';

export const fakeObject = ({ fields = {}, fake = true, suffix = '' } = {}) => (
  Object.entries(fields).reduce((acc, [key, value]) => {
    acc[key] = fake ? Faker.fake(`{{${value}}}`) : `${value}${suffix}`;

    return acc;
  }, {})
);
