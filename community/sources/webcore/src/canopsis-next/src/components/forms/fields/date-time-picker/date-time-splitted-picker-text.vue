<template>
  <c-layout class="date-time-splitted-field" column>
    <c-layout :style="rowStyle" gap="2">
      <div class="date-time-splitted-field__date">
        <date-picker-text
          :value="dateValue"
          :label="!reverse || fullDay ? label : ''"
          hide-details
        />
      </div>
      <time-picker-text
        v-if="!fullDay"
        :value="timeValue"
        :label="reverse ? label : ''"
        class="date-time-splitted-field__time"
        hide-details
      />
    </c-layout>
  </c-layout>
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import DatePickerText from './date-picker-text.vue';
import TimePickerText from './time-picker-text.vue';

export default {
  components: { DatePickerText, TimePickerText },
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
    reverse: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    dateValue() {
      return convertDateToString(this.value, DATETIME_FORMATS.vuetifyDatePicker, null);
    },

    timeValue() {
      return convertDateToString(this.value, DATETIME_FORMATS.timePicker, null);
    },

    rowStyle() {
      return {
        flexDirection: this.reverse ? 'row-reverse' : undefined,
      };
    },
  },
};
</script>
