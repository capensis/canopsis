<template lang="pug">
  div
    modal-base(:modals="modals")
</template>

<script>

/**
 * Wrapper for all modal windows
 */
export default {
  computed: {
    modals() {
      return this.$store.getters[`${this.$modals.moduleName}/modals`];
    },
  },
  watch: {
    $route: {
      handler() {
        if (this.modals && this.modals.length) {
          this.modals.map(modal => this.$modals.hide({ id: modal.id }));
        }
      },
      deep: true,
    },
  },
};
</script>
