<template lang="pug">
  modal-wrapper(:modal="modal")
    component(:is="modal.name", :modal="modal")
</template>

<script>
import ModalWrapper from './modal-wrapper.vue';

export default {
  name: 'modal-base',
  components: { ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
};
</script>
