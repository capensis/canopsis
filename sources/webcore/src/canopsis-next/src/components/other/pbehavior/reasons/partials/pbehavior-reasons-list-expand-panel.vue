<template lang="pug">
  v-tabs(color="secondary lighten-1", dark, centered, slider-color="primary")
    v-tab {{ $t('common.description') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(column)
        v-layout(row)
          v-text-field.theme--dark(
            :value="pbehaviorReason.description",
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
