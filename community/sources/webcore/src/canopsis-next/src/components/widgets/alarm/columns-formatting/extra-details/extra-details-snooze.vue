<template lang="pug">
  v-tooltip.c-extra-details(top)
    v-icon.pink.white--text.c-extra-details__badge(slot="activator", small) {{ icon }}
    div.text-md-center
      strong {{ $t('alarm.actions.iconsTitles.snooze') }}
      div {{ $t('common.by') }} : {{ snooze.a }}
      div {{ $t('common.date') }} : {{ date }}
      div {{ $t('common.end') }} : {{ end }}
      div(v-if="snooze.initiator") {{ $t('common.initiator') }} : {{ snooze.initiator }}
      div.c-extra-details__message(v-if="snooze.m") {{ $tc('common.comment') }} : {{ snooze.m }}
</template>

<script>
import { EVENT_ENTITY_STYLE, EVENT_ENTITY_TYPES } from '@/constants';

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
      return EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.snooze].icon;
    },
  },
};
</script>
