<template lang="pug">
  div
    v-layout(row)
      v-flex(xs8)
        c-id-field.mr-3(
          v-field="form._id",
          :disabled="isDisabledIdField",
          :help-text="$t('eventFilter.idHelp')"
        )
      v-flex(xs4)
        c-event-filter-type-field.ml-3(v-field="form.type")
    c-description-field(v-field="form.description", required)
    v-layout(row, justify-space-between)
      c-enabled-field(v-field="form.enabled")
      c-priority-field(v-field="form.priority")
    c-information-block(:title="$t('eventFilter.duringPeriod')")
      event-filter-drop-intervals-field(v-field="form")
    pbehavior-recurrence-rule-field.mb-3(v-field="form")
    c-patterns-field(v-field="form.patterns", with-entity, with-event, some-required, entity-counters-type)

    template(v-if="hasAdditionalOptions")
      v-divider.my-3
      c-information-block(
        :title="isEnrichmentType ? $t('eventFilter.enrichmentOptions') : $t('eventFilter.changeEntityOptions')"
      )
        c-collapse-panel.mb-2(:title="$t('externalData.title')")
          external-data-form(v-field="form.external_data", :variables="externalDataVariables")

        event-filter-enrichment-form(
          v-if="isEnrichmentType",
          v-field="form",
          :template-variables="actionsDataVariables"
        )
        event-filter-change-entity-form(
          v-else-if="isChangeEntityType",
          v-field="form.config",
          :variables="actionsDataVariables"
        )
</template>

<script>
import {
  EVENT_FILTER_TYPES,
  EXTERNAL_DATA_DEFAULT_CONDITION_VALUES,
  EXTERNAL_DATA_PAYLOADS_VARIABLES,
} from '@/constants';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import PbehaviorRecurrenceRuleField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-recurrence-rule-field.vue';
import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

import EventFilterEnrichmentForm from './fields/event-filter-enrichment-form.vue';
import EventFilterChangeEntityForm from './fields/event-filter-change-entity-form.vue';
import EventFilterDropIntervalsField from './fields/event-filter-drop-intervals-field.vue';

export default {
  inject: ['$validator'],
  components: {
    ExternalDataForm,
    EventFilterDropIntervalsField,
    DateTimePickerField,
    PbehaviorRecurrenceRuleField,
    EventFilterEnrichmentForm,
    EventFilterChangeEntityForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isEnrichmentType() {
      return this.form.type === EVENT_FILTER_TYPES.enrichment;
    },

    isChangeEntityType() {
      return this.form.type === EVENT_FILTER_TYPES.changeEntity;
    },

    hasAdditionalOptions() {
      return this.isEnrichmentType || this.isChangeEntityType;
    },

    regexpVariables() {
      return [{
        value: EXTERNAL_DATA_PAYLOADS_VARIABLES.regexp,
        text: this.$t('common.regexp'),
      }];
    },

    externalDataVariables() {
      return [
        ...EXTERNAL_DATA_DEFAULT_CONDITION_VALUES.map(({ value, text }) => ({
          value,
          text: this.$t(`externalData.conditionValues.${text}`),
        })),
        ...this.regexpVariables,
      ];
    },

    referencesVariables() {
      return this.form.external_data.length
        ? this.form.external_data.map(({ reference }) => ({
          value: EXTERNAL_DATA_PAYLOADS_VARIABLES.externalData.replace('%reference%', reference),
          text: `${this.$t('externalData.title')}: ${reference}`,
        }))
        : [{
          value: EXTERNAL_DATA_PAYLOADS_VARIABLES.externalData,
          text: this.$t('externalData.title'),
        }];
    },

    actionsDataVariables() {
      return [
        ...this.externalDataVariables,
        ...this.referencesVariables,
      ];
    },
  },
};
</script>
