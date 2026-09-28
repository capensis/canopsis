<template>
  <div>
    <c-simple-tooltip :content="tooltipContent" html top>
      <template #activator="{ props: activatorProps }">
        <c-alarm-extra-details-chip v-bind="activatorProps" :icon="icon" theme-color="snooze">
          <template #default="{ iconColor }">
            <v-icon :color="iconColor" size="small">
              {{ icon }}
            </v-icon>
            <strong v-if="hasInactivePbehavior" :style="{ color: iconColor }" class="ml-2">
              !
            </strong>
          </template>
        </c-alarm-extra-details-chip>
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
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

    return {
      icon,
      tooltipContent,
    };
  },
};
</script>
