<template>
  <c-clickable-tooltip
    class="c-extra-details"
    top
  >
    <template #activator="">
      <c-alarm-extra-details-chip :icon="icon" :color="color" />
    </template>
    <div class="text-md-center">
      <strong>{{ $t('alarm.actions.iconsTitles.comment') }}</strong>
      <div>{{ $t('common.by') }}: {{ lastComment.a }}</div>
      <div>{{ $t('common.date') }}: {{ date }}</div>
      <div class="c-extra-details__message">
        {{ $tc('common.comment') }}:&nbsp;
        <div v-html="sanitizedLastComment" />
      </div>
    </div>
  </c-clickable-tooltip>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';
import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import { sanitizeHtml, linkifyHtml } from '@/helpers/html';
import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

export default {
  props: {
    lastComment: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const date = computed(() => convertDateToStringWithFormatForToday(props.lastComment.t));
    const sanitizedLastComment = computed(() => sanitizeHtml(linkifyHtml(String(props.lastComment?.m ?? ''))));
    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.comment);
    const color = COLORS.alarmExtraDetails.comment;

    return {
      date,
      sanitizedLastComment,
      icon,
      color,
    };
  },
};
</script>
