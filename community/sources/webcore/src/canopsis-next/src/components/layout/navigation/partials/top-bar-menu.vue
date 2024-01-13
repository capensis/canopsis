<template>
  <v-menu
    v-if="preparedLinks.length"
    bottom
    offset-y
  >
    <template #activator="{ on }">
      <v-btn
        class="white--text"
        v-on="on"
        text
      >
        {{ title }}
      </v-btn>
    </template>
    <v-list class="py-0">
      <top-bar-menu-link
        v-for="link in preparedLinks"
        :key="link.title"
        :link="link"
      />
    </v-list>
  </v-menu>
</template>

<script>
import { layoutNavigationTopBarMenuMixin } from '@/mixins/layout/navigation/top-bar-menu';

import TopBarMenuLink from './top-bar-menu-link.vue';

export default {
  components: { TopBarMenuLink },
  mixins: [layoutNavigationTopBarMenuMixin],
  props: {
    title: {
      type: String,
      default: '',
    },
    links: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    preparedLinks() {
      return this.prepareLinks(this.links);
    },
  },
};
</script>
