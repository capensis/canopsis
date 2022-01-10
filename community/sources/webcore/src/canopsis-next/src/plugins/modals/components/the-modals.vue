<template lang="pug">
  div
    modal-base(
      v-for="modal in modals",
      :key="modal.id",
      :modal="modal"
    )
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
          this.modals.forEach(modal => !modal.minimized && this.$modals.hide({ id: modal.id }));
        }
      },
      deep: true,
    },
  },
};
</script>
