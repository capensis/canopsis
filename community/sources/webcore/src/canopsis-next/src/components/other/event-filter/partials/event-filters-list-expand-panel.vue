<template>
  <v-tabs
    v-model="activeTab"
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab :href="`#${tabs.summary}`">
      {{ $t('common.summary') }}
    </v-tab>
    <v-tab :href="`#${tabs.pattern}`">
      {{ $tc('common.pattern', 2) }}
    </v-tab>
    <template v-if="isEnrichment">
      <v-tab :href="`#${tabs.action}`">
        {{ $tc('common.action', 2) }}
      </v-tab>
      <v-tab :href="`#${tabs.externalData}`" :disabled="!externalDataForm.length">
        {{ $t('externalData.title') }}
      </v-tab>
    </template>
    <v-tab v-if="eventFilter.failures_count" :href="`#${tabs.errors}`">
      {{ $tc('common.error', 2) }}
    </v-tab>

    <v-tabs-items
      :value="activeTab"
      mandatory
    >
      <v-tab-item :value="tabs.summary">
        <v-layout
          class="py-3 secondary lighten-2"
          justify-center
        >
          <v-flex xs11>
            <v-card>
              <v-card-text>
                <v-flex
                  xs12
                  md8
                  offset-md2
                  lg6
                  offset-lg3
                >
                  <event-filters-rule-summary :event-filter="eventFilter" />
                </v-flex>
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
      <v-tab-item :value="tabs.pattern">
        <v-layout
          class="pa-3 secondary lighten-2"
          justify-center
        >
          <v-flex xs10>
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
          </v-flex>
        </v-layout>
      </v-tab-item>
      <template v-if="isEnrichment">
        <v-tab-item :value="tabs.action">
          <v-layout
            class="py-3 secondary lighten-2"
            justify-center
          >
            <v-flex xs11>
              <v-data-table
                :items="eventFilter.config.actions"
                :headers="headers"
              />
            </v-flex>
          </v-layout>
        </v-tab-item>
        <v-tab-item :value="tabs.externalData">
          <v-layout
            class="py-3 secondary lighten-2"
            justify-center
          >
            <v-flex xs11>
              <external-data-form
                :form="externalDataForm"
                disabled
                optionally
              />
            </v-flex>
          </v-layout>
        </v-tab-item>
      </template>
      <v-tab-item v-if="eventFilter.failures_count" :value="tabs.errors">
        <v-layout
          class="py-3 secondary lighten-2"
          justify-center
        >
          <v-flex xs11>
            <v-card>
              <v-card-text>
                <event-filter-failures
                  :event-filter="eventFilter"
                  @refresh="$emit('refresh')"
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import { ref, computed } from 'vue';

import { EVENT_FILTER_EXPAND_PANEL_TABS } from '@/constants';

import { externalDataToForm } from '@/helpers/entities/shared/external-data/form';
import { eventFilterPatternToForm } from '@/helpers/entities/event-filter/rule/form';
import { isEnrichmentEventFilterRuleType } from '@/helpers/entities/event-filter/rule/entity';

import { useI18n } from '@/hooks/i18n';

import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

import EventFiltersRuleSummary from './event-filters-rule-summary.vue';
import EventFilterFailures from './event-filter-failures.vue';

export default {
  components: {
    EventFiltersRuleSummary,
    EventFilterFailures,
    ExternalDataForm,
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
