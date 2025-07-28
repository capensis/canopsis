<template>
  <c-advanced-data-table
    :items="pbehaviors"
    :options="options"
    :loading="pending"
    :headers="headers"
    :total-items="totalItems"
    :select-all="removable || enablable || disablable"
    :advanced-search-fields="advancedSearchFields"
    advanced-search
    advanced-pagination
    expand
    @update:options="$emit('update:options', $event)"
  >
    <template v-if="shownUserTimezone" #toolbar="">
      <v-layout justify-end>
        <v-flex xs3>
          <c-timezone-field
            v-model="timezone"
            server
          />
        </v-flex>
      </v-layout>
    </template>
    <template #mass-actions="{ selected, clearSelected }">
      <pbehaviors-mass-actions-panel
        :items="selected"
        :removable="removable"
        :enablable="enablable"
        :disablable="disablable"
        @clear:items="clearSelected"
      />
    </template>
    <template #name="{ item }">
      <c-ellipsis :text="item.name" />
    </template>
    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>
    <template #tstart="{ item }">
      {{ formatIntervalDate(item, 'tstart') }}
    </template>
    <template #tstop="{ item }">
      {{ formatIntervalDate(item, 'tstop') }}
    </template>
    <template #rrule_end="{ item }">
      {{ formatRruleEndDate(item) }}
    </template>
    <template #last_alarm_date="{ item }">
      {{ item.last_alarm_date | timezone(timezone) }}
    </template>
    <template #created="{ item }">
      {{ item.created | timezone(timezone) }}
    </template>
    <template #updated="{ item }">
      {{ item.updated | timezone(timezone) }}
    </template>
    <template #rrule="{ item }">
      <v-icon>{{ item.rrule ? 'check' : 'clear' }}</v-icon>
    </template>
    <template #type.icon_name="{ item }">
      <v-icon color="primary">
        {{ item.type.icon_name }}
      </v-icon>
    </template>
    <template #is_active_status="{ item }">
      <v-icon :color="item.is_active_status ? 'success' : 'error'">
        $vuetify.icons.settings_sync
      </v-icon>
    </template>
    <template #actions="{ item }">
      <pbehavior-actions
        :pbehavior="item"
        :removable="removable"
        :updatable="updatable"
        :duplicable="duplicable"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #expand="{ item }">
      <pbehaviors-list-expand-item :pbehavior="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import { usePbehaviorDateFormat } from '@/components/other/pbehavior/pbehaviors/hooks/pbehavior-date-format';

import { ADVANCED_SEARCH_DATE_CONDITIONS } from '@/constants/advanced-search';
import { PBEHAVIOR_LIST_FIELDS } from '@/constants/pbehavior';

import PbehaviorsMassActionsPanel from './actions/pbehaviors-mass-actions-panel.vue';
import PbehaviorActions from './partials/pbehavior-actions.vue';
import PbehaviorsListExpandItem from './partials/pbehaviors-list-expand-item.vue';

export default {
  inject: ['$system'],
  components: {
    PbehaviorActions,
    PbehaviorsListExpandItem,
    PbehaviorsMassActionsPanel,
  },
  props: {
    pbehaviors: {
      type: Array,
      required: true,
    },
    options: {
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
  setup() {
    const { t, tc } = useI18n();
    const { timezone, shownUserTimezone, formatIntervalDate, formatRruleEndDate } = usePbehaviorDateFormat();

    const headers = computed(() => [
      { text: t('common.name'), value: PBEHAVIOR_LIST_FIELDS.name },
      { text: t('common.author'), value: PBEHAVIOR_LIST_FIELDS.author },
      { text: t('pbehavior.isEnabled'), value: PBEHAVIOR_LIST_FIELDS.enabled },
      { text: t('pbehavior.begins'), value: PBEHAVIOR_LIST_FIELDS.begins },
      { text: t('pbehavior.ends'), value: PBEHAVIOR_LIST_FIELDS.ends },
      { text: t('pbehavior.rruleEnd'), value: PBEHAVIOR_LIST_FIELDS.rruleEnd, sortable: false },
      { text: t('common.recurrence'), value: PBEHAVIOR_LIST_FIELDS.rrule },
      { text: t('common.type'), value: PBEHAVIOR_LIST_FIELDS.type },
      { text: t('common.reason'), value: PBEHAVIOR_LIST_FIELDS.reason },
      { text: t('common.created'), value: PBEHAVIOR_LIST_FIELDS.created },
      { text: t('common.updated'), value: PBEHAVIOR_LIST_FIELDS.updated },
      { text: t('pbehavior.lastAlarmDate'), value: PBEHAVIOR_LIST_FIELDS.lastAlarmDate },
      { text: t('pbehavior.alarmCount'), value: PBEHAVIOR_LIST_FIELDS.alarmCount, sortable: false },
      { text: tc('common.icon', 1), value: PBEHAVIOR_LIST_FIELDS.typeIcon },
      { text: t('common.status'), value: PBEHAVIOR_LIST_FIELDS.status, sortable: false },
      { text: t('common.actionsLabel'), value: PBEHAVIOR_LIST_FIELDS.actions, sortable: false },
    ]);

    const notSearchableFields = [
      PBEHAVIOR_LIST_FIELDS.rruleEnd,
      PBEHAVIOR_LIST_FIELDS.lastAlarmDate,
      PBEHAVIOR_LIST_FIELDS.alarmCount,
      PBEHAVIOR_LIST_FIELDS.typeIcon,
      PBEHAVIOR_LIST_FIELDS.status,
      PBEHAVIOR_LIST_FIELDS.actions,
    ];

    const dateSearchableFields = [
      PBEHAVIOR_LIST_FIELDS.begins,
      PBEHAVIOR_LIST_FIELDS.ends,
      PBEHAVIOR_LIST_FIELDS.created,
      PBEHAVIOR_LIST_FIELDS.updated,
    ];

    const advancedSearchFields = computed(() => (
      headers.value.filter(header => !notSearchableFields.includes(header.value))
        .map(header => (
          dateSearchableFields.includes(header.value)
            ? { ...header, conditions: ADVANCED_SEARCH_DATE_CONDITIONS }
            : header
        ))
    ));

    return {
      timezone,
      shownUserTimezone,
      headers,
      advancedSearchFields,

      formatIntervalDate,
      formatRruleEndDate,
    };
  },
};
</script>
