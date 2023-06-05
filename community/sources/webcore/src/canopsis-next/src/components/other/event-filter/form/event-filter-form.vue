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
    c-patterns-field(v-field="form.patterns", with-entity, with-event, some-required)

    template(v-if="hasAdditionalOptions")
      v-divider.my-3
      c-information-block(v-if="isEnrichmentType", :title="$t('eventFilter.enrichmentOptions')")
        event-filter-enrichment-form(v-field="form")
      c-information-block(v-else-if="isChangeEntityType", :title="$t('eventFilter.changeEntityOptions')")
        event-filter-change-entity-form(v-field="form.config")
</template>

<script>
import { EVENT_FILTER_TYPES } from '@/constants';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import PbehaviorRecurrenceRuleField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-recurrence-rule-field.vue';

import EventFilterEnrichmentForm from './partials/event-filter-enrichment-form.vue';
import EventFilterChangeEntityForm from './partials/event-filter-change-entity-form.vue';
import EventFilterDropIntervalsField from './partials/fields/event-filter-drop-intervals-field.vue';

export default {
  inject: ['$validator'],
  components: {
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
  },
};
</script>
