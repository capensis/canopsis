<template>
  <v-layout column="column">
    <v-checkbox
      v-for="chip in chips"
      :key="chip.value"
      :input-value="chip.active"
      :label="chip.text"
      color="primary"
      hide-details="hide-details"
      @change="updateActive(chip.value)"
    />
  </v-layout>
</template>

<script>
import { filterValue } from '@/helpers/array';

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
