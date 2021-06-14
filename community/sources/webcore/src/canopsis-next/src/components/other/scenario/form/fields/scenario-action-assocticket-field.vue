<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="value.ticket",
        v-validate="'required'",
        :error-messages="errors.collect(ticketFieldName)",
        :name="ticketFieldName",
        :label="$t('common.ticket')"
      )
    scenario-action-note-field(v-field="value")
</template>

<script>
import ScenarioActionNoteField from './scenario-action-note-field.vue';

export default {
  inject: ['$validator'],
  components: { ScenarioActionNoteField },
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
