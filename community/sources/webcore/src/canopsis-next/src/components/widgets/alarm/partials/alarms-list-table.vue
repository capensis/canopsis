<template lang="pug">
  v-flex(
    v-on="wrapperListeners",
    v-resize="changeHeaderPositionOnResize"
  )
    c-empty-data-table-columns(v-if="!columns.length")
    div(v-else)
      v-layout.alarms-list-table__top-pagination.px-4.position-relative(
        v-if="totalItems && (densable || !hideActions || !hidePagination)",
        ref="actions",
        row,
        align-center
      )
        v-flex.alarms-list-table__top-pagination--left(v-if="densable || !hideActions", xs6)
          v-layout(row, align-center, justify-start)
            c-density-btn-toggle(v-if="densable", :value="dense", @change="$emit('update:dense', $event)")
            v-fade-transition
              v-flex.px-1(v-show="unresolvedSelected.length")
                mass-actions-panel(
                  v-if="!hideActions",
                  :items="unresolvedSelected",
                  :widget="widget",
                  :refresh-alarms-list="refreshAlarmsList",
                  @clear:items="clearSelected"
                )
        v-flex.alarms-list-table__top-pagination--center-absolute(xs4)
          c-pagination(
            v-if="!hidePagination",
            :page="pagination.page",
            :limit="pagination.rowsPerPage",
            :total="totalItems",
            type="top",
            @input="updateQueryPage"
          )
      v-data-table.alarms-list-table(
        ref="dataTable",
        v-model="selected",
        :class="vDataTableClass",
        :items="alarms",
        :headers="headers",
        :total-items="totalItems",
        :pagination="pagination",
        :select-all="selectable",
        :loading="loading || columnsFiltersPending",
        :expand="expandable",
        :dense="isMediumHeight",
        :ultra-dense="isSmallHeight",
        item-key="_id",
        hide-actions,
        multi-sort,
        @update:pagination="updatePaginationHandler"
      )
        template(#progress="")
          v-fade-transition
            v-progress-linear(color="primary", height="2", indeterminate)
        template(#headerCell="{ header }")
          alarm-header-cell(
            :header="header",
            :selected-tag="selectedTag",
            @clear:tag="$emit('clear:tag')"
          )
        template(#items="props")
          alarms-list-row(
            v-model="props.selected",
            v-on="rowListeners",
            :ref="`row${props.item._id}`",
            :key="props.item._id",
            :selectable="selectable",
            :expandable="expandable",
            :row="props",
            :widget="widget",
            :columns="preparedColumns",
            :parent-alarm="parentAlarm",
            :is-tour-enabled="checkIsTourEnabledForAlarmByIndex(props.index)",
            :refresh-alarms-list="refreshAlarmsList",
            :selecting="selecting",
            :selected-tag="selectedTag",
            :hide-actions="hideActions",
            :medium="isMediumHeight",
            :small="isSmallHeight",
            @select:tag="$emit('select:tag', $event)"
          )
        template(#expand="{ item, index }")
          alarms-expand-panel(
            :alarm="item",
            :parent-alarm-id="parentAlarmId",
            :widget="widget",
            :hide-children="hideChildren",
            :is-tour-enabled="checkIsTourEnabledForAlarmByIndex(index)"
          )
    c-table-pagination(
      v-if="!hidePagination",
      :total-items="totalItems",
      :rows-per-page="pagination.rowsPerPage",
      :page="pagination.page",
      @update:page="updateQueryPage",
      @update:rows-per-page="updateRecordsPerPage"
    )
    component(
      v-bind="additionalComponent.props",
      v-on="additionalComponent.on",
      :is="additionalComponent.is"
    )
</template>

<script>
import { get, intersectionBy, throttle } from 'lodash';

import { TOP_BAR_HEIGHT } from '@/config';
import { ALARM_DENSE_TYPES, ALARMS_LIST_HEADER_OPACITY_DELAY } from '@/constants';

import { isActionAvailableForAlarm, calculateAlarmLinksColumnWidth } from '@/helpers/entities';

import featuresService from '@/services/features';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { widgetColumnsAlarmMixin } from '@/mixins/widget/columns/alarm';

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
    pagination: {
      type: Object,
      default: () => ({}),
    },
    columns: {
      type: Array,
      default: () => [],
    },
    isTourEnabled: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    expandable: {
      type: Boolean,
      default: false,
    },
    stickyHeader: {
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
  },
  data() {
    return {
      prevEvent: null,
      selecting: false,
      selected: [],
    };
  },

  computed: {
    wrapperListeners() {
      return this.selectable
        ? { mousemove: this.throttledMousemoveHandler }
        : {};
    },

    topBarHeight() {
      return this.shownHeader ? TOP_BAR_HEIGHT : 0;
    },

    unresolvedSelected() {
      return this.selected.filter(item => isActionAvailableForAlarm(item));
    },

    expanded() {
      return this.$refs.dataTable.expanded;
    },

    hasInstructionsAlarms() {
      return this.alarms.some(alarm => alarm.assigned_instructions?.length);
    },

    headers() {
      const headers = this.preparedColumns.map((column) => {
        if (column.linksInRowCount) {
          const linksCounts = this.alarms.map(alarm => Object.values(get(alarm, column.value, {})).flat().length);
          const maxLinksCount = Math.max(...linksCounts);
          const actualInlineLinksCount = maxLinksCount > column.inlineLinksCount
            ? column.inlineLinksCount + 1
            : maxLinksCount;

          const width = calculateAlarmLinksColumnWidth(
            this.dense,
            Math.max(Math.min(actualInlineLinksCount, column.linksInRowCount), 1),
          );

          return { ...column, width };
        }

        return column;
      });

      if (!this.hideActions) {
        headers.push({ text: this.$t('common.actionsLabel'), sortable: false });
      }

      if ((this.expandable || this.hasInstructionsAlarms) && !this.selectable) {
        /**
         * We need it for the expand panel open button
         */
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

      const { label } = Object.values(COLUMNS_SIZES_VALUES)
        .find(({ min, max }) => columnsLength >= min && columnsLength <= max);

      return {
        [`columns-${label}`]: true,
        'alarms-list-table__selecting': this.selecting,
        'alarms-list-table__selecting--text-unselectable': this.selectingMousePressed,
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
        return featuresService.call('components.alarmListTable.computed.additionalComponent', this);
      }

      return {};
    },

    tableHeader() {
      return this.$el.querySelector('.v-table__overflow > table > thead');
    },

    tableBody() {
      return this.$el.querySelector('.v-table__overflow > table > tbody');
    },

    isMediumHeight() {
      return this.dense === ALARM_DENSE_TYPES.medium;
    },

    isSmallHeight() {
      return this.dense === ALARM_DENSE_TYPES.small;
    },

    parentAlarmId() {
      return this.parentAlarm?._id;
    },

    selectingMousePressed() {
      return this.selecting && !!this.prevEvent;
    },
  },

  watch: {
    alarms(alarms) {
      this.selected = intersectionBy(alarms, this.selected, '_id');
    },

    stickyHeader(stickyHeader) {
      if (stickyHeader) {
        this.calculateHeaderOffsetPosition();
        this.setHeaderPosition();
        this.addShadowToHeader();

        window.addEventListener('scroll', this.changeHeaderPosition);
      } else {
        window.removeEventListener('scroll', this.changeHeaderPosition);

        this.resetHeaderPosition();
      }
    },
  },

  created() {
    this.actionsTranslateY = 0;
    this.translateY = 0;
    this.previousTranslateY = 0;
    this.throttledMousemoveHandler = throttle(this.mousemoveHandler, 50);
  },

  async mounted() {
    if (this.stickyHeader) {
      window.addEventListener('scroll', this.changeHeaderPosition);
    }

    if (this.selectable) {
      window.addEventListener('keydown', this.enableSelecting);
      window.addEventListener('keyup', this.disableSelecting);
      window.addEventListener('mousedown', this.mousedownHandler);
      window.addEventListener('mouseup', this.mouseupHandler);
    }
  },
  updated() {
    if (this.selecting) {
      this.calculateRowsPositions();
    }
  },
  beforeDestroy() {
    window.removeEventListener('scroll', this.changeHeaderPosition);
    window.removeEventListener('keydown', this.enableSelecting);
    window.removeEventListener('keyup', this.disableSelecting);
    window.removeEventListener('mousedown', this.mousedownHandler);
    window.removeEventListener('mouseup', this.mouseupHandler);
  },

  methods: {
    calculateRowsPositions() {
      this.rowsPositions = Object.entries(this.$refs).reduce((acc, [key, value]) => {
        if (!key.startsWith('row') || !value) {
          return acc;
        }

        const position = value.$el.getBoundingClientRect();

        acc.push({
          position: {
            x1: position.x,
            x2: position.x + position.width,
            y1: position.y,
            y2: position.y + position.height,
          },
          row: value.$options.propsData.row,
        });

        return acc;
      }, []);
    },

    getIntersectRowsByPosition(newX, newY, prevX, prevY) {
      return this.rowsPositions?.reduce((acc, { position, row }) => {
        if (
          (prevX >= position.x1 && prevX <= position.x2 && prevY >= position.y1 && prevY <= position.y2)
          || (newX < position.x1 && prevX < position.x1)
          || (newX > position.x2 && prevX > position.x2)
          || (newY < position.y1 && prevY < position.y1)
          || (newY > position.y2 && prevY > position.y2)
        ) {
          return acc;
        }

        acc.push(row);

        return acc;
      }, []) ?? [];
    },

    mousedownHandler(event) {
      this.prevEvent = event;
    },

    mouseupHandler() {
      this.prevEvent = null;
    },

    mousemoveHandler(event) {
      if (!event.ctrlKey || !event.buttons || !this.prevEvent) {
        return;
      }

      const rows = this.getIntersectRowsByPosition(
        event.clientX,
        event.clientY,
        this.prevEvent.clientX,
        this.prevEvent.clientY,
      );

      this.prevEvent = event;

      rows.forEach(row => this.toggleSelected(row.item));
    },

    toggleSelected(alarm) {
      const index = this.selected.findIndex(({ _id: id }) => id === alarm._id);

      if (index === -1) {
        this.selected.push(alarm);

        return;
      }

      this.selected.splice(index, 1);
    },

    clearSelected() {
      this.selected = [];
    },

    updateRecordsPerPage(limit) {
      this.$emit('update:rows-per-page', limit);
    },

    updateQueryPage(page) {
      this.$emit('update:page', page);
    },

    enableSelecting({ key }) {
      if (key === 'Control') {
        this.selecting = true;
      }
    },

    disableSelecting({ key }) {
      if (key === 'Control') {
        this.selecting = false;
      }
    },

    startScrolling() {
      if (this.translateY !== this.previousTranslateY) {
        this.tableHeader.style.opacity = '0';

        if (this.$refs.actions) {
          this.$refs.actions.style.opacity = '0';
        }
      }

      this.scrooling = true;
    },

    finishScrolling() {
      if (!Number(this.tableHeader.style.opacity)) {
        this.tableHeader.style.opacity = '1.0';

        if (this.$refs.actions) {
          this.$refs.actions.style.opacity = '1.0';
        }
      }

      this.scrooling = false;
    },

    clearFinishTimer() {
      if (this.finishTimer) {
        clearTimeout(this.finishTimer);
      }
    },

    setHeaderPosition() {
      this.tableHeader.style.transform = `translateY(${this.translateY}px)`;

      if (this.$refs.actions) {
        this.$refs.actions.style.transform = `translateY(${this.actionsTranslateY}px)`;
      }
    },

    calculateHeaderOffsetPosition() {
      const { top: headerTop } = this.tableHeader.getBoundingClientRect();
      const { height: bodyHeight } = this.tableBody.getBoundingClientRect();
      const { top: actionsTop = 0, height: actionsHeight = 0 } = this.$refs.actions?.getBoundingClientRect() ?? {};

      const offset = headerTop - this.translateY - this.topBarHeight - actionsHeight;
      const actionsOffset = actionsTop - this.actionsTranslateY - this.topBarHeight;

      this.previousTranslateY = this.actionsTranslateY;
      this.translateY = Math.min(bodyHeight, Math.max(0, -offset));
      this.actionsTranslateY = Math.min(bodyHeight, Math.max(0, -actionsOffset));
    },

    addShadowToHeader() {
      this.tableHeader.classList.add('head-shadow');
    },

    removeShadowFromHeader() {
      this.tableHeader.classList.remove('head-shadow');
    },

    changeHeaderPosition() {
      this.clearFinishTimer();

      this.calculateHeaderOffsetPosition();
      this.setHeaderPosition();

      if (!this.actionsTranslateY || !this.translateY) {
        this.removeShadowFromHeader();
        this.finishScrolling();

        return;
      }

      if (!this.scrooling) {
        this.addShadowToHeader();
        this.startScrolling();
      }

      this.finishTimer = setTimeout(this.finishScrolling, ALARMS_LIST_HEADER_OPACITY_DELAY);
    },

    resetHeaderPosition() {
      this.translateY = 0;
      this.actionsTranslateY = 0;
      this.previousTranslateY = 0;

      this.setHeaderPosition();
      this.clearFinishTimer();
      this.removeShadowFromHeader();
    },

    changeHeaderPositionOnResize() {
      if (this.stickyHeader) {
        this.changeHeaderPosition();
      }

      if (this.selecting) {
        this.calculateRowsPositions();
      }
    },

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
  &__top-pagination {
    position: relative;
    min-height: 48px;
    background: var(--v-background-base);
    z-index: 2;
    transition: .3s cubic-bezier(.25, .8, .5,1);
    transition-property: opacity, background-color;

    .theme--dark & {
      background: #424242;
    }

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
    & > .v-table__overflow > table > tbody > .alarm-list-row:after {
      pointer-events: auto;
      opacity: 1;
    }

    &--text-unselectable {
      * {
        user-select: none;
      }
    }
  }

  table {
    max-width: unset;
    min-width: 100%;
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
        border-bottom: none !important;
        box-shadow: 0 1px 10px 0 rgba(0, 0, 0, 0.12) !important;
      }
    }

    tr {
      background: white;
      transition: background-color .3s cubic-bezier(.25,.8,.5,1);

      .theme--dark & {
        background: #424242;
      }

      th {
        transition: none;
      }
    }
  }

  &.columns-lg .v-table {
    &:not(.v-datatable--dense) {
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
}
</style>
