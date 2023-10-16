<template>
  <v-tooltip
    :disabled="!text"
    right
  >
    <template #activator="{ on }">
      <div
        class="color-indicator"
        v-on="on"
        :class="{ 'color-indicator--invalid': !text }"
        :style="{ backgroundColor: color }"
      >
        <slot>{{ value }}</slot>
      </div>
    </template>
    <span>{{ text }}</span>
  </v-tooltip>
</template>

<script>
import { COLORS } from '@/config';
import { COLOR_INDICATOR_TYPES } from '@/constants';

import { getEntityStateColor, getImpactStateColor } from '@/helpers/entities/entity/color';

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
  computed: {
    isImpactState() {
      return this.type === COLOR_INDICATOR_TYPES.impactState;
    },

    impactLevel() {
      return this.entity.impact_level ?? 0;
    },

    state() {
      return this.alarm?.v?.state?.val
        ?? this.entity?.state
        ?? 0;
    },

    impactState() {
      return this.entity?.impact_state
        ?? this.alarm?.impact_state
        ?? this.state * this.impactLevel;
    },

    value() {
      return this.isImpactState
        ? this.impactState
        : this.state;
    },

    color() {
      const color = this.isImpactState
        ? getImpactStateColor(this.impactState)
        : getEntityStateColor(this.state);

      return color ?? 'black';
    },

    text() {
      if (this.isImpactState) {
        return this.$t('common.countOfTotal', { count: this.impactState, total: COLORS.impactState.length - 1 });
      }

      const key = `common.stateTypes.${this.state}`;

      return this.$te(key) && this.$t(key);
    },
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
