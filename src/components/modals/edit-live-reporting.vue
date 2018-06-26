<template lang="pug">
  v-card
    v-card-title.blue.darken-4.white--text(primary-title)
      v-layout(justify-space-between, align-center)
        h2 {{ $t('modals.liveReporting.editLiveReporting') }}
        v-btn(@click="hideModal", icon, small)
          v-icon(color="white") close
    v-divider
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
          date-time-picker(v-model="tstart", clearable, label="tstart", name="tstart", v-validate="'required'")
        v-flex(xs12)
          date-time-picker(
            v-model="tstop",
            clearable, label="tstop",
            name="tstop",
            v-validate="'required|after:tstart'")
      v-btn(@click="submit", color="green darken-4 white--text", small) {{ $t('common.apply') }}
</template>

<script>
import DateTimePicker from '@/components/forms/date-time-picker.vue';
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  name: MODALS.editLiveReporting,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    DateTimePicker,
  },
  mixins: [modalMixin],
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
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.isCustomRangeEnabled) {
          this.$router.push({
            query: {
              ...this.$route.query,
              interval: this.selectedInterval.value,
              tstart: this.tstart.getTime() / 1000,
              tstop: this.tstop.getTime() / 1000,
            },
          });
        } else {
          this.$router.push({
            query: {
              ...this.$route.query,
              interval: this.selectedInterval.value,
            },
          });
        }
        this.hideModal();
      }
    },
  },
};
</script>
