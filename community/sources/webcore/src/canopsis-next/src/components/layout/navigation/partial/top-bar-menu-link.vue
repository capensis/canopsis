<template lang="pug">
  v-list-tile.top-bar-menu-link(:to="link.route", active-class="")
    v-list-tile-title
      v-layout(justify-space-between)
        span {{ link.title }}
        v-icon.ml-2 {{ link.icon }}
</template>

<script>
export default {
  props: {
    link: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

<style lang="scss" scoped>
.top-bar-menu-link ::v-deep a {
  text-decoration: none;
  color: inherit;
}
</style>
