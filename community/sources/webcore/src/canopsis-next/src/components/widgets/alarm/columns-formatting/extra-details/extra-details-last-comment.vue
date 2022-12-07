<template lang="pug">
  c-clickable-tooltip.c-extra-details(top)
    v-icon.white--text.purple.lighten-2.c-extra-details__badge(
      slot="activator",
      small
    ) {{ icon }}
    div.text-md-center
      strong {{ $t('alarm.actions.iconsTitles.comment') }}
      div {{ $t('common.by') }}: {{ lastComment.a }}
      div {{ $t('common.date') }}: {{ date }}
      div.c-extra-details__message
        | {{ $tc('common.comment') }}:&nbsp;
        span(v-html="lastComment.m")
</template>

<script>
import { EVENT_ENTITY_STYLE, EVENT_ENTITY_TYPES } from '@/constants';

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
      return EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.comment].icon;
    },
  },
};
</script>
