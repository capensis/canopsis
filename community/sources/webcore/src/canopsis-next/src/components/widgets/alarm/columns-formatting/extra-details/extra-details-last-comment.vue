<template lang="pug">
  c-clickable-tooltip.c-extra-details(top)
    template(#activator="")
      span.c-extra-details__badge.purple.lighten-2
        v-icon(color="white", small) {{ icon }}
    div.text-md-center
      strong {{ $t('alarm.actions.iconsTitles.comment') }}
      div {{ $t('common.by') }}: {{ lastComment.a }}
      div {{ $t('common.date') }}: {{ date }}
      div.c-extra-details__message
        | {{ $tc('common.comment') }}:&nbsp;
        span(v-html="lastComment.m")
</template>

<script>
import { EVENT_ENTITY_TYPES } from '@/constants';

import { getEntityEventIcon } from '@/helpers/icon';
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
  },
};
</script>
