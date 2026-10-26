<template>
  <v-tabs
    v-model="activeTab"
    class="bg-secondary-lighten-1"
    slider-color="primary"
    color="white"
    align-tabs="center"
  >
    <v-tab :value="tabs.summary">
      {{ $t('common.summary') }}
    </v-tab>
    <v-tab :value="tabs.pattern">
      {{ $t('common.pattern', 2) }}
    </v-tab>
    <template v-if="isEnrichment">
      <v-tab :value="tabs.action">
        {{ $t('common.action', 2) }}
      </v-tab>
      <v-tab :value="tabs.externalData" :disabled="!externalDataForm.length">
        {{ $t('externalData.title') }}
      </v-tab>
    </template>
    <v-tab v-if="eventFilter.failures_count" :value="tabs.errors">
      {{ $t('common.error', 2) }}
    </v-tab>
  </v-tabs>

  <v-window v-model="activeTab" class="bg-secondary-lighten-2" mandatory>
    <c-window-expand-item :value="tabs.summary" cols="11">
      <v-row justify="center">
        <v-col cols="12" md="8" lg="6">
          <event-filters-rule-summary :event-filter="eventFilter" />
        </v-col>
      </v-row>
    </c-window-expand-item>
    <c-window-expand-item :value="tabs.pattern" cols="8">
      <c-patterns-field
        :value="patterns"
        readonly
        with-entity
        with-event
      />
    </c-window-expand-item>
    <template v-if="isEnrichment">
      <c-window-expand-item :value="tabs.action" cols="11">
        <v-data-table :items="eventFilter.config.actions" :headers="headers" />
      </c-window-expand-item>
      <c-window-expand-item :value="tabs.externalData" cols="11">
        <external-data-field
          :external-data="externalDataForm"
          disabled
          optionally
          with-on-mismatch
        />
      </c-window-expand-item>
    </template>
    <c-window-expand-item
      v-if="eventFilter.failures_count"
      :value="tabs.errors"
      cols="11"
    >
      <event-filter-failures :event-filter="eventFilter" @refresh="$emit('refresh')" />
    </c-window-expand-item>
  </v-window>
</template>

<script>
import { ref, computed } from 'vue';

import { EVENT_FILTER_EXPAND_PANEL_TABS } from '@/constants';

import { isEnrichmentEventFilterRuleType } from '@/helpers/entities/event-filter/rule/entity';
import { eventFilterPatternToForm } from '@/helpers/entities/event-filter/rule/form';
import { externalDataToForm } from '@/helpers/entities/shared/external-data/form';

import { useI18n } from '@/hooks/i18n';

import ExternalDataField from '@/components/forms/external-data/external-data-field.vue';

import EventFilterFailures from './event-filter-failures.vue';
import EventFiltersRuleSummary from './event-filters-rule-summary.vue';

export default {
  components: {
    EventFiltersRuleSummary,
    EventFilterFailures,
    ExternalDataField,
  },
  props: {
    eventFilter: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ['refresh'],
  setup(props) {
    const { t } = useI18n();

    const tabs = EVENT_FILTER_EXPAND_PANEL_TABS;
    const activeTab = ref(EVENT_FILTER_EXPAND_PANEL_TABS.summary);

    const headers = computed(() => [
      { value: 'type', text: t('common.type'), sortable: false },
      { value: 'name', text: t('common.name'), sortable: false },
      { value: 'value', text: t('common.value'), sortable: false },
      { value: 'description', text: t('common.description'), sortable: false },
    ]);

    const patterns = computed(() => eventFilterPatternToForm(props.eventFilter));
    const isEnrichment = computed(() => isEnrichmentEventFilterRuleType(props.eventFilter.type));
    const externalDataForm = computed(() => externalDataToForm(props.eventFilter.external_data, true));

    return {
      tabs,
      activeTab,

      headers,
      patterns,
      isEnrichment,
      externalDataForm,
    };
  },
};
</script>
