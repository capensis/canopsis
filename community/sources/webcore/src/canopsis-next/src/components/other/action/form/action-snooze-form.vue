<template lang="pug">
  div
    action-note-form(v-field="value")
    c-duration-field(
      v-field="value.duration",
      :units="snoozeUnits",
      :name="`${name}.duration`",
      required
    )
</template>

<script>
import ActionNoteForm from './action-note-form.vue';
import { SNOOZE_DURATION_UNITS } from '@/constants';

export default {
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
    snoozeUnits() {
      return Object.values(SNOOZE_DURATION_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.value.duration.value || 0),
      }));
    },
  },
};
</script>
