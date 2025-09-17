<template>
  <div>
    <v-layout>
      <v-flex xs8>
        <c-id-field
          v-field="form._id"
          :disabled="isDisabledIdField"
          :help-text="$t('eventFilter.idHelp')"
          class="mr-3"
          autofocus
        />
      </v-flex>
      <v-flex xs4>
        <c-event-filter-type-field
          v-field="form.type"
          :autofocus="isDisabledIdField"
          class="ml-3"
        />
      </v-flex>
    </v-layout>
    <c-description-field
      v-field="form.description"
      required
    />
    <v-layout justify-space-between>
      <c-enabled-field
        v-field="form.enabled"
        class="mr-3"
      />
      <c-priority-field v-field="form.priority" />
    </v-layout>
    <c-information-block :title="$t('eventFilter.duringPeriod')">
      <event-filter-drop-intervals-field v-field="form" />
    </c-information-block>
    <pbehavior-recurrence-rule-field
      v-field="form"
      class="mb-3"
    />
    <c-patterns-field
      v-field="form.patterns"
      :some-required="!isChangeEntityType"
      :required="isChangeEntityType"
      :with-entity="!isChangeEntityType"
      with-event
      entity-counters-type
    />
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
          :set-tags-items="setTagsItems"
        />
        <event-filter-change-entity-form
          v-else-if="isChangeEntityType"
          v-field="form.config"
          :variables="templateVars.config"
        />
      </c-information-block>
    </template>
  </div>
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
