<template>
  <patterns-list
    :options="options"
    :patterns="corporatePatterns"
    :total-items="corporatePatternsMeta.total_count"
    :pending="corporatePatternsPending"
    corporate
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
import { useCorporatePattern } from '@/hooks/store/modules/pattern-corporate';

import PatternsList from './patterns-list.vue';

export default {
  components: {
    PatternsList,
  },

  emits: ['edit', 'remove', 'remove-selected'],

  setup() {
    const {
      corporatePatternsMeta,
      corporatePatternsPending,
      corporatePatterns,
      fetchCorporatePatternsList,
    } = useCorporatePattern();

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useLocalQueryWithOptions({
      onUpdate: fetchQuery =>
        fetchCorporatePatternsList({
          params: convertQueryToRequest(fetchQuery),
        }),
    });

    onMounted(fetchList);

    return {
      corporatePatternsMeta,
      corporatePatternsPending,
      corporatePatterns,
      options,
      updateOptions,
    };
  },
};
</script>
