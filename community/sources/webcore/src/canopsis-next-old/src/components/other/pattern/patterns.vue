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
import { onMounted } from 'vue';

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

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useLocalQueryWithOptions({
      onUpdate: fetchQuery => fetchPatternsList({
        params: {
          ...convertQueryToRequest(fetchQuery),
          corporate: false,
        },
      }),
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
