<template lang="pug">
  div
    sidebar-base(:sidebar="sidebar")
</template>

<script>

/**
 * Wrapper for all sidebars
 */
export default {
  computed: {
    sidebar() {
      return this.$store.getters[`${this.$sidebar.moduleName}/sidebar`];
    },
  },
};
</script>
