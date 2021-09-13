<template lang="pug">
  c-clickable-tooltip.c-extra-details(top)
    v-icon.white--text.purple.lighten-2.c-extra-details__badge(
      slot="activator",
      small
    ) {{ icon }}
    div.text-md-center
      strong {{ $t('alarmList.actions.iconsTitles.comment') }}
      div {{ $t('common.by') }}: {{ lastComment.a }}
      div {{ $t('common.date') }}: {{ lastComment.t | dateWithToday('long') }}
      div.c-extra-details__message
        | {{ $tc('common.comment') }}:&nbsp;
        span(v-html="lastComment.m")
</template>

<script>
import { EVENT_ENTITY_STYLE, EVENT_ENTITY_TYPES } from '@/constants';

export default {
  props: {
    lastComment: {
      type: Object,
      required: true,
    },
  },
  computed: {
    icon() {
      return EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.comment].icon;
    },
  },
};
</script>
