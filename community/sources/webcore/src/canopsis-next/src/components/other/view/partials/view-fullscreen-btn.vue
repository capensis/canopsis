<template>
  <v-tooltip
    class="view-fullscreen-btn"
    :left="leftTooltip"
    :top="topTooltip"
  >
    <v-btn
      class="view-fullscreen-btn__button"
      slot="activator"
      :value="value"
      :small="small"
      fab="fab"
      dark="dark"
      @click.stop="toggleFullScreen"
    >
      <v-icon>fullscreen</v-icon>
      <v-icon>fullscreen_exit</v-icon>
    </v-btn>
    <div>
      {{ $t('view.fullScreen') }}
      <div class="font-italic text-caption ml-1">
        ({{ $t('view.fullScreenShortcut') }})
      </div>
    </div>
  </v-tooltip>
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    leftTooltip: {
      type: Boolean,
      default: false,
    },
    topTooltip: {
      type: Boolean,
      default: false,
    },
    toggleFullScreen: {
      type: Function,
      default: () => () => {},
    },
  },
};
</script>

<style lang="scss">
.view-fullscreen-btn {
  &__button  {
    border-color: #212121 !important;
    background-color: #212121 !important;

    .theme--dark & {
      border-color: #979797 !important;
      background-color: #979797 !important;
    }
  }
}
</style>
