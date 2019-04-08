<template lang="pug">
  v-dialog(v-model="isOpen", v-bind="dialogProps")
    // @slot use this slot default
    slot
</template>

<script>
import modalMixin from '@/mixins/modal';

/**
 * Wrapper for each modal window
 *
 * @prop {Object} modal - The current modal object
 * @prop {number} index - The current modal index in the store
 * @prop {Object} [dialogProps={}] - Properties for vuetify v-dialog
 */
export default {
  mixins: [modalMixin],
  props: {
    modal: {
      type: Object,
      required: true,
    },
    index: {
      type: Number,
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
        this.hideModal({ id: this.modal.id });
      },
    },
  },
  mounted() {
    this.ready = true;
  },
};
</script>
