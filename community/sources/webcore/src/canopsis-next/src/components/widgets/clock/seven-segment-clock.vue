<template>
  <time :datetime="dateTime" :aria-label="accessibleTime" class="seven-segment-clock">
    <span class="seven-segment-clock__digits" aria-hidden="true">
      <template v-for="(token, index) in tokens" :key="`${token}-${index}`">
        <span v-if="token === ':'" class="seven-segment-clock__colon">
          <i />
          <i />
        </span>
        <span v-else class="seven-segment-clock__digit">
          <i
            v-for="segment in segments"
            :key="segment"
            :class="[
              `seven-segment-clock__segment--${segment}`,
              { 'seven-segment-clock__segment--active': activeSegments[token].includes(segment) },
            ]"
            class="seven-segment-clock__segment"
          />
        </span>
      </template>
    </span>
    <span v-if="period" class="seven-segment-clock__period">{{ period }}</span>
  </time>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  date: {
    type: Object,
    required: true,
  },
  use12HourFormat: {
    type: Boolean,
    default: false,
  },
  showSeconds: {
    type: Boolean,
    default: true,
  },
});

const segments = ['a', 'b', 'c', 'd', 'e', 'f', 'g'];
const activeSegments = {
  0: ['a', 'b', 'c', 'd', 'e', 'f'],
  1: ['b', 'c'],
  2: ['a', 'b', 'd', 'e', 'g'],
  3: ['a', 'b', 'c', 'd', 'g'],
  4: ['b', 'c', 'f', 'g'],
  5: ['a', 'c', 'd', 'f', 'g'],
  6: ['a', 'c', 'd', 'e', 'f', 'g'],
  7: ['a', 'b', 'c'],
  8: ['a', 'b', 'c', 'd', 'e', 'f', 'g'],
  9: ['a', 'b', 'c', 'd', 'f', 'g'],
};
const format = computed(() => {
  if (props.use12HourFormat) {
    return props.showSeconds ? 'hh:mm:ss' : 'hh:mm';
  }

  return props.showSeconds ? 'HH:mm:ss' : 'HH:mm';
});
const time = computed(() => props.date.format(format.value));
const tokens = computed(() => time.value.split(''));
const period = computed(() => props.use12HourFormat ? props.date.format('A') : '');
const accessibleTime = computed(() => [time.value, period.value].filter(Boolean).join(' '));
const dateTime = computed(() => props.date.format());
</script>

<style lang="scss" scoped>
.seven-segment-clock {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  color: inherit;
  white-space: nowrap;

  &__digits {
    display: flex;
    align-items: center;
    gap: calc(var(--clock-segment-width) * 0.16);
  }

  &__digit {
    position: relative;
    display: inline-block;
    width: var(--clock-segment-width);
    height: calc(var(--clock-segment-width) * 1.8);
  }

  &__segment {
    position: absolute;
    display: block;
    background: currentColor;
    opacity: 0.08;

    &--active {
      opacity: 1;
    }

    &--a,
    &--d,
    &--g {
      left: 14%;
      width: 72%;
      height: 9%;
      clip-path: polygon(8% 0, 92% 0, 100% 50%, 92% 100%, 8% 100%, 0 50%);
    }

    &--a {
      top: 0;
    }

    &--g {
      top: 45.5%;
    }

    &--d {
      bottom: 0;
    }

    &--b,
    &--c,
    &--e,
    &--f {
      width: 16%;
      height: 42%;
      clip-path: polygon(50% 0, 100% 10%, 100% 90%, 50% 100%, 0 90%, 0 10%);
    }

    &--b,
    &--c {
      right: 0;
    }

    &--e,
    &--f {
      left: 0;
    }

    &--b,
    &--f {
      top: 5%;
    }

    &--c,
    &--e {
      bottom: 5%;
    }
  }

  &__colon {
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    width: calc(var(--clock-segment-width) * 0.24);
    height: calc(var(--clock-segment-width) * 1.2);

    i {
      width: calc(var(--clock-segment-width) * 0.15);
      height: calc(var(--clock-segment-width) * 0.15);
      background: currentColor;
      border-radius: 50%;
    }
  }

  &__period {
    margin-left: 0.35em;
    font-size: calc(var(--clock-segment-width) * 0.42);
    font-weight: 500;
    line-height: 1;
  }
}
</style>
