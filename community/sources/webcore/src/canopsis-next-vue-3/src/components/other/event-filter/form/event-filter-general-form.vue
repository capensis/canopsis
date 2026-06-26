<template>
  <c-layout gap="3" column>
    <c-form-block>
      <c-form-block-row :label="$t('common.type')">
        <c-event-filter-type-field v-model="typeField" />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.priority')">
        <c-priority-field v-model="priorityField" />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.description')">
        <c-description-field
          v-model="descriptionField"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('eventFilter.duringPeriod')">
        <event-filter-drop-intervals-field
          v-model:start="startField"
          v-model:stop="stopField"
          :required="hasRRule"
          :rrule="rruleField"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.recurrence')" indented>
        <pbehavior-recurrence-rule-field
          v-model:rrule="rruleField"
          v-model:exdates="exdatesField"
          v-model:exceptions="exceptionsField"
          :start="startField"
        />
      </c-form-block-row>
    </c-form-block>

    <event-filter-enrichment-form
      v-if="isEnrichmentType"
      :template-vars="templateVars"
      :copy-vars="copyVars"
    />

    <event-filter-change-entity-form
      v-else-if="isChangeEntityType"
      :template-vars="templateVars"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import {
  isEnrichmentEventFilterRuleType,
  isChangeEntityEventFilterRuleType,
} from '@/helpers/entities/event-filter/rule/entity';

import { useFormModels } from '@/hooks/form';

import PbehaviorRecurrenceRuleField from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-field.vue';

import EventFilterEnrichmentForm from './fields/event-filter-enrichment-form.vue';
import EventFilterChangeEntityForm from './fields/event-filter-change-entity-form.vue';
import EventFilterDropIntervalsField from './fields/event-filter-drop-intervals-field.vue';

export default {
  components: {
    EventFilterDropIntervalsField,
    PbehaviorRecurrenceRuleField,
    EventFilterEnrichmentForm,
    EventFilterChangeEntityForm,
  },
  props: {
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
  setup() {
    const {
      typeField,
      priorityField,
      descriptionField,
      rruleField,
      exdatesField,
      exceptionsField,
      startField,
      stopField,
    } = useFormModels({
      typeField: 'type',
      priorityField: 'priority',
      descriptionField: 'description',
      rruleField: 'rrule',
      exdatesField: 'exdates',
      exceptionsField: 'exceptions',
      startField: 'start',
      stopField: 'stop',
    });

    const hasRRule = computed(() => !!rruleField.value);
    const isEnrichmentType = computed(() => isEnrichmentEventFilterRuleType(typeField.value));
    const isChangeEntityType = computed(() => isChangeEntityEventFilterRuleType(typeField.value));

    return {
      typeField,
      priorityField,
      descriptionField,
      rruleField,
      exdatesField,
      exceptionsField,
      startField,
      stopField,
      hasRRule,
      isEnrichmentType,
      isChangeEntityType,
    };
  },
};
</script>
