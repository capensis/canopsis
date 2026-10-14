import Faker from 'faker';

export const fakeMeta = ({ count, limit = 10, page = 1 } = {}) => ({
  meta: {
    page,
    per_page: limit,
    page_count: Math.floor(count / limit),
    total_count: count,
  },
});

export const fakeParams = ({ limit, page } = {}) => ({
  meta: {
    limit: limit ?? Faker.datatype.number(),
    page: page ?? Faker.datatype.number(),
  },
});
