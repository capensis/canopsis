<template lang="pug" functional>
  tr(:data-test="`tableRow-${props.row.item._id}`")
    td(data-test="rowCheckbox")
      v-layout(row, align-center)
        v-checkbox-functional(
          v-if="!props.isResolvedAlarm",
          v-model="props.row.selected",
          hide-details
        )
        v-checkbox-functional(
          v-else,
          :value="false",
          disabled,
          hide-details
        )
        v-btn.ma-0(icon, small, @click="props.row.expanded = !props.row.expanded")
          v-icon {{ props.row.expanded ? 'keyboard_arrow_down' : 'keyboard_arrow_up' }}
    td(v-for="column in props.columns")
      alarm-column-value(
        :alarm="props.row.item",
        :column="column",
        :columnFiltersMap="props.columnFiltersMap",
        :widget="props.widget"
      )
    td
      actions-panel(
        :item="props.row.item",
        :widget="props.widget",
        :isResolvedAlarm="props.isResolvedAlarm",
        :isEditingMode="props.isEditingMode"
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
  props: {
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
