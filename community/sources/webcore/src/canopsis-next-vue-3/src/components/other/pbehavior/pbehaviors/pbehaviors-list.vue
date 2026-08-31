<template>
  <c-advanced-data-table
    :items="pbehaviors"
    :options="options"
    :loading="pending"
    :headers="headers"
    :total-items="totalItems"
    :select-all="removable || updatable"
    :advanced-search-attributes="advancedSearchAttributes"
    expand
    advanced-search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #toolbar="">
      <c-layout gap="3" justify-end align-center>
        <pbehaviors-list-alarm-filtering-btn class="ma-0" @refresh="$emit('refresh')" />
      </c-layout>
    </template>

    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        :enablable="updatable"
        :disablable="updatable"
        pbehavior
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>

    <template #name="{ item }">
      <span class="text-truncate">{{ item.name }}</span>
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
      {{ timezoneDateToString(item.last_alarm_date) }}
    </template>

    <template #pattern_ms="{ item }">
      {{ durationToString(item.pattern_ms) }}
    </template>

    <template #pattern_exec_at="{ item }">
      {{ timezoneDateToString(item.pattern_exec_at) }}
    </template>

    <template #created="{ item }">
      {{ timezoneDateToString(item.created) }}
    </template>

    <template #updated="{ item }">
      {{ timezoneDateToString(item.updated) }}
    </template>

    <template #rrule="{ item }">
      <v-icon>{{ item.rrule ? 'check' : 'clear' }}</v-icon>
    </template>

    <template #type.icon_name="{ item }">
      <v-icon color="primary">
        {{ item.type?.icon_name }}
      </v-icon>
    </template>

    <template #is_active_status="{ item }">
      <v-icon
        :color="item.is_active_status ? 'success' : 'error'"
        icon="$vuetify.icons.settings_sync"
      />
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

import {
  DATETIME_FORMATS,
  PBEHAVIOR_LIST_FIELDS,
  PATTERN_DURATION_FORMAT,
  TIME_UNITS,
} from '@/constants';

import { convertDateToTimezoneDateString } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';

import { usePbehaviorAdvancedSearchAttributes } from '@/components/common/search/hooks/advanced-search';
import { usePbehaviorDateFormat } from '@/components/other/pbehavior/pbehaviors/hooks/pbehavior-date-format';

import PbehaviorActions from './partials/pbehavior-actions.vue';
import PbehaviorsListAlarmFilteringBtn from './partials/pbehaviors-list-alarm-filtering-btn.vue';
import PbehaviorsListExpandItem from './partials/pbehaviors-list-expand-item.vue';

export default {
  components: {
    PbehaviorActions,
    PbehaviorsListExpandItem,
    PbehaviorsListAlarmFilteringBtn,
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
  },

  emits: ['update:options', 'refresh'],

  setup() {
    const durationFormat = PATTERN_DURATION_FORMAT;
    const millisecondUnit = TIME_UNITS.millisecond;

    const { t } = useI18n();
    const { timezone, formatIntervalDate, formatRruleEndDate } = usePbehaviorDateFormat();
    const { attributes: advancedSearchAttributes } = usePbehaviorAdvancedSearchAttributes();

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
      { text: t('pbehavior.alarmFilteringTime'), value: PBEHAVIOR_LIST_FIELDS.patternMs },
      { text: t('pbehavior.lastFilteringDate'), value: PBEHAVIOR_LIST_FIELDS.patternExecAt },
      { text: t('common.icon', 1), value: PBEHAVIOR_LIST_FIELDS.typeIcon },
      { text: t('common.status'), value: PBEHAVIOR_LIST_FIELDS.status, sortable: false },
      { text: t('common.actionsLabel'), value: PBEHAVIOR_LIST_FIELDS.actions, sortable: false },
    ]);

    const timezoneDateToString = date => convertDateToTimezoneDateString(date, timezone.value, DATETIME_FORMATS.long);
    const durationToString = duration => convertDurationToString(duration, durationFormat, millisecondUnit);

    return {
      advancedSearchAttributes,
      headers,
      formatIntervalDate,
      formatRruleEndDate,
      timezoneDateToString,
      durationToString,
    };
  },
};
</script>
