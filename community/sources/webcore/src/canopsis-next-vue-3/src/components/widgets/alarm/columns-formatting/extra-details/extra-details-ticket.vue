<template>
  <c-tooltip :content="tooltipContent" html top>
    <template #activator="{ props: activatorProps }">
      <v-badge
        :model-value="!!failedTicket"
        class="time-line-flag"
        color="transparent"
        offset-y="15px"
      >
        <template #badge>
          <v-icon class="extra-details-ticket__badge-icon" color="error" size="14">
            error
          </v-icon>
        </template>

        <c-alarm-extra-details-chip v-bind="activatorProps" :icon="icon" theme-color="tickets" />
      </v-badge>
    </template>
  </c-tooltip>
</template>

<script>
import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';

import { useExtraDetailsTicketTooltip } from '../../hooks/extra-details-tooltips';

export default {
  props: {
    tickets: {
      type: Array,
      required: true,
    },
    failedTicket: {
      type: Object,
      required: false,
    },
    limit: {
      type: Number,
      default: 5,
    },
  },
  setup(props) {
    const { tooltipContent } = useExtraDetailsTicketTooltip(props);

    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.declareTicket);
    return {
      tooltipContent,
      icon,
    };
  },
};
</script>

<style lang="scss">
.extra-details-ticket {
  &__badge-icon {
    background: white;
    border-radius: 50%;
  }

  &__list {
    gap: 10px;
  }
}
</style>
