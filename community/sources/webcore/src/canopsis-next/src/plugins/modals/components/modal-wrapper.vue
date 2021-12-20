<template lang="pug">
  v-dialog(v-model="isOpen", v-bind="dialogProps")
    // @slot use this slot default
    slot
</template>

<script>
/**
 * Wrapper for each modal window
 *
 * @prop {Object} modal - The current modal object
 * @prop {Object} [dialogProps={}] - Properties for vuetify v-dialog
 */
export default {
  inject: ['$clickOutside'],
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      ready: false,
    };
  },
  computed: {
    isOpen: {
      get() {
        return !this.modal.hidden && this.ready;
      },
      set() {
        this.$modals.hide({ id: this.modal.id });
      },
    },
    dialogProps() {
      const defaultDialogProps = { maxWidth: 700, lazy: true };
      const { dialogPropsMap = {} } = this.$modals;
      const { name, dialogProps, minimized } = this.modal;

      const props = {
        ...defaultDialogProps,
        ...dialogPropsMap[name],
        ...dialogProps,

        customCloseConditional: (...args) => this.$clickOutside.call(...args),
      };

      return {
        ...props,

        hideOverlay: props.hideOverlay || minimized,
        ignoreClickOutside: props.ignoreClickOutside || minimized,
        contentClass: minimized ? `v-dialog--minimized ${props.contentClass}` : props.contentClass,
      };
    },
  },
  mounted() {
    this.ready = true;
  },
};
</script>

<style lang="scss">
$minimizedDialogMaxWidth: 360px;

.v-dialog {
  .v-card__title {
    .headline {
      word-break: break-word;
    }
  }

  &.v-dialog--minimized {
    position: fixed;
    bottom: 0;
    max-width: $minimizedDialogMaxWidth !important;
    margin-bottom: 0 !important;
    transition: all .1s linear;
    top: auto;
    left: auto;
    right: auto;
    height: auto;

    .v-card__title {
      padding: 0 10px;
      transition: all .1s linear;

      .headline {
        font-size: 16px !important;
      }
    }

    .v-card.fill-min-height {
      min-height: auto;
    }
  }
}
</style>
