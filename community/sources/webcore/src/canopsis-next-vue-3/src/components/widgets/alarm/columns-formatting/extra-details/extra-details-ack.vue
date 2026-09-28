<template>
  <div>
    <c-simple-tooltip :content="tooltipContent" html top>
      <template #activator="{ props: activatorProps }">
        <c-alarm-extra-details-chip
          v-bind="activatorProps"
          :icon="icon"
          theme-color="ack"
        />
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';

import { useExtraDetailsAckTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    ack: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { tooltipContent } = useExtraDetailsAckTooltip(props);

    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.ack);
    return {
      icon,
      tooltipContent,
    };
  },
};
</script>
