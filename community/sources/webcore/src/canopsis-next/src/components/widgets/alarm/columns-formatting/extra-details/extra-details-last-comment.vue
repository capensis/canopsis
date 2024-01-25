<template>
  <c-clickable-tooltip
    class="c-extra-details"
    top
  >
    <template #activator="">
      <span class="c-extra-details__badge purple lighten-2">
        <v-icon
          color="white"
          small
        >
          {{ icon }}
        </v-icon>
      </span>
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
import { EVENT_ENTITY_TYPES } from '@/constants';

import { sanitizeHtml, linkifyHtml } from '@/helpers/html';
import { getEntityEventIcon } from '@/helpers/entities/entity/icons';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

export default {
  props: {
    lastComment: {
      type: Object,
      required: true,
    },
  },
  computed: {
    date() {
      return convertDateToStringWithFormatForToday(this.lastComment.t);
    },

    icon() {
      return getEntityEventIcon(EVENT_ENTITY_TYPES.comment);
    },

    sanitizedLastComment() {
      return sanitizeHtml(linkifyHtml(String(this.lastComment?.m ?? '')));
    },
  },
};
</script>
