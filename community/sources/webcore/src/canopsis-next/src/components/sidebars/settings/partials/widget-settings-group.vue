<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ title }}
    v-list.grey.lighten-4.px-2.py-0(expand)
      slot
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  props: {
    title: {
      type: String,
      required: true,
    },
  },
};
</script>
