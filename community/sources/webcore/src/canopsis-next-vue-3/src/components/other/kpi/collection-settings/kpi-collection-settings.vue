<template>
  <v-form @submit.prevent="submit">
    <v-row>
      <v-col cols="6">
        <kpi-collection-settings-basic-metrics-list />
      </v-col>
      <v-col cols="6">
        <kpi-collection-settings-optional-metrics-form
          v-if="form"
          :value="form"
          @update:value="form = $event"
        />
        <c-layout v-else class="my-2" justify-center>
          <v-progress-circular color="primary" indeterminate />
        </c-layout>
      </v-col>
    </v-row>
    <v-divider class="mt-3" />
    <c-layout class="mt-3" justify-end>
      <v-btn
        :disabled="!form || isFormNotChanged"
        :loading="submitting"
        class="mr-0"
        color="primary"
        type="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </c-layout>
  </v-form>
</template>

<script>
import { cloneDeep, isEqual } from 'lodash';

import { metricsSettingsToForm } from '@/helpers/entities/metrics-settings/form';

import { usePendingHandler } from '@/hooks/query/pending';
import { useMetricsSettings } from '@/hooks/store/modules/metrics-settings';

import KpiCollectionSettingsOptionalMetricsForm from './form/kpi-collection-settings-optional-metrics-form.vue';
import KpiCollectionSettingsBasicMetricsList from './partials/kpi-collection-settings-basic-metrics-list.vue';

export default {
  components: {
    KpiCollectionSettingsOptionalMetricsForm,
    KpiCollectionSettingsBasicMetricsList,
  },

  setup() {
    const { fetchMetricsSettingsWithoutStore, updateMetricsSettings } = useMetricsSettings();

    const { pending: submitting, handler: submitWithPending } = usePendingHandler(handler =>
      handler(),
    );

    return {
      fetchMetricsSettingsWithoutStore,
      updateMetricsSettings,
      submitting,
      submitWithPending,
    };
  },

  data() {
    return {
      initialForm: null,
      form: null,
    };
  },

  computed: {
    isFormNotChanged() {
      return isEqual(this.form, this.initialForm);
    },
  },

  mounted() {
    this.fetchMetricsSettings();
  },

  methods: {
    setInitialForm() {
      this.initialForm = cloneDeep(this.form);
    },

    async fetchMetricsSettings() {
      const metricsSettings = await this.fetchMetricsSettingsWithoutStore();

      this.form = metricsSettingsToForm(metricsSettings);
      this.setInitialForm();
    },

    async submit() {
      await this.submitWithPending(async () => {
        await this.updateMetricsSettings({ data: this.form });

        this.setInitialForm();

        this.$popups.success({ text: this.$t('success.default') });
      });
    },
  },
};
</script>
