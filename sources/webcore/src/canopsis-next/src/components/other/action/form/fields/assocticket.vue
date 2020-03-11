<template lang="pug">
  div.mt-4
    v-layout(row)
      v-text-field(
        v-field="value.ticket",
        :label="$t('modals.createAction.fields.ticket')",
        data-test="createAckEventTicket"
      )
    v-layout(row)
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
