<template lang="pug">
  v-layout.tab-panel-content(xs12)
    v-flex.tab-content-block.secondary.pa-2.white--text(xs4) {{ group.title }}
    v-flex.tab-content-block.secondary.lighten-1.pa-2.white--text(xs4) {{ view.title }}
    v-flex.tab-content-block.pa-2.white--text(xs4) {{ tab.title }}
</template>

<script>
export default {
  props: {
    tab: {
      type: Object,
      required: true,
    },
    view: {
      type: Object,
      required: true,
    },
    group: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.tab-panel-content {
  height: 100%;
}

.tab-content-block {
  display: flex;
  align-items: center;
}
</style>
