<template lang="pug">
  v-data-iterator(:items="exceptions", hide-actions)
    template(#item="{ item }")
      v-flex
        v-card
          v-card-title {{ item.name }}
        v-divider
    template(#no-data="")
      v-flex
        v-card
          v-card-title {{ $t('common.noData') }}
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
