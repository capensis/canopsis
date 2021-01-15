<template lang="pug">
  div
    v-menu(v-if="!asList", @click.native.stop, :disabled="isDisabled")
      v-btn(slot="activator", :disabled="isDisabled", depressed, small) {{ $t('common.links') }}
      v-list(dark, dense)
        category-links-list(:categories="filteredCategories")
    category-links-list(v-else, :categories="filteredCategories", :limit="limit")
</template>

<script>
import { BUSINESS_USER_RIGHTS_ACTIONS_MAP, WIDGETS_ACTIONS_TYPES } from '@/constants';

import { harmonizeCategories } from '@/helpers/links';

import authMixin from '@/mixins/auth';

import CategoryLinksList from './category-links/category-links-list.vue';

export default {
  components: { CategoryLinksList },
  mixins: [authMixin],
  props: {
    links: {
      type: Object,
      default: () => ({
        Category1: [
          { label: 'Procédure', link: 'http://uneurl.local/?composant=feeder2_0&message=' },
          { label: 'Procédure2', link: 'http://uneurl.local/?composant=feeder2_0&message=' },
        ],
        Category2: [{ label: 'Procédure', link: 'http://uneurl.local/?composant=feeder2_0&message=' }],
        Category3: [{ label: 'Procédure', link: 'http://uneurl.local/?composant=feeder2_0&message=' }],
        Category4: [{ label: 'Procédure', link: 'http://uneurl.local/?composant=feeder2_0&message=' }],
        Category5: [{ label: 'Procédure', link: 'http://uneurl.local/?composant=feeder2_0&message=' }],
      }),
    },
    asList: {
      type: Boolean,
      default: false,
    },
    limit: {
      type: Number,
      required: false,
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
