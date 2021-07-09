<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    v-textarea(
      v-field="form.description",
      v-validate="'required'",
      :label="$t('common.description')",
      :error-messages="errors.collect('description')",
      name="description"
    )
    v-layout(row, justify-space-between)
      v-flex(xs7)
        c-duration-field(v-field="form.duration", required)
      v-flex(xs3)
        c-priority-field(v-field="form.priority", required)
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
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
  },
};
</script>
