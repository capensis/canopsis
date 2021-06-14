<template lang="pug">
  v-data-iterator(:items="exceptions", hide-actions)
    v-flex(slot="item", slot-scope="props")
      v-card
        v-card-title {{ props.item.name }}
      v-divider
    v-flex(slot="no-data")
      v-card
        v-card-title {{ $t('tables.noData') }}
</template>

<script>
export default {
  props: {
    exceptions: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
