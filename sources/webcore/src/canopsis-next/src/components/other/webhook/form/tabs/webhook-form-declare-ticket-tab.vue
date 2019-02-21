<template lang="pug">
  div
    h2 {{ $t('webhook.tabs.declareTicket.title') }}
    text-pairs(
    :textLabel="$t('webhook.tabs.declareTicket.fields.text')",
    :valueLabel="$t('webhook.tabs.declareTicket.fields.value')",
    :items="declareTicket",
    :disabled="disabled",
    @input="$emit('input', $event)"
    )
</template>

<script>
import TextPairs from '@/components/forms/fields/text-pairs.vue';

export default {
  components: { TextPairs },
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
  },
};
</script>
