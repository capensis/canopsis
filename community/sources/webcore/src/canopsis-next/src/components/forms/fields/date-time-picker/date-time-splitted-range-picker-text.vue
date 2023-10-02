<template lang="pug">
  v-layout.date-time-splitted-range-picker-field(row)
    v-flex.date-time-splitted-range-picker-field__start
      date-time-splitted-picker-text(
        :value="start",
        :full-day="fullDay",
        :label="startLabel"
      )
    template(v-if="!noEnding")
      div.date-time-splitted-range-picker-field__time-dash –
      v-flex.date-time-splitted-range-picker-field__end
        date-time-splitted-picker-text(
          :value="end",
          :full-day="fullDay",
          :label="endLabel",
          reverse
        )
</template>

<script>
import DateTimeSplittedPickerText from './date-time-splitted-picker-text.vue';

export default {
  components: { DateTimeSplittedPickerText },
  props: {
    start: {
      type: Date,
      default: null,
    },
    end: {
      type: Date,
      default: null,
    },
    startLabel: {
      type: String,
      required: false,
    },
    endLabel: {
      type: String,
      required: false,
    },
    noEnding: {
      type: Boolean,
      default: false,
    },
    fullDay: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
