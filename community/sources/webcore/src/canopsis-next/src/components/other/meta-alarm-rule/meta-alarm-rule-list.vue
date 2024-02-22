<template>
  <c-advanced-data-table
    :items="rules"
    :headers="headers"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-expandable-item="hasRulePatterns"
    :select-all="removable"
    expand
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
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
      </v-layout>
    </template>
    <template #expand="{ item }">
      <meta-alarm-rule-list-expand-panel :meta-alarm-rule="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { isMetaAlarmRuleTypeHasPatterns } from '@/helpers/entities/meta-alarm/rule/form';

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
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.id'),
          value: '_id',
        },
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.type'),
          value: 'type',
        },
        {
          text: this.$t('metaAlarmRule.autoResolve'),
          value: 'auto_resolve',
          sortable: false,
        },
        {
          text: this.$t('metaAlarmRule.thresholdRate'),
          value: 'config.threshold_rate',
          sortable: false,
        },
        {
          text: this.$t('metaAlarmRule.thresholdCount'),
          value: 'config.threshold_count',
          sortable: false,
        },
        {
          text: this.$t('metaAlarmRule.timeInterval'),
          value: 'config.time_interval',
          sortable: false,
        },
        {
          text: this.$t('common.author'),
          value: 'author.display_name',
        },
        {
          text: this.$t('common.created'),
          value: 'created',
        },
        {
          text: this.$t('common.updated'),
          value: 'updated',
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
  methods: {
    hasRulePatterns({ type }) {
      return isMetaAlarmRuleTypeHasPatterns(type);
    },
  },
};
</script>
