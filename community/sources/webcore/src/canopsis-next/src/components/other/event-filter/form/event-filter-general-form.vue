<template>
  <v-layout class="gap-3" column>
    <c-form-block>
      <c-form-block-row :label="$t('common.type')">
        <c-event-filter-type-field v-field="form.type" />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.priority')">
        <c-priority-field v-field="form.priority" />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.description')">
        <c-description-field v-field="form.description" required />
      </c-form-block-row>

      <c-form-block-row :label="$t('eventFilter.duringPeriod')">
        <event-filter-drop-intervals-field v-field="form" :required="hasRRule" />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.recurrence')" indented>
        <c-label class="mb-3">{{ $t('common.recurrence') }}</c-label>
        <pbehavior-recurrence-rule-field v-field="form" />
      </c-form-block-row>
    </c-form-block>

    <event-filter-enrichment-form
      v-if="isEnrichmentType"
      v-field="form"
      :template-vars="templateVars"
      :copy-vars="copyVars"
    />

    <event-filter-change-entity-form
      v-else-if="isChangeEntityType"
      v-field="form"
      :template-vars="templateVars"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import {
  isEnrichmentEventFilterRuleType,
  isChangeEntityEventFilterRuleType,
} from '@/helpers/entities/event-filter/rule/entity';

import PbehaviorRecurrenceRuleField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-recurrence-rule-field.vue';

import EventFilterEnrichmentForm from './fields/event-filter-enrichment-form.vue';
import EventFilterChangeEntityForm from './fields/event-filter-change-entity-form.vue';
import EventFilterDropIntervalsField from './fields/event-filter-drop-intervals-field.vue';

export default {
  inject: ['$validator'],
  components: {
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
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    copyVars: {
      type: Object,
      default: () => ({}),
    },
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const hasRRule = computed(() => !!props.form.rrule);
    const isEnrichmentType = computed(() => isEnrichmentEventFilterRuleType(props.form.type));
    const isChangeEntityType = computed(() => isChangeEntityEventFilterRuleType(props.form.type));

    return {
      hasRRule,
      isEnrichmentType,
      isChangeEntityType,
    };
  },
};
</script>
