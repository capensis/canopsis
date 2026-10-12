<template>
  <v-layout
    class="date-time-splitted-field"
    column
  >
    <v-layout :reverse="reverse">
      <v-flex class="date-time-splitted-field__date">
        <date-picker-text
          :value="value | date('vuetifyDatePicker', null)"
          :label="!reverse || fullDay ? label : ''"
          hide-details
        />
      </v-flex>
      <time-picker-text
        v-if="!fullDay"
        :value="value | date('timePicker', null)"
        :label="reverse ? label : ''"
        class="date-time-splitted-field__time"
        hide-details
      />
    </v-layout>
  </v-layout>
</template>

<script>
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
};
</script>
