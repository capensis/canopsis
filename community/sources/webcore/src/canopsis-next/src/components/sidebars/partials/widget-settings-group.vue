<template lang="pug">
  v-list-group.widget-settings-group
    template(#activator="")
      v-list-tile
        div(:class="validationHeaderClass") {{ title }}
    v-list.px-2.py-0.background.darken-1(expand)
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
