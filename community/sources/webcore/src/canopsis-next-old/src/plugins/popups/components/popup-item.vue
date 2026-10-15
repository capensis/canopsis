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
      :context="context"
      :template-props="context"
      class="pre-line"
    />
  </v-alert>
</template>

<script>
import {
  ref,
  computed,
  watch,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { VUETIFY_ANIMATION_DELAY, POPUP_TICK_DELAY } from '@/config';

import { usePopups } from '@/hooks/popups';

/**
 * Popup component
 *
 * @prop {String} [id] - Id of the popup
 * @prop {String} [type] - Type of the popup (info, error, ...)
 * @prop {String} [text] - Text displayed in the popup
 * @prop {Object} [context] - Context object for template compilation
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
    context: {
      type: Object,
      default: () => ({}),
    },
    autoClose: {
      type: [Number, Boolean],
      required: true,
    },
  },
  setup(props) {
    const animationTimeout = ref(null);
    const closeInterval = ref(null);
    const closeValue = ref(props.autoClose);
    const isPaused = ref(false);
    const visible = ref(false);

    const popups = usePopups();

    /**
     * Decrements the auto-close countdown timer by POPUP_TICK_DELAY milliseconds.
     */
    const progressTick = () => {
      if (closeValue.value <= 0) {
        visible.value = false;
      } else {
        closeValue.value -= POPUP_TICK_DELAY;
      }
    };

    /**
     * Starts or resumes the auto-close countdown timer for the popup.
     */
    const playProgress = () => {
      closeInterval.value = setInterval(progressTick, POPUP_TICK_DELAY);
      isPaused.value = false;
    };

    /**
     * Completely stops the auto-close timer and resets the countdown to its original value.
     */
    const stopProgress = () => {
      clearInterval(closeInterval.value);
      closeInterval.value = undefined;
      closeValue.value = props.autoClose;
    };

    /**
     * Temporarily pauses the auto-close countdown timer without resetting the value.
     */
    const pauseProgress = () => {
      clearInterval(closeInterval.value);
      isPaused.value = true;
    };

    /**
     * Initiates the popup removal process with a delay to allow exit animations.
     */
    const removeWithTimeout = () => {
      stopProgress();
      animationTimeout.value = setTimeout(() => popups.remove({ id: props.id }), VUETIFY_ANIMATION_DELAY);
    };

    const progressLineStyle = computed(() => ({
      animationDuration: `${props.autoClose / 1000}s`,
    }));

    const progressLineClass = computed(() => ({
      'progress-line--active': visible.value,
      'progress-line--paused': isPaused.value,
    }));

    const alertListeners = computed(() => {
      if (props.autoClose) {
        return {
          mouseover: pauseProgress,
          mouseout: playProgress,
        };
      }

      return {};
    });

    watch(visible, (value) => {
      if (!value) {
        removeWithTimeout();
      }
    });

    onMounted(() => {
      visible.value = true;

      if (props.autoClose) {
        playProgress();
      }
    });

    onBeforeUnmount(() => {
      clearInterval(closeInterval.value);
      clearTimeout(animationTimeout.value);
    });

    return {
      visible,
      progressLineStyle,
      progressLineClass,
      alertListeners,
    };
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
