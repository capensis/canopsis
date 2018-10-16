<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.tstop') }}
    v-container
      template(v-if="durationUnit === 'm'")
        v-menu(
        ref="menu"
        :close-on-content-click="false",
        transition="scale-transition",
        )
          v-text-field(slot="activator", :value="dateObject | date('MM/YYYY', true)", readonly)
          v-date-picker(@input="updateValue", :allowed-dates="allowedDates")
      template(v-else)
        date-time-picker(:value="dateObject", @input="updateValue", roundHours)
</template>

<script>
import moment from 'moment-timezone';
import DateTimePicker from '@/components/forms/date-time-picker.vue';

export default {
  components: {
    DateTimePicker,
  },
  props: {
    value: {
      type: Number,
      required: true,
    },
    duration: {
      type: String,
      required: true,
    },
  },
  computed: {
    dateObject() {
      return moment.unix(this.value).toDate();
    },
    durationUnit() {
      return this.duration.slice(-1);
    },
  },
  watch: {
    // Each time the duration change on settings, we need to check if the duration unit is 'm' (for 'Month')
    // If it's 'm', we need to put the date on the 1st day of the month, at 00:00 (UTC)
    durationUnit(value, oldValue) {
      if (value !== oldValue && value === 'm') {
        const date = moment.tz(this.value * 1000, moment.tz.guess()).startOf('month');
        this.$emit('input', date.add(date.utcOffset(), 'm').unix());
      }
    },
  },
  methods: {
    updateValue(value) {
      /**
       * If the duration is 'm' (for 'Month'), we need to put the date on the first day of the month,
       * at 00:00 (UTC)
       */
      if (this.durationUnit === 'm') {
        // Get the value's date object, and put it on the first day of the month, at 00:00 (local date)
        const date = moment.tz(value, moment.tz.guess()).startOf('month');
        // Add the difference between the local date, and the UTC one.
        this.$emit('input', date.add(date.utcOffset(), 'm').unix());
      } else {
        this.$emit('input', moment(value).unix());
      }
    },
    allowedDates: val => parseInt(val.split('-')[2], 10) === 1,
  },
};
</script>
