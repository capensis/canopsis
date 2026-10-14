<template>
  <article class="clock-card">
    <h3 class="clock-card__label">
      {{ label }}
    </h3>

    <div class="clock-card__display">
      <analog-clock v-if="display === CLOCK_DISPLAY_TYPES.analog" :date="date" />
      <seven-segment-clock
        v-else-if="display === CLOCK_DISPLAY_TYPES.sevenSegments"
        :date="date"
        :use12-hour-format="use12HourFormat"
        :show-seconds="showSeconds"
      />
      <digital-clock
        v-else
        :date="date"
        :use12-hour-format="use12HourFormat"
        :show-seconds="showSeconds"
      />
    </div>

    <div v-if="showDate || metadata" class="clock-card__details">
      <div v-if="showDate">
        {{ formattedDate }}
      </div>
      <div v-if="metadata">
        {{ metadata }}
      </div>
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue';

import { CLOCK_DISPLAY_TYPES } from '@/constants';

import { formatTimezoneLabel } from '@/helpers/date/date';

import AnalogClock from './analog-clock.vue';
import DigitalClock from './digital-clock.vue';
import SevenSegmentClock from './seven-segment-clock.vue';

const props = defineProps({
  clock: {
    type: Object,
    required: true,
  },
  date: {
    type: Object,
    required: true,
  },
  display: {
    type: String,
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
  showDate: {
    type: Boolean,
    default: true,
  },
  showUtcOffset: {
    type: Boolean,
    default: true,
  },
  showTimezone: {
    type: Boolean,
    default: true,
  },
});

const timezoneLabel = computed(() => formatTimezoneLabel(props.clock.timezone));
const label = computed(() => {
  if (props.clock.label?.trim()) {
    return props.clock.label.trim();
  }

  return timezoneLabel.value.split('/').at(-1);
});
const formattedDate = computed(() => props.date.format('LL'));
const metadata = computed(() => [
  props.showTimezone ? timezoneLabel.value : '',
  props.showUtcOffset ? `UTC${props.date.format('Z')}` : '',
].filter(Boolean).join(' · '));
</script>

<style lang="scss" scoped>
.clock-card {
  display: flex;
  flex: 0 0 var(--clock-card-width);
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-width: 0;
  min-height: var(--clock-card-height);
  padding: var(--clock-card-padding);
  color: rgb(var(--v-theme-on-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.58);
  border-radius: 6px;

  &__label {
    overflow: hidden;
    width: 100%;
    margin: 0 0 var(--clock-label-spacing);
    font-size: var(--clock-label-font-size);
    font-weight: 600;
    line-height: 1.25;
    text-align: center;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__display {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: var(--clock-display-height);
  }

  &__details {
    min-height: var(--clock-details-height);
    margin-top: var(--clock-details-spacing);
    font-size: var(--clock-details-font-size);
    line-height: 1.45;
    text-align: center;
  }
}
</style>
