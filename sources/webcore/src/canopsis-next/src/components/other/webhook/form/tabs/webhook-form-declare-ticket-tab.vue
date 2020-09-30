<template lang="pug">
  div
    v-switch(
      v-if="!disabled",
      :input-value="emptyResponse",
      :label="$t('webhook.tabs.declareTicket.emptyResponse')",
      @change="$emit('update:emptyResponse', $event)",
      color="primary"
    )
    text-pairs(
      :textLabel="$t('webhook.tabs.declareTicket.fields.text')",
      :valueLabel="$t('webhook.tabs.declareTicket.fields.value')",
      :items="declareTicket",
      :disabled="disabled",
      valueValidationRules="required",
      mixed,
      @input="$emit('input', $event)"
    )
</template>

<script>
import TextPairs from '@/components/forms/fields/text-pairs.vue';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

export default {
  inject: ['$validator'],
  components: { TextPairs },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'declareTicket',
    event: 'input',
  },
  props: {
    declareTicket: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    emptyResponse: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
