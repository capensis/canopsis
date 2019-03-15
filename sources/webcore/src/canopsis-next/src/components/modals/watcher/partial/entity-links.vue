<template lang="pug">
  div.mt-1
    div(v-for="(link, index) in filteredLinks", :key="`links-${index}`")
      div(v-for="(item, itemIndex) in link.links", :key="`links-item-${index}-${itemIndex}`")
        v-divider(light)
        div.pa-2.text-xs-right
          span.category.mr-2 {{ link.cat_name }} {{ itemIndex + 1 }}
          a(:href="item", target="_blank") {{ $t('common.link') }}
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

<style lang="scss" scoped>
  .category {
    display: inline-block;

    &:first-letter {
      text-transform: uppercase;
    }
  }
</style>
