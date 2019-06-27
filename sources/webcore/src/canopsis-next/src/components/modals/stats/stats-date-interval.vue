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
          v-select.pt-0(
          v-model="form.periodUnit",
          :items="periodUnits",
          :label="$t('modals.statsDateInterval.fields.periodUnit')"
          )
        v-alert.mb-2(
        v-if="form.periodUnit === $constants.STATS_DURATION_UNITS.month", type="info", value="true"
        ) {{ $t('settings.statsDateInterval.monthPeriodInfo') }}
        stats-date-selector.my-1(v-model="form", :periodUnit="form.periodUnit", @input="resetValidation")
      v-alert(
      value="errors",
      type="error",
      v-for="error in localErrors",
      :key="error",
      ) {{ error }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit", :disabled="localErrors.length !== 0") {{ $t('common.submit') }}
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
      form: cloneDeep(interval || defaultInterval),
      periodUnits: [
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
      ],
      localErrors: [],
    };
  },
  computed: {
    hiddenFields() {
      return this.modal.config.hiddenFields || [];
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
    resetValidation() {
      this.localErrors = [];
    },

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
