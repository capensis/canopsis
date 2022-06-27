<template lang="pug">
  v-layout(column)
    c-duration-field(v-field="form.duration", required)
    c-description-field(
      v-field="form.output",
      :label="$t('common.note')",
      :required="isNoteRequired",
      name="output"
    )
</template>

<script>
export default {
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
