<template lang="pug">
  div
    v-tooltip(v-if="isEnabled", :disabled="!formattedData.text", right)
      div.color-indicator.white--text(
        slot="activator",
        :style="{ backgroundColor: formattedData.color }"
      )
        slot {{ value }}
      span {{ formattedData.text }}
    slot(v-else)
</template>

<script>
import { get, isUndefined } from 'lodash';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import { formatState, formatImpactState } from '@/helpers/formatting';

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
    isEnabled() {
      return !!this.type;
    },

    impactLevel() {
      return get(this.entity, 'impact_level', 0);
    },

    state() {
      return get(this.alarm, 'v.state.val', 0);
    },

    impactState() {
      const impactState = get(this.entity, 'impact_state');

      if (!isUndefined(impactState)) {
        return impactState;
      }

      return get(this.alarm, 'impact_state', this.state * this.impactLevel);
    },

    value() {
      return {
        [COLOR_INDICATOR_TYPES.state]: this.state,
        [COLOR_INDICATOR_TYPES.impactState]: this.impactState,
      }[this.type];
    },

    formattedData() {
      const formatter = {
        [COLOR_INDICATOR_TYPES.state]: formatState,
        [COLOR_INDICATOR_TYPES.impactState]: formatImpactState,
      }[this.type];

      return formatter ? formatter(this.value) : {};
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
