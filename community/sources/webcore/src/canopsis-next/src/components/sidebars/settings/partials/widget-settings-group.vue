<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile
        div(:class="validationHeaderClass") {{ title }}
    v-list.grey.px-2.py-0(:class="{ 'darken-2': $vuetify.dark, 'lighten-4': !$vuetify.dark }", expand)
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
