<template lang="pug">
  modal-wrapper(:modal="modal")
    component(:is="modal.name", :modal="modal")
</template>

<script>
import ModalWrapper from './modal-wrapper.vue';

export default {
  provide() {
    return {
      $clickOutside: this.$clickOutside,
      $modal: this.modal,
    };
  },
  components: { ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  beforeCreate() {
    this.$clickOutside = {
      handlers: [],

      register(handler) {
        this.handlers.push(handler);
      },

      unregister(handler) {
        this.handlers = this.handlers.filter(h => h !== handler);
      },

      call(...args) {
        return this.handlers.reduce((acc, handler) => acc && handler(...args), true);
      },
    };
  },

};
</script>
