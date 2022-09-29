<template lang="pug">
  v-flex
    date-time-picker-field(
      v-validate="startRules",
      v-field="form.start",
      :label="$t('common.start')",
      :error-message="errors.collect('start')",
      name="start"
    )
    date-time-picker-field.ml-2(
      v-validate="stopRules",
      v-field="form.stop",
      :label="$t('common.stop')",
      :error-message="errors.collect('stop')",
      name="stop"
    )
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    startRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    stopRules() {
      return {
        required: true,
        after: [convertDateToString(this.form.start, DATETIME_FORMATS.dateTimePicker)],
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },
  },
};
</script>
