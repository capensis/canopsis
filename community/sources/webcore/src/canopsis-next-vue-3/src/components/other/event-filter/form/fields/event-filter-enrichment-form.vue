<template>
  <c-layout gap="3" column>
    <span class="text-subtitle-1">
      {{ $t('eventFilter.enrichmentOptions') }}
    </span>

    <c-form-block>
      <c-form-block-row :label="$t('externalData.title')" indented>
        <external-data-field
          v-model:external-data="externalDataField"
          :variables="templateVars.external_data"
          optionally
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.actionsLabel')" indented>
        <c-label :label="$t('common.actionsLabel')" class="mb-3" />
        <event-filter-enrichment-actions-field
          v-model:actions="actionsField"
          :variables="templateVars.config"
          :copy-variables="copyVars.config"
          :set-tags-items="setTagsItems"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('eventFilter.onSuccessAndFailure')">
        <event-filter-enrichment-after-outcome-fields
          v-model:on-success="configField.on_success"
          v-model:on-failure="configField.on_failure"
        />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { getSetTagsItemsFromPattern } from '@/helpers/entities/event-filter/rule/entity';

import { useFormModels } from '@/hooks/form';

import ExternalDataField from '@/components/forms/external-data/external-data-field.vue';

import EventFilterEnrichmentActionsField from './event-filter-enrichment-actions-field.vue';
import EventFilterEnrichmentAfterOutcomeFields from './event-filter-enrichment-after-outcome-fields.vue';

export default {
  components: {
    ExternalDataField,
    EventFilterEnrichmentActionsField,
    EventFilterEnrichmentAfterOutcomeFields,
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
  },
  setup() {
    const { externalDataField, configField, actionsField, patternsField } = useFormModels({
      externalDataField: 'external_data',
      configField: 'config',
      actionsField: 'actions',
      patternsField: 'patterns',
    });

    const setTagsItems = computed(() => getSetTagsItemsFromPattern(patternsField.value?.event_pattern));

    return {
      setTagsItems,

      externalDataField,
      configField,
      actionsField,
    };
  },
};
</script>
