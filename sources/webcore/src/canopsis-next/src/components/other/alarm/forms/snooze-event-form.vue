<template lang="pug">
  div
    duration-field(v-field="form.duration")
    v-layout(row)
      v-textarea(
        v-field="form.output",
        v-validate="isNoteRequired ? 'required' : ''",
        :label="$t('common.output')",
        :error-messages="errors.collect('output')",
        name="output"
      )
</template>

<script>
import DurationField from '@/components/forms/fields/duration.vue';

export default {
  inject: ['$validator'],
  components: { DurationField },
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
