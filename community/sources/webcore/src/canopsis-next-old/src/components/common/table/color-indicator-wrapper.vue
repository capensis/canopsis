<template>
  <v-tooltip
    :disabled="!text"
    right
  >
    <template #activator="{ on }">
      <div
        :class="{ 'color-indicator--invalid': !text }"
        :style="{ backgroundColor: color }"
        class="color-indicator"
        v-on="on"
        @click="handleClick"
      >
        <slot>{{ value }}</slot>
      </div>
    </template>
    <span>{{ text }}</span>
  </v-tooltip>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';
import { COLOR_INDICATOR_TYPES } from '@/constants';

import { getAlarmStateColor, getAlarmImpactStateColor } from '@/helpers/entities/alarm/color';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
    alarm: {
      type: Object,
      default: () => ({}),
    },
    type: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit, listeners }) {
    const { t, te } = useI18n();

    const isImpactState = computed(() => props.type === COLOR_INDICATOR_TYPES.impactState);

    const impactLevel = computed(() => props.entity.impact_level ?? 0);

    const state = computed(() => props.alarm?.v?.state?.val
      ?? props.entity?.state
      ?? 0);

    const impactState = computed(() => props.entity?.impact_state
      ?? props.alarm?.impact_state
      ?? state.value * impactLevel.value);

    const value = computed(() => (isImpactState.value
      ? impactState.value
      : state.value));

    const color = computed(() => {
      const colorValue = isImpactState.value
        ? getAlarmImpactStateColor(impactState.value)
        : getAlarmStateColor(state.value);

      return colorValue ?? 'black';
    });

    const text = computed(() => {
      if (isImpactState.value) {
        return t('common.countOfTotal', { count: impactState.value, total: COLORS.impactState.length - 1 });
      }

      const key = `common.stateTypes.${state.value}`;

      return te(key) && t(key);
    });

    const handleClick = (event) => {
      if (listeners?.click) {
        listeners.click(event);
      }
      emit('click', event);
    };

    return {
      isImpactState,
      impactLevel,
      state,
      impactState,
      value,
      color,
      text,
      handleClick,
    };
  },
};
</script>

<style lang="scss" scoped>
.color-indicator {
  display: inline-block;
  border-radius: 10px;
  padding: 3px 7px;
  color: black;

  &--invalid {
    color: white;
  }
}
</style>
