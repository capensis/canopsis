<template lang="pug">
  div
    div(v-for="(link, index) in filteredLinks", :key="`links-${index}`")
      div(v-for="(item, itemIndex) in link.links", :key="`links-item-${index}-${itemIndex}`")
        span {{ link.cat_name }} {{ itemIndex + 1 }}
        a(:href="item", target="_blank") link
</template>

<script>
export default {
  props: {
    links: {
      type: Array,
      default: () => [],
    },
    category: {
      type: String,
      default: null,
    },
  },
  computed: {
    filteredLinks() {
      if (this.category) {
        return this.links.filter(({ cat_name: catName }) => catName === this.category);
      }

      return this.links;
    },
  },
};
</script>

