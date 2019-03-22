<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.liveReporting.editLiveReporting') }}
    v-card-text
      h3 {{ $t('modals.liveReporting.dateInterval') }}
      v-layout(wrap)
        v-radio-group(v-model="selectedInterval")
          v-radio(
          v-for="interval in dateIntervals",
          :label="interval.text",
          :value="interval.value",
          :key="interval.value"
          )
      v-layout(wrap, v-if="isCustomRangeEnabled")
        v-flex(xs12)
          date-time-picker-field(
          v-model="tstart",
          v-validate="'required'"
          :label="$t('modals.liveReporting.tstart')",
          name="tstart",
          clearable
          )
        v-flex(xs12)
          date-time-picker-field(
          v-model="tstop",
          v-validate="tstopRules",
          :label="$t('modals.liveReporting.tstop')",
          name="tstop",
          clearable
          )
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit", :disabled="errors.any()") {{ $t('common.apply') }}
</template>

<script>
import moment from 'moment';

import { MODALS, LIVE_REPORTING_INTERVALS, DATETIME_FORMATS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

/**
   * Modal to add a time filter on alarm-list
   */
export default {
  name: MODALS.editLiveReporting,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    DateTimePickerField,
  },
  mixins: [modalInnerMixin],
  data() {
    const { config } = this.modal;

    return {
      selectedInterval: config.interval || '',
      dateIntervals: Object.values(LIVE_REPORTING_INTERVALS).map(value => ({
        value,
        text: this.$t(`modals.liveReporting.${value}`),
      })),
      tstart: config.tstart ? moment.unix(config.tstart).toDate() : new Date(),
      tstop: config.tstop ? moment.unix(config.tstop).toDate() : new Date(),
    };
  },
  computed: {
    isCustomRangeEnabled() {
      return this.selectedInterval === LIVE_REPORTING_INTERVALS.custom;
    },
    tstopRules() {
      const rules = { required: true };

      if (this.tstart) {
        rules.after = [moment(this.tstart).format(DATETIME_FORMATS.dateTimePicker)];
        rules.date_format = DATETIME_FORMATS.dateTimePicker;
      }

      return rules;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          const params = {
            interval: this.selectedInterval,
          };

          if (this.isCustomRangeEnabled) {
            params.tstart = this.tstart.getTime() / 1000;
            params.tstop = this.tstop.getTime() / 1000;
          }

          await this.config.action(params);
        }

        this.hideModal();
      }
    },
  },
};
</script>
