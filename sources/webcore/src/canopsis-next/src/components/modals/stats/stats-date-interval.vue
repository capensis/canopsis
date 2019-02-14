<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Stats - Date interval
    v-card-text
      v-container
        v-select.pt-0(
        :items="periodUnits",
        v-model="periodUnit",
        label="Period",
        )
        stats-date-selector(v-model="dateForm", :periodUnit="periodUnit")
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import StatsDateSelector from '@/components/forms/stats-date-selector.vue';

export default {
  name: MODALS.statsDateInterval,
  components: {
    StatsDateSelector,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      periodUnit: 'h',
      dateForm: {
        tstart: 'now+1d',
        tstop: 'now+2d',
      },
      periodUnits: [
        {
          text: this.$tc('common.times.hour'),
          value: 'h',
        },
        {
          text: this.$tc('common.times.day'),
          value: 'd',
        },
        {
          text: this.$tc('common.times.week'),
          value: 'w',
        },
        {
          text: this.$tc('common.times.month'),
          value: 'm',
        },
      ],
    };
  },
  mounted() {
    if (this.config.interval) {
      const { periodUnit, tstart, tstop } = this.config.interval;

      this.periodUnit = periodUnit;
      this.dateForm = { ...this.dateForm, tstart, tstop };
    }
  },
  methods: {
    async submit() {
      if (this.config.action) {
        this.config.action({ periodUnit: this.periodUnit, tstart: this.dateForm.tstart, tstop: this.dateForm.tstop });
      }

      this.hideModal();
    },
  },
};
</script>
