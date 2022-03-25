<template lang="pug">
  div
    v-tabs(v-model="activeTab", color="secondary lighten-1", dark, centered, slider-color="primary")
      v-tab {{ $tc('common.pattern', 2) }}
    v-layout.pa-3.secondary.lighten-2
      v-flex(xs12)
        v-card.pa-3
          v-tabs-items.pt-2(v-model="activeTab")
            v-tab-item(lazy)
              v-flex(offset-xs2, xs8)
                c-patterns-field(
                  :value="filter",
                  disabled,
                  with-entity
                )
</template>

<script>
export default {
  props: {
    filter: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      activeTab: 0,
    };
  },
};
</script>
