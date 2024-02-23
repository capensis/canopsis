<template>
  <v-layout
    class="panzoom"
    justify-center
    @mousewheel="handleScale"
    @mousemove="handleMove"
  >
    <div :style="containerStyles">
      <slot />
    </div>
    <v-layout
      class="panzoom__top-left-actions"
      column
    >
      <v-btn
        :disabled="scale === maxZoom"
        class="secondary ma-0 mb-1"
        icon
        dark
        @click="handleScaleIn"
      >
        <v-icon>add</v-icon>
      </v-btn>
      <v-btn
        :disabled="scale === minZoom"
        class="secondary ma-0"
        dark
        icon
        @click="handleScaleOut"
      >
        <v-icon>remove</v-icon>
      </v-btn>
    </v-layout>
    <v-layout
      v-if="helpText"
      class="panzoom__bottom-right-actions"
      column
    >
      <v-tooltip top>
        <template #activator="{ on }">
          <v-icon
            class="panzoom__help-icon"
            color="secondary"
            size="32"
            v-on="on"
          >
            help
          </v-icon>
        </template>
        <div
          v-html="helpText"
          class="pre-wrap"
        />
      </v-tooltip>
    </v-layout>
  </v-layout>
</template>

<script>
export default {
  props: {
    zoomStep: {
      type: Number,
      default: 0.1,
    },
    maxZoom: {
      type: Number,
      default: 4,
    },
    minZoom: {
      type: Number,
      default: 0.1,
      validator: value => value >= 0,
    },
    moveStep: {
      type: Number,
      default: 20,
    },
    helpText: {
      type: String,
      required: false,
    },
  },
  data() {
    return {
      scale: 1,
      translateX: 0,
      translateY: 0,
    };
  },
  computed: {
    containerStyles() {
      return {
        transform: `scale(${this.scale}) translateX(${this.translateX}px) translateY(${this.translateY}px)`,
      };
    },
  },
  methods: {
    reset() {
      this.scale = 1;
      this.translateX = 0;
      this.translateY = 0;
    },

    handleMove(event) {
      if (event.ctrlKey && event.buttons === 1) {
        this.translateX += event.movementX / this.scale;
        this.translateY += event.movementY / this.scale;

        event.preventDefault();
      }
    },

    handleScaleIn() {
      this.scale = Math.min(this.maxZoom, this.scale + this.zoomStep);
    },

    handleScaleOut() {
      this.scale = Math.max(this.minZoom, this.scale - this.zoomStep);
    },

    handleScale(event) {
      if (event.ctrlKey) {
        if (event.deltaY < 0) {
          this.handleScaleIn();
        } else {
          this.handleScaleOut();
        }

        event.preventDefault();
      }

      if (event.shiftKey) {
        if (event.deltaY < 0) {
          this.translateX += this.moveStep;
        } else {
          this.translateX -= this.moveStep;
        }

        event.preventDefault();
      }

      if (event.altKey) {
        if (event.deltaY < 0) {
          this.translateY += this.moveStep;
        } else {
          this.translateY -= this.moveStep;
        }

        event.preventDefault();
      }
    },
  },
};
</script>

<style lang="scss">
.panzoom {
  position: relative;
  overflow: hidden;

  &__top-left-actions {
    position: absolute;
    top: 10px;
    left: 10px;
  }

  &__bottom-right-actions {
    position: absolute;
    bottom: 10px;
    right: 10px;
  }

  &__help-icon {
    cursor: pointer;
  }
}
</style>
