<template lang="pug">
  v-tooltip(:left="leftTooltip", :top="topTooltip")
    v-btn(
      slot="activator",
      :value="value",
      :small="small",
      fab,
      dark,
      @click.stop="toggleFullScreen"
    )
      v-icon fullscreen
      v-icon fullscreen_exit
    div {{ $t('view.fullScreen') }}
      div.font-italic.caption.ml-1 ({{ $t('view.fullScreenShortcut') }})
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
