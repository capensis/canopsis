<template>
  <c-no-events-icon
    v-if="isNoEventsStatus"
    :value="idleSince"
    :size="iconSize"
    color="error"
    top
  />
  <v-tooltip
    v-else
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
    <span>{{ $t(`common.statusTypes.${statusValue}`) }}</span>
  </v-tooltip>
</template>

<script>
import { ENTITIES_STATUSES } from '@/constants';

import { formatState, formatStatus } from '@/helpers/entities/entity/formatting';

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
      return this.statusValue === ENTITIES_STATUSES.noEvents;
    },

    isOngoingStatus() {
      return this.statusValue === ENTITIES_STATUSES.ongoing;
    },

    idleSince() {
      return this.alarm.entity.idle_since;
    },

    status() {
      return formatStatus(this.statusValue);
    },

    state() {
      return formatState(this.alarm.v.state.val);
    },

    statusColor() {
      return this.isOngoingStatus ? this.state.color : this.status.color;
    },
  },
};
</script>
