<template>
  <div>
    <c-simple-tooltip
      :content="tooltipContent"
      top
    >
      <template #activator="{ on }">
        <c-alarm-extra-details-chip :color="color" :icon="icon" v-on="on" />
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { COLORS } from '@/config';
import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';

import { useExtraDetailsSnoozeTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    snooze: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { tooltipContent } = useExtraDetailsSnoozeTooltip(props);

    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.snooze);
    const color = COLORS.alarmExtraDetails.snooze;

    return {
      tooltipContent,
      icon,
      color,
    };
  },
};
</script>
