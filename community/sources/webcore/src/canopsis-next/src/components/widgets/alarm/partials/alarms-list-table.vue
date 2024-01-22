<template>
  <v-flex
    v-resize="changeHeaderPositionOnResize"
    v-on="wrapperListeners"
  >
    <c-empty-data-table-columns v-if="!columns.length" />
    <div v-else>
      <v-layout
        v-if="shownTopPagination"
        ref="actions"
        class="alarms-list-table__top-pagination px-4 position-relative"
        align-center
      >
        <v-flex
          v-if="densable || !hideActions"
          class="alarms-list-table__top-pagination--left"
          xs6
        >
          <v-layout
            align-center
            justify-start
          >
            <c-density-btn-toggle
              v-if="densable"
              :value="dense"
              @change="$emit('update:dense', $event)"
            />
            <v-fade-transition v-if="!hideActions">
              <v-flex
                v-show="unresolvedSelected.length"
                class="px-1"
              >
                <mass-actions-panel
                  :items="unresolvedSelected"
                  :widget="widget"
                  :refresh-alarms-list="refreshAlarmsList"
                  @clear:items="clearSelected"
                />
              </v-flex>
            </v-fade-transition>
          </v-layout>
        </v-flex>
        <v-flex
          v-if="!hidePagination"
          class="alarms-list-table__top-pagination--center-absolute"
          xs4
        >
          <c-pagination
            :page="options.page"
            :limit="options.itemsPerPage"
            :total="totalItems"
            type="top"
            @input="updatePage"
          />
        </v-flex>
        <v-flex
          v-if="resizableColumn || draggableColumn"
          class="alarms-list-table__top-pagination--right-absolute"
        >
          <c-action-btn
            v-if="isColumnsChanging"
            :tooltip="$t('alarm.tooltips.resetChangeColumns')"
            icon="$vuetify.icons.restart_alt"
            @click="resetColumnsSettings"
          />
          <c-action-btn
            :icon="isColumnsChanging ? 'lock_open' : 'lock_outline'"
            :tooltip="$t(`alarm.tooltips.${isColumnsChanging ? 'finishChangeColumns' : 'startChangeColumns'}`)"
            @click="toggleColumnEditingMode"
          />
        </v-flex>
      </v-layout>
      <v-data-table
        v-model="selected"
        ref="dataTable"
        :class="vDataTableClass"
        :style="vDataTableStyle"
        :items="alarms"
        :headers="headersWithWidth"
        :server-items-length="totalItems"
        :options="options"
        :show-select="selectable"
        :loading="loading || columnsFiltersPending"
        :dense="isMediumDense"
        :ultra-dense="isSmallDense"
        class="alarms-list-table"
        item-key="_id"
        loader-height="2"
        hide-default-footer
        multi-sort
        @update:options="$emit('update:options', $event)"
      >
        <template #progress="">
          <v-fade-transition>
            <v-progress-linear
              color="primary"
              height="2"
              indeterminate
            />
          </v-fade-transition>
        </template>
        <template
          v-for="item in headers"
          #[`header.${item.value}`]="{ header }"
        >
          <alarm-header-cell
            :key="`header.${item.value}`"
            :header="header"
            :selected-tag="selectedTag"
            :resizing="resizingMode"
            @clear:tag="$emit('clear:tag')"
          />
          <template>
            <span
              v-if="draggingMode"
              :key="`header.${item.value}.drag`"
              class="alarms-list-table__dragging-handler"
              @click.stop=""
            />
            <span
              v-if="resizingMode"
              :key="`header.cell.${item.value}.resize`"
              class="alarms-list-table__resize-handler"
              @mousedown.stop.prevent="startColumnResize(header.value)"
              @click.stop=""
            />
          </template>
        </template>
        <template #item="{ isSelected, isExpanded, item, select, expand }">
          <alarms-list-row
            :ref="`row${item._id}`"
            :key="item._id"
            :selected="isSelected"
            :selectable="selectable"
            :expandable="expandable"
            :expanded="isExpanded"
            :alarm="item"
            :widget="widget"
            :headers="headers"
            :parent-alarm="parentAlarm"
            :refresh-alarms-list="refreshAlarmsList"
            :selecting="selecting"
            :selected-tag="selectedTag"
            :medium="isMediumDense"
            :small="isSmallDense"
            :resizing="resizingMode"
            :search="search"
            :wrap-actions="resizableColumn"
            :show-instruction-icon="hasInstructionsAlarms"
            v-on="rowListeners"
            @start:resize="startColumnResize"
            @select:tag="$emit('select:tag', $event)"
            @click:state="openRootCauseDiagram"
            @expand="expand"
            @input="select"
          />
        </template>
        <template #expanded-item="{ item }">
          <alarms-expand-panel
            :alarm="item"
            :parent-alarm-id="parentAlarmId"
            :widget="widget"
            :search="search"
            :hide-children="hideChildren"
            @select:tag="$emit('select:tag', $event)"
          />
        </template>
      </v-data-table>
    </div>
    <c-table-pagination
      v-if="!hidePagination"
      :total-items="totalItems"
      :items-per-page="options.itemsPerPage"
      :page="options.page"
      @update:page="updatePage"
      @update:items-per-page="updateItemsPerPage"
    />
    <component
      v-bind="additionalComponent.props"
      :is="additionalComponent.is"
      v-on="additionalComponent.on"
    />
  </v-flex>
</template>

<script>
import { get, intersectionBy } from 'lodash';

import { ALARM_DENSE_TYPES, ALARMS_RESIZING_CELLS_CONTENTS_BEHAVIORS, MODALS } from '@/constants';

import featuresService from '@/services/features';

import { isActionAvailableForAlarm } from '@/helpers/entities/alarm/form';
import { calculateAlarmLinksColumnWidth } from '@/helpers/entities/alarm/list';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { widgetColumnsAlarmMixin } from '@/mixins/widget/columns/alarm';
import { widgetRowsSelectingAlarmMixin } from '@/mixins/widget/rows/alarm-selecting';
import { widgetColumnResizingAlarmMixin } from '@/mixins/widget/columns/alarm-resizing';
import { widgetColumnDraggingAlarmMixin } from '@/mixins/widget/columns/alarm-dragging';
import { widgetHeaderStickyAlarmMixin } from '@/mixins/widget/rows/alarm-sticky-header';

import AlarmHeaderCell from '../headers-formatting/alarm-header-cell.vue';
import AlarmsExpandPanel from '../expand-panel/alarms-expand-panel.vue';
import MassActionsPanel from '../actions/mass-actions-panel.vue';

import AlarmsListRow from './alarms-list-row.vue';

/**
 * Alarm-list-table component
 *
 * @module alarm
 */
export default {
  components: {
    MassActionsPanel,
    AlarmHeaderCell,
    AlarmsExpandPanel,
    AlarmsListRow,
  },
  mixins: [
    entitiesInfoMixin,
    widgetColumnsAlarmMixin,
    widgetHeaderStickyAlarmMixin,
    widgetRowsSelectingAlarmMixin,
    widgetColumnResizingAlarmMixin,
    widgetColumnDraggingAlarmMixin,

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
    totalItems: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    columns: {
      type: Array,
      default: () => [],
    },
    loading: {
      type: Boolean,
      default: false,
    },
    expandable: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Number,
      default: ALARM_DENSE_TYPES.large,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
    selectedTag: {
      type: String,
      default: '',
    },
    hideChildren: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    hidePagination: {
      type: Boolean,
      default: false,
    },
    densable: {
      type: Boolean,
      default: false,
    },
    resizableColumn: {
      type: Boolean,
      default: false,
    },
    draggableColumn: {
      type: Boolean,
      default: false,
    },
    columnsSettings: {
      type: Object,
      default: () => ({}),
    },
    cellsContentBehavior: {
      type: String,
      required: false,
    },
    search: {
      type: String,
      default: '',
    },
  },

  computed: {
    shownTopPagination() {
      return this.totalItems && (this.densable || !this.hideActions || !this.hidePagination);
    },

    wrapperListeners() {
      return this.selectable
        ? { mousemove: this.throttledMousemoveHandler }
        : {};
    },

    unresolvedSelected() {
      return this.selected.filter(item => isActionAvailableForAlarm(item));
    },

    expanded() {
      return this.$refs.dataTable.expansion;
    },

    isColumnsChanging() {
      return this.resizingMode || this.draggingMode;
    },

    hasInstructionsAlarms() {
      return this.alarms.some(alarm => alarm.assigned_instructions?.length);
    },

    isCellContentWrapped() {
      return this.cellsContentBehavior === ALARMS_RESIZING_CELLS_CONTENTS_BEHAVIORS.wrap;
    },

    isCellContentTruncated() {
      return this.cellsContentBehavior === ALARMS_RESIZING_CELLS_CONTENTS_BEHAVIORS.truncate;
    },

    needToAddLeftActionsCell() {
      return (this.expandable || this.hasInstructionsAlarms) && !this.selectable;
    },

    hasLeftActions() {
      return this.selectable || this.needToAddLeftActionsCell;
    },

    headers() {
      const headers = this.preparedColumns.map((column) => {
        const header = {
          ...column,
          class: this.draggableClass,
        };

        if (column.linksInRowCount) {
          const linksCounts = this.alarms.map(alarm => Object.values(get(alarm, column.value, {})).flat().length);
          const maxLinksCount = Math.max(...linksCounts);
          const actualInlineLinksCount = maxLinksCount > column.inlineLinksCount
            ? column.inlineLinksCount + 1
            : maxLinksCount;

          header.width = calculateAlarmLinksColumnWidth(
            this.dense,
            Math.max(Math.min(actualInlineLinksCount, column.linksInRowCount), 1),
          );
        }

        return header;
      });

      if (!this.hideActions) {
        headers.push({
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
          class: this.draggableClass,
        });
      }

      if (this.needToAddLeftActionsCell) {
        /**
         * We need it for the expand panel open button
         */
        headers.unshift({ sortable: false, width: 100 });
      }

      return this.draggableColumn
        ? headers.sort((a, b) => this.getColumnPositionByField(a.value) - this.getColumnPositionByField(b.value))
        : headers;
    },

    headersWithWidth() {
      if (this.resizableColumn) {
        return this.headers.map((header) => {
          const width = this.getColumnWidthByField(header.value);

          return {
            ...header,
            width: header.width
              ? header.width
              : width && `${width}%`,
          };
        });
      }

      return this.headers;
    },

    vDataTableClass() {
      const columnsLength = this.headers.length;
      const COLUMNS_SIZES_VALUES = {
        sm: { min: 0, max: 10, label: 'sm' },
        md: { min: 11, max: 12, label: 'md' },
        lg: { min: 13, max: Number.MAX_VALUE, label: 'lg' },
      };

      const { label } = Object.values(COLUMNS_SIZES_VALUES)
        .find(({ min, max }) => columnsLength >= min && columnsLength <= max);

      return {
        [`columns-${label}`]: true,
        'alarms-list-table__selecting': this.selecting,
        'alarms-list-table__selecting--text-unselectable': this.selectingMousePressed,
        'alarms-list-table__grid': this.isColumnsChanging,
        'alarms-list-table__dragging': this.draggingMode,
        'alarms-list-table--wrapped': this.isCellContentWrapped,
        'alarms-list-table--truncated': this.isCellContentTruncated,
        'alarms-list-table--fixed': this.resizableColumn || this.draggableColumn,
      };
    },

    leftActionsWidth() {
      /**
       * left expand/instruction icon/select actions width
       */
      return this.isMediumDense || this.isSmallDense ? 100 : 120;
    },

    vDataTableStyle() {
      if (this.resizableColumn) {
        const actionsWidth = this.hasLeftActions ? this.leftActionsWidth : 0;

        return {
          '--alarms-list-table-width': `calc(${actionsWidth}px + ${this.sumOfColumnsWidth}%)`,
        };
      }

      return {};
    },

    rowListeners() {
      if (featuresService.has('components.alarmListTable.computed.rowListeners')) {
        return featuresService.call('components.alarmListTable.computed.rowListeners', this, {});
      }

      return {};
    },

    additionalComponent() {
      if (featuresService.has('components.alarmListTable.computed.additionalComponent')) {
        return featuresService.call('components.alarmListTable.computed.additionalComponent', this);
      }

      return {};
    },

    isMediumDense() {
      return this.dense === ALARM_DENSE_TYPES.medium;
    },

    isSmallDense() {
      return this.dense === ALARM_DENSE_TYPES.small;
    },

    parentAlarmId() {
      return this.parentAlarm?._id;
    },
  },

  watch: {
    alarms(alarms) {
      this.selected = intersectionBy(alarms, this.selected, '_id');
    },

    columns() {
      if (this.isColumnsChanging) {
        this.updateColumnsSettings();

        this.disableDraggingMode();
        this.disableResizingMode();
      }
    },

    columnsSettings: {
      immediate: true,
      deep: true,
      handler() {
        if (!this.draggingMode && this.columnsSettings?.columns_position) {
          this.setColumnsPosition(this.columnsSettings?.columns_position);
        }

        if (!this.resizingMode && this.columnsSettings?.columns_width) {
          this.setColumnsWidth(this.columnsSettings?.columns_width);
        }
      },
    },
  },

  methods: {
    updateColumnsSettings() {
      const settings = {};

      if (this.resizingMode) {
        settings.columns_width = this.columnsWidthByField;
      }

      if (this.draggingMode) {
        settings.columns_position = this.columnsPositionByField;
      }

      this.$emit('update:columns-settings', settings);
    },

    toggleColumnEditingMode() {
      if (this.isColumnsChanging) {
        this.updateColumnsSettings();
      }

      if (this.resizableColumn) {
        this.toggleResizingMode();
      }

      if (this.draggableColumn) {
        this.toggleDraggingMode();
      }
    },

    resetColumnsSettings() {
      if (this.resizableColumn) {
        this.setColumnsPosition({});
      }

      if (this.draggableColumn) {
        this.setColumnsWidth({});
        this.$nextTick(this.calculateColumnsWidths);
      }
    },

    updateItemsPerPage(limit) {
      this.$emit('update:items-per-page', limit);
    },

    updatePage(page) {
      this.$emit('update:page', page);
    },

    changeHeaderPositionOnResize() {
      if (this.stickyHeader) {
        this.changeHeaderPosition();
      }

      if (this.selecting) {
        this.calculateRowsPositions();
      }
    },

    openRootCauseDiagram(entity) {
      this.$modals.show({
        name: MODALS.entitiesRootCauseDiagram,
        config: {
          entity,
        },
      });
    },
  },
};
</script>

<style lang="scss">
.alarms-list-table {
  .theme--light & {
    --alarms-list-table-border-color: rgba(0, 0, 0, 0.12);
  }

  .theme--dark & {
    --alarms-list-table-border-color: rgba(255, 255, 255, 0.12);
  }

  &__top-pagination {
    position: relative;
    min-height: 48px;
    background: var(--v-background-base);
    z-index: 2;
    transition: .3s cubic-bezier(.25, .8, .5,1);
    transition-property: opacity, background-color;

    &:after {
      content: ' ';
      width: 100%;
      height: 2px;
      background: inherit;
      position: absolute;
      left: 0;
      right: 0;
      bottom: -1px;
    }

    &--left {
      padding-right: 80px;
    }

    &--center-absolute {
      position: absolute;
      left: 50%;
      transform: translate(-50%, 0);
    }

    &--right-absolute {
      position: absolute;
      right: 0;
    }
  }

  &__resize-handler {
    cursor: col-resize;

    display: flex;
    justify-content: center;

    width: 15px;

    position: absolute;
    right: -7px;
    top: 0;

    height: 100%;

    z-index: 2;
  }

  &__dragging-handler {
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: grab;
    z-index: 1;

    &:after {
      content: ' ';
      position: absolute;
      transition: .3s cubic-bezier(.25, .8, .5,1);
      top: 0;
      right: 0;
      bottom: 0;
      left: 0;
      background: var(--v-secondary-base);
      opacity: 0.0;
    }

    &:hover:after {
      opacity: 0.1;
    }
  }

  .alarm-list-row {
    position: relative;

    &:after {
      content: '';
      position: absolute;
      width: 100%;
      height: 100%;
      top: 0;
      left: 0;
      opacity: 0;
      pointer-events: none;
      background: rgba(200, 220, 200, .3);
      transition: opacity linear .3s;
    }
  }

  &__selecting {
    & > .v-data-table__wrapper > table > tbody > .alarm-list-row:after {
      pointer-events: auto;
      opacity: 1;
    }

    &--text-unselectable {
      * {
        user-select: none;
      }
    }
  }

  &__grid {
    & > .v-data-table__wrapper > table {
      & > tbody > tr > td,
      & > thead > tr > th {
        position: relative;

        &:after {
          content: ' ';
          background: rgba(0, 0, 0, 0.12);
          position: absolute;
          right: -1px;
          top: 0;
          width: 1px;
          height: 100%;
        }
      }
    }
  }

  &--fixed {
    & > .v-data-table__wrapper > table {
      table-layout: fixed;
      /**
       * TODO: Should be used v-bind later. We should update compiler.
       * Current compiler cannot to handle script setup and v-bind
       */
      width: var(--alarms-list-table-width);
      max-width: unset;
      min-width: 100%;

      & > thead > tr > th {
        word-break: break-all;
        white-space: pre-wrap;
      }
    }
  }

  &--wrapped {
    & > .v-data-table__wrapper > table > tbody > tr > td:not(:last-of-type) {
      word-break: break-all;
      word-wrap: break-word;
    }
  }

  &--truncated {
    .color-indicator {
      max-width: 100%;
    }

    .alarms-column-cell__layout .alarm-column-cell__text {
      display: grid;
    }

    .alarm-list-row__cell {
      .alarm-column-cell__text > span,
      .alarm-column-value {
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        display: block;
      }
    }
  }

  tbody {
    position: relative;
  }

  thead {
    position: relative;
    transition: .3s cubic-bezier(.25, .8, .5,1);
    transition-property: opacity, background-color;
    z-index: 1;

    &.head-shadow {
      tr:first-child {
        box-shadow: 0 1px 10px 0 rgba(0, 0, 0, 0.12) !important;

        &:after {
          content: unset;
        }
      }
    }

    tr {
      background: var(--v-table-background-base);
      transition: background-color .3s cubic-bezier(.25,.8,.5,1);

      .theme--dark & {
        background: var(--v-table-background-base);
      }

      th {
        position: relative;
        transition: none;
      }
    }
  }

  tr:not(.v-data-table__expanded) th:first-of-type {
    width: 120px !important;
  }

  &.v-data-table--dense,
  &.v-data-table--ultra-dense {
    thead tr:not(.v-data-table__expanded) th:first-of-type {
      width: 100px !important;
    }
  }

  &.columns-lg table {
    &:not(.v-data-table--dense) {
      td, th {
        padding: 0 8px;
      }
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

  &.columns-md .v-table {
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

  &.columns-sm .v-table {
    td, th {
      padding: 0 12px;
    }
  }
}
</style>
