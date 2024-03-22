<template>
  <c-advanced-data-table
    :headers="headers"
    :items="availabilities"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    expand
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #value="{ item }">
      <availability-list-column-value
        :availability="item"
        :display-parameter="displayParameter"
        :show-trend="showTrend"
        :show-type="showType"
      />
    </template>
    <template v-for="column of columns" #[column.value]="{ item }">
      <entity-column-cell
        :key="column.value"
        :entity="item.entity"
        :column="column"
      />
    </template>
    <template #expand="{ item }">
      <availability-list-expand-panel
        :availability="item"
        :active-alarms-columns="activeAlarmsColumns"
        :resolved-alarms-columns="resolvedAlarmsColumns"
        :interval="interval"
        :default-show-type="showType"
        :display-parameter="displayParameter"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import AvailabilityListColumnValue from '@/components/other/availability/partials/availability-list-column-value.vue';
import EntityColumnCell from '@/components/widgets/context/columns-formatting/entity-column-cell.vue';

import AvailabilityListExpandPanel from './partials/availability-list-expand-panel.vue';

export default {
  components: { EntityColumnCell, AvailabilityListColumnValue, AvailabilityListExpandPanel },
  props: {
    availabilities: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
    activeAlarmsColumns: {
      type: Array,
      default: () => [],
    },
    resolvedAlarmsColumns: {
      type: Array,
      default: () => [],
    },
    displayParameter: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
    showType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
    showTrend: {
      type: Boolean,
      default: false,
    },
    interval: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const isUptimeParameter = computed(() => props.displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime);
    const headers = computed(() => [
      {
        text: t(`common.${isUptimeParameter.value ? 'uptime' : 'downtime'}`),
        value: 'value',
        sortable: false,
      },
      ...props.columns,
    ]);

    return {
      headers,
    };
  },
};
</script>
