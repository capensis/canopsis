<template lang="pug">
  div.d-inline-flex
    date-time-picker-field(
      v-validate,
      v-field="value.from",
      :label="$t('common.from')",
      :disabled="disabled",
      :name="fromFieldName",
      :hide-details="hideDetails"
    )
    date-time-picker-field.pl-3(
      v-validate,
      v-field="value.to",
      :label="$t('common.to')",
      :disabled="disabled",
      :name="toFieldName",
      :hide-details="hideDetails"
    )
</template>

<script>
import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimePickerField,
  },
  model: {
    event: 'input',
    prop: 'value',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'interval',
    },
    disabled: {
      type: Boolean,
      required: false,
    },
    hideDetails: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    fromFieldName() {
      return `${this.name}.from`;
    },

    toFieldName() {
      return `${this.name}.to`;
    },
  },
};
</script>
