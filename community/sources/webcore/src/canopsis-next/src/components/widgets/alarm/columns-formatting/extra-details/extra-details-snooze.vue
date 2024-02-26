<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <span
          class="c-extra-details__badge pink"
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
        <strong>{{ $t('alarm.actions.iconsTitles.snooze') }}</strong>
        <div>{{ $t('common.by') }} : {{ snooze.a }}</div>
        <div>{{ $t('common.date') }} : {{ date }}</div>
        <div>{{ $t('common.end') }} : {{ end }}</div>
        <div v-if="snooze.initiator">
          {{ $t('common.initiator') }} : {{ snooze.initiator }}
        </div>
        <div
          v-if="snooze.m"
          class="c-extra-details__message"
        >
          {{ $tc('common.comment') }} : {{ snooze.m }}
        </div>
      </div>
    </v-tooltip>
  </div>
</template>

<script>
import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

export default {
  props: {
    snooze: {
      type: Object,
      required: true,
    },
  },
  computed: {
    date() {
      return convertDateToStringWithFormatForToday(this.snooze.t);
    },

    end() {
      return convertDateToStringWithFormatForToday(this.snooze.val);
    },

    icon() {
      return getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.snooze);
    },
  },
};
</script>
