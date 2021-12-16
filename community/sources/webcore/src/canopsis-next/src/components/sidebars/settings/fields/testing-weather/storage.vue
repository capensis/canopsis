<template lang="pug">
  v-list-group(:disabled="disabled", hide-on-disabled)
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass") {{ title }}
    v-container.pr-2
      c-storage-field(
        v-field="storage",
        v-on="$listeners",
        :required="required",
        :disabled="disabled"
      )
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'storage',
    event: 'input',
  },
  props: {
    storage: {
      type: String,
      default: '',
    },
    title: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
