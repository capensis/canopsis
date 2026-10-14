<template>
  <c-layout gap="2" column>
    <c-lazy-search-field
      :value="form.ticket"
      :items="associatedTickets"
      :label="t('modals.removeAssociatedTicketEvent.associatedTicketLabel')"
      :hint="t('modals.removeAssociatedTicketEvent.associatedTicketHint')"
      :disabled="ticketDisabled"
      item-text="text"
      item-value="value"
      name="ticket"
      with-type
      persistent-hint
      required
      autocomplete
      @update:value="updateField('ticket', $event)"
    />
    <c-description-field
      :model-value="form.comment"
      :label="t('common.reason')"
      :max-length="FORM_MAX_STRING_LENGTH"
      name="reason"
      required
      @update:model-value="updateField('comment', $event)"
    />
  </c-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import { FORM_MAX_STRING_LENGTH } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    form: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    ticketDisabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
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
      t,
      FORM_MAX_STRING_LENGTH,
      associatedTickets,
      updateField,
    };
  },
};
</script>
