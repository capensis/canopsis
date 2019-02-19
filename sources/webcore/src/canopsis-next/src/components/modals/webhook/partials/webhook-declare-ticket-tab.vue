<template lang="pug">
  div
    h2 Declare ticket
    v-layout
      v-flex
        v-textarea(v-model="input")
</template>

<script>
import formMixin from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'declareTicket',
    event: 'input',
  },
  props: {
    declareTicket: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      input: JSON.stringify(this.declareTicket),
    };
  },
};
</script>
