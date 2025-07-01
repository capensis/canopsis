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
        :style="iconStyle"
        v-on="on"
      >
        {{ status.icon }}
      </v-icon>
    </template>
  </c-simple-tooltip>
</template>

<script>
import { computed } from 'vue';

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
  setup(props) {
    const statusValue = computed(() => props.alarm.v.status.val);
    const isNoEventsStatus = computed(() => statusValue.value === ALARM_STATUSES.noEvents);
    const isOngoingStatus = computed(() => statusValue.value === ALARM_STATUSES.ongoing);
    const idleSince = computed(() => props.alarm.entity.idle_since);
    const status = computed(() => formatAlarmStatus(statusValue.value));
    const state = computed(() => formatAlarmState(props.alarm.v.state.val));
    const statusColor = computed(() => (isOngoingStatus.value ? state.value.color : status.value.color));
    const iconSize = computed(() => (props.small ? 24 : undefined));
    const iconStyle = computed(() => ({ color: statusColor.value, caretColor: statusColor.value }));

    return {
      statusValue,
      isNoEventsStatus,
      isOngoingStatus,
      idleSince,
      status,
      state,
      iconSize,
      iconStyle,
    };
  },
};
</script>
