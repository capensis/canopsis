<template>
  <div
    class="modals-wrapper"
    :class="{ 'modals-wrapper--active': modals.length > 0 }"
  >
    <modal-base
      v-for="modal in modals"
      :key="modal.id"
      :modal="modal"
    />
  </div>
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

<style lang="scss" scoped>
$minimizedDialogMaxWidth: 360px;

.modals-wrapper {
  display: none;

  &--active {
    position: fixed;
    display: flex;
    bottom: 0;
    width: 100%;
    height: 100%;
    flex-wrap: wrap-reverse;
    pointer-events: none;
    justify-content: center;
    align-items: flex-end;
    align-content: flex-start;
    z-index: 300;

    & ::v-deep {
      .v-menu__content {
        pointer-events: auto;
      }

      .v-dialog__content {
        &--minimized {
          margin: 8px 8px 0 8px;
          position: relative;
          height: auto;
          max-width: $minimizedDialogMaxWidth;
          pointer-events: all;
          z-index: inherit !important;

          .v-dialog {
            margin: 0;
            box-shadow: none;
            transition: none;

            .v-card {
              &.fill-min-height {
                min-height: auto;
              }

              &__title {
                padding: 0 10px;

                .headline {
                  font-size: 16px !important;
                }
              }
            }
          }
        }
      }
    }
  }
}
</style>
