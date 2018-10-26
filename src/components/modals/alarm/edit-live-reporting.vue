<template lang="pug">
  v-card
    v-card-title
      v-layout(justify-space-between, align-center)
        h2 {{ $t('modals.liveReporting.editLiveReporting') }}
        v-btn(@click="hideModal", icon, small)
          v-icon close
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
          date-time-picker(v-model="tstart",
          clearable,
          :label="$t('modals.liveReporting.tstart')",
          name="tstart",
          :rules="'required'")
        v-flex(xs12)
          date-time-picker(
          v-model="tstop",
          clearable,
          :label="$t('modals.liveReporting.tstop')",
          name="tstop",
          :rules="tstopRules")
      v-btn(@click="submit", color="green darken-4 white--text", small) {{ $t('common.apply') }}
</template>

<script>
import moment from 'moment';

import { MODALS } from '@/constants';

import DateTimePicker from '@/components/forms/date-time-picker.vue';
import modalInnerMixin from '@/mixins/modal/modal-inner';

/**
   * Modal to add a time filter on alarm-list
   */
export default {
  name: MODALS.editLiveReporting,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    DateTimePicker,
  },
  mixins: [modalInnerMixin],
  data() {
    const { config } = this.modal;

    return {
      selectedInterval: config.interval || '',
      dateIntervals: Object.values(this.$constants.LIVE_REPORTING_INTERVALS).map(value => ({
        value,
        text: this.$t(`modals.liveReporting.${value}`),
      })),
      tstart: config.tstart ? moment.unix(config.tstart).toDate() : new Date(),
      tstop: config.tstop ? moment.unix(config.tstop).toDate() : new Date(),
    };
  },
  computed: {
    isCustomRangeEnabled() {
      return this.selectedInterval === this.$constants.LIVE_REPORTING_INTERVALS.custom;
    },
    tstopRules() {
      return {
        required: true,
        after: [moment(this.tstart).format('DD/MM/YYYY HH:mm')],
        date_format: 'DD/MM/YYYY HH:mm',
      };
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
