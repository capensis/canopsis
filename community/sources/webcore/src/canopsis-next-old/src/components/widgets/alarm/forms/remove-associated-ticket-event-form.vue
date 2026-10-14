<template>
  <v-layout column>
    <c-lazy-search-field
      v-field="form.ticket"
      :items="associatedTickets"
      :label="$t('modals.removeAssociatedTicketEvent.associatedTicketLabel')"
      :hint="$t('modals.removeAssociatedTicketEvent.associatedTicketHint')"
      item-text="text"
      item-value="value"
      name="ticket"
      with-type
      persistent-hint
      required
      autocomplete
    />
    <c-description-field
      v-field="form.comment"
      :label="$t('common.reason')"
      :max-length="255"
      name="reason"
      required
    />
  </v-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const associatedTickets = computed(() => {
      if (!props.items.length) {
        return [];
      }

      const allTicketsMap = props.items.reduce((acc, alarm) => {
        const tickets = alarm.v?.tickets ?? [];

        tickets.forEach((ticket) => {
          if (ticket.ticket && !acc[ticket.ticket]) {
            acc[ticket.ticket] = {
              text: ticket.ticket,
              value: ticket.ticket,
              type: ticket.ticket_system_name,
            };
          }
        });

        return acc;
      }, {});

      return Object.values(allTicketsMap);
    });

    onMounted(() => {
      if (associatedTickets.value.length === 1 && !props.form.ticket) {
        updateField('ticket', associatedTickets.value[0].value);
      }
    });

    return {
      associatedTickets,
    };
  },
};
</script>
