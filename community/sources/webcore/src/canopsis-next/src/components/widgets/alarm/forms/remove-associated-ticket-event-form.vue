<template>
  <v-layout column>
    <c-select-field
      v-field="form.ticket"
      :items="associatedTicketsOptions"
      :label="$t('modals.removeAssociatedTicketEvent.associatedTicketLabel')"
      :hint="associatedTicketsOptions.length ? $t('modals.removeAssociatedTicketEvent.associatedTicketHint') : ''"
      name="ticket"
      persistent-hint
      required
    />
    <c-description-field
      v-field="form.reason"
      :label="$t('modals.removeAssociatedTicketEvent.reasonLabel')"
      :max-length="255"
      name="reason"
      required
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

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
  setup(props) {
    const associatedTicketsOptions = computed(() => {
      if (!props.items.length) {
        return [];
      }

      const allTickets = [];
      props.items.forEach((alarm) => {
        const tickets = alarm.v?.tickets ?? [];
        tickets.forEach((ticket) => {
          if (ticket.ticket && !allTickets.find(t => t.value === ticket.ticket)) {
            allTickets.push({
              text: ticket.ticket,
              value: ticket.ticket,
            });
          }
        });
      });

      return allTickets;
    });

    return {
      associatedTicketsOptions,
    };
  },
};
</script>
