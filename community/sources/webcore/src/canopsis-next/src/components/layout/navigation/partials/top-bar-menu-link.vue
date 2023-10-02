<template lang="pug">
  v-list-tile.top-bar-menu-link(:to="link.route", active-class="", @click="handleClick")
    v-list-tile-title
      v-layout(justify-space-between)
        span {{ link.title }}
    v-list-tile-avatar.ml-2(size="24")
      v-icon(size="24") {{ link.icon }}
</template>

<script>
export default {
  props: {
    link: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    handleClick() {
      this.link.handler?.();
    },
  },
};
</script>

<style lang="scss" scoped>
.top-bar-menu-link ::v-deep a {
  text-decoration: none;
  color: inherit;

  .v-list__tile__avatar {
    min-width: unset;
  }
}
</style>
