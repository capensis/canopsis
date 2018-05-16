<template lang="pug">
  v-dialog(v-model="isOpened", v-bind="dialogProps")
    slot(v-if="isActive")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('modal');

export default {
  props: {
    name: {
      type: String,
      required: true,
    },
    dialogProps: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    ...mapGetters({
      modalName: 'name',
      modalAnimationPending: 'animationPending',
    }),
    isActive() {
      return this.modalName === this.name;
    },
    isOpened: {
      get() {
        return this.modalName === this.name && !this.modalAnimationPending;
      },
      set() {
        this.hideModal();
      },
    },
  },
  methods: {
    ...mapActions({
      hideModal: 'hide',
    }),
  },
};
</script>
