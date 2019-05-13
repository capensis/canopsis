<template lang="pug">
  v-menu(@click.native.stop, :disabled="isDisabled")
    v-btn(slot="activator", :disabled="isDisabled", depressed, small) {{ $t('common.links') }}
    v-list(dark, dense)
      v-list-tile-content(
      v-for="(categoryLinks, category) in linkList",
      :key="category"
      )
        template(v-if="hasAccessToCategory(category)")
          v-list-tile-title.px-2.font-weight-bold.category {{ category }}
          v-list-tile(
          v-for="(link, index) in categoryLinks",
          :key="index"
          )
            alarm-column-value-link(:link="link")
</template>

<script>
import { BUSINESS_USER_RIGHTS_ACTIONS_MAP, WIDGETS_ACTIONS_TYPES } from '@/constants';

import linksMixin from '@/mixins/links';
import authMixin from '@/mixins/auth';

import AlarmColumnValueLink from './alarm-column-value-link.vue';

export default {
  components: { AlarmColumnValueLink },
  mixins: [linksMixin, authMixin],
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
    linkList() {
      return Object.entries(this.links).reduce((acc, [category, categoryLinks]) => {
        if (categoryLinks.length) {
          acc[category] = this.harmonizeLinks(categoryLinks, category);
        }

        return acc;
      }, {});
    },

    isDisabled() {
      return Object.values(this.links).every(element => !element.length);
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
