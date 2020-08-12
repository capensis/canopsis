<template lang="pug">
  v-layout.pa-3.secondary.lighten-1(column)
    v-layout(row)
      v-text-field.theme--dark(
        :value="pbehaviorReason.description",
        :label="$t('common.description')",
        disabled
      )
</template>

<script>
export default {
  props: {
    pbehaviorReason: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
