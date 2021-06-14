<template lang="pug">
  v-list-tile.top-bar-menu-link(@click="")
    v-list-tile-title
      router-link(:to="link.route")
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
.top-bar-menu-link /deep/ a {
  text-decoration: none;
  color: inherit;
}
</style>
