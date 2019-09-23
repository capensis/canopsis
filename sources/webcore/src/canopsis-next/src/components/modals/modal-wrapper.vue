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
  props: {
    modal: {
      type: Object,
      required: true,
    },
    dialogProps: {
      type: Object,
      default: () => ({}),
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
  },
  mounted() {
    this.ready = true;
  },
};
</script>
