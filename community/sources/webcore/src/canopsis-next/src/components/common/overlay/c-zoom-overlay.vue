<template>
  <div
    class="zoom-overlay__wrapper"
    @wheel="wheelListener"
  >
    <v-fade-transition>
      <div
        class="zoom-overlay"
        v-if="shown"
      >
        <span class="zoom-overlay__text">{{ $t('common.ctrlZoom') }}</span>
      </div>
    </v-fade-transition>
    <slot />
  </div>
</template>

<script>
import { ZOOM_OVERLAY_DELAY } from '@/config';

export default {
  props: {
    skipAlt: {
      type: Boolean,
      default: false,
    },
    skipShift: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      shown: false,
    };
  },
  methods: {
    wheelListener(event) {
      if (this.timer) {
        clearTimeout(this.timer);
      }

      if (event.ctrlKey || (this.skipShift && event.shiftKey) || (this.skipAlt && event.altKey)) {
        event.preventDefault();
        this.shown = false;

        return;
      }

      this.shown = true;
      this.timer = setTimeout(() => this.shown = false, ZOOM_OVERLAY_DELAY);
    },
  },
};
</script>

<style lang="scss">
.zoom-overlay {
  position: absolute;
  display: grid;
  top: 0;
  left: 0;
  height: 100%;
  width: 100%;
  z-index: 2;
  background: rgba(0, 0, 0, .5);
  align-content: center;
  justify-content: center;
  pointer-events: none;

  &__text {
    color: white;
    font-size: 2em;
    pointer-events: none;
  }

  &__wrapper {
    position: relative;
    width: 100%;
    height: 100%;
  }
}
</style>
