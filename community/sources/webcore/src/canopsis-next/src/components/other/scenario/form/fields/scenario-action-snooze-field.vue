<template lang="pug">
  div
    scenario-action-note-field(v-field="value")
    c-duration-field(v-field="value.duration", :name="`${name}.duration`", required)
</template>

<script>
import ScenarioActionNoteField from './scenario-action-note-field.vue';

export default {
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
};
</script>
