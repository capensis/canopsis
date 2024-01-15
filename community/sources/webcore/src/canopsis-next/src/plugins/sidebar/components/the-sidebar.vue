<template>
  <div>
    <sidebar-base :sidebar="sidebar" />
  </div>
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
