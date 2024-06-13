<template lang="pug">
  c-advanced-data-table(
    :items="eventFilters",
    :headers="headers",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :select-all="removable",
    advanced-pagination,
    search,
    expand,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn(
        v-if="removable",
        type="delete",
        @click="$emit('remove-selected', selected)"
      )
    template(#priority="{ item }") {{ item.priority || '-' }}
    template(#type="{ item }") {{ $t(`eventFilter.types.${item.type}`) }}
    template(#enabled="{ item }")
      c-enabled(:value="item | get('enabled', true)")
    template(#unread_failures_count="{ item }")
      c-circle-badge.error(v-if="item.unread_failures_count") {{ item.unread_failures_count }}
      template(v-else) -
    template(#created="{ item }") {{ item.created | date }}
    template(#updated="{ item }") {{ item.updated | date }}
    template(#calendar="{ item }")
      v-icon {{ isCalendarRule(item) ? 'check' : 'clear' }}
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
        pbehaviors-create-action-btn(:entity-id="item._id")
        pbehaviors-list-action-btn(:entity-id="item._id")
    template(#expand="{ item }")
      event-filters-list-expand-panel(:event-filter="item", @refresh="$emit('refresh')")
</template>

<script>
import { OLD_PATTERNS_FIELDS } from '@/constants';

import { isOldPattern } from '@/helpers/entities/pattern/form';

import PbehaviorsCreateActionBtn from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-create-action-btn.vue';
import PbehaviorsListActionBtn from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-list-action-btn.vue';

import EventFiltersListExpandPanel from './partials/event-filters-list-expand-panel.vue';

export default {
  components: {
    EventFiltersListExpandPanel,
    PbehaviorsCreateActionBtn,
    PbehaviorsListActionBtn,
  },
  props: {
    eventFilters: {
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
        { text: this.$t('common.id'), value: '_id' },
        { text: this.$t('common.type'), value: 'type', sortable: false },
        { text: this.$t('common.priority'), value: 'priority' },
        { text: this.$t('common.enabled'), value: 'enabled', sortable: false },
        { text: this.$tc('common.error', 2), value: 'unread_failures_count', sortable: false },
        { text: this.$t('common.author'), value: 'author.display_name' },
        { text: this.$t('common.created'), value: 'created' },
        { text: this.$t('common.updated'), value: 'updated' },
        { text: this.$t('common.calendar'), value: 'calendar', sortable: false },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
    },
  },

  methods: {
    isCalendarRule(item) {
      return item.start
        || item.rrule
        || item.exdates?.length
        || item.exceptions?.length;
    },

    isOldPattern(item) {
      return isOldPattern(item, [OLD_PATTERNS_FIELDS.patterns]);
    },
  },
};
</script>
