<template>
  <div>
    <c-simple-tooltip
      :content="tooltipContent"
      top
      disable-resize
    >
      <template #activator="{ on }">
        <c-alarm-extra-details-chip :color="color" icon="center_focus_weak" v-on="on" />
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { COLORS } from '@/config';

import { useExtraDetailsParentsTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    rules: {
      type: Array,
      default: () => [],
    },
    total: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const { tooltipContent } = useExtraDetailsParentsTooltip(props);

    const color = COLORS.alarmExtraDetails.parent;

    return {
      tooltipContent,
      color,
    };
  },
};
</script>

<style lang="scss" scoped>
.rule-name:nth-of-type(odd) {
  color: #b5b5b5;
}
</style>
