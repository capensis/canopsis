<template lang="pug">
  v-tooltip.c-extra-details(top)
    v-icon.purple.white--text.c-extra-details__badge(slot="activator", small) {{ icon }}
    div.text-md-center
      strong {{ $t('alarmList.actions.iconsTitles.ack') }}
      div {{ $t('common.by') }} : {{ ack.a }}
      div {{ $t('common.date') }} : {{ ack.t | date('long') }}
      div(v-if="ack.initiator") {{ $t('common.initiator') }} : {{ ack.initiator }}
      div.c-extra-details__message(v-if="ack.m") {{ $tc('common.comment') }} : {{ ack.m }}
</template>

<script>
import { EVENT_ENTITY_STYLE, EVENT_ENTITY_TYPES } from '@/constants';

export default {
  props: {
    ack: {
      type: Object,
      required: true,
    },
  },
  computed: {
    icon() {
      return EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon;
    },
  },
};
</script>
