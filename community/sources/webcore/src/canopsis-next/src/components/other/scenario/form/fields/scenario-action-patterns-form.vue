<template lang="pug">
  div.mt-4
    c-patterns-field(
      v-field="patterns",
      :name="name",
      some-required,
      with-alarm,
      with-entity
    )
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

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
    name: {
      type: String,
      default: 'patterns',
    },
  },
};
</script>
