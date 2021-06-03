<template lang="pug">
  v-flex.white
    v-flex.px-3(v-show="selectedIds.length", xs12)
      mass-actions-panel(:itemsIds="selectedIds", :widget="widget")
    c-empty-data-table-columns(v-if="!hasColumns")
    div(v-else)
      v-data-table.alarms-list-table(
        v-model="selected",
        :class="vDataTableClass",
        :items="alarms",
        :headers="headers",
        :total-items="totalItems",
        :pagination="pagination",
        :select-all="selectable",
        :loading="loading || columnsFiltersPending",
        :expand="expandable",
        data-test="tableWidget",
        ref="dataTable",
        item-key="_id",
        hide-actions,
        @update:pagination="updatePaginationHandler"
      )
        template(slot="progress")
          v-fade-transition
            v-progress-linear(height="2", indeterminate, color="primary")
        template(slot="headerCell", slot-scope="props")
          alarm-header-cell(:header="props.header")
        template(slot="items", slot-scope="props")
          alarms-list-row(
            v-model="props.selected",
            v-on="rowListeners",
            :selectable="selectable",
            :expandable="expandable",
            :row="props",
            :widget="widget",
            :columns="columns",
            :columns-filters="columnsFilters",
            :hide-groups="hideGroups",
            :parent-alarm="parentAlarm",
            :is-tour-enabled="checkIsTourEnabledForAlarmByIndex(props.index)"
          )
        template(slot="expand", slot-scope="props")
          alarms-expand-panel(
            :alarm="props.item",
            :widget="widget",
            :hide-groups="hideGroups",
            :is-tour-enabled="checkIsTourEnabledForAlarmByIndex(props.index)"
          )
    slot
    component(v-bind="additionalComponent.props", v-on="additionalComponent.on", :is="additionalComponent.is")
</template>

<script>
import { isResolvedAlarm } from '@/helpers/entities';

import Observer from '@/services/observer';
import featuresService from '@/services/features';

import { entitiesAlarmColumnsFiltersMixin } from '@/mixins/entities/associative-table/alarm-columns-filters';

import AlarmHeaderCell from '../headers-formatting/alarm-header-cell.vue';
import MassActionsPanel from '../actions/mass-actions-panel.vue';
import AlarmsExpandPanel from '../partials/alarms-expand-panel.vue';
import AlarmsListRow from '../partials/alarms-list-row.vue';

/**
   * Alarm-list-table component
   *
   * @module alarm
   */
export default {
  inject: {
    $periodicRefresh: {
      default: new Observer(),
    },
  },
  components: {
    AlarmsListRow,
    AlarmsExpandPanel,
    MassActionsPanel,
    AlarmHeaderCell,

    ...featuresService.get('components.alarmListTable.components', {}),
  },
  mixins: [
    entitiesAlarmColumnsFiltersMixin,

    ...featuresService.get('components.alarmListTable.mixins', []),
  ],
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
    hideGroups: {
      type: Boolean,
      default: false,
    },
    expandable: {
      type: Boolean,
      default: false,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
  },
  data() {
    const data = featuresService.has('components.alarmListTable.data')
      ? featuresService.call('components.alarmListTable.data', this, {})
      : {};

    return {
      selected: [],
      columnsFilters: [],
      columnsFiltersPending: false,

      ...data,
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

    hasInstructionsAlarms() {
      return this.alarms.some(alarm => alarm.assigned_instructions.length);
    },

    headers() {
      if (!this.hasColumns) {
        return [];
      }

      const headers = [...this.columns, { text: this.$t('common.actionsLabel'), sortable: false }];

      if ((this.expandable || this.hasInstructionsAlarms) && !this.selectable) {
        // We need it for the expand panel open button
        headers.unshift({ sortable: false });
      }

      return headers;
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

    rowListeners() {
      if (featuresService.has('components.alarmListTable.computed.rowListeners')) {
        return featuresService.call('components.alarmListTable.computed.rowListeners', this, {});
      }

      return {};
    },

    additionalComponent() {
      if (featuresService.has('components.alarmListTable.computed.additionalComponent')) {
        return featuresService.call('components.alarmListTable.computed.additionalComponent', this, {});
      }

      return {};
    },
  },
  watch: {
    ...featuresService.get('components.alarmListTable.watch', {}),
  },
  async mounted() {
    if (featuresService.has('components.alarmListTable.mounted')) {
      featuresService.call('components.alarmListTable.mounted', this, {});
    }

    this.columnsFiltersPending = true;
    this.columnsFilters = await this.fetchAlarmColumnsFiltersList();
    this.columnsFiltersPending = false;
  },
  beforeDestroy() {
    if (featuresService.has('components.alarmListTable.beforeDestroy')) {
      featuresService.call('components.alarmListTable.beforeDestroy', this, {});
    }
  },

  methods: {
    ...featuresService.get('components.alarmListTable.methods', {}),

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