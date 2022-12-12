<template lang="pug">
  v-list.pa-0(:class="{ 'transparent': limit }", :dark="!limit", dense)
    category-links(
      v-for="category in categoriesLinks",
      :links="category.links",
      :label="category.label",
      :key="category.label"
    )
    v-menu(v-if="dropDownShown", full-width)
      template(#activator="{ on }")
        v-btn.ma-0(v-on="on", small, block, flat) ...
      categories-list(:categories="dropDownCategories")
</template>

<script>
import CategoryLinks from './category-links.vue';

export default {
  name: 'CategoriesList',
  components: { CategoryLinks },
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
