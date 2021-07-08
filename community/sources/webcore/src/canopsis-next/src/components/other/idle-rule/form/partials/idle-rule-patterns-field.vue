<template lang="pug">
  div
    patterns-list(
      v-if="isEntityType",
      v-field="form.entity_patterns",
      v-validate="'required'",
      name="entity_patterns"
    )
    c-patterns-field(
      v-else,
      v-field="form",
      alarm,
      entity
    )
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

import PatternsList from '@/components/common/patterns-list/patterns-list.vue';

export default {
  inject: ['$validator'],
  components: { PatternsList },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isEntityType: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
