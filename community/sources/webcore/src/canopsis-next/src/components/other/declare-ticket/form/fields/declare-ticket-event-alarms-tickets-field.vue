<template lang="pug">
  v-layout(column)
    v-data-table(:headers="headers", :items="alarms", hide-actions)
      template(#items="{ item, index }")
        td.text-xs-left {{ item.v.connector_name }}
        td.text-xs-left {{ item.v.connector }}
        td.text-xs-left {{ item.v.component }}
        td.text-xs-left {{ item.v.resource }}
        td
          declare-ticket-event-tickets-chips-field(
            :value="value[item._id]",
            :tickets="ticketsByAlarms[item._id]",
            @input="updateTickets(item._id, $event)"
          )
    v-divider
    c-alert(v-if="!hasTickets", type="info") {{ $t('declareTicket.noRulesForAlarms') }}
    c-alert(v-if="hasErrors", type="error") {{ $t('declareTicket.errors.ticketRequired') }}
</template>

<script>
import { formMixin } from '@/mixins/form';

import DeclareTicketEventTicketsChipsField from './declare-ticket-event-tickets-chips-field.vue';

export default {
  inject: ['$validator'],
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
      required: () => ({}),
    },
    alarms: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'tickets_by_alarms',
    },
  },
  computed: {
    hasTickets() {
      return Object.values(this.ticketsByAlarms).filter(tickets => tickets.length).length;
    },

    hasErrors() {
      return this.errors.has(this.name);
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
        {
          text: this.$tc('common.ticket', 2),
          sortable: false,
        },
      ];
    },
  },
  created() {
    this.attachMinValueRule();
  },
  beforeDestroy() {
    this.detachRules();
  },
  methods: {
    attachMinValueRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'min_value:1',
        getter: () => Object.values(this.value).filter(tickets => tickets.length).length,
        vm: this,
      });
    },

    detachRules() {
      this.$validator.detach(this.name);
    },

    updateTickets(alarmId, tickets) {
      this.updateField(alarmId, tickets);
    },
  },
};
</script>
