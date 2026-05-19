<template>
  <v-layout class="gap-3" column>
    <span class="text-subtitle-1">
      {{ $t('eventFilter.enrichmentOptions') }}
    </span>

    <c-form-block>
      <c-form-block-row :label="$t('externalData.title')" indented>
        <external-data-form
          v-field="form.external_data"
          :variables="templateVars.external_data"
          optionally
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.actionsLabel')" indented>
        <c-label :label="$t('common.actionsLabel')" class="mb-3" />
        <event-filter-enrichment-actions-form
          v-field="form.config.actions"
          :variables="templateVars.config"
          :copy-variables="copyVars.config"
          :name="name"
          :set-tags-items="setTagsItems"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('eventFilter.onSuccessAndFailure')">
        <event-filter-enrichment-after-outcome-fields v-field="form" />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { getSetTagsItemsFromPattern } from '@/helpers/entities/event-filter/rule/entity';

import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

import EventFilterEnrichmentActionsForm from './event-filter-enrichment-actions-form.vue';
import EventFilterEnrichmentAfterOutcomeFields from './event-filter-enrichment-after-outcome-fields.vue';

export default {
  inject: ['$validator'],
  components: {
    ExternalDataForm,
    EventFilterEnrichmentActionsForm,
    EventFilterEnrichmentAfterOutcomeFields,
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
    name: {
      type: String,
      default: 'config.actions',
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    copyVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const setTagsItems = computed(() => getSetTagsItemsFromPattern(props.form.patterns?.event_pattern));

    return {
      setTagsItems,
    };
  },
};
</script>
