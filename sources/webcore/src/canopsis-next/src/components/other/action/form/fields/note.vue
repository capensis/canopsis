<template lang="pug">
  v-layout.mt-4(row)
    v-textarea(
      v-field="value.output",
      :label="$t('modals.createAction.fields.output')"
    )
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>
