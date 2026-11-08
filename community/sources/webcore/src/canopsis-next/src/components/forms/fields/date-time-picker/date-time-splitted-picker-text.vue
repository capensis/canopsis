<template>
  <v-row :style="rowStyle" class="date-time-splitted-field">
    <v-col class="date-time-splitted-field__date">
      <date-picker-text
        :value="dateValue"
        :label="!reverse || fullDay ? label : ''"
        hide-details
      />
    </v-col>
    <v-col v-if="!fullDay" class="date-time-splitted-field__time" cols="auto">
      <time-picker-text
        :value="timeValue"
        :label="reverse ? label : ''"
        hide-details
      />
    </v-col>
  </v-row>
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
