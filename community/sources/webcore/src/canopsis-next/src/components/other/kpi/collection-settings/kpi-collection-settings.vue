<template>
  <v-form @submit.prevent="submit">
    <v-layout>
      <v-flex xs6>
        <kpi-collection-settings-basic-metrics-list />
      </v-flex>
      <v-flex xs6>
        <kpi-collection-settings-optional-metrics-form
          v-if="form"
          v-model="form"
        />
        <v-layout
          class="my-2"
          v-else
          justify-center
        >
          <v-progress-circular
            indeterminate
            color="primary"
          />
        </v-layout>
      </v-flex>
    </v-layout>
    <v-divider class="mt-3" />
    <v-layout
      class="mt-3"
      justify-end
    >
      <v-btn
        class="primary mr-0"
        :disabled="isDisabled || isFormNotChanged"
        :loading="submitting"
        type="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </v-layout>
  </v-form>
</template>

<script>
import { cloneDeep, isEqual } from 'lodash';

import { metricsSettingsToForm } from '@/helpers/entities/metrics-settings/form';

import { entitiesMetricsSettingsMixin } from '@/mixins/entities/metrics-settings';
import { submittableMixinCreator } from '@/mixins/submittable';

import KpiCollectionSettingsBasicMetricsList from './partials/kpi-collection-settings-basic-metrics-list.vue';
import KpiCollectionSettingsOptionalMetricsForm from './form/kpi-collection-settings-optional-metrics-form.vue';

export default {
  inject: ['$validator'],
  components: {
    KpiCollectionSettingsOptionalMetricsForm,
    KpiCollectionSettingsBasicMetricsList,
  },
  mixins: [
    entitiesMetricsSettingsMixin,
    submittableMixinCreator(),
  ],
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
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.updateMetricsSettings({ data: this.form });

        this.setInitialForm();

        this.$popups.success({ text: this.$t('success.default') });
      }
    },
  },
};
</script>
