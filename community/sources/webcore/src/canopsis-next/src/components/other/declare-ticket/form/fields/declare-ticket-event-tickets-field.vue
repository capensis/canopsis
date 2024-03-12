<template>
  <v-layout column>
    <v-checkbox
      v-for="chip in chips"
      :key="chip.value"
      :input-value="chip.active"
      :label="chip.text"
      color="primary"
      hide-details
      @change="updateActive(chip.value)"
    >
      <template #append="">
        <extra-details-ticket
          v-if="chip.assignedTickets.length"
          :tickets="chip.assignedTickets"
          class="ml-2"
        />
      </template>
    </v-checkbox>
  </v-layout>
</template>

<script>
import { groupBy } from 'lodash';

import { EVENT_ENTITY_TYPES } from '@/constants';

import { filterValue } from '@/helpers/array';

import { formMixin } from '@/mixins/form';

import ExtraDetailsTicket from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-ticket.vue';

export default {
  components: { ExtraDetailsTicket },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    tickets: {
      type: Array,
      default: () => [],
    },
    alarmTickets: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    successAlarmTickets() {
      return this.alarmTickets
        .filter(ticket => [EVENT_ENTITY_TYPES.declareTicket, EVENT_ENTITY_TYPES.assocTicket].includes(ticket._t));
    },

    successAlarmTicketsByTicketId() {
      return groupBy(this.successAlarmTickets, 'ticket_rule_id');
    },

    chips() {
      return this.tickets.map(({ _id: id, name }) => ({
        active: this.isActiveTicket(id),
        text: name,
        value: id,
        assignedTickets: this.successAlarmTicketsByTicketId[id] ?? [],
      }));
    },
  },
  methods: {
    isActiveTicket(ticketId) {
      return this.value.includes(ticketId);
    },

    updateActive(id) {
      this.updateModel(
        this.isActiveTicket(id)
          ? filterValue(this.value, id)
          : [...this.value, id],
      );
    },
  },
};
</script>
