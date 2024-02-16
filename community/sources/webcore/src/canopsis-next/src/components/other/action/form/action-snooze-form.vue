<template>
  <div>
    <action-note-form v-field="value" />
    <c-duration-field
      v-field="value.duration"
      :name="`${name}.duration`"
      required
    />
  </div>
</template>

<script>
import ActionNoteForm from './action-note-form.vue';

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
};
</script>
