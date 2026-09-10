<template>
  <v-fade-transition>
    <c-alert
      v-if="visible"
      :type="type"
      class="alert-without-border"
      closable
      @click:close="visible = false"
      @mouseenter="pauseProgress"
      @mouseleave="playProgress"
    >
      <div v-if="hasAutoClose" class="progress">
        <div :style="progressLineStyle" :class="progressLineClass" class="progress-line" />
      </div>

      <div class="pre-line">
        <c-compiled-template
          :template="text"
          :context="context"
          :template-props="context"
        />
      </div>
    </c-alert>
  </v-fade-transition>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import { VUETIFY_ANIMATION_DELAY, POPUP_TICK_DELAY } from '@/config';

import { usePopupsStore } from '../stores';

const props = defineProps({
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
});

const popupsStore = usePopupsStore();

const animationTimeout = ref();
const closeInterval = ref();
const closeValue = ref(Number(props.autoClose) || 0);
const isPaused = ref(false);
const visible = ref(false);

const hasAutoClose = computed(() => !!props.autoClose);

const progressLineStyle = computed(() => ({
  animationDuration: `${Number(props.autoClose) / 1000}s`,
}));

const progressLineClass = computed(() => ({
  'progress-line--active': visible.value,
  'progress-line--paused': isPaused.value,
}));

/**
 * Decrements the remaining auto-close time.
 * Closes the popup when the countdown reaches zero.
 */
const progressTick = () => {
  if (closeValue.value <= 0) {
    visible.value = false;
    return;
  }

  closeValue.value -= POPUP_TICK_DELAY;
};

/**
 * Starts or resumes the auto-close timer.
 */
const playProgress = () => {
  if (!hasAutoClose.value) {
    return;
  }

  clearInterval(closeInterval.value);

  closeInterval.value = setInterval(progressTick, POPUP_TICK_DELAY);
  isPaused.value = false;
};

/**
 * Stops the auto-close timer and resets the countdown value.
 */
const stopProgress = () => {
  clearInterval(closeInterval.value);
  closeInterval.value = undefined;
  closeValue.value = Number(props.autoClose) || 0;
};

/**
 * Pauses the auto-close timer without resetting the countdown value.
 */
const pauseProgress = () => {
  clearInterval(closeInterval.value);
  isPaused.value = true;
};

/**
 * Removes the popup after animation delay
 * to allow the close transition to finish.
 */
const removeWithTimeout = () => {
  stopProgress();

  animationTimeout.value = setTimeout(() => {
    popupsStore.remove({ id: props.id });
  }, VUETIFY_ANIMATION_DELAY);
};

watch(visible, (value) => {
  if (!value) {
    removeWithTimeout();
  }
});

onMounted(() => {
  visible.value = true;

  if (hasAutoClose.value) {
    playProgress();
  }
});

onBeforeUnmount(() => {
  clearInterval(closeInterval.value);
  clearTimeout(animationTimeout.value);
});
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

:deep(.alert-without-border.v-alert) {
  border: 0;
  margin-left: 0;
  margin-right: 0;
  width: fit-content;
  max-width: 100%;
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
