<template>
  <div>
    <c-simple-tooltip :content="tooltipContent" html top>
      <template #activator="{ props: activatorProps }">
        <c-alarm-extra-details-chip v-bind="activatorProps" :icon="icon" :color="color" />
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
      color,
      icon,
      tooltipContent,
    };
  },
};
</script>
