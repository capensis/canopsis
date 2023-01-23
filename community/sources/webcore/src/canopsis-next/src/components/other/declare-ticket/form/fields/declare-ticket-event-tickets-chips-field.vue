<template lang="pug">
  v-layout.declare-ticket-event-tickets-chips-field(row, wrap)
    v-chip.declare-ticket-event-tickets-chips-field__chip(
      v-for="chip in chips",
      :key="chip.value",
      :selected="chip.active",
      :color="chip.active ? 'primary' : undefined",
      :text-color="chip.active ? 'white' : undefined",
      small,
      @click="updateActive(chip.value)"
    ) {{ chip.text }}
</template>

<script>
import { filterValue } from '@/helpers/entities';

import { formMixin } from '@/mixins/form';

export default {
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
  &__chip .v-chip__content {
    cursor: pointer;
  }
}
</style>
