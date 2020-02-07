<template lang="pug">
  tr(:data-test="`tableRow-${row.item._id}`")
    td.pr-0(data-test="rowCheckbox")
      v-layout(row, align-center)
        v-checkbox-functional(
          v-if="!isResolvedAlarm",
          v-field="selected",
          hide-details
        )
        v-checkbox-functional(
          v-else,
          :value="false",
          disabled,
          hide-details
        )
        v-layout.ml-2(align-center)
          v-btn.ma-0(icon, small, @click="row.expanded = !row.expanded")
            v-icon {{ row.expanded ? 'keyboard_arrow_up' : 'keyboard_arrow_down' }}
    td(v-for="column in columns")
      alarm-column-value(
        :alarm="row.item",
        :column="column",
        :columnFiltersMap="columnFiltersMap",
        :widget="widget"
      )
    td
      actions-panel(
        :item="row.item",
        :widget="widget",
        :isResolvedAlarm="isResolvedAlarm",
        :isEditingMode="isEditingMode"
      )
</template>

<script>
import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm/actions/mass-actions-panel.vue';
import TimeLine from '@/components/other/alarm/time-line/time-line.vue';
import AlarmListSearch from '@/components/other/alarm/search/alarm-list-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';
import NoColumnsTable from '@/components/tables/no-columns.vue';
import FilterSelector from '@/components/other/filter/selector/filter-selector.vue';

export default {
  components: {
    AlarmListSearch,
    RecordsPerPage,
    TimeLine,
    MassActionsPanel,
    ActionsPanel,
    AlarmColumnValue,
    NoColumnsTable,
    FilterSelector,
  },
  model: {
    prop: 'selected',
    event: 'input',
  },
  props: {
    selected: {
      type: Boolean,
      required: false,
      default: false,
    },
    row: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    columnFiltersMap: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    isResolvedAlarm: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
