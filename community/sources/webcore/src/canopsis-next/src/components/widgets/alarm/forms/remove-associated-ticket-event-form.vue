<template>
  <v-layout column>
    <c-lazy-search-field
      v-field="form.ticket"
      :items="associatedTicketsOptions"
      :label="$t('modals.removeAssociatedTicketEvent.associatedTicketLabel')"
      :hint="associatedTicketsOptions.length ? $t('modals.removeAssociatedTicketEvent.associatedTicketHint') : ''"
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

    return {
      associatedTicketsOptions,
    };
  },
};
</script>
