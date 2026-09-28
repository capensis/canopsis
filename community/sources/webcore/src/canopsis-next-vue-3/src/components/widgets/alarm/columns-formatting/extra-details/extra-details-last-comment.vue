<template>
  <c-simple-tooltip :content="tooltipContent" html top>
    <template #activator="{ props: activatorProps }">
      <c-alarm-extra-details-chip v-bind="activatorProps" :icon="icon" theme-color="comment" />
    </template>
  </c-simple-tooltip>
</template>

<script>
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
    return {
      icon,
      tooltipContent,
    };
  },
};
</script>
