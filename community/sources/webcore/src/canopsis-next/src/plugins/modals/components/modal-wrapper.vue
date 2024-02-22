<template>
  <v-dialog
    v-model="isOpen"
    v-bind="dialogProps"
  >
    <!-- @slot use this slot default-->
    <slot />
  </v-dialog>
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
      const defaultDialogProps = {
        maxWidth: 700,
        attach: '.modals-wrapper',
        absolute: true,
        retainFocus: false,
      };
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
        contentWrapperClass: minimized ? 'v-dialog__content--minimized' : '',
      };
    },
  },
  mounted() {
    this.ready = true;
  },
};
</script>

<style lang="scss">
.v-dialog .v-card__title {
  .headline {
    word-break: break-word;
  }
}
</style>
