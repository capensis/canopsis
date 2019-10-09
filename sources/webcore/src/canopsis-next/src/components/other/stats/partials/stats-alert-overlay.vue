<template lang="pug">
  alert-overlay(:value="value", :opacity="1")
    v-layout.alert-content(align-center, justify-center)
      v-tooltip(top, v-if="hasErrorMessage && !editionError")
        v-btn.ml-2.mr-0(slot="activator", icon, small)
          v-icon.error--text error
        span {{ errorMessage }}
      v-icon.error--text(v-if="editionError") error
      p.error--text.mb-0.ml-1 {{ message }}
</template>

<script>
import AlertOverlay from '@/components/layout/alert/alert-overlay.vue';

export default {
  components: {
    AlertOverlay,
  },
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    message: {
      type: String,
      default: null,
    },
    errorMessage: {
      type: String,
      default: null,
    },
    editionError: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      default: 'error',
    },
  },
  computed: {
    hasErrorMessage() {
      return Boolean(this.errorMessage);
    },
  },
};
</script>

<style>
  .alert-content {
    z-index: 2;
    opacity: 0.7;
  }
</style>
