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
          date-time-picker(v-model="tstop", clearable, label="tstop")
      v-btn(@click="handleSubmit", color="green darken-4 white--text", small) Apply
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import DateTimePicker from '@/components/forms/date-time-picker.vue';
import dateIntervals from '@/helpers/date-intervals';
import ModalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

const { mapActions: alarmsMapActions } = createNamespacedHelpers('alarm');
const { mapActions: alarmsListMapActions } = createNamespacedHelpers('alarmsList');

export default {
  name: MODALS.editLiveReporting,
  components: {
    DateTimePicker,
  },
  mixins: [ModalMixin],
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
      tstop: new Date(),
    };
  },
  computed: {
    isCustomRangeEnabled() {
      return this.selectedInterval.value === 'custom';
    },
  },
  methods: {
    ...alarmsMapActions({
      fetchAlarmList: 'fetchList',
    }),
    ...alarmsListMapActions(['addLiveReportingFilter']),

    handleIntervalClick(interval) {
      try {
        const { tstart, tstop } = dateIntervals[interval.value]();
        this.tstart = tstart;
        this.tstop = tstop;
      } catch (err) {
        console.error(err);
      }
    },
    handleSubmit() {
      this.addLiveReportingFilter(this.selectedInterval);
      this.fetchAlarmList({ params: { tstart: this.tstart.getTime() / 1000, tstop: this.tstop.getTime() / 1000 } });
      this.hideModal();
    },
  },
};
</script>
