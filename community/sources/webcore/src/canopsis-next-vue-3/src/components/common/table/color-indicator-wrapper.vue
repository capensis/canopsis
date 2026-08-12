<template>
  <c-simple-tooltip :content="text" top>
    <template #activator="{ props: activatorProps }">
      <div
        v-bind="activatorProps"
        :class="{ 'color-indicator--invalid': !text }"
        :style="{ backgroundColor: color }"
        class="color-indicator"
        @click="$emit('click', $event)"
      >
        <slot>{{ value }}</slot>
      </div>
    </template>
  </c-simple-tooltip>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';
import { COLOR_INDICATOR_TYPES } from '@/constants';
import { getAlarmImpactStateColor, getAlarmStateColor } from '@/helpers/entities/alarm/color';
import { useI18n } from '@/hooks/i18n';

export default {
  emits: ['click'],

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

  setup(props) {
    const { t, te } = useI18n();

    const isImpactState = computed(() => props.type === COLOR_INDICATOR_TYPES.impactState);
    const impactLevel = computed(() => props.entity.impact_level ?? 0);
    const state = computed(() => props.alarm?.v?.state?.val ?? props.entity?.state ?? 0);
    const impactState = computed(
      () => props.entity?.impact_state ?? props.alarm?.impact_state ?? state.value * impactLevel.value,
    );
    const value = computed(() => (isImpactState.value ? impactState.value : state.value));
    const color = computed(() => {
      const colorValue = isImpactState.value
        ? getAlarmImpactStateColor(impactState.value)
        : getAlarmStateColor(state.value);

      return colorValue ?? 'black';
    });
    const text = computed(() => {
      if (isImpactState.value) {
        return t('common.countOfTotal', {
          count: impactState.value,
          total: COLORS.impactState.length - 1,
        });
      }

      const key = `common.stateTypes.${state.value}`;

      return te(key) && t(key);
    });

    return {
      color,
      text,
      value,
    };
  },
};
</script>

<style lang="scss" scoped>
.color-indicator {
  display: inline-block;
  padding: 3px 7px;
  color: black;
  border-radius: 10px;

  &--invalid {
    color: white;
  }
}
</style>
