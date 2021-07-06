<template lang="pug">
  v-list-group(data-test="columnNames")
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass") {{ label }}
    c-columns-field(v-field="columns", :with-html="withHtml", :with-color-indicator="withColorIndicator")
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import ColorIndicatorField from '../partials/color-indicator-field.vue';

export default {
  inject: ['$validator'],
  components: { ColorIndicatorField },
  mixins: [
    formValidationHeaderMixin,
  ],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      required: true,
    },
    columns: {
      type: [Array, Object],
      default: () => [],
    },
    withHtml: {
      type: Boolean,
      default: false,
    },
    withColorIndicator: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
