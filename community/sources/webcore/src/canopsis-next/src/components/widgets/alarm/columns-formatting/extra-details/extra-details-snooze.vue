<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <c-alarm-extra-details-chip :color="color" :icon="icon" v-on="on" />
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
import { computed } from 'vue';

import { COLORS } from '@/config';
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
  setup(props) {
    const date = computed(() => convertDateToStringWithFormatForToday(props.snooze.t));
    const end = computed(() => convertDateToStringWithFormatForToday(props.snooze.val));
    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.snooze);
    const color = COLORS.alarmExtraDetails.snooze;

    return {
      date,
      end,
      icon,
      color,
    };
  },
};
</script>
