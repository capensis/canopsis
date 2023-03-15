<template lang="pug">
  c-advanced-data-table(
    :items="pbehaviors",
    :pagination="pagination",
    :loading="pending",
    :headers="headers",
    :total-items="totalItems",
    :search-tooltip="$t('pbehavior.searchHelp')",
    :select-all="removable || enablable || disablable",
    advanced-search,
    advanced-pagination,
    expand,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected, clearSelected }")
      pbehaviors-mass-actions-panel(
        :items="selected",
        :removable="removable",
        :enablable="enablable",
        :disablable="disablable",
        @clear:items="clearSelected"
      )
    template(#name="{ item }")
      c-ellipsis(:text="item.name")
    template(#enabled="{ item }")
      c-enabled(:value="item.enabled")
    template(#tstart="{ item }") {{ item.tstart | timezone($system.timezone) }}
    template(#tstop="{ item }") {{ item.tstop | timezone($system.timezone) }}
    template(#last_alarm_date="{ item }") {{ item.last_alarm_date | timezone($system.timezone) }}
    template(#created="{ item }") {{ item.created | date }}
    template(#updated="{ item }") {{ item.updated | date }}
    template(#rrule="{ item }")
      v-icon {{ item.rrule ? 'check' : 'clear' }}
    template(#type.icon_name="{ item }")
      v-icon(color="primary") {{ item.type.icon_name }}
    template(#is_active_status="{ item }")
      v-icon(:color="item.is_active_status ? 'primary' : 'error'") $vuetify.icons.settings_sync
    template(#actions="{ item }")
      v-layout(row)
        c-action-btn(
          v-if="updatable",
          :tooltip="item.editable ? $t('common.edit') : $t('pbehavior.notEditable')",
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
      pbehaviors-list-expand-item(:pbehavior="item")
</template>

<script>
import { isOldPattern } from '@/helpers/pattern';

import PbehaviorsMassActionsPanel from './actions/pbehaviors-mass-actions-panel.vue';
import PbehaviorsListExpandItem from './partials/pbehaviors-list-expand-item.vue';

export default {
  inject: ['$system'],
  components: {
    PbehaviorsListExpandItem,
    PbehaviorsMassActionsPanel,
  },
  props: {
    pbehaviors: {
      type: Array,
      required: true,
    },
    pagination: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
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
    enablable: {
      type: Boolean,
      default: false,
    },
    disablable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.author'), value: 'author.name' },
        { text: this.$t('pbehavior.isEnabled'), value: 'enabled' },
        { text: this.$t('pbehavior.begins'), value: 'tstart' },
        { text: this.$t('pbehavior.ends'), value: 'tstop' },
        { text: this.$t('common.type'), value: 'type.name' },
        { text: this.$t('common.reason'), value: 'reason.name' },
        { text: this.$t('common.created'), value: 'created' },
        { text: this.$t('common.updated'), value: 'updated' },
        { text: this.$t('pbehavior.lastAlarmDate'), value: 'last_alarm_date' },
        { text: this.$t('common.recurrence'), value: 'rrule' },
        { text: this.$t('common.icon'), value: 'type.icon_name' },
        { text: this.$t('common.status'), value: 'is_active_status', sortable: false },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
    },
  },
  methods: {
    isOldPattern(item) {
      return isOldPattern(item);
    },
  },
};
</script>
