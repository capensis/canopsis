<template>
  <patterns-list
    :options="options"
    :patterns="patterns"
    :total-items="patternsMeta.total_count"
    :pending="patternsPending"
    @update:options="updateOptions"
    @edit="$emit('edit', $event)"
    @remove="$emit('remove', $event)"
    @remove-selected="$emit('remove-selected', $event)"
  />
</template>

<script>
import { onMounted } from 'vue';

import { convertQueryToRequest } from '@/helpers/query';

import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { usePattern } from '@/hooks/store/modules/pattern';

import PatternsList from './patterns-list.vue';

export default {
  components: {
    PatternsList,
  },

  emits: ['edit', 'remove', 'remove-selected'],

  setup() {
    const { patternsMeta, patternsPending, patterns, fetchPatternsList } = usePattern();

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useLocalQueryWithOptions({
      onUpdate: fetchQuery =>
        fetchPatternsList({
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
