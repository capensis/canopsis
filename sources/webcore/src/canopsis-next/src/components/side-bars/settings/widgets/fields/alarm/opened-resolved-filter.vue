<template lang="pug">
  v-list-group(data-test="filterOnOpenResolved")
    v-list-tile(slot="activator") {{$t('settings.filterOnOpenResolved')}}
    v-container
      v-layout
        v-checkbox(
          v-field="value.opened",
          :label="$t('settings.open')",
          data-test="openFilter",
          color="primary",
          hide-details
        )
        v-checkbox(
          v-field="value.resolved",
          :label="$t('settings.resolved')",
          data-test="resolvedFilter",
          color="primary",
          hide-details
        )
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
