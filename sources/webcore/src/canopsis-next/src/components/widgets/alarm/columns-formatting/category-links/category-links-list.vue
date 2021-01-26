<template lang="pug">
  v-list.pa-0(:class="{ 'category-links-list': limit }", :dark="!limit", dense)
    v-list-tile-content(
      v-for="category in categoriesLinks",
      :key="category.label"
    )
      v-list-tile-title.px-2.font-weight-bold.category {{ category.label }}
      v-list-tile.category-list-tile.py-1(
        v-for="(link, index) in category.links",
        :key="index"
      )
        category-link(:link="link")
    v-menu(v-if="dropDownShown", full-width)
      v-btn.ma-0(slot="activator", small, block, flat) ...
      category-links-list(:categories="dropDownCategories")
</template>

<script>
import CategoryLink from './category-link.vue';

export default {
  name: 'category-links-list',
  components: { CategoryLink },
  props: {
    categories: {
      type: Array,
      default: () => [],
    },
    limit: {
      type: Number,
      required: false,
    },
  },
  computed: {
    dropDownShown() {
      return this.limit && this.categories.length > this.limit;
    },

    categoriesLinks() {
      return this.limit
        ? this.categories.slice(0, this.limit)
        : this.categories;
    },

    dropDownCategories() {
      return this.categories.slice(this.limit);
    },
  },
};
</script>

<style lang="scss" scoped>
  .category {
    display: inline-block;

    &-links-list {
      background: transparent !important;
    }

    &-list-tile /deep/ .v-list__tile {
      height: unset !important;
    }

    &:first-letter {
      text-transform: uppercase;
    }
  }
</style>
