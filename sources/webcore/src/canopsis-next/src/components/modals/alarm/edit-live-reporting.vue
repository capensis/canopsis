<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.liveReporting.editLiveReporting') }}
    v-card-text
      h3 {{ $t('modals.liveReporting.dateInterval') }}
      stats-date-selector(
      v-model="form",
      :getDateObjectPreparer="getDateObjectPreparer"
      )
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit", :disabled="errors.any()") {{ $t('common.apply') }}
</template>

<script>
import moment from 'moment';

import { MODALS, DATETIME_FORMATS } from '@/constants';

import { dateParse } from '@/helpers/date-intervals';

import modalInnerMixin from '@/mixins/modal/inner';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';
import StatsDateSelector from '@/components/forms/stats-date-selector.vue';

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
    StatsDateSelector,
  },
  mixins: [modalInnerMixin],
  data() {
    const { config } = this.modal;

    return {
      form: {
        tstart: config.tstart || '',
        tstop: config.tstop || '',
      },
    };
  },
  computed: {
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
            return momentDate.toDate();
          }
        }

        return null;
      };
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.hideModal();
      }
    },
  },
};
</script>
