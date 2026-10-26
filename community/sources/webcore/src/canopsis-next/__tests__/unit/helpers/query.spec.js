import { convertQueryToRequest } from '@/helpers/query';

describe('helpers/query', () => {
  it('keeps a zero webhook rule type in request parameters', () => {
    expect(convertQueryToRequest({
      page: 1,
      itemsPerPage: 10,
      rule_type: 0,
    })).toEqual({
      page: 1,
      limit: 10,
      rule_type: 0,
    });
  });
});
