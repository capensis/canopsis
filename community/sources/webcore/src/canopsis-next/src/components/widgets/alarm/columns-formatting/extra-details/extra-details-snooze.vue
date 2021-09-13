<template lang="pug">
  v-tooltip.c-extra-details(top)
    v-icon.pink.white--text.c-extra-details__badge(slot="activator", small) {{ icon }}
    div.text-md-center
      strong {{ $t('alarmList.actions.iconsTitles.snooze') }}
      div {{ $t('common.by') }} : {{ snooze.a }}
      div {{ $t('common.date') }} : {{ snooze.t | dateWithToday('long') }}
      div {{ $t('common.end') }} : {{ snooze.val | dateWithToday('long') }}
      div(v-if="snooze.initiator") {{ $t('common.initiator') }} : {{ snooze.initiator }}
      div.c-extra-details__message(v-if="snooze.m") {{ $tc('common.comment') }} : {{ snooze.m }}
</template>

<script>
import { EVENT_ENTITY_STYLE, EVENT_ENTITY_TYPES } from '@/constants';

export default {
  props: {
    snooze: {
      type: Object,
      required: true,
    },
  },
  computed: {
    icon() {
      return EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.snooze].icon;
    },
  },
};
</script>
