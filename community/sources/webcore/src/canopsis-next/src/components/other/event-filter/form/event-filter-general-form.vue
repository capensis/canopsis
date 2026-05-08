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
        <event-filter-drop-intervals-field v-field="form" />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.recurrence')">
        <div class="py-3">
          <pbehavior-recurrence-rule-field v-field="form" />
        </div>
      </c-form-block-row>
    </c-form-block>

    <span class="text-subtitle-1">
      {{ $t('eventFilter.enrichmentOptions') }}
    </span>

    <template v-if="hasAdditionalOptions">
      <v-divider class="my-3" />
      <c-information-block
        :title="
          isEnrichmentType ? $t('eventFilter.enrichmentOptions') : $t('eventFilter.changeEntityOptions')
        "
      >
        <c-collapse-panel
          :title="$t('externalData.title')"
          class="mb-2"
        >
          <external-data-form
            v-field="form.external_data"
            :variables="templateVars.external_data"
            optionally
          />
        </c-collapse-panel>
        <event-filter-enrichment-form
          v-if="isEnrichmentType"
          v-field="form"
          :template-variables="templateVars.config"
          :copy-variables="copyVars.config"
          :set-tags-items="setTagsItems"
        />
        <event-filter-change-entity-form
          v-else-if="isChangeEntityType"
          v-field="form.config"
          :variables="templateVars.config"
        />
      </c-information-block>
    </template>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

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
    const isEnrichmentType = computed(() => isEnrichmentEventFilterRuleType(props.form.type));

    const isChangeEntityType = computed(() => isChangeEntityEventFilterRuleType(props.form.type));

    const hasAdditionalOptions = computed(() => isEnrichmentType.value || isChangeEntityType.value);

    const setTagsItems = computed(() => getSetTagsItemsFromPattern(props.form.patterns?.event_pattern));

    return {
      isEnrichmentType,
      isChangeEntityType,
      hasAdditionalOptions,
      setTagsItems,
    };
  },
};
</script>
