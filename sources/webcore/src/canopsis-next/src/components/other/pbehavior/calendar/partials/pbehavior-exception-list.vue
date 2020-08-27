<template lang="pug">
  v-list.pa-0
    template(v-if="!exceptions || !exceptions.length")
      v-list-tile
        v-list-tile-content
          v-list-tile-title {{ $t('tables.noData') }}
      v-divider
    template(v-for="(exception) in exceptions")
      v-list-tile(:key="exception._id")
        v-list-tile-content
          v-list-tile-title {{ exception.name }}
      v-divider
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
