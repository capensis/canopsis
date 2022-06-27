<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile
        div(:class="validationHeaderClass") {{ title }}
          span.font-italic.caption.ml-1(v-if="optional") ({{ $t('common.optional') }})
    v-container
      enabled-limit-field(v-field="value", :label="label")
</template>

<script>
import { formMixin, formValidationHeaderMixin } from '@/mixins/form';

import EnabledLimitField from '@/components/forms/fields/enabled-limit-field.vue';

export default {
  inject: ['$validator'],
  components: { EnabledLimitField },
  mixins: [formMixin, formValidationHeaderMixin],
  props: {
    value: {
      type: Object,
      default: () => ({
        enabled: false,
        limit: 1,
      }),
    },
    title: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    optional: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
