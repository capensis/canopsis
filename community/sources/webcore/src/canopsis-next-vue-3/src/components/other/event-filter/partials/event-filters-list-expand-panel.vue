<template>
  <v-tabs
    v-model="activeTab"
    class="bg-secondary-lighten-1"
    slider-color="primary"
    color="white"
    centered
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

  <v-window v-model="activeTab" mandatory>
    <v-window-item :value="tabs.summary">
      <c-layout class="py-3 secondary lighten-2" justify-center>
        <div :style="gridFlexStyle(11)">
          <v-card>
            <v-card-text>
              <div>
                <event-filters-rule-summary :event-filter="eventFilter" />
              </div>
            </v-card-text>
          </v-card>
        </div>
      </c-layout>
    </v-window-item>
    <v-window-item :value="tabs.pattern">
      <c-layout class="pa-3 secondary lighten-2" justify-center>
        <div :style="gridFlexStyle(10)">
          <v-card>
            <v-card-text>
              <c-patterns-field
                :value="patterns"
                readonly
                with-entity
                with-event
              />
            </v-card-text>
          </v-card>
        </div>
      </c-layout>
    </v-window-item>
    <template v-if="isEnrichment">
      <v-window-item :value="tabs.action">
        <c-layout class="py-3 secondary lighten-2" justify-center>
          <div :style="gridFlexStyle(11)">
            <v-data-table :items="eventFilter.config.actions" :headers="headers" />
          </div>
        </c-layout>
      </v-window-item>
      <v-window-item :value="tabs.externalData">
        <c-layout class="py-3 secondary lighten-2" justify-center>
          <div :style="gridFlexStyle(11)">
            <external-data-field :form="externalDataForm" disabled optionally />
          </div>
        </c-layout>
      </v-window-item>
    </template>
    <v-window-item v-if="eventFilter.failures_count" :value="tabs.errors">
      <c-layout class="py-3 secondary lighten-2" justify-center>
        <div :style="gridFlexStyle(11)">
          <v-card>
            <v-card-text>
              <event-filter-failures :event-filter="eventFilter" @refresh="$emit('refresh')" />
            </v-card-text>
          </v-card>
        </div>
      </c-layout>
    </v-window-item>
  </v-window>
</template>

<script>
import { ref, computed } from 'vue';

import ExternalDataField from '@/components/forms/external-data/external-data-field.vue';
import { EVENT_FILTER_EXPAND_PANEL_TABS } from '@/constants';
import { isEnrichmentEventFilterRuleType } from '@/helpers/entities/event-filter/rule/entity';
import { eventFilterPatternToForm } from '@/helpers/entities/event-filter/rule/form';
import { externalDataToForm } from '@/helpers/entities/shared/external-data/form';
import { gridFlexStyle } from '@/helpers/layout';
import { useI18n } from '@/hooks/i18n';

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
    const externalDataForm = computed(() => externalDataToForm(props.eventFilter.external_data));

    return {
      gridFlexStyle,
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
