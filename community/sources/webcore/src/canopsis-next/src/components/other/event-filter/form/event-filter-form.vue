<template lang="pug">
  div
    c-id-field(
      v-field="form._id",
      :disabled="isDisabledIdField",
      :help-text="$t('eventFilter.idHelp')"
    )
    c-event-filter-type-field(v-field="form.type")
    c-description-field(v-field="form.description", required)
    c-priority-field(v-field="form.priority")
    c-enabled-field(v-field="form.enabled")
    template(v-if="isDropType")
      event-filter-drop-intervals-field(v-field="form")
      pbehavior-recurrence-rule-field.mb-1(v-field="form")
    c-patterns-field(v-field="form.patterns", with-entity, with-event, some-required)

    template(v-if="isChangeEntityType || isEnrichmentType")
      v-divider.my-3
      c-information-block(:title="$t('eventFilter.configuration')")
        template(v-if="isChangeEntityType")
          event-filter-change-entity-form(v-field="form.config")
        template(v-if="isEnrichmentType")
          event-filter-enrichment-form(v-field="form.config")
</template>

<script>
import { DATETIME_FORMATS, EVENT_FILTER_TYPES } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import EventFilterEnrichmentForm from '@/components/other/event-filter/form/event-filter-enrichment-form.vue';
import EventFilterChangeEntityForm from '@/components/other/event-filter/form/event-filter-change-entity-form.vue';
import PbehaviorRecurrenceRuleField from '@/components/other/pbehavior/calendar/partials/pbehavior-recurrence-rule-field.vue';

import EventFilterDropIntervalsField from './fields/event-filter-drop-intervals-field.vue';

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
    startRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    stopRules() {
      return {
        required: true,
        after: [convertDateToString(this.form.start, DATETIME_FORMATS.dateTimePicker)],
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    isEnrichmentType() {
      return this.form.type === EVENT_FILTER_TYPES.enrichment;
    },

    isDropType() {
      return this.form.type === EVENT_FILTER_TYPES.drop;
    },

    isChangeEntityType() {
      return this.form.type === EVENT_FILTER_TYPES.changeEntity;
    },
  },
};
</script>
