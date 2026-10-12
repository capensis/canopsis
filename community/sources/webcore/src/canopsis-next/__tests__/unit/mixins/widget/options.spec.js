import { describe, expect, it } from 'vitest';

import { widgetOptionsMixin } from '@/mixins/widget/options';

describe('widgetOptionsMixin', () => {
  it('preserves descending sort from canonical table options', () => {
    const context = {
      query: {
        page: 1,
        itemsPerPage: 10,
        sortBy: [],
        sortDesc: [],
      },
    };

    widgetOptionsMixin.computed.options.set.call(context, {
      page: 1,
      itemsPerPage: 10,
      sortBy: ['name'],
      sortDesc: [true],
    });

    expect(context.query.sortBy).toEqual(['name']);
    expect(context.query.sortDesc).toEqual([true]);
  });

  it('converts Vuetify sort options to canonical query options', () => {
    const context = {
      query: {
        page: 1,
        itemsPerPage: 10,
        sortBy: [],
        sortDesc: [],
      },
    };

    widgetOptionsMixin.computed.options.set.call(context, {
      page: 1,
      itemsPerPage: 10,
      sortBy: [{ key: 'name', order: 'desc' }],
    });

    expect(context.query.sortBy).toEqual(['name']);
    expect(context.query.sortDesc).toEqual([true]);
  });
});
