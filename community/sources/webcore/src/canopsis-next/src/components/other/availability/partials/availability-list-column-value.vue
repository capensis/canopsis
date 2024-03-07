<template>
  <span class="availability-list-column-value">
    {{ value }}

    <c-help-icon
      v-if="isTrendEnabled"
      :text="trendTooltipText"
      :icon-class="{
        'availability-list-column-value__trend': true,
        'availability-list-column-value__trend--up': trendUp
      }"
      :color="trendUp ? 'success' : 'error'"
      icon="arrow_downward"
      size="16"
      top
    />
  </span>
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE } from '@/constants';

import { convertDurationToString } from '@/helpers/date/duration';
import {
  getAvailabilityFieldByDisplayParameterAndShowType,
  getAvailabilityTrendFieldByDisplayParameter,
} from '@/helpers/entities/availability/entity';

export default {
  props: {
    availability: {
      type: Object,
      required: true,
    },
    displayParameter: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
    showType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
    showTrend: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const isPercentType = computed(() => props.showType === AVAILABILITY_SHOW_TYPE.percent);
    const valueField = computed(() => getAvailabilityFieldByDisplayParameterAndShowType(
      props.displayParameter,
      props.showType,
    ));
    const trendValueField = computed(() => getAvailabilityTrendFieldByDisplayParameter(props.displayParameter));

    const targetValue = computed(() => props.availability[valueField.value]);
    const targetTrendValue = computed(() => props.availability[trendValueField.value]);

    const trendUp = computed(
      () => targetValue.value > targetTrendValue.value,
    );
    const isTrendEnabled = computed(
      () => props.showTrend
        && isPercentType.value
        && targetTrendValue.value
        && targetValue.value !== targetTrendValue.value,
    );

    const trendTooltipText = computed(
      () => (isPercentType.value
        ? `${targetTrendValue.value}%`
        : convertDurationToString(targetTrendValue.value)),
    );
    const value = computed(
      () => (isPercentType.value
        ? `${targetValue.value}%`
        : convertDurationToString(targetValue.value)),
    );

    return {
      value,
      isPercentType,
      isTrendEnabled,
      trendUp,
      trendTooltipText,
    };
  },
};
</script>

<style lang="scss">
.availability-list-column-value {
  &__trend {
    transform: rotate(0deg);
    transition: transform linear .3s;

    &--up {
      transform: rotate(180deg);
    }
  }
}
</style>
