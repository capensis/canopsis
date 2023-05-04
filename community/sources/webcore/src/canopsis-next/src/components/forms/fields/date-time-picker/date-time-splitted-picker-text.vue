<template lang="pug">
  v-layout.date-time-splitted-field(column)
    v-layout(row, :reverse="reverse")
      v-flex.date-time-splitted-field__date
        date-picker-text(
          :value="value | date('vuetifyDatePicker', null)",
          :label="!reverse || fullDay ? label : ''",
          hide-details
        )
      time-picker-text.date-time-splitted-field__time(
        v-if="!fullDay",
        :value="value | date('timePicker', null)",
        :label="reverse ? label : ''",
        hide-details
      )
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
