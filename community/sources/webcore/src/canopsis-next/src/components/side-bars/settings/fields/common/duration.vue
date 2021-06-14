<template lang="pug">
  v-list-group
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass")  {{ label }}
    v-container
      c-duration-field(
        v-field="value",
        :units-label="$t('common.unit')",
        :name="name",
        clearable
      )
</template>

<script>
import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

export default {
  inject: ['$validator'],
  mixins: [formMixin, formValidationHeaderMixin],
  props: {
    value: {
      type: Object,
      required: false,
    },
    label: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: false,
    },
  },
};
</script>
