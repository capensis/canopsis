<template>
  <div
    :class="[
      `clock-widget--size-${parameters.size}`,
      `clock-widget--spacing-${parameters.spacing}`,
    ]"
    class="clock-widget"
  >
    <clock-card
      v-for="clock in preparedClocks"
      :key="clock.key"
      :clock="clock"
      :date="clock.date"
      :display="parameters.display"
      :use12-hour-format="parameters.is_12_hour_format"
      :show-seconds="parameters.show_seconds"
      :show-date="parameters.show_date"
      :show-utc-offset="parameters.show_utc_offset"
      :show-timezone="parameters.show_timezone"
    />
  </div>
</template>

<script setup>
import moment from 'moment-timezone';
import { computed, onBeforeUnmount, ref, watch } from 'vue';

import {
  CLOCK_DISPLAY_TYPES,
  DEFAULT_CLOCK_WIDGET_PARAMETERS,
} from '@/constants';

import ClockCard from './clock-card.vue';

const props = defineProps({
  widget: {
    type: Object,
    required: true,
  },
  visible: {
    type: Boolean,
    default: true,
  },
});

const now = ref(Date.now());
let timer;

const parameters = computed(() => ({
  ...DEFAULT_CLOCK_WIDGET_PARAMETERS,
  ...props.widget.parameters,
}));

const updatesEverySecond = computed(() =>
  parameters.value.display === CLOCK_DISPLAY_TYPES.analog || parameters.value.show_seconds,
);

const updateInterval = computed(() => updatesEverySecond.value ? 1000 : 60000);
const preparedClocks = computed(() => (parameters.value.clocks ?? []).reduce((clocks, clock, index) => {
  if (!moment.tz.zone(clock.timezone)) {
    return clocks;
  }

  clocks.push({
    ...clock,
    key: `${clock.timezone}-${index}`,
    date: moment(now.value).tz(clock.timezone),
  });

  return clocks;
}, []));

const clearTimer = () => {
  clearTimeout(timer);
  timer = undefined;
};

const scheduleUpdate = () => {
  clearTimer();

  if (!props.visible) {
    return;
  }

  now.value = Date.now();

  const delay = updateInterval.value - now.value % updateInterval.value;
  timer = setTimeout(scheduleUpdate, delay);
};

watch([() => props.visible, updateInterval], scheduleUpdate, { immediate: true });

onBeforeUnmount(clearTimer);
</script>

<style lang="scss" scoped>
.clock-widget {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  align-items: flex-start;
  width: 100%;
  padding: 20px;

  &--spacing-automatic {
    justify-content: space-between;
  }

  &--spacing-fixed {
    justify-content: flex-start;
  }

  &--size-s {
    --clock-card-width: 220px;
    --clock-card-height: 126px;
    --clock-card-padding: 12px;
    --clock-display-height: 40px;
    --clock-details-height: 30px;
    --clock-label-spacing: 6px;
    --clock-label-font-size: 12px;
    --clock-time-font-size: 34px;
    --clock-segment-width: 17px;
    --clock-face-size: 80px;
    --clock-details-spacing: 6px;
    --clock-details-font-size: 10px;
  }

  &--size-m {
    --clock-card-width: 280px;
    --clock-card-height: 166px;
    --clock-card-padding: 16px;
    --clock-display-height: 56px;
    --clock-details-height: 36px;
    --clock-label-spacing: 8px;
    --clock-label-font-size: 14px;
    --clock-time-font-size: 46px;
    --clock-segment-width: 22px;
    --clock-face-size: 112px;
    --clock-details-spacing: 8px;
    --clock-details-font-size: 11px;
  }

  &--size-l {
    --clock-card-width: 340px;
    --clock-card-height: 218px;
    --clock-card-padding: 20px;
    --clock-display-height: 78px;
    --clock-details-height: 42px;
    --clock-label-spacing: 10px;
    --clock-label-font-size: 16px;
    --clock-time-font-size: 60px;
    --clock-segment-width: 28px;
    --clock-face-size: 156px;
    --clock-details-spacing: 10px;
    --clock-details-font-size: 12px;
  }
}

@media (max-width: 599px) {
  .clock-widget {
    padding: 12px;

    :deep(.clock-card) {
      flex-basis: min(100%, var(--clock-card-width));
    }
  }
}
</style>
