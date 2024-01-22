<template>
  <div>
    <v-tooltip
      class="c-extra-details extra-details-ticket"
      top
    >
      <template #activator="{ on }">
        <v-badge
          class="time-line-flag"
          :value="isLastFailed"
          color="transparent"
          overlap
        >
          <template #badge="">
            <v-icon
              class="extra-details-ticket__badge-icon"
              color="error"
              size="14"
            >
              error
            </v-icon>
          </template>
          <span
            class="c-extra-details__badge blue"
            v-on="on"
          >
            <v-icon
              color="white"
              small
            >
              {{ icon }}
            </v-icon>
          </span>
        </v-badge>
      </template>
      <v-layout
        class="extra-details-ticket__list"
        column
      >
        <div
          class="text-md-center"
          v-for="(ticket, index) in shownTickets"
          :key="index"
        >
          <strong>{{ ticket.ticket_rule_name }} {{ getTicketStatusText(ticket) }}</strong>
          <div>{{ $t('common.by') }} : {{ ticket.a }}</div>
          <div>{{ $t('common.date') }} : {{ convertDateWithToday(ticket.t) }}</div>
          <div v-if="ticket.ticket">
            {{ $t('alarm.actions.iconsFields.ticketNumber') }} : {{ ticket.ticket }}
          </div>
          <div v-if="ticket.ticket_comment">
            {{ $tc('common.comment') }} : {{ ticket.ticket_comment }}
          </div>
        </div>
      </v-layout>
      <div class="mt-2">
        <i v-if="tickets.length > limit">{{ $t('alarm.otherTickets') }}</i>
      </div>
    </v-tooltip>
  </div>
</template>

<script>
import { last } from 'lodash';

import { EVENT_ENTITY_TYPES } from '@/constants';

import { getEntityEventIcon } from '@/helpers/entities/entity/icons';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

export default {
  props: {
    tickets: {
      type: Array,
      required: true,
    },
    limit: {
      type: Number,
      default: 5,
    },
  },
  computed: {
    shownTickets() {
      return this.tickets.slice(0, this.limit);
    },

    isLastFailed() {
      return !this.isSuccessTicketDeclaration(last(this.tickets));
    },

    icon() {
      return getEntityEventIcon(EVENT_ENTITY_TYPES.declareTicket);
    },
  },
  methods: {
    isSuccessTicketDeclaration(ticket = {}) {
      return [EVENT_ENTITY_TYPES.declareTicket, EVENT_ENTITY_TYPES.assocTicket].includes(ticket?._t);
    },

    getTicketStatusText(ticket) {
      return this.$t(`common.${this.isSuccessTicketDeclaration(ticket) ? 'ok' : 'failed'}`);
    },

    convertDateWithToday(date) {
      return convertDateToStringWithFormatForToday(date);
    },
  },
};
</script>

<style lang="scss">
.extra-details-ticket {
  &__badge-icon {
    background: white;
    border-radius: 50%;
  }

  &__list {
    gap: 10px;
  }
}
</style>
