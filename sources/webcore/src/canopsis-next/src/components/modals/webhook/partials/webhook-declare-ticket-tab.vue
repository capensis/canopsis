<template lang="pug">
  div
    h2 Declare ticket
    text-pairs(
    :textLabel="'Text'",
    :valueLabel="'Value'",
    :items="declareTicket",
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
  },
};
</script>
