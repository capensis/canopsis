<template lang="pug">
  c-collapse-panel(:error="hasChildrenError", color="grey")
    template(#header="")
      span.white--text {{ title }}
    v-card
      v-card-text
        slot
</template>

<script>
import { validationChildrenMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [validationChildrenMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    title: {
      type: String,
      required: true,
    },
  },
};
</script>
