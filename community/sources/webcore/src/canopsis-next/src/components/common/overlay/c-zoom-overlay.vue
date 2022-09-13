<template lang="pug">
  div.zoom-overlay__wrapper(@wheel="throttledWheelListener")
    v-fade-transition
      div.zoom-overlay(v-if="shown")
        span.zoom-overlay__text Use ctrl + mouse wheel for zoom
    slot
</template>

<script>
import { throttle } from 'lodash';

export default {
  data() {
    return {
      shown: false,
    };
  },
  created() {
    this.throttledWheelListener = throttle(this.wheelListener, 100);
  },
  methods: {
    wheelListener({ ctrlKey }) {
      if (this.timer) {
        clearTimeout(this.timer);
      }

      if (ctrlKey) {
        this.shown = false;
        return;
      }

      this.shown = true;
      this.timer = setTimeout(() => this.shown = false, 2000); // TODO: move to constants
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
  z-index: 10;
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
