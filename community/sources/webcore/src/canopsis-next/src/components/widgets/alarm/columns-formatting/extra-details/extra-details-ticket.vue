<template>
  <div>
    <v-tooltip
      class="c-extra-details extra-details-ticket"
      top
    >
      <template #activator="{ on }">
        <v-badge
          :value="isLastFailed"
          class="time-line-flag"
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
          <c-alarm-extra-details-chip :color="color" :icon="icon" v-on="on" />
        </v-badge>
      </template>
      <v-layout
        class="extra-details-ticket__list"
        column
      >
        <div
          v-for="(ticket, index) in shownTickets"
          :key="index"
          class="text-md-center"
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
import { computed } from 'vue';

import { COLORS } from '@/config';
import { ALARM_LIST_ACTIONS_TYPES, ALARM_LIST_STEPS } from '@/constants';

import { getAlarmActionIcon } from '@/helpers/entities/alarm/icons';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

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
  setup(props) {
    const { t } = useI18n();

    const isSuccessTicketDeclaration = ticket => [
      ALARM_LIST_STEPS.declareTicket,
      ALARM_LIST_STEPS.assocTicket,
    ].includes(ticket?._t);
    const getTicketStatusText = ticket => t(`common.${isSuccessTicketDeclaration(ticket) ? 'ok' : 'failed'}`);
    const convertDateWithToday = date => convertDateToStringWithFormatForToday(date);

    const shownTickets = computed(() => props.tickets.slice(0, props.limit));
    const isLastFailed = computed(() => isSuccessTicketDeclaration(last(props.tickets)));
    const icon = getAlarmActionIcon(ALARM_LIST_ACTIONS_TYPES.declareTicket);
    const color = COLORS.alarmExtraDetails.ticket;

    return {
      getTicketStatusText,
      convertDateWithToday,
      shownTickets,
      isLastFailed,
      icon,
      color,
    };
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
