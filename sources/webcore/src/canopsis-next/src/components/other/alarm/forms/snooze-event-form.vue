<template lang="pug">
  div
    duration-field(v-field="form.duration")
    output-field(v-field="form.output", :required="isNoteRequired")
</template>

<script>
import DurationField from '@/components/forms/fields/duration.vue';
import OutputField from '@/components/forms/fields/output.vue';

export default {
  components: { DurationField, OutputField },
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
