<template>
  <c-advanced-search
    :searches="searches"
    :attributes="attributes"
    :pending="pending"
    allow-or
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
    const { attributes, pending } = useEntityAdvancedSearchAttributes();

    return {
      attributes,
      pending,
    };
  },
};
</script>
