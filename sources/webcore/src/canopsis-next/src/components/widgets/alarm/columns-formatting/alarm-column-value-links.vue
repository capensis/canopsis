<template lang="pug">
  v-menu(@click.native.stop, :disabled="isDisabled")
    v-btn(slot="activator", :disabled="isDisabled", depressed, small) {{ $t('common.links') }}
    v-list(dark, dense)
      v-list-tile-content(
        v-for="category in filteredCategories",
        :key="category.label"
      )
        v-list-tile-title.px-2.font-weight-bold.category {{ category.label }}
        v-list-tile(
          v-for="(link, index) in category.links",
          :key="index"
        )
          alarm-column-value-link(:link="link")
</template>

<script>
import { BUSINESS_USER_RIGHTS_ACTIONS_MAP, WIDGETS_ACTIONS_TYPES } from '@/constants';

import { harmonizeCategories } from '@/helpers/links';

import authMixin from '@/mixins/auth';

import AlarmColumnValueLink from './alarm-column-value-link.vue';

export default {
  components: { AlarmColumnValueLink },
  mixins: [authMixin],
  props: {
    links: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      isOpen: false,
    };
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
      return this.checkAccess(BUSINESS_USER_RIGHTS_ACTIONS_MAP.alarmsList[WIDGETS_ACTIONS_TYPES.alarmsList.links]);
    },
  },
  methods: {
    checkAccessForSpecialCategory(category) {
      const rightPrefix = BUSINESS_USER_RIGHTS_ACTIONS_MAP.alarmsList[WIDGETS_ACTIONS_TYPES.alarmsList.links];
      const right = `${rightPrefix}_${category}`;

      return this.hasAccessToLinks || this.checkAccess(right);
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
