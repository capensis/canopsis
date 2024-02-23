<template>
  <v-tooltip
    :left="leftTooltip"
    :top="topTooltip"
    class="view-fullscreen-btn"
  >
    <template #activator="{ on }">
      <v-btn
        :value="value"
        class="view-fullscreen-btn__button"
        small
        fab
        dark
        v-on="on"
        @click.stop="toggleFullScreen"
      >
        <v-icon small>
          fullscreen
        </v-icon>
        <v-icon small>
          fullscreen_exit
        </v-icon>
      </v-btn>
    </template>
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
