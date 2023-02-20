<template lang="pug">
  div.alarm-column-value-categories
    categories-list(v-if="asList", :categories="filteredCategories", :limit="limit")
    v-menu(
      v-else,
      :disabled="isDisabled",
      @input="$emit('activate', $event)",
      @click.native.stop=""
    )
      template(#activator="{ on }")
        v-btn.ma-0.alarm-column-value-categories__button(
          v-on="on",
          :disabled="isDisabled",
          depressed,
          small
        ) {{ $tc('common.link', 2) }}
      v-list(dark, dense)
        categories-list(:categories="filteredCategories")
</template>

<script>
import { ALARM_LIST_ACTIONS_TYPES, BUSINESS_USER_PERMISSIONS_ACTIONS_MAP } from '@/constants';

import { harmonizeCategories } from '@/helpers/links';

import { authMixin } from '@/mixins/auth';

import CategoriesList from './category-links/categories-list.vue';

export default {
  components: { CategoriesList },
  mixins: [authMixin],
  props: {
    links: {
      type: Object,
      default: () => ({}),
    },
    asList: {
      type: Boolean,
      default: false,
    },
    limit: {
      type: Number,
      required: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    filteredCategories() {
      return harmonizeCategories(this.links, this.checkAccessForSpecialCategory);
    },

    /**
     * If there are no links, in every category
     */
    isDisabled() {
      return Object.values(this.links).every(element => !element.length);
    },

    /**
     * Check if user has access to all links/categories
     */
    hasAccessToLinks() {
      return this.checkAccess(BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links]);
    },
  },
  methods: {
    checkAccessForSpecialCategory(category) {
      const permissionPrefix = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList[ALARM_LIST_ACTIONS_TYPES.links];
      const permission = `${permissionPrefix}_${category}`;

      return this.hasAccessToLinks || this.checkAccess(permission);
    },
  },
};
</script>

<style lang="scss">
.alarm-column-value-categories {
  &__button {
    height: 20px;
  }
}
</style>
