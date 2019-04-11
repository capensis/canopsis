<template lang="pug">
  v-menu(@click.native.stop, :disabled="isDisabled")
    v-btn(slot="activator", :disabled="isDisabled", depressed, small) {{ $t('common.links') }}
    v-list(dark)
      v-list-tile-content(
      v-for="(category, key) in linkList",
      :key="key"
      )
        template(v-if="category.length")
          v-list-tile-title.px-2.font-weight-bold.category {{ key }}
          v-list-tile(
          v-for="(link, index) in category",
          :key="index"
          )
            alarm-column-value-link(:link="link")
</template>

<script>
import { cloneDeep } from 'lodash';

import AlarmColumnValueLink from './alarm-column-value-link.vue';

export default {
  components: { AlarmColumnValueLink },
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
      const links = cloneDeep(this.links);

      return Object.keys(links).reduce((acc, category) => {
        const categoryLinks = links[category].map((link, index) => {
          if (typeof link === 'object' && link.link && link.label) {
            return link;
          }

          return {
            label: `${category} - ${index}`,
            link,
          };
        });

        acc[category] = categoryLinks;

        return acc;
      }, {});
    },

    isDisabled() {
      return Object.values(this.links).every(element => !element.length);
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
