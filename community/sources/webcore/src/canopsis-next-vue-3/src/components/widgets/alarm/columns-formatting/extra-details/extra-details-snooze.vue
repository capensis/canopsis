<template>
  <div>
    <c-simple-tooltip :content="tooltipContent" html top>
      <template #activator="{ props: activatorProps }">
        <c-alarm-extra-details-chip v-bind="activatorProps" :color="color" :icon="icon">
          <v-icon color="white" size="small">
            {{ icon }}
          </v-icon>
          <strong v-if="hasInactivePbehavior" class="ml-2 text-white">
            !
          </strong>
        </c-alarm-extra-details-chip>
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
    hasInactivePbehavior: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const { tooltipContent } = useExtraDetailsSnoozeTooltip(props);

    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.snooze);
    const color = COLORS.alarmExtraDetails.snooze;

    return {
      color,
      icon,
      tooltipContent,
    };
  },
};
</script>
