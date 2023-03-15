<template lang="pug">
  v-layout.declare-ticket-event-tickets-chips-field(row, wrap)
    declare-ticket-event-tickets-chip-field(
      v-for="chip in chips",
      :key="chip.value",
      :value="chip.active",
      :disabled="disabled",
      @input="updateActive(chip.value)"
    ) {{ chip.text }}
</template>

<script>
import { filterValue } from '@/helpers/entities';

import { formMixin } from '@/mixins/form';

import DeclareTicketEventTicketsChipField from './declare-ticket-event-tickets-chip-field.vue';

export default {
  components: { DeclareTicketEventTicketsChipField },
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    chips() {
      return this.tickets.map(({ _id: id, name }) => ({
        active: this.isActiveTicket(id),
        text: name,
        value: id,
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

<style lang="scss">
.declare-ticket-event-tickets-chips-field {
  margin: 5px;
  gap: 5px;
}
</style>
