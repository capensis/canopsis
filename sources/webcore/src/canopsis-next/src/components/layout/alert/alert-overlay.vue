<template lang="pug" functional>
  v-fade-transition
    v-layout.alert(v-show="props.value")
      div.overlay(:class="props.backgroundColor", :style="{ opacity: props.opacity }")
      div.content
        slot
          v-alert(type="error", :value="true") {{ parent.$t('errors.default') }}
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
  },
};
</script>

<style lang="scss" scoped>
  .alert {
    z-index: 2;

    &, .overlay {
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
