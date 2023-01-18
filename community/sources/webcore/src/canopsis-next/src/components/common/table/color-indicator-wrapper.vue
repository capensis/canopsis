<template lang="pug">
  v-tooltip(:disabled="!text", right)
    div.color-indicator.white--text(
      slot="activator",
      :style="{ backgroundColor: color }"
    )
      slot {{ value }}
    span {{ text }}
</template>

<script>
import { COLORS } from '@/config';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import { getEntityStateColor, getImpactStateColor } from '@/helpers/color';

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
      return this.alarm?.v?.state?.val ?? 0;
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
        return this.$t('common.countOfMax', { count: this.impactState, total: COLORS.impactState.length - 1 }); // TODO: add translations
      }

      return this.$t(`common.stateTypes.${this.state}`);
    },
  },
};
</script>

<style lang="scss" scoped>
.color-indicator {
  display: inline-block;
  border-radius: 10px;
  padding: 3px 7px;
}
</style>
