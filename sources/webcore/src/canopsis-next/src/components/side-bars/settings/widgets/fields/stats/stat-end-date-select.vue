<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.tstop') }}
    v-container
      template(v-if="isDurationUnitEqualMonth")
        v-menu(
        ref="menu"
        :close-on-content-click="false",
        v-model="menu",
        transition="scale-transition",
        )
          v-text-field(slot="activator", :value="dateObject | date('MM/YYYY', true)", readonly)
          date-picker(
          ref="picker",
          :value="dateString",
          @input="updateValue",
          @change="save",
          type="month",
          year-first,
          )
      template(v-else)
        date-time-picker-field(:value="dateObject", roundHours, @input="updateValue")
</template>

<script>
import moment from 'moment-timezone';

import DatePicker from '@/components/forms/fields/date-picker/date-picker.vue';
import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  components: {
    DatePicker,
    DateTimePickerField,
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
  data() {
    return {
      menu: false,
    };
  },
  computed: {
    dateObject() {
      return moment.unix(this.value).toDate();
    },
    dateString() {
      return moment.unix(this.value).format('YYYY-MM-DD');
    },
    durationUnit() {
      return this.duration.slice(-1);
    },
    isDurationUnitEqualMonth() {
      return this.durationUnit === this.$constants.STATS_DURATION_UNITS.month;
    },
  },
  watch: {
    /**
     * Each time the duration change on settings, we need to check if the duration unit is 'm' (for 'Month')
     * If it's 'm', we need to put the date on the 1st day of the month, at 00:00 (UTC)
    */
    durationUnit(value) {
      if (value && value === this.$constants.STATS_DURATION_UNITS.month) {
        const date = moment.tz(this.value * 1000, moment.tz.guess()).startOf('month');
        this.$emit('input', date.add(date.utcOffset(), 'm').unix());
      } else {
        this.$emit('input', moment().startOf('hour').unix());
      }
    },
    menu(value) {
      if (value) {
        setTimeout(() => {
          this.$refs.picker.$refs.years.scrollToActive();
        }, 50);
      }
    },
  },
  methods: {
    updateValue(value) {
      /**
       * If the duration is 'm' (for 'Month'), we need to put the date on the first day of the month,
       * at 00:00 (UTC)
       */
      if (this.durationUnit === this.$constants.STATS_DURATION_UNITS.month) {
        // Get the value's date object, and put it on the first day of the month, at 00:00 (local date)
        const date = moment.tz(value, moment.tz.guess()).startOf('month');
        // Add the difference between the local date, and the UTC one.
        this.$emit('input', date.add(date.utcOffset(), 'm').unix());
      } else {
        this.$emit('input', moment(value).unix());
      }
    },
    save(date) {
      this.$refs.menu.save(date);
    },
  },
};
</script>
