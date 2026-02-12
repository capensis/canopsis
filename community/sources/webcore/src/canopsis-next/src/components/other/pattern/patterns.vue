<template>
  <patterns-list
    :options="options"
    :patterns="patterns"
    :total-items="patternsMeta.total_count"
    :pending="patternsPending"
    @update:options="updateOptions"
    @edit="$listeners.edit"
    @remove="$listeners.remove"
    @remove-selected="$listeners['remove-selected']"
  />
</template>

<script>
import { isEqual } from 'lodash';
import { onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { convertQueryToRequest } from '@/helpers/query';

import { usePattern } from '@/hooks/store/modules/pattern';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';

import PatternsList from './patterns-list.vue';

export default {
  components: {
    PatternsList,
  },
  setup() {
    const {
      patternsMeta,
      patternsPending,
      patterns,
      fetchPatternsList,
    } = usePattern();

    const initialQuery = {
      page: 1,
      itemsPerPage: PAGINATION_LIMIT,
      search: '',
      sortBy: [],
      sortDesc: [],
    };

    const fetchList = (fetchQuery = initialQuery) => fetchPatternsList({
      params: {
        ...convertQueryToRequest(fetchQuery),
        corporate: false,
      },
    });

    const {
      options,
      updateOptions,
    } = useLocalQueryWithOptions({
      initialQuery,
      comparator: isEqual,
      onUpdate: fetchList,
    });

    onMounted(fetchList);

    return {
      patternsMeta,
      patternsPending,
      patterns,
      options,
      updateOptions,
    };
  },
};
</script>
