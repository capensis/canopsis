<template lang="pug">
  modal-wrapper(
    v-if="modal",
    :key="modal.id",
    :modal="modal"
  )
    component(:is="modal.name", :modal="modal")
    modal-base(:modals="modals", :index="index + 1")
</template>

<script>
import ModalWrapper from './modal-wrapper.vue';

export default {
  name: 'modal-base',
  components: { ModalWrapper },
  provide() {
    return {
      $clickOutside: this.$clickOutside,
    };
  },
  props: {
    modals: {
      type: Array,
      default: () => [],
    },
    index: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    modal() {
      return this.modals[this.index];
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
        return this.handlers.reduce((acc, handler) => handler(...args), true);
      },
    };
  },
};
</script>
