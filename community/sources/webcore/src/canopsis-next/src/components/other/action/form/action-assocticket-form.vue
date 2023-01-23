<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="value.ticket",
        v-validate="'required'",
        :error-messages="errors.collect(ticketFieldName)",
        :name="ticketFieldName",
        :label="$tc('common.ticket')"
      )
    action-note-form(v-field="value")
</template>

<script>
import ActionNoteForm from './action-note-form.vue';

export default {
  inject: ['$validator'],
  components: { ActionNoteForm },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: 'parameters',
    },
  },
  computed: {
    ticketFieldName() {
      return `${this.name}.ticket`;
    },

    outputFieldName() {
      return `${this.name}.output`;
    },
  },
};
</script>
