<template>
  <v-menu
    v-if="preparedLinks.length"
    v-bind="$attrs"
    bottom
    offset-y
  >
    <template #activator="{ on }">
      <v-btn
        class="white--text"
        text
        v-on="on"
      >
        <slot name="title">
          {{ title }}
        </slot>
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
import { computed } from 'vue';

import { useTopBarMenu } from './hooks/top-bar-menu';
import TopBarMenuLink from './top-bar-menu-link.vue';

export default {
  components: { TopBarMenuLink },
  inheritAttrs: false,
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
  setup(props) {
    const { prepareLinks } = useTopBarMenu();

    const preparedLinks = computed(() => prepareLinks(props.links));

    return {
      preparedLinks,
    };
  },
};
</script>
