<template lang="pug">
  v-card
    v-card-title(primary-title, class="blue darken-4 white--text")
      v-layout(justify-space-between, align-center)
        h2 Edit live reporting
        v-btn(@click="hideModal", icon, small)
          v-icon(color="white") close
    v-divider
    v-card-text
      h3 Interval de date
      v-layout(wrap)
        v-radio-group(@change="handleIntervalClick", v-model="selectedInterval")
          v-radio(
            v-for="interval in dateIntervals",
            :label="interval.text",
            :value="interval",
            :key="interval.value"
          )
      v-layout(wrap, v-if="isCustomRangeEnabled")
        v-flex(xs12)
          date-time-picker(v-model="tstart", clearable, label="tstart")
        v-flex(xs12)
          date-time-picker(v-model="tend", clearable, label="tend")
      v-btn(@click="handleSubmit", color="green darken-4 white--text", small) Apply
</template>

<script>
import moment from 'moment';
import { createNamespacedHelpers } from 'vuex';
import DateTimePicker from '@/components/forms/date-time-picker.vue';

const { mapActions: modalMapActions } = createNamespacedHelpers('modal');
const { mapActions: alarmsMapActions } = createNamespacedHelpers('alarm');
const { mapActions: alarmsListMapActions } = createNamespacedHelpers('alarmsList');

export default {
  name: 'edit-live-reporting',
  components: {
    DateTimePicker,
  },
  data() {
    return {
      selectedInterval: '',
      dateIntervals: [
        {
          text: 'Today',
          value: 'today',
        },
        {
          text: 'Yesterday',
          value: 'yesterday',
        },
        {
          text: 'Last 7 days',
          value: 'last7Days',
        },
        {
          text: 'Last 30 days',
          value: 'last30Days',
        },
        {
          text: 'This month',
          value: 'thisMonth',
        },
        {
          text: 'Last month',
          value: 'lastMonth',
        },
        {
          text: 'Custom range',
          value: 'custom',
        },
      ],
      tstart: new Date(),
      tend: new Date(),
    };
  },
  computed: {
    isCustomRangeEnabled() {
      return this.selectedInterval.value === 'custom';
    },
  },
  methods: {
    ...modalMapActions({
      hideModal: 'hide',
    }),
    ...alarmsMapActions({
      fetchAlarmList: 'fetchList',
    }),
    ...alarmsListMapActions(['addLiveReportingFilter']),

    handleIntervalClick(interval) {
      if (interval.value === 'today') {
        this.tstart = moment().startOf('day').toDate();
        this.tend = moment().toDate();
      }

      if (interval.value === 'yesterday') {
        this.tstart = moment().subtract(1, 'day').startOf('day').toDate();
        this.tend = moment().subtract(1, 'day').endOf('day').toDate();
      }

      if (interval.value === 'last7Days') {
        this.tstart = moment().subtract(7, 'day').toDate();
        this.tend = moment().toDate();
      }

      if (interval.value === 'last30Days') {
        this.tstart = moment().subtract(30, 'day').toDate();
        this.tend = moment().toDate();
      }

      if (interval.value === 'thisMonth') {
        this.tstart = moment().startOf('month').toDate();
        this.tend = moment().toDate();
      }

      if (interval.value === 'lastMonth') {
        this.tstart = moment().subtract(1, 'month').startOf('month').toDate();
        this.tend = moment().startOf('month').toDate();
      }
    },
    handleSubmit() {
      this.addLiveReportingFilter(this.selectedInterval);
      this.fetchAlarmList({ params: { tstart: this.tstart.getTime() / 1000, tstop: this.tend.getTime() / 1000 } });
      this.hideModal();
    },
  },
};
</script>
