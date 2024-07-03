<template>
  <c-simple-tooltip
    :content="tooltipContent"
    top
    clickable
  >
    <template #activator="{ on }">
      <c-alarm-extra-details-chip :icon="icon" :color="color" v-on="on" />
    </template>
  </c-simple-tooltip>
</template>

<script>
import { COLORS } from '@/config';
import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';

import { useExtraDetailsLastCommentTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    lastComment: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { tooltipContent } = useExtraDetailsLastCommentTooltip(props);

    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.comment);
    const color = COLORS.alarmExtraDetails.comment;

    return {
      tooltipContent,
      icon,
      color,
    };
  },
};
</script>
