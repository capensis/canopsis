<template>
  <div>
    <c-simple-tooltip
      :content="tooltipContent"
      top
      disable-resize
    >
      <template #activator="{ on }">
        <c-alarm-extra-details-chip :icon="icon" :color="color" v-on="on" />
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { COLORS } from '@/config';
import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';

import { useExtraDetailsCanceledTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    canceled: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { tooltipContent } = useExtraDetailsCanceledTooltip(props);

    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.fastCancel);
    const color = COLORS.alarmExtraDetails.canceled;

    return {
      tooltipContent,
      icon,
      color,
    };
  },
};
</script>
