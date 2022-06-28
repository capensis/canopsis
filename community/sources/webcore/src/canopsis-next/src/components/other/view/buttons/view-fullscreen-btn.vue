<template lang="pug">
  v-tooltip(:left="leftTooltip", :top="topTooltip")
    v-btn(
      slot="activator",
      v-model="fullscreen",
      :small="small",
      fab,
      dark,
      @click.stop="toggleFullScreenMode"
    )
      v-icon fullscreen
      v-icon fullscreen_exit
    div {{ $t('view.fullScreen') }}
      div.font-italic.caption.ml-1 ({{ $t('view.fullScreenShortcut') }})
</template>

<script>
export default {
  props: {
    activeTab: {
      type: Object,
      required: false,
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
  },
  data() {
    return {
      fullscreen: false,
    };
  },
  created() {
    document.addEventListener('keydown', this.keyDownListener);
  },
  beforeDestroy() {
    this.$fullscreen.exit();
    document.removeEventListener('keydown', this.keyDownListener);
  },
  methods: {
    toggleFullScreenMode() {
      if (!this.activeTab) {
        this.$popups.warning({ text: this.$t('view.errors.emptyTabs') });
        return;
      }

      const element = document.getElementById(`view-tab-${this.activeTab._id}`);

      if (!element) {
        return;
      }

      this.$fullscreen.toggle(element, {
        fullscreenClass: 'full-screen',
        background: 'white',
        callback: value => this.fullscreen = value,
      });
    },

    keyDownListener(event) {
      if (event.key === 'Enter' && event.altKey) {
        this.toggleFullScreenMode();
        event.preventDefault();
      }
    },
  },
};
</script>
