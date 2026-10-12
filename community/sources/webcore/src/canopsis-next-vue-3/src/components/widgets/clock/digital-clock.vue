<template>
  <time :datetime="dateTime" :aria-label="accessibleTime" class="digital-clock">
    <span>{{ time }}</span>
    <span v-if="period" class="digital-clock__period">{{ period }}</span>
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

const format = computed(() => {
  if (props.use12HourFormat) {
    return props.showSeconds ? 'hh:mm:ss' : 'hh:mm';
  }

  return props.showSeconds ? 'HH:mm:ss' : 'HH:mm';
});

const time = computed(() => props.date.format(format.value));
const period = computed(() => props.use12HourFormat ? props.date.format('A') : '');
const accessibleTime = computed(() => [time.value, period.value].filter(Boolean).join(' '));
const dateTime = computed(() => props.date.format());
</script>

<style lang="scss" scoped>
.digital-clock {
  display: flex;
  align-items: baseline;
  justify-content: center;
  color: inherit;
  font-size: var(--clock-time-font-size);
  font-variant-numeric: tabular-nums;
  font-weight: 400;
  line-height: 1;
  white-space: nowrap;

  &__period {
    margin-left: 0.25em;
    font-size: 0.32em;
    font-weight: 500;
  }
}
</style>
