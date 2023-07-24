<template lang="pug">
  div
    c-collapse-panel.mb-2(:title="$t('externalData.title')")
      external-data-form(v-field="form.external_data", :variables="externalDataVariables")
    c-collapse-panel.mb-2(:title="$t('eventFilter.editActions')")
      event-filter-enrichment-actions-form(v-field="form.config.actions", :variables="actionsDataVariables")
    v-layout(row)
      v-select.mr-3(
        v-field="form.config.on_success",
        :label="$t('eventFilter.onSuccess')",
        :items="eventFilterAfterTypes"
      )
      v-select.ml-3(
        v-field="form.config.on_failure",
        :label="$t('eventFilter.onFailure')",
        :items="eventFilterAfterTypes"
      )
    v-alert(:value="errors.has(name)", type="error") {{ $t('eventFilter.actionsRequired') }}
</template>

<script>
import {
  EVENT_FILTER_ENRICHMENT_AFTER_TYPES,
  EXTERNAL_DATA_DEFAULT_CONDITION_VALUES,
  EXTERNAL_DATA_PAYLOADS_VARIABLES,
} from '@/constants';

import { formMixin } from '@/mixins/form';

import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

import EventFilterEnrichmentActionsForm from './event-filter-enrichment-actions-form.vue';

export default {
  inject: ['$validator'],
  components: {
    ExternalDataForm,
    EventFilterEnrichmentActionsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'actions',
    },
  },
  computed: {
    eventFilterAfterTypes() {
      return Object.values(EVENT_FILTER_ENRICHMENT_AFTER_TYPES);
    },

    externalDataVariables() {
      return EXTERNAL_DATA_DEFAULT_CONDITION_VALUES.map(({ value, text }) => ({
        value,
        text: this.$t(`externalData.conditionValues.${text}`),
      }));
    },

    actionsDataVariables() {
      const referencesVariables = this.form.external_data.length
        ? this.form.external_data.map(({ reference }) => ({
          value: EXTERNAL_DATA_PAYLOADS_VARIABLES.externalData.replace('%reference%', reference),
          text: `${this.$t('externalData.title')}: ${reference}`,
        }))
        : [{
          value: EXTERNAL_DATA_PAYLOADS_VARIABLES.externalData,
          text: this.$t('externalData.title'),
        }];

      return [
        ...this.externalDataVariables,
        ...referencesVariables,
      ];
    },
  },
  watch: {
    'form.config.actions': function validateActions() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.attachActionsRequiredRule();
  },
  beforeDestroy() {
    this.detachActionsRequiredRule();
  },
  methods: {
    attachActionsRequiredRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'required:true',
        getter: () => this.form.config.actions,
      });
    },

    detachActionsRequiredRule() {
      this.$validator.detach(this.name);
    },
  },
};
</script>
