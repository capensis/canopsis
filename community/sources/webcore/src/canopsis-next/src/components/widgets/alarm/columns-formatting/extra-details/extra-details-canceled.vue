<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
      disable-resize
    >
      <template #activator="{ on }">
        <c-alarm-extra-details-chip :icon="icon" :color="color" v-on="on" />
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
import { computed } from 'vue';

import { COLORS } from '@/config';
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
  setup(props) {
    const date = computed(() => convertDateToStringWithFormatForToday(props.canceled.t));
    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.fastCancel);
    const color = COLORS.alarmExtraDetails.canceled;

    return {
      date,
      icon,
      color,
    };
  },
};
</script>
