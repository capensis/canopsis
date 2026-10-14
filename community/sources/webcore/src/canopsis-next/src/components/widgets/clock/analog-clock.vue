<template>
  <time :datetime="dateTime" :aria-label="accessibleTime" class="analog-clock">
    <svg
      :aria-label="accessibleTime"
      class="analog-clock__face"
      viewBox="0 0 200 200"
      role="img"
    >
      <circle
        class="analog-clock__outline"
        cx="100"
        cy="100"
        r="92"
      />

      <line
        v-for="tick in ticks"
        :key="tick"
        :class="{ 'analog-clock__tick--hour': tick % 5 === 0 }"
        :transform="`rotate(${tick * 6} 100 100)`"
        :y2="tick % 5 === 0 ? 22 : 17"
        class="analog-clock__tick"
        x1="100"
        y1="10"
        x2="100"
      />

      <text x="100" y="37">12</text>
      <text x="165" y="105">3</text>
      <text x="100" y="173">6</text>
      <text x="35" y="105">9</text>

      <line
        :transform="`rotate(${hourAngle} 100 100)`"
        class="analog-clock__hand analog-clock__hand--hour"
        x1="100"
        y1="104"
        x2="100"
        y2="58"
      />
      <line
        :transform="`rotate(${minuteAngle} 100 100)`"
        class="analog-clock__hand analog-clock__hand--minute"
        x1="100"
        y1="106"
        x2="100"
        y2="34"
      />
      <line
        :transform="`rotate(${secondAngle} 100 100)`"
        class="analog-clock__hand analog-clock__hand--second"
        x1="100"
        y1="112"
        x2="100"
        y2="26"
      />
      <circle
        class="analog-clock__center"
        cx="100"
        cy="100"
        r="4"
      />
    </svg>
  </time>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  date: {
    type: Object,
    required: true,
  },
});

const ticks = Array.from({ length: 60 }, (_, index) => index);
const seconds = computed(() => props.date.seconds() + props.date.milliseconds() / 1000);
const minutes = computed(() => props.date.minutes() + seconds.value / 60);
const hours = computed(() => props.date.hours() % 12 + minutes.value / 60);
const secondAngle = computed(() => seconds.value * 6);
const minuteAngle = computed(() => minutes.value * 6);
const hourAngle = computed(() => hours.value * 30);
const accessibleTime = computed(() => props.date.format('HH:mm:ss'));
const dateTime = computed(() => props.date.format());
</script>

<style lang="scss" scoped>
.analog-clock {
  display: flex;
  justify-content: center;
  color: inherit;

  &__face {
    width: var(--clock-face-size);
    height: var(--clock-face-size);
    overflow: visible;
  }

  &__outline,
  &__tick,
  &__hand {
    vector-effect: non-scaling-stroke;
  }

  &__outline {
    fill: transparent;
    stroke: currentColor;
    stroke-width: 2;
  }

  &__tick {
    stroke: currentColor;
    stroke-width: 1;

    &--hour {
      stroke-width: 2;
    }
  }

  text {
    fill: currentColor;
    font-size: 12px;
    text-anchor: middle;
  }

  &__hand {
    stroke: currentColor;
    stroke-linecap: round;

    &--hour {
      stroke-width: 5;
    }

    &--minute {
      stroke-width: 4;
    }

    &--second {
      opacity: 0.55;
      stroke-width: 1.5;
    }
  }

  &__center {
    fill: currentColor;
  }
}
</style>
