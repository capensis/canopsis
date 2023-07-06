<template lang="pug">
  v-layout(column)
    v-layout.secondary.lighten-2.tab-panel-item(align-center)
      slot(name="title")
        span.white--text {{ tab.title }}
    v-divider(dark)
</template>

<script>
export default {
  props: {
    tab: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
  .tab-panel-item {
    display: flex;
    align-items: center;
    height: 48px;
  }
</style>
