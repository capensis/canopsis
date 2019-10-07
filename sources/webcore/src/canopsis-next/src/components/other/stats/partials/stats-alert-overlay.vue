<template lang="pug">
  alert-overlay(:value="value")
    v-alert(:value="value", :type="type")
      v-layout(align-center)
        v-flex
          span {{ message }}
        v-flex(v-show="hasErrorMessage")
          v-tooltip.alert-tooltip(top)
            v-btn.ml-2.mr-0(slot="activator", icon, small)
              v-icon.alert-tooltip-icon help
            span {{ errorMessage }}
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
  .alert-tooltip .alert-tooltip-icon {
    color: rgba(0, 0, 0, .3) !important;
  }
</style>
