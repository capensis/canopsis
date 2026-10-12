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

import { useExtraDetailsCloseDelayTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    closeDelay: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { tooltipContent } = useExtraDetailsCloseDelayTooltip(props);

    const icon = '$vuetify.icons.close_with_delay';
    const color = COLORS.alarmExtraDetails.closeDelay;

    return {
      tooltipContent,
      icon,
      color,
    };
  },
};
</script>
