<template lang="pug">
  div
    c-duration-field(v-field="form.duration", required)
    v-layout(row)
      v-textarea(
        v-field="form.output",
        v-validate="isNoteRequired ? 'required' : ''",
        :label="$t('common.note')",
        :error-messages="errors.collect('output')",
        name="output"
      )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    isNoteRequired: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
