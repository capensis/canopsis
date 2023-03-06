<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $t('common.description') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(column)
        v-flex(xs12)
          v-card
            v-card-text
              pre {{ pbehaviorType.description }}
    v-tab {{ $t('common.type') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(column)
        v-flex(xs12)
          v-card
            v-card-text
              pre {{ pbehaviorType.type }}
</template>

<script>
export default {
  props: {
    pbehaviorType: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
