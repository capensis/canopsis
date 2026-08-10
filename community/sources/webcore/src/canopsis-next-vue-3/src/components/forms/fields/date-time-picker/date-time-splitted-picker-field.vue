<template>
  <c-layout class="date-time-splitted-field" column>
    <c-layout :style="rowStyle">
      <div class="date-time-splitted-field__date">
        <c-date-picker-field
          :value="dateValue"
          :label="!reverse || fullDay ? label : ''"
          :placeholder="placeholder"
          :error="!!errorMessage"
          :disabled="disabled"
          :control-variant="controlVariant"
          :min="min"
          :max="max"
          header-color="primary"
          hide-details
          @update:value="updateDate"
        />
      </div>
      <time-picker-field
        v-if="!fullDay"
        :value="timeValue"
        :label="reverse ? label : ''"
        :error="!!errorMessage"
        :disabled="disabled"
        class="date-time-splitted-field__time"
        hide-details
        @update:value="updateTime"
      />
    </c-layout>
    <div v-if="errorMessage" class="v-text-field__details mt-2">
      <v-messages :messages="[errorMessage]" color="error" active />
    </div>
  </c-layout>
</template>

<script>
import TimePickerField from '@/components/forms/fields/time-picker/time-picker-field.vue';
import { DATETIME_FORMATS, TIME_UNITS } from '@/constants';
import { convertDateToStartOfUnitDateObject, convertDateToString } from '@/helpers/date/date';
import { getDateObjectByTime, getDateObjectByDate } from '@/helpers/date/date-time-picker';
import { formBaseMixin } from '@/mixins/form';

import { createDateTimeValidationRules } from './validation-rules';

export default {
  emits: ['update:value'],
  components: { TimePickerField },
  mixins: [formBaseMixin],
  props: {
    value: {
      type: Date,
      default: null,
    },
    fullDay: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'date',
    },
    reverse: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    min: {
      type: String,
      required: false,
    },
    max: {
      type: String,
      required: false,
    },
    rules: {
      type: Object,
      default: () => ({}),
    },
    controlVariant: {
      type: String,
      default: undefined,
    },
  },
  computed: {
    dateValue() {
      return convertDateToString(this.value, DATETIME_FORMATS.vuetifyDatePicker, null);
    },

    timeValue() {
      return convertDateToString(this.value, DATETIME_FORMATS.timePicker, null);
    },

    errorMessage() {
      if (this.rules.required && !this.value) {
        return 'Field is required';
      }

      const rule = createDateTimeValidationRules(this.rules).find((validationRule) => {
        const result = validationRule(this.validatedValue);

        return result !== true;
      });

      return rule ? rule(this.validatedValue) : '';
    },

    rowStyle() {
      return {
        flexDirection: this.reverse ? 'row-reverse' : undefined,
      };
    },

    validatedValue() {
      if (!this.value) {
        return this.value;
      }

      return convertDateToString(
        convertDateToStartOfUnitDateObject(this.value, TIME_UNITS.minute),
        DATETIME_FORMATS.dateTimePicker,
      );
    },
  },
  methods: {
    updateTime(time) {
      this.updateModel(getDateObjectByTime(this.value, time));
    },

    updateDate(date) {
      this.updateModel(getDateObjectByDate(this.value, date));
    },
  },
};
</script>

<style lang="scss">
.date-time-splitted-field {
  & > .layout,
  & > .c-layout {
    gap: 4px;
  }

  &__date {
    flex: 1;
    min-width: 90px;
  }

  &__time {
    width: 50px;
    max-width: 50px;
  }
}
</style>
