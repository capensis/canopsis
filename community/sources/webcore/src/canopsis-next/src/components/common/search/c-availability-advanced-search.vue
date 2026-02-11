<template>
  <c-advanced-search
    :searches="searches"
    :attributes="attributes"
    allow-or
    with-history
    v-on="$listeners"
  />
</template>

<script>
import { useAvailabilityAdvancedSearchAttributes } from './hooks/advanced-search';

export default {
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    const { attributes } = useAvailabilityAdvancedSearchAttributes();

    return {
      attributes,
    };
  },
};
</script>
