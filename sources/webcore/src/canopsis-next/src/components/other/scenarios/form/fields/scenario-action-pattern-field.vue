<template lang="pug">
  div
    c-patterns-field(
      v-field="patterns",
      :disabled="disabled",
      :name="name",
      required,
      alarm,
      entity
    )
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'patterns',
    },
  },
};
</script>
