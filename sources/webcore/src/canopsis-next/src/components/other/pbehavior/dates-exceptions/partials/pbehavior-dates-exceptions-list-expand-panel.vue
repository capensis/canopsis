<template lang="pug">
  v-tabs(color="secondary lighten-1", dark, centered, slider-color="primary")
    v-tab {{ $t('common.description') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(column)
        v-flex(xs12)
          v-card
            v-card-text
              pre {{ pbehaviorDateException.description }}
</template>

<script>
export default {
  props: {
    pbehaviorDateException: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
