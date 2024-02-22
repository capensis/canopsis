<template>
  <v-layout
    align-center
    justify-space-around
  >
    <span class="handler">
      <v-icon
        v-show="!disabled"
        :class="dragClass"
        class="draggable"
      >
        drag_indicator
      </v-icon>
    </span>
    <v-avatar
      :color="color"
      class="white--text"
      size="32"
    >
      <slot />
    </v-avatar>
  </v-layout>
</template>

<script>
export default {
  props: {
    dragClass: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'primary',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
  .handler {
    height: 24px !important;
    width: 18px !important;
  }
</style>
