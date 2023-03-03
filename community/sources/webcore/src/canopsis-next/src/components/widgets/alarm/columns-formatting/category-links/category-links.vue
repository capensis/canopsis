<template lang="pug">
  v-list-tile-content
    v-list-tile-title.px-2.font-weight-bold.category(v-show="label") {{ label }}
    v-list-tile.category-list-tile.py-1(
      v-for="(link, index) in links",
      :key="index"
    )
      a(
        :href="link.link",
        target="_blank",
        @click.stop=""
      ) {{ link.label }}
</template>

<script>
export default {
  props: {
    links: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      required: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.category {
  display: inline-block;

  &-list-tile ::v-deep .v-list__tile {
    height: unset !important;
  }

  &:first-letter {
    text-transform: uppercase;
  }
}
</style>
