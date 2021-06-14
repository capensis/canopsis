<template lang="pug">
  v-menu(v-show="links.length", bottom, offset-y)
    v-btn.white--text(slot="activator", flat) {{ title }}
    v-list.py-0
      top-bar-menu-link(
        v-for="link in links",
        :key="link.title",
        :link="link"
      )
</template>

<script>
import TopBarMenuLink from './top-bar-menu-link.vue';

export default {
  components: { TopBarMenuLink },
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
};
</script>
