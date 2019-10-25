<template lang="pug">
  v-fade-transition
    v-layout.alert(v-show="value")
      div.overlay(:class="backgroundColor", :style="{ opacity: opacity }")
      div.content
        slot
          v-alert(type="error", :value="true") {{ errorMessage }}
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    opacity: {
      type: Number,
      default: 0.5,
    },
    backgroundColor: {
      type: String,
      default: 'white',
    },
    message: {
      type: String,
      default: '',
    },
  },

  computed: {
    errorMessage() {
      return this.message || this.$t('errors.default');
    },
  },
};
</script>

<style lang="scss" scoped>
  .alert {
    z-index: 2;

    &, .overlay {
      min-height: 100px;
      position: absolute;
      top: 0;
      left: 0;
      bottom: 0;
      right: 0;
    }

    .content {
      width: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }
</style>
