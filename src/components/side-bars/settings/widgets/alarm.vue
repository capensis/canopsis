<template lang="pug">
  div
    v-list.pt-0(expand)
      field-row-grid-size(
      :rowId.sync="rowId",
      :size.sync="settings.widget.size",
      :availableRows="getWidgetAvailableRows(config.widget._id)",
      @createRow="createRow"
      :rowForCreation.sync="rowForCreation"
      )
      v-divider
      field-title(v-model="settings.widget.title")
      v-divider
      field-default-sort-column(v-model="settings.widget.parameters.sort")
      v-divider
      field-columns(v-model="settings.widget.parameters.widgetColumns")
      v-divider
      field-periodic-refresh(v-model="settings.widget.parameters.periodicRefresh")
      v-divider
      field-default-elements-per-page(v-model="settings.widget_preferences.itemsPerPage")
      v-divider
      field-opened-resolved-filter(v-model="settings.widget.parameters.alarmsStateFilter")
      v-divider
      field-filters(
      v-model="settings.widget_preferences.mainFilter",
      :filters.sync="settings.widget_preferences.viewFilters"
      )
      v-divider
      field-info-popup(v-model="settings.widget.parameters.infoPopups")
      v-divider
      field-more-info(v-model="settings.widget.parameters.moreInfoTemplate")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import get from 'lodash/get';
import cloneDeep from 'lodash/cloneDeep';
import { normalize, denormalize } from 'normalizr';

import { viewSchema } from '@/store/schemas';

import { PAGINATION_LIMIT } from '@/config';
import { SIDE_BARS } from '@/constants';
import widgetSettingsMixin from '@/mixins/widget/settings';
import entitiesViewMixin from '@/mixins/entities/view/view';

import FieldRowGridSize from '../partial/fields/row-grid-size.vue';
import FieldTitle from '../partial/fields/title.vue';
import FieldDefaultSortColumn from '../partial/fields/default-sort-column.vue';
import FieldColumns from '../partial/fields/columns.vue';
import FieldPeriodicRefresh from '../partial/fields/periodic-refresh.vue';
import FieldDefaultElementsPerPage from '../partial/fields/default-elements-per-page.vue';
import FieldOpenedResolvedFilter from '../partial/fields/opened-resolved-filter.vue';
import FieldFilters from '../partial/fields/filters.vue';
import FieldInfoPopup from '../partial/fields/info-popup.vue';
import FieldMoreInfo from '../partial/fields/more-info.vue';

/**
 * Component to regroup the alarms list settings fields
 */
export default {
  name: SIDE_BARS.alarmSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldPeriodicRefresh,
    FieldDefaultElementsPerPage,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldInfoPopup,
    FieldMoreInfo,
  },
  mixins: [entitiesViewMixin, widgetSettingsMixin],
  data() {
    const { widget } = this.config;

    return {
      entities: {
        view: {},
        viewRow: {},
        widget: {},
      },
      rowId: get(widget, '_embedded.parentId', null),
      rowForCreation: null,
      settings: {
        widget: cloneDeep(widget),
        widget_preferences: {
          itemsPerPage: PAGINATION_LIMIT,
          viewFilters: [],
          mainFilter: {},
        },
      },
    };
  },
  computed: {
    availableRows() {
      return this.getWidgetAvailableRows(this.config.widget._id);
    },
    localView() {
      return denormalize(this.view._id, viewSchema, this.entities) || { rows: [] };
    },
    getWidgetAvailableRows() {
      return widgetId => this.localView.rows.map((row) => {
        const availableSize = row.widgets.reduce((acc, widget) => {
          if (widget._id !== widgetId) {
            acc.sm -= widget.size.sm;
            acc.md -= widget.size.md;
            acc.lg -= widget.size.lg;
          }

          return acc;
        }, { sm: 12, md: 12, lg: 12 });

        return {
          _id: row._id,
          title: row.title,

          availableSize,
        };
      }).filter(({ availableSize }) =>
        availableSize.sm >= 3 &&
        availableSize.md >= 3 &&
        availableSize.lg >= 3);
    },
  },
  mounted() {
    const { itemsPerPage, viewFilters, mainFilter } = this.userPreference.widget_preferences;
    const { entities } = normalize(this.view, viewSchema);

    this.settings.widget_preferences = {
      itemsPerPage,
      viewFilters,
      mainFilter,
    };
    this.entities = entities;
  },
  methods: {
    createRow(row) {
      const { rows } = this.entities.view[this.view._id];

      this.$set(this.entities.viewRow, row._id, row);
      this.$set(this.entities.view[this.view._id], 'rows', [...rows, row._id]);
    },

    prefixFormatter(value) {
      return value.replace('alarm.', 'v.');
    },

    prepareSettingsWidget() {
      const { widget } = this.settings;

      return {
        ...widget,
        parameters: {
          ...widget.parameters,

          widgetColumns: widget.parameters.widgetColumns.map(v => ({
            ...v,
            value: this.prefixFormatter(v.value),
          })),

          infoPopups: widget.parameters.infoPopups.map(v => ({
            ...v,
            column: this.prefixFormatter(v.column),
          })),
        },
      };
    },
  },
};
</script>
