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
          :value="interval",
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
import DateTimePicker from '@/components/forms/date-time-picker.vue';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

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
    return {
      selectedInterval: '',
      dateIntervals: [
        {
          text: this.$t('modals.liveReporting.today'),
          value: 'today',
        },
        {
          text: this.$t('modals.liveReporting.yesterday'),
          value: 'yesterday',
        },
        {
          text: this.$t('modals.liveReporting.last7Days'),
          value: 'last7Days',
        },
        {
          text: this.$t('modals.liveReporting.last30Days'),
          value: 'last30Days',
        },
        {
          text: this.$t('modals.liveReporting.thisMonth'),
          value: 'thisMonth',
        },
        {
          text: this.$t('modals.liveReporting.lastMonth'),
          value: 'lastMonth',
        },
        {
          text: this.$t('modals.liveReporting.custom'),
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

      if (isFormValid && this.config.updateQuery) {
        const params = {
          interval: this.selectedInterval.value,
        };

        if (this.isCustomRangeEnabled) {
          params.tstart = this.tstart.getTime() / 1000;
          params.tstop = this.tstop.getTime() / 1000;
        }

        this.config.updateQuery(params);
        this.hideModal();
      }
    },
  },
};
</script>
