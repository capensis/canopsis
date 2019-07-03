<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.statsDateInterval.title') }}
    v-card-text
      v-container
        v-layout
          v-flex(xs3, v-if="!hiddenFields.includes('periodValue')")
            v-text-field.pt-0(
            type="number",
            v-model="form.periodValue",
            :label="$t('modals.statsDateInterval.fields.periodValue')"
            )
          v-flex
            v-select.pt-0(
            v-model="form.periodUnit",
            :items="periodUnits",
            :label="$t('modals.statsDateInterval.fields.periodUnit')"
            )
        v-alert.mb-2(
        :value="isPeriodMonth",
        type="info",
        ) {{ $t('settings.statsDateInterval.monthPeriodInfo') }}
        stats-date-selector.my-1(
        v-model="form",
        :periodUnit="form.periodUnit",
        @update:tstartObjectValue="objectValues.tstart = $event",
        @update:tstopObjectValue="objectValues.tstop = $event"
        )
        v-alert.mb-2(
        :value="isPeriodMonth"
        type="info",
        ) {{ objectValues.tstart | date('long', true) }} - {{ objectValues.tstop | date('long', true) }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit", :disabled="errors.any()") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, DATETIME_FORMATS, STATS_DURATION_UNITS, STATS_QUICK_RANGES } from '@/constants';

import { dateParse } from '@/helpers/date-intervals';

import modalInnerMixin from '@/mixins/modal/inner';

import StatsDateSelector from '@/components/forms/stats-date-selector.vue';

export default {
  name: MODALS.statsDateInterval,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    StatsDateSelector,
  },
  mixins: [modalInnerMixin],
  data() {
    const { interval } = this.modal.config;
    const defaultInterval = {
      periodValue: 1,
      periodUnit: STATS_DURATION_UNITS.hour,
      tstart: STATS_QUICK_RANGES.thisMonthSoFar.start,
      tstop: STATS_QUICK_RANGES.thisMonthSoFar.stop,
    };

    return {
      objectValues: {
        tstart: null,
        tstop: null,
      },
      form: cloneDeep(interval || defaultInterval),
    };
  },
  computed: {
    periodUnits() {
      return [
        {
          text: this.$tc('common.times.hour'),
          value: STATS_DURATION_UNITS.hour,
        },
        {
          text: this.$tc('common.times.day'),
          value: STATS_DURATION_UNITS.day,
        },
        {
          text: this.$tc('common.times.week'),
          value: STATS_DURATION_UNITS.week,
        },
        {
          text: this.$tc('common.times.month'),
          value: STATS_DURATION_UNITS.month,
        },
      ];
    },

    hiddenFields() {
      return this.modal.config.hiddenFields || [];
    },

    isPeriodMonth() {
      return this.form.periodUnit === STATS_DURATION_UNITS.month;
    },
  },
  created() {
    this.$validator.extend('after_custom', {
      getMessage: () => this.$t('modals.statsDateInterval.errors.endDateLessOrEqualStartDate'),
      validate: (value, [otherValue]) => {
        try {
          const convertedStop = dateParse(value, 'stop', DATETIME_FORMATS.dateTimePicker);
          const convertedStart = dateParse(otherValue, 'start', DATETIME_FORMATS.dateTimePicker);

          return !convertedStop.isSameOrBefore(convertedStart);
        } catch (err) {
          return true; // TODO: problem with i18n: https://github.com/baianat/vee-validate/issues/2025
        }
      },
    }, {
      hasTarget: true,
    });
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          this.config.action(this.form);
        }

        this.hideModal();
      }
    },
  },
};
</script>
