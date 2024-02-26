<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <span
          class="c-extra-details__badge blue-grey"
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
        <strong>{{ $t('alarm.actions.iconsTitles.canceled') }}</strong>
        <div>{{ $t('common.by') }} : {{ canceled.a }}</div>
        <div>{{ $t('common.date') }} : {{ date }}</div>
        <div
          v-if="canceled.m"
          class="c-extra-details__message"
        >
          {{ $tc('common.comment') }} : {{ canceled.m }}
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
    canceled: {
      type: Object,
      required: true,
    },
  },
  computed: {
    date() {
      return convertDateToStringWithFormatForToday(this.canceled.t);
    },

    icon() {
      return getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.fastCancel);
    },
  },
};
</script>
