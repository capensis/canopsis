import { describe, expect, it } from 'vitest';
import { ref } from 'vue';

import { useQueryOptions } from '@/hooks/query/options';

describe('useQueryOptions', () => {
  it('preserves descending sort from canonical table options', () => {
    const query = ref({
      page: 1,
      itemsPerPage: 10,
      sortBy: [],
      sortDesc: [],
    });
    const updateQuery = (value) => {
      query.value = value;
    };

    const { updateOptions } = useQueryOptions(query, updateQuery);

    updateOptions({
      page: 1,
      itemsPerPage: 10,
      sortBy: ['name'],
      sortDesc: [true],
    });

    expect(query.value.sortBy).toEqual(['name']);
    expect(query.value.sortDesc).toEqual([true]);
  });

  it('converts Vuetify sort options to canonical query options', () => {
    const query = ref({
      page: 1,
      itemsPerPage: 10,
      sortBy: [],
      sortDesc: [],
    });
    const updateQuery = (value) => {
      query.value = value;
    };

    const { updateOptions } = useQueryOptions(query, updateQuery);

    updateOptions({
      page: 1,
      itemsPerPage: 10,
      sortBy: [{ key: 'name', order: 'desc' }],
    });

    expect(query.value.sortBy).toEqual(['name']);
    expect(query.value.sortDesc).toEqual([true]);
  });
});
