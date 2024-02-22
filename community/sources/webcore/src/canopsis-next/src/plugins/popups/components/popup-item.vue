<template>
  <v-alert
    v-model="visible"
    :type="type"
    class="alert-without-border"
    transition="fade-transition"
    dismissible
    v-on="alertListeners"
  >
    <div
      v-if="autoClose"
      class="progress"
    >
      <div
        :style="progressLineStyle"
        :class="progressLineClass"
        class="progress-line"
      />
    </div>
    <c-compiled-template
      :template="text"
      class="pre-line"
    />
  </v-alert>
</template>

<script>
import { VUETIFY_ANIMATION_DELAY, POPUP_TICK_DELAY } from '@/config';

/**
 * Popup component
 *
 * @prop {String} [id] - Id of the popup
 * @prop {String} [type] - Type of the popup (info, error, ...)
 * @prop {String} [text] - Text displayed in the popup
 * @prop {Number,Boolean} [autoClose] - Auto close delay
 */
export default {
  props: {
    id: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      default: 'error',
    },
    text: {
      type: String,
      default: '',
    },
    autoClose: {
      type: [Number, Boolean],
      required: true,
    },
  },
  data() {
    return {
      animationTimeout: null,
      closeInterval: null,
      closeValue: this.autoClose,
      isPaused: false,
      visible: false,
    };
  },
  computed: {
    progressLineStyle() {
      return { animationDuration: `${this.autoClose / 1000}s` };
    },
    progressLineClass() {
      return {
        'progress-line--active': this.visible,
        'progress-line--paused': this.isPaused,
      };
    },
    alertListeners() {
      if (this.autoClose) {
        return {
          mouseover: this.pauseProgress,
          mouseout: this.playProgress,
        };
      }

      return {};
    },
  },
  watch: {
    visible(value) {
      if (!value) {
        this.removeWithTimeout();
      }
    },
  },
  mounted() {
    this.visible = true;

    if (this.autoClose) {
      this.playProgress();
    }
  },
  beforeDestroy() {
    clearInterval(this.closeInterval);
    clearTimeout(this.animationTimeout);
  },
  methods: {
    playProgress() {
      this.closeInterval = setInterval(this.progressTick, POPUP_TICK_DELAY);
      this.isPaused = false;
    },

    stopProgress() {
      clearInterval(this.closeInterval);
      this.closeInterval = undefined;
      this.closeValue = this.autoClose;
    },

    pauseProgress() {
      clearInterval(this.closeInterval);
      this.isPaused = true;
    },

    progressTick() {
      if (this.closeValue <= 0) {
        this.visible = false;
      } else {
        this.closeValue -= POPUP_TICK_DELAY;
      }
    },

    removeWithTimeout() {
      this.stopProgress();
      this.animationTimeout = setTimeout(() => this.$popups.remove({ id: this.id }), VUETIFY_ANIMATION_DELAY);
    },
  },
};
</script>

<style lang="scss" scoped>
  @keyframes progress {
    from {
      width: 0;
    }

    to {
      width: 100%;
    }
  }

  .alert-without-border {
    border: 0;

    &.v-alert {
      margin-left: 0;
      margin-right: 0;
    }
  }

  .progress {
    height: 5px;
    position: absolute;
    width: 100%;
    top: 0;
    margin: 0;
    padding: 0;
    left: 0;

    &-line {
      animation-play-state: paused;
      animation: progress linear;
      display: block;
      height: 100%;
      background: black;
      opacity: 0.2;

      &--active {
        animation-play-state: running;
      }

      &--paused {
        animation-play-state: paused;
      }
    }
  }
</style>
