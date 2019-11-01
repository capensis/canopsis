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
};
</script>
