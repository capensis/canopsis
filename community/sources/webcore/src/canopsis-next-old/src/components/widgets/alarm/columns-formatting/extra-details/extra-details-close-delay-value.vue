<template>
  <div>
    <c-simple-tooltip
      :content="tooltipContent"
      top
    >
      <template #activator="{ on }">
        <c-alarm-extra-details-chip :icon="icon" :color="color" v-on="on" />
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { COLORS } from '@/config';

import { useExtraDetailsCloseDelayValueTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    closeDelayValue: {
      type: Number,
      required: true,
    },
  },
  setup(props) {
    const { tooltipContent } = useExtraDetailsCloseDelayValueTooltip(props);

    const icon = '$vuetify.icons.remove_clock';
    const color = COLORS.alarmExtraDetails.closeDelay;

    return {
      tooltipContent,
      icon,
      color,
    };
  },
};
</script>
