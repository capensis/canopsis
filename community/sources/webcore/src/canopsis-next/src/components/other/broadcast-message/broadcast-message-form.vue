<template lang="pug">
  div
    v-layout(row)
      v-flex(xs12)
        text-editor-field(
          v-field="form.message",
          v-validate="'required'",
          :label="$t('common.message')",
          :error-messages="errors.collect('message')",
          name="message",
          public
        )
    v-layout(row)
      c-color-picker-field(v-field="form.color")
    v-layout(row)
      date-time-picker-field(
        v-validate="startRules",
        :value="form.start",
        :label="$t('common.start')",
        :error-message="errors.collect('start')",
        name="start",
        @input="updateField('start', $event)"
      )
    v-layout(row)
      date-time-picker-field(
        v-validate="endRules",
        :value="form.end",
        :label="$t('common.end')",
        :error-message="errors.collect('end')",
        name="end",
        @input="updateField('end', $event)"
      )
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { formMixin } from '@/mixins/form';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  inject: ['$validator'],
  components: { TextEditorField, DateTimePickerField },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  computed: {
    startRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    endRules() {
      return {
        required: true,
        after: [convertDateToString(this.form.start, DATETIME_FORMATS.dateTimePicker)],
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },
  },
};
</script>
