<template>
  <v-layout column>
    <c-duration-field
      v-field="form.duration"
      required
    />
    <c-description-field
      v-field="form.comment"
      :label="$t('common.note')"
      :required="isNoteRequired"
      name="comment"
    />
  </v-layout>
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
