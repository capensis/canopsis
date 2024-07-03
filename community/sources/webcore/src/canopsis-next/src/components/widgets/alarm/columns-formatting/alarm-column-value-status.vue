<template>
  <c-no-events-icon
    v-if="isNoEventsStatus"
    :value="idleSince"
    :size="iconSize"
    color="error"
    top
  />
  <c-simple-tooltip
    v-else
    :content="$t(`common.statusTypes.${statusValue}`)"
    top
  >
    <template #activator="{ on }">
      <v-icon
        :size="iconSize"
        :style="{ color: statusColor, caretColor: statusColor }"
        v-on="on"
      >
        {{ status.icon }}
      </v-icon>
    </template>
  </c-simple-tooltip>
</template>

<script>
import { ALARM_STATUSES } from '@/constants';

import { formatAlarmState, formatAlarmStatus } from '@/helpers/entities/alarm/formatting';

export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    iconSize() {
      return this.small ? 24 : undefined;
    },

    statusValue() {
      return this.alarm.v.status.val;
    },

    isNoEventsStatus() {
      return this.statusValue === ALARM_STATUSES.noEvents;
    },

    isOngoingStatus() {
      return this.statusValue === ALARM_STATUSES.ongoing;
    },

    idleSince() {
      return this.alarm.entity.idle_since;
    },

    status() {
      return formatAlarmStatus(this.statusValue);
    },

    state() {
      return formatAlarmState(this.alarm.v.state.val);
    },

    statusColor() {
      return this.isOngoingStatus ? this.state.color : this.status.color;
    },
  },
};
</script>
