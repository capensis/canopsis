<template lang="pug">
  div
    v-tooltip.c-extra-details(top, lazy)
      template(#activator="{ on }")
        span.c-extra-details__badge.purple(v-on="on")
          v-icon(color="white", small) {{ icon }}
      div.text-md-center
        strong {{ $t('alarm.actions.iconsTitles.ack') }}
        div {{ $t('common.by') }} : {{ ack.a }}
        div {{ $t('common.date') }} : {{ date }}
        div(v-if="ack.initiator") {{ $t('common.initiator') }} : {{ ack.initiator }}
        div.c-extra-details__message(v-if="ack.m") {{ $tc('common.comment') }} : {{ ack.m }}
</template>

<script>
import { EVENT_ENTITY_TYPES } from '@/constants';

import { getEntityEventIcon } from '@/helpers/entities/entity/icons';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

export default {
  props: {
    ack: {
      type: Object,
      required: true,
    },
  },
  computed: {
    date() {
      return convertDateToStringWithFormatForToday(this.ack.t);
    },

    icon() {
      return getEntityEventIcon(EVENT_ENTITY_TYPES.ack);
    },
  },
};
</script>
