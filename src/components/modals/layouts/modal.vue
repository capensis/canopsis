<template lang="pug">
  v-dialog(v-model="isOpened", v-bind="dialogProps")
    // @slot use this slot default
    slot(v-if="isActive")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('modal');

/**
 * Wrapper for each modal window
 *
 * @prop {string} name - Name of the modal
 * @prop {Object} [dialogProps={}] - Properties for vuetify v-dialog
 */
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
