<template>
  <div>
    <c-simple-tooltip :content="tooltipContent" html top>
      <template #activator="{ props: activatorProps }">
        <c-alarm-extra-details-chip
          v-bind="activatorProps"
          :color="color"
          icon="center_focus_strong"
        />
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { COLORS } from '@/config';

import { useExtraDetailsChildrenTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    rule: {
      type: Object,
      required: true,
    },
    total: {
      type: Number,
      default: 0,
    },
    opened: {
      type: Number,
      default: 0,
    },
    closed: {
      type: Number,
      default: 0,
    },
  },

  setup(props) {
    const { tooltipContent } = useExtraDetailsChildrenTooltip(props);

    const color = COLORS.alarmExtraDetails.children;

    return {
      color,
      tooltipContent,
    };
  },
};
</script>
