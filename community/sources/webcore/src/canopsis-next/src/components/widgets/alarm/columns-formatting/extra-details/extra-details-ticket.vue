<template>
  <div>
    <c-simple-tooltip
      :content="tooltipContent"
      top
    >
      <template #activator="{ on }">
        <v-badge
          :value="!!failedTicket"
          class="time-line-flag"
          color="transparent"
          offset-y="15px"
          overlap
        >
          <template #badge="">
            <v-icon
              class="extra-details-ticket__badge-icon"
              color="error"
              size="14"
            >
              error
            </v-icon>
          </template>
          <c-alarm-extra-details-chip :color="color" :icon="icon" v-on="on" />
        </v-badge>
      </template>
    </c-simple-tooltip>
  </div>
</template>

<script>
import { COLORS } from '@/config';
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
    const color = COLORS.alarmExtraDetails.ticket;

    return {
      tooltipContent,
      icon,
      color,
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
