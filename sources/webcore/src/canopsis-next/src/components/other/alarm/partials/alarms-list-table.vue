<template lang="pug">
  v-flex.white
    v-flex.px-3(v-show="selectedIds.length", xs12)
      mass-actions-panel(:itemsIds="selectedIds", :widget="widget")
    no-columns-table(v-if="!hasColumns")
    div(v-else)
      v-data-table.alarms-list-table(
        v-model="selected",
        :class="vDataTableClass",
        :items="alarms",
        :headers="headers",
        :total-items="totalItems",
        :pagination="pagination",
        :select-all="selectable",
        :loading="loading || alarmColumnFiltersPending",
        data-test="tableWidget",
        ref="dataTable",
        item-key="_id",
        expand,
        hide-actions,
        @update:pagination="updatePaginationHandler"
      )
        template(slot="progress")
          v-fade-transition
            v-progress-linear(height="2", indeterminate, color="primary")
        template(slot="headerCell", slot-scope="props")
          span {{ props.header.text }}
        template(slot="items", slot-scope="props")
          alarms-list-row(
            v-model="props.selected",
            :selectable="selectable",
            :row="props",
            :isEditingMode="isEditingMode",
            :widget="widget",
            :columns="columns",
            :columnFiltersMap="columnFiltersMap",
            :isTourEnabled="checkIsTourEnabledForAlarmByIndex(props.index)"
          )
        template(slot="expand", slot-scope="props")
          alarms-expand-panel(
            :alarm="props.item",
            :widget="widget",
            :isTourEnabled="checkIsTourEnabledForAlarmByIndex(props.index)"
          )
    slot
</template>

<script>
import { isResolvedAlarm } from '@/helpers/entities';
import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm/actions/mass-actions-panel.vue';
import AlarmsExpandPanel from '@/components/other/alarm/partials/alarms-expand-panel.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';
import NoColumnsTable from '@/components/tables/no-columns.vue';
import AlarmsListRow from '@/components/other/alarm/partials/alarms-list-row.vue';

import alarmColumnFilters from '@/mixins/entities/alarm-column-filters';

/**
   * Alarm-list-table component
   *
   * @module alarm
   */
export default {
  components: {
    AlarmsListRow,
    RecordsPerPage,
    AlarmsExpandPanel,
    MassActionsPanel,
    ActionsPanel,
    AlarmColumnValue,
    NoColumnsTable,
  },
  mixins: [alarmColumnFilters],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    alarms: {
      type: Array,
      required: true,
    },
    columns: {
      type: Array,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pagination: {
      type: Object,
      required: false,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    isTourEnabled: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    hasColumns: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      selected: [],
    };
  },

  computed: {
    expanded() {
      return this.$refs.dataTable.expanded;
    },

    selectedIds() {
      return this.selected
        .filter(item => !isResolvedAlarm(item))
        .map(item => item._id);
    },

    headers() {
      if (!this.hasColumns) {
        return [];
      }

      const filters = [...this.columns, { text: this.$t('common.actionsLabel'), sortable: false }];

      if (!this.selectable) {
        filters.unshift({ sortable: false });
      }

      return filters;
    },

    vDataTableClass() {
      const columnsLength = this.headers.length;
      const COLUMNS_SIZES_VALUES = {
        sm: { min: 0, max: 10, label: 'sm' },
        md: { min: 11, max: 12, label: 'md' },
        lg: { min: 13, max: Number.MAX_VALUE, label: 'lg' },
      };

      const { label = COLUMNS_SIZES_VALUES.sm.label } = Object.values(COLUMNS_SIZES_VALUES)
        .find(({ min, max }) => columnsLength >= min && columnsLength <= max);

      return {
        [`columns-${label}`]: true,
      };
    },
  },

  mounted() {
    this.fetchAlarmColumnFilters();
  },

  methods: {
    checkIsTourEnabledForAlarmByIndex(index) {
      return this.isTourEnabled && index === 0;
    },

    updatePaginationHandler(data) {
      this.$emit('update:pagination', data);
    },
  },
};
</script>

<style lang="scss">
  .alarms-list-table {
    &.columns-lg {
      table.v-table {
        tbody, thead {
          td, th {
            padding: 0 8px;
          }

          @media screen and (max-width: 1600px) {
            td, th {
              padding: 0 4px;
            }
          }

          @media screen and (max-width: 1450px) {
            td, th {
              font-size: 0.85em;
            }

            .badge {
              font-size: inherit;
            }
          }
        }
      }
    }

    &.columns-md {
      table.v-table {
        tbody, thead {
          @media screen and (max-width: 1700px) {
            td, th {
              padding: 0 12px;
            }
          }

          @media screen and (max-width: 1250px) {
            td, th {
              padding: 0 8px;
            }
          }

          @media screen and (max-width: 1150px) {
            td, th {
              font-size: 0.85em;
              padding: 0 4px;
            }

            .badge {
              font-size: inherit;
            }
          }
        }
      }
    }
  }
</style>
