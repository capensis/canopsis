<template>
  <v-layout column>
    <c-information-block
      v-if="!singleMode"
      :title="$t('declareTicket.applyRules')"
    >
      <declare-ticket-event-form
        :form="commonValue"
        :alarms="alarms"
        :tickets-by-alarms="ticketsByAlarms"
        :hide-ticket-resource="hideTicketResource"
        hide-row-select
        @input="updateCommonValue"
      />
    </c-information-block>
    <template v-else>
      <c-information-block
        v-for="group in groups"
        :key="group.ticketId"
        :title="$t('declareTicket.applyRules')"
      >
        <v-layout class="mt-2">
          <v-checkbox
            :input-value="group.everyEnabled"
            :indeterminate="group.someEnabled"
            :label="group.name"
            color="primary"
            hide-details
            @change="updateEnabledTickets(group.ticketId)"
          />
        </v-layout>
        <declare-ticket-event-form
          :form="group.value"
          :alarms="group.alarms"
          :tickets-by-alarms="group.ticketsByAlarms"
          :hide-ticket-resource="hideTicketResource"
          hide-tickets
          @input="updateGroup(group.ticketId, $event)"
        />
      </c-information-block>
      <c-information-block
        v-if="alarmsWithoutTickets.length"
        :title="$t('declareTicket.noRulesForAlarms')"
      >
        <declare-ticket-event-alarms-tickets-field
          :alarms="alarmsWithoutTickets"
          hide-tickets
          hide-row-select
        />
      </c-information-block>
    </template>
    <c-alert
      v-if="hasErrors"
      type="error"
    >
      {{ $t('declareTicket.errors.ticketRequired') }}
    </c-alert>
  </v-layout>
</template>

<script>
import { keyBy, pick } from 'lodash';

import { formMixin } from '@/mixins/form';

import DeclareTicketEventForm from './declare-ticket-event-form.vue';
import DeclareTicketEventAlarmsTicketsField from './fields/declare-ticket-event-alarms-tickets-field.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketEventAlarmsTicketsField, DeclareTicketEventForm },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    alarms: {
      type: Array,
      default: () => [],
    },
    ticketsByAlarms: {
      type: Object,
      required: () => ({}),
    },
    alarmsByTickets: {
      type: Object,
      required: () => ({}),
    },
    tickets: {
      type: Array,
      default: () => [],
    },
    hideTicketResource: {
      type: Boolean,
      default: false,
    },
    singleMode: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    declareTicketFieldName() {
      return 'declare_event';
    },

    hasErrors() {
      return this.errors.has(this.declareTicketFieldName);
    },

    alarmsById() {
      return keyBy(this.alarms, '_id');
    },

    alarmsByTicketId() {
      return Object.entries(this.alarmsByTickets).reduce((acc, [ticketId, { alarms: alarmsIds }]) => {
        acc[ticketId] = alarmsIds.map(id => this.alarmsById[id]);

        return acc;
      }, {});
    },

    alarmsWithoutTickets() {
      return Object.keys(this.ticketsByAlarms).reduce((acc, alarmId) => {
        if (this.ticketsByAlarms[alarmId].length === 0) {
          acc.push(this.alarmsById[alarmId]);
        }

        return acc;
      }, []);
    },

    commonValue() {
      return {
        alarms_by_tickets: this.form.alarms_by_tickets,
        comment: Object.values(this.form.comments_by_tickets)[0] ?? '',
        ticket_resources: Object.values(this.form.ticket_resources_by_tickets)[0] ?? false,
      };
    },

    groups() {
      return Object.entries(this.alarmsByTickets).map(([ticketId, { alarms: alarmsIds, name }]) => {
        const selectedAlarms = this.form.alarms_by_tickets[ticketId];

        const everyEnabled = selectedAlarms.length === alarmsIds.length;
        const someEnabled = !everyEnabled && selectedAlarms.length > 0;
        const activeTickets = [{ _id: ticketId, name }];

        return {
          everyEnabled,
          someEnabled,
          value: {
            alarms_by_tickets: pick(this.form.alarms_by_tickets, [ticketId]),
            comment: this.form.comments_by_tickets[ticketId],
            ticket_resources: this.form.ticket_resources_by_tickets[ticketId],
          },
          name,
          ticketId,
          alarms: this.alarmsByTicketId[ticketId],
          ticketsByAlarms: alarmsIds.reduce((acc, id) => {
            acc[id] = activeTickets;

            return acc;
          }, {}),
        };
      });
    },
  },
  watch: {
    'form.alarms_by_tickets': {
      handler() {
        if (this.hasErrors) {
          this.$validator.validate(this.name);
        }
      },
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
        name: this.declareTicketFieldName,
        rules: 'min_value:1',
        getter: () => Object.values(this.form.alarms_by_tickets).filter(alarms => alarms.length).length,
        vm: this,
      });
    },

    detachRules() {
      this.$validator.detach(this.declareTicketFieldName);
    },

    updateEnabledTickets(ticketId) {
      const { alarms } = this.alarmsByTickets[ticketId];
      const selectedAlarms = this.form.alarms_by_tickets[ticketId];

      this.updateField(
        'alarms_by_tickets',
        {
          ...this.form.alarms_by_tickets,
          [ticketId]: selectedAlarms.length === alarms.length ? [] : alarms,
        },
      );
    },

    updateCommonValue({ comment, alarms_by_tickets: alarmsByTickets, ticket_resources: ticketResources }) {
      this.updateModel(
        {
          alarms_by_tickets: alarmsByTickets,
          ...Object.keys(this.form.comments_by_tickets).reduce((acc, ticketId) => {
            acc.comments_by_tickets[ticketId] = comment;
            acc.ticket_resources_by_tickets[ticketId] = ticketResources;

            return acc;
          }, {
            comments_by_tickets: {},
            ticket_resources_by_tickets: {},
          }),
        },
      );
    },

    updateGroup(ticketId, { comment, alarms_by_tickets: alarmsByTickets, ticket_resources: ticketResources }) {
      this.updateModel(
        {
          alarms_by_tickets: {
            ...this.form.alarms_by_tickets,
            ...alarmsByTickets,
          },
          comments_by_tickets: {
            ...this.form.comments_by_tickets,
            [ticketId]: comment,
          },
          ticket_resources_by_tickets: {
            ...this.form.ticket_resources_by_tickets,
            [ticketId]: ticketResources,
          },
        },
      );
    },
  },
};
</script>
