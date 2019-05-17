<template lang="pug">
  div.mt-1
    div(v-if="hasAccessToLinks")
      div(v-for="category in linkList", :key="category.cat_name")
        template(v-if="category.links.length && hasAccessToCategory(category.cat_name)")
          span.category.mr-2 {{ category.cat_name }}
          v-divider(light)
          div(v-for="(link, index) in category.links", :key="`links-${index}`")
            div.pa-2.text-xs-right
              a(:href="link.link", target="_blank") {{ link.label }}
</template>

<script>
import { BUSINESS_USER_RIGHTS_ACTIONS_MAP, WIDGETS_ACTIONS_TYPES } from '@/constants';

import linksMixin from '@/mixins/links';
import authMixin from '@/mixins/auth';

export default {
  mixins: [linksMixin, authMixin],
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
      return this.category ?
        this.links.filter(({ cat_name: catName }) => catName === this.category) :
        this.links;
    },

    linkList() {
      return this.filteredLinks.map((category) => {
        const categoryLinks = this.harmonizeLinks(category.links, category.cat_name);

        return {
          cat_name: category.cat_name,
          links: categoryLinks,
        };
      });
    },

    hasAccessToLinks() {
      return this.checkAccess(BUSINESS_USER_RIGHTS_ACTIONS_MAP.weather[WIDGETS_ACTIONS_TYPES.weather.entityLinks]);
    },

    hasAccessToCategory() {
      return category => this.checkAccess(`${
        BUSINESS_USER_RIGHTS_ACTIONS_MAP
          .weather[WIDGETS_ACTIONS_TYPES.weather.entityLinks]}_${category}`);
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
