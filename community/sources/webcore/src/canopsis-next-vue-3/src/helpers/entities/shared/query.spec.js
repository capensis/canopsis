import { describe, expect, it } from 'vitest';

import { convertDataTableSortToQuery } from './query';

describe('convertDataTableSortToQuery', () => {
  it('converts Vuetify sort array to query sort', () => {
    expect(convertDataTableSortToQuery([{ key: 'name', order: 'desc' }])).toEqual({
      sortBy: ['name'],
      sortDesc: [true],
    });
  });

  it('converts canonical options object to query sort', () => {
    expect(convertDataTableSortToQuery({ sortBy: ['name'], sortDesc: [true] })).toEqual({
      sortBy: ['name'],
      sortDesc: [true],
    });
  });
});
