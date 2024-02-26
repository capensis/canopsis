<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <span
          class="c-extra-details__badge purple"
          v-on="on"
        >
          <v-icon
            color="white"
            small
          >
            {{ icon }}
          </v-icon>
        </span>
      </template>
      <div class="text-md-center">
        <strong>{{ $t('alarm.actions.iconsTitles.ack') }}</strong>
        <div>{{ $t('common.by') }} : {{ ack.a }}</div>
        <div>{{ $t('common.date') }} : {{ date }}</div>
        <div v-if="ack.initiator">
          {{ $t('common.initiator') }} : {{ ack.initiator }}
        </div>
        <div
          v-if="ack.m"
          class="c-extra-details__message"
        >
          {{ $tc('common.comment') }} : {{ ack.m }}
        </div>
      </div>
    </v-tooltip>
  </div>
</template>

<script>
import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';

export default {
  props: {
    ack: {
      type: Object,
      required: true,
    },
  },
  computed: {
    date() {
      return convertDateToStringWithFormatForToday(this.ack.t);
    },

    icon() {
      return getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.ack);
    },
  },
};
</script>
