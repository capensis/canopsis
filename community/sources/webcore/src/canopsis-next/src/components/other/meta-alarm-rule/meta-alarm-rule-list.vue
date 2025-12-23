<template>
  <c-advanced-data-table
    :items="rules"
    :headers="headers"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-expandable-item="hasRulePatterns"
    :select-all="removable || updatable"
    expand
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        :enablable="updatable"
        :disablable="updatable"
        meta-alarm-rule
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #auto_resolve="{ item }">
      <c-enabled :value="item.auto_resolve" />
    </template>
    <template #config.threshold_rate="{ item }">
      {{ item | get('config.threshold_rate') | percentage }}
    </template>
    <template #config.threshold_count="{ item }">
      {{ item | get('config.threshold_count') }}
    </template>
    <template #config.time_interval="{ item }">
      {{ item | get('config.time_interval') | duration }}
    </template>
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #updated="{ item }">
      {{ item.updated | date }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="duplicable"
          type="duplicate"
          @click="$emit('duplicate', item)"
        />
        <c-action-btn
          v-if="removable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
        <c-db-export-btn :id="item._id" meta-alarm-rule />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <meta-alarm-rule-list-expand-panel :meta-alarm-rule="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { isMetaAlarmRuleTypeHasPatterns } from '@/helpers/entities/meta-alarm/rule/form';

import { useI18n } from '@/hooks/i18n';

import MetaAlarmRuleListExpandPanel from './partials/meta-alarm-rule-list-expand-panel.vue';

export default {
  components: {
    MetaAlarmRuleListExpandPanel,
  },
  props: {
    rules: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { text: t('common.id'), value: '_id' },
      { text: t('common.name'), value: 'name' },
      { text: t('common.type'), value: 'type' },
      { text: t('metaAlarmRule.autoResolve'), value: 'auto_resolve', sortable: false },
      { text: t('metaAlarmRule.thresholdRate'), value: 'config.threshold_rate', sortable: false },
      { text: t('metaAlarmRule.thresholdCount'), value: 'config.threshold_count', sortable: false },
      { text: t('metaAlarmRule.timeInterval'), value: 'config.time_interval', sortable: false },
      { text: t('common.author'), value: 'author.display_name' },
      { text: t('common.created'), value: 'created' },
      { text: t('common.updated'), value: 'updated' },
      { text: t('common.actionsLabel'), value: 'actions', sortable: false },
    ]);

    const hasRulePatterns = item => isMetaAlarmRuleTypeHasPatterns(item.type);

    return {
      headers,
      hasRulePatterns,
    };
  },
};
</script>
