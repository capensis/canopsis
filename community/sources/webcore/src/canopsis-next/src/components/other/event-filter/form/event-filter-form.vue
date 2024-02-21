<template>
  <div>
    <v-layout>
      <v-flex xs8>
        <c-id-field
          class="mr-3"
          v-field="form._id"
          :disabled="isDisabledIdField"
          :help-text="$t('eventFilter.idHelp')"
        />
      </v-flex>
      <v-flex xs4>
        <c-event-filter-type-field
          class="ml-3"
          v-field="form.type"
        />
      </v-flex>
    </v-layout>
    <c-description-field
      v-field="form.description"
      required
    />
    <v-layout justify-space-between>
      <c-enabled-field
        class="mr-3"
        v-field="form.enabled"
      />
      <c-priority-field v-field="form.priority" />
    </v-layout>
    <c-information-block :title="$t('eventFilter.duringPeriod')">
      <event-filter-drop-intervals-field v-field="form" />
    </c-information-block>
    <pbehavior-recurrence-rule-field
      class="mb-3"
      v-field="form"
    />
    <c-patterns-field
      v-field="form.patterns"
      with-entity
      with-event
      some-required
      entity-counters-type
    />
    <template v-if="hasAdditionalOptions">
      <v-divider class="my-3" />
      <c-information-block
        :title="isEnrichmentType ? $t('eventFilter.enrichmentOptions') : $t('eventFilter.changeEntityOptions')"
      >
        <c-collapse-panel
          class="mb-2"
          :title="$t('externalData.title')"
        >
          <external-data-form
            v-field="form.external_data"
            :variables="externalDataVariables"
          />
        </c-collapse-panel>
        <event-filter-enrichment-form
          v-if="isEnrichmentType"
          v-field="form"
          :template-variables="actionsDataVariables"
          :set-tags-items="setTagsItems"
        />
        <event-filter-change-entity-form
          v-else-if="isChangeEntityType"
          v-field="form.config"
          :variables="actionsDataVariables"
        />
      </c-information-block>
    </template>
  </div>
</template>

<script>
import { EXTERNAL_DATA_DEFAULT_CONDITION_VALUES, EXTERNAL_DATA_PAYLOADS_VARIABLES } from '@/constants';

import {
  isEnrichmentEventFilterRuleType,
  isChangeEntityEventFilterRuleType,
  getSetTagsItemsFromPattern,
} from '@/helpers/entities/event-filter/rule/entity';

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
      return isEnrichmentEventFilterRuleType(this.form.type);
    },

    isChangeEntityType() {
      return isChangeEntityEventFilterRuleType(this.form.type);
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

    setTagsItems() {
      return getSetTagsItemsFromPattern(this.form.patterns.event_pattern);
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
