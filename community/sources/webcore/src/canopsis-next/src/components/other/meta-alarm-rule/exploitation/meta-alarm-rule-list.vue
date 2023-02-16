<template lang="pug">
  c-advanced-data-table(
    :items="rules",
    :headers="headers",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-expandable-item="hasRulePatterns",
    :select-all="removable",
    expand,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn(
        v-if="removable",
        type="delete",
        @click="$emit('remove-selected', selected)"
      )
    template(#auto_resolve="{ item }")
      c-enabled(:value="item.auto_resolve")
    template(#config.threshold_rate="{ item }") {{ item | get('config.threshold_rate') | percentage }}
    template(#config.threshold_count="{ item }") {{ item | get('config.threshold_count') }}
    template(#config.time_interval="{ item }") {{ item | get('config.time_interval') | duration }}
    template(#created="{ item }") {{ item.created | date }}
    template(#updated="{ item }") {{ item.updated | date }}
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          :badge-value="isOldPattern(item)",
          :badge-tooltip="$t('pattern.oldPatternTooltip')",
          type="edit",
          @click="$emit('edit', item)"
        )
        c-action-btn(
          v-if="duplicable",
          type="duplicate",
          @click="$emit('duplicate', item)"
        )
        c-action-btn(
          v-if="removable",
          type="delete",
          @click="$emit('remove', item._id)"
        )
    template(#expand="{ item }")
      meta-alarm-rule-list-expand-panel(:meta-alarm-rule="item")
</template>

<script>
import { OLD_PATTERNS_FIELDS } from '@/constants';

import { isOldPattern } from '@/helpers/pattern';
import { isMetaAlarmRuleTypeHasPatterns } from '@/helpers/forms/meta-alarm-rule';

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
    pagination: {
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
          value: 'author.name',
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
    isOldPattern(item) {
      return isOldPattern(item, [
        OLD_PATTERNS_FIELDS.entity,
        OLD_PATTERNS_FIELDS.alarm,
        OLD_PATTERNS_FIELDS.totalEntity,
      ]);
    },

    hasRulePatterns({ type }) {
      return isMetaAlarmRuleTypeHasPatterns(type);
    },
  },
};
</script>
