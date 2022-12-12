<template lang="pug">
  div
    v-tooltip.c-extra-details(top)
      v-icon.c-extra-details__badge.blue.white--text(slot="activator", small) {{ icon }}
      div.text-md-center
        strong {{ $t('alarm.actions.iconsTitles.declareTicket') }}
        div {{ $t('common.by') }} : {{ ticket.a }}
        div {{ $t('common.date') }} : {{ date }}
        div(
          v-if="ticket.val"
        ) {{ $t('alarm.actions.iconsFields.ticketNumber') }} : {{ ticket.val }}
</template>

<script>
import { EVENT_ENTITY_STYLE, EVENT_ENTITY_TYPES } from '@/constants';

import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

export default {
  props: {
    ticket: {
      type: Object,
      required: true,
    },
  },
  computed: {
    date() {
      return convertDateToStringWithFormatForToday(this.ticket.t);
    },

    icon() {
      return EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.declareTicket].icon;
    },
  },
};
</script>
