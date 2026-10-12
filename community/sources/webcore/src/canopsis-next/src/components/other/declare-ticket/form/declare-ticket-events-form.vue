<template>
  <c-layout column>
    <c-information-block v-if="!singleMode" :title="$t('declareTicket.applyRules')">
      <declare-ticket-event-form
        :form="commonValue"
        :alarms="alarms"
        :tickets-by-alarms="ticketsByAlarms"
        :hide-ticket-resource="hideTicketResource"
        autofocus
        hide-row-select
        @update:form="updateCommonValue"
      />
    </c-information-block>
    <template v-else>
      <c-information-block
        v-for="(group, index) in groups"
        :key="group.ticketId"
        :title="$t('declareTicket.applyRules')"
      >
        <c-layout class="mt-2">
          <c-checkbox
            :model-value="group.everyEnabled"
            :indeterminate="group.someEnabled"
            :label="group.name"
            color="primary"
            hide-details
            @update:model-value="updateEnabledTickets(group.ticketId)"
          />
        </c-layout>
        <declare-ticket-event-form
          :form="group.value"
          :alarms="group.alarms"
          :tickets-by-alarms="group.ticketsByAlarms"
          :hide-ticket-resource="hideTicketResource"
          :autofocus="!index"
          hide-tickets
          @update:form="updateGroup(group.ticketId, $event)"
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
    <c-alert v-if="hasErrors" type="error">
      {{ $t('declareTicket.errors.ticketRequired') }}
    </c-alert>
  </c-layout>
</template>

<script>
import { keyBy, pick } from 'lodash';
import { computed, onBeforeUnmount, watch } from 'vue';

import { useFormErrorsMap } from '@/hooks/form';
import { useModelField } from '@/hooks/form/model-field';
import { useValidationAttachRule } from '@/hooks/validator/validation-attach-rule';

import DeclareTicketEventForm from './declare-ticket-event-form.vue';
import DeclareTicketEventAlarmsTicketsField from './fields/declare-ticket-event-alarms-tickets-field.vue';

const DECLARE_TICKET_FIELD_NAME = 'declare_event';

export default {
  components: { DeclareTicketEventAlarmsTicketsField, DeclareTicketEventForm },
  emits: ['update:form'],
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
  setup(props, { emit }) {
    const { hasError } = useFormErrorsMap();
    const { updateField, updateModel } = useModelField(props, emit);
    const { attachRule, detachRule, validateRule } = useValidationAttachRule(
      DECLARE_TICKET_FIELD_NAME,
      'min_value:1',
      () => Object.values(props.form.alarms_by_tickets).filter(alarms => alarms.length).length,
    );

    const declareTicketFieldName = DECLARE_TICKET_FIELD_NAME;
    const hasErrors = computed(() => hasError(declareTicketFieldName));
    const alarmsById = computed(() => keyBy(props.alarms, '_id'));
    const alarmsByTicketId = computed(() =>
      Object.entries(props.alarmsByTickets).reduce(
        (acc, [ticketId, { alarms: alarmsIds }]) => {
          acc[ticketId] = alarmsIds.map(id => alarmsById.value[id]);

          return acc;
        },
        {},
      ),
    );

    const alarmsWithoutTickets = computed(() =>
      Object.keys(props.ticketsByAlarms).reduce((acc, alarmId) => {
        if (props.ticketsByAlarms[alarmId].length === 0) {
          acc.push(alarmsById.value[alarmId]);
        }

        return acc;
      }, []),
    );

    const commonValue = computed(() => ({
      alarms_by_tickets: props.form.alarms_by_tickets,
      comment: Object.values(props.form.comments_by_tickets)[0] ?? '',
      ticket_resources: Object.values(props.form.ticket_resources_by_tickets)[0] ?? false,
    }));

    const groups = computed(() =>
      Object.entries(props.alarmsByTickets).map(([ticketId, { alarms: alarmsIds, name }]) => {
        const selectedAlarms = props.form.alarms_by_tickets[ticketId];

        const everyEnabled = selectedAlarms.length === alarmsIds.length;
        const someEnabled = !everyEnabled && selectedAlarms.length > 0;
        const activeTickets = [{ _id: ticketId, name }];

        return {
          everyEnabled,
          someEnabled,
          value: {
            alarms_by_tickets: pick(props.form.alarms_by_tickets, [ticketId]),
            comment: props.form.comments_by_tickets[ticketId],
            ticket_resources: props.form.ticket_resources_by_tickets[ticketId],
          },
          name,
          ticketId,
          alarms: alarmsByTicketId.value[ticketId],
          ticketsByAlarms: alarmsIds.reduce((acc, id) => {
            acc[id] = activeTickets;

            return acc;
          }, {}),
        };
      }),
    );

    const updateEnabledTickets = (ticketId) => {
      const { alarms } = props.alarmsByTickets[ticketId];
      const selectedAlarms = props.form.alarms_by_tickets[ticketId];

      updateField('alarms_by_tickets', {
        ...props.form.alarms_by_tickets,
        [ticketId]: selectedAlarms.length === alarms.length ? [] : alarms,
      });
    };

    const updateCommonValue = ({
      comment,
      alarms_by_tickets: alarmsByTickets,
      ticket_resources: ticketResources,
    }) => {
      updateModel({
        alarms_by_tickets: alarmsByTickets,
        ...Object.keys(props.form.comments_by_tickets).reduce(
          (acc, ticketId) => {
            acc.comments_by_tickets[ticketId] = comment;
            acc.ticket_resources_by_tickets[ticketId] = ticketResources;

            return acc;
          },
          {
            comments_by_tickets: {},
            ticket_resources_by_tickets: {},
          },
        ),
      });
    };

    const updateGroup = (
      ticketId,
      { comment, alarms_by_tickets: alarmsByTickets, ticket_resources: ticketResources },
    ) => {
      updateModel({
        alarms_by_tickets: {
          ...props.form.alarms_by_tickets,
          ...alarmsByTickets,
        },
        comments_by_tickets: {
          ...props.form.comments_by_tickets,
          [ticketId]: comment,
        },
        ticket_resources_by_tickets: {
          ...props.form.ticket_resources_by_tickets,
          [ticketId]: ticketResources,
        },
      });
    };

    attachRule();

    watch(
      () => props.form.alarms_by_tickets,
      () => {
        if (hasErrors.value) {
          validateRule();
        }
      },
    );

    onBeforeUnmount(detachRule);

    return {
      declareTicketFieldName,
      hasErrors,
      alarmsById,
      alarmsByTicketId,
      alarmsWithoutTickets,
      commonValue,
      groups,
      updateEnabledTickets,
      updateCommonValue,
      updateGroup,
    };
  },
};
</script>
