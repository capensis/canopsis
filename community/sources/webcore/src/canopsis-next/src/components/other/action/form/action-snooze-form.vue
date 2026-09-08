<template>
  <div>
    <action-note-form v-field="value" :template-vars="templateVars" :depth="depth" />

    <c-form-block-row :label="$t('common.duration')" :depth="depth">
      <c-duration-field
        v-field="value.duration"
        :name="`${name}.duration`"
        required
      />
    </c-form-block-row>
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
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    depth: {
      type: Number,
      default: 0,
    },
  },
};
</script>
