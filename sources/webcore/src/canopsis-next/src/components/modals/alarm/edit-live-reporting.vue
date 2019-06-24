<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.liveReporting.editLiveReporting') }}
    v-card-text
      h3 {{ $t('modals.liveReporting.dateInterval') }}
      v-layout(wrap)
        v-radio-group(v-model="form.interval")
          v-radio(
          v-for="interval in dateIntervals",
          :label="interval.text",
          :value="interval.value",
          :key="interval.value"
          )
      v-layout(wrap, v-if="isCustomRangeEnabled")
        v-flex(xs6)
          v-layout(align-center)
            date-time-picker-text-field(
            v-model="form.tstart",
            :label="$t('common.startDate')",
            :dateObjectPreparer="getDateObjectPreparer('start')",
            name="tstart"
            )
          v-layout(align-center)
            date-time-picker-text-field(
            v-model="form.tstop",
            :label="$t('common.endDate')",
            :dateObjectPreparer="getDateObjectPreparer('stop')",
            name="tstop"
            )
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit", :disabled="errors.any()") {{ $t('common.apply') }}
</template>

<script>
import moment from 'moment';

import { MODALS, LIVE_REPORTING_INTERVALS, DATETIME_FORMATS } from '@/constants';

import { dateParse } from '@/helpers/date-intervals';

import modalInnerMixin from '@/mixins/modal/inner';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

/**
   * Modal to add a time filter on alarm-list
   */
export default {
  name: MODALS.editLiveReporting,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    DateTimePickerTextField,
  },
  mixins: [modalInnerMixin],
  data() {
    const { config } = this.modal;

    return {
      form: {
        interval: config.interval || '',
        tstart: config.tstart || '',
        tstop: config.tstop || '',
      },
    };
  },
  computed: {
    dateIntervals() {
      return Object.values(LIVE_REPORTING_INTERVALS).map(value => ({
        value,
        text: this.$t(`modals.liveReporting.${value}`),
      }));
    },

    isCustomRangeEnabled() {
      return this.form.interval === LIVE_REPORTING_INTERVALS.custom;
    },

    tstartRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    tstopRules() {
      const rules = { required: true };

      if (this.tstart) {
        rules.after = [moment(this.tstart).format(DATETIME_FORMATS.dateTimePicker)];
        rules.date_format = DATETIME_FORMATS.veeValidateDateTimeFormat;
      }

      return rules;
    },
  },
  methods: {
    getDateObjectPreparer(type) {
      return (date) => {
        if (date) {
          const momentDate = dateParse(date, type, DATETIME_FORMATS.dateTimePicker);

          if (momentDate.isValid()) {
            return momentDate.startOf('hour').toDate();
          }
        }

        return null;
      };
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          const data = this.isCustomRangeEnabled ? this.form : {
            interval: this.form.interval,
          };

          await this.config.action(data);
        }

        this.hideModal();
      }
    },
  },
};
</script>
