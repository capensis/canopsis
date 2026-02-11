<template>
  <c-advanced-search
    :searches="searches"
    :attributes="attributes"
    with-history
    v-on="$listeners"
  />
</template>

<script>
import { useEntityAdvancedSearchAttributes } from './hooks/advanced-search';

export default {
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    const { attributes } = useEntityAdvancedSearchAttributes();

    return {
      attributes,
    };
  },
};
</script>
