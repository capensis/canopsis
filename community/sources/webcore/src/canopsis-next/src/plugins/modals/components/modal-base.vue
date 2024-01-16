<template>
  <modal-wrapper :modal="modal">
    <component
      :is="modal.name"
      :modal="modal"
    />
  </modal-wrapper>
</template>

<script>
import ClickOutside from '@/services/click-outside';

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
    this.$clickOutside = new ClickOutside();
  },
};
</script>
