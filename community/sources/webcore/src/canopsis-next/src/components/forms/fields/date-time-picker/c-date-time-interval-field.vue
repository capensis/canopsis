<template>
  <div class="d-inline-flex">
    <date-time-picker-field
      v-validate="rules"
      v-field="value.from"
      :label="$t('common.from')"
      :disabled="disabled"
      :name="fromFieldName"
      :hide-details="hideDetails"
    />
    <date-time-picker-field
      class="pl-3"
      v-validate="rules"
      v-field="value.to"
      :label="$t('common.to')"
      :disabled="disabled"
      :name="toFieldName"
      :hide-details="hideDetails"
    />
  </div>
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
    required: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    fromFieldName() {
      return `${this.name}.from`;
    },

    toFieldName() {
      return `${this.name}.to`;
    },
  },
};
</script>
