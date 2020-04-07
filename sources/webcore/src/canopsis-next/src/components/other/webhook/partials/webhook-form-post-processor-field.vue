<template lang="pug">
  v-card.pa-4.mt-2
    v-switch(
      v-if="!disabled",
      v-field="postProcessor.emptyResponse",
      :label="$t('webhook.tabs.declareTicket.emptyResponse')",
      color="primary"
    )
    v-text-field(
      v-field="postProcessor.ticketId",
      v-validate="'required'",
      :label="$t('webhook.tabs.declareTicket.ticketId')",
      :disabled="disabled",
      :error-messages="errors.collect(getFieldName('ticketId'))",
      :name="getFieldName('ticketId')"
    )
    text-pairs(
      v-field="postProcessor.fields",
      :textLabel="$t('webhook.tabs.declareTicket.fields.text')",
      :valueLabel="$t('webhook.tabs.declareTicket.fields.value')",
      :disabled="disabled",
      :name="getFieldName('fields')",
      valueValidationRules="required",
      mixed
    )
</template>

<script>
import formMixin from '@/mixins/form';

import TextPairs from '@/components/forms/fields/text-pairs.vue';

export default {
  inject: ['$validator'],
  components: { TextPairs },
  mixins: [formMixin],
  model: {
    prop: 'postProcessor',
    event: 'input',
  },
  props: {
    postProcessor: {
      type: Object,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'postProcessor',
    },
  },
  methods: {
    getFieldName(name) {
      return `${this.name}.${name}`;
    },
  },
};
</script>
