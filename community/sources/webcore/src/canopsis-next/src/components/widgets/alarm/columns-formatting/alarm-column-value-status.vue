<template lang="pug">
  c-no-events-icon.mr-1(v-if="isNoEventsStatus", :value="idleSince", color="red", top)
  v-tooltip(v-else, top)
    template(#activator="{ on }")
      v-icon(v-on="on", :color="statusColor") {{ status.icon }}
    span {{ $t(`common.statusTypes.${statusValue}`) }}
</template>

<script>
import { ENTITIES_STATUSES } from '@/constants';

import { formatState, formatStatus } from '@/helpers/formatting';

export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  computed: {
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
