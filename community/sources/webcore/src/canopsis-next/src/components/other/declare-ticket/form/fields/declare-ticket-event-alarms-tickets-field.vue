<template lang="pug">
  v-layout(column)
    v-data-table(:headers="headers", :items="alarms", hide-actions)
      template(#items="{ item, index }")
        td.text-xs-left {{ item.v.connector_name }}
        td.text-xs-left {{ item.v.connector }}
        td.text-xs-left {{ item.v.component }}
        td.text-xs-left {{ item.v.resource }}
        td(v-if="!hideTickets")
          v-layout(row, align-center)
            declare-ticket-event-tickets-chips-field(
              :value="activeTicketsByAlarms[item._id]",
              :tickets="ticketsByAlarms[item._id]",
              :disabled="disableTickets",
              @input="updateTickets(item._id, $event)"
            )
            c-action-btn(
              v-if="!hideRemove",
              :disabled="!hasActiveTickets(item._id)",
              type="delete",
              @click="removeTickets(item._id)"
            )
    v-divider
</template>

<script>
import { difference } from 'lodash';
import { filterValue, revertGroupBy } from '@/helpers/entities';

import { formMixin } from '@/mixins/form';

import DeclareTicketEventTicketsChipsField from './declare-ticket-event-tickets-chips-field.vue';

export default {
  components: { DeclareTicketEventTicketsChipsField },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    ticketsByAlarms: {
      type: Object,
      default: () => ({}),
    },
    alarms: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'alarms_by_tickets',
    },
    disableTickets: {
      type: Boolean,
      default: false,
    },
    hideRemove: {
      type: Boolean,
      default: false,
    },
    hideTickets: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    activeTicketsByAlarms() {
      return revertGroupBy(this.value);
    },

    headers() {
      return [
        {
          text: this.$t('common.connectorName'),
          sortable: false,
        },
        {
          text: this.$t('common.connector'),
          sortable: false,
        },
        {
          text: this.$t('common.component'),
          sortable: false,
        },
        {
          text: this.$t('common.resource'),
          sortable: false,
        },
        !this.hideTickets && {
          text: this.$tc('common.ticket', 2),
          sortable: false,
        },
      ].filter(Boolean);
    },
  },
  methods: {
    hasActiveTickets(alarmId) {
      return !!this.activeTicketsByAlarms[alarmId]?.length;
    },

    updateTickets(alarmId, tickets) {
      const oldTickets = this.activeTicketsByAlarms[alarmId] ?? [];

      const removedTickets = difference(oldTickets, tickets);
      const addedTickets = difference(tickets, oldTickets);

      const newValue = { ...this.value };

      addedTickets.forEach((ticketId) => {
        newValue[ticketId] = [...this.value[ticketId], alarmId];
      });

      removedTickets.forEach((ticketId) => {
        newValue[ticketId] = filterValue(this.value[ticketId], alarmId);
      });

      this.updateModel(newValue);
    },

    removeTickets(alarmId) {
      const tickets = Object.entries(this.value).reduce((acc, [ticketId, alarms]) => {
        acc[ticketId] = filterValue(alarms, alarmId);

        return acc;
      }, {});

      this.updateModel(tickets);
    },
  },
};
</script>
