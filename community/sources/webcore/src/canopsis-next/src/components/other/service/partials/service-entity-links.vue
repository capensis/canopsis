<template lang="pug">
  div.mt-1(v-if="hasAccessToLinks")
    div(v-for="(categoryLinks, category) in preparedLinks", :key="category")
      span.category.mr-2 {{ category }}
      v-divider(light)
      div(v-for="(link, index) in categoryLinks", :key="`links-${index}`")
        div.pa-2.text-xs-right
          a(:href="link.url", target="_blank") {{ link.label }}
</template>

<script>
import { BUSINESS_USER_PERMISSIONS_ACTIONS_MAP, WEATHER_ACTIONS_TYPES } from '@/constants';

import { harmonizeCategoryLinks, harmonizeCategoriesLinks } from '@/helpers/links';

import { authMixin } from '@/mixins/auth';

export default {
  mixins: [authMixin],
  props: {
    links: {
      type: Object,
      default: () => ({}),
    },
    category: {
      type: String,
      default: null,
    },
  },
  computed: {
    hasAccessToLinks() {
      return this.checkAccess(BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.weather[WEATHER_ACTIONS_TYPES.entityLinks]);
    },

    preparedLinks() {
      return this.category
        ? { [this.category]: harmonizeCategoryLinks(this.links, this.category) }
        : harmonizeCategoriesLinks(this.links);
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
