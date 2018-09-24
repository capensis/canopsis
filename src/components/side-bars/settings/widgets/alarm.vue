<template lang="pug">
  div
    v-list.pt-0(expand)
      field-row-grid-size(
      :rowId.sync="settings.rowId",
      :size.sync="settings.widget.size",
      :availableRows="availableRows",
      @createRow="createRow"
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

import { PAGINATION_LIMIT } from '@/config';
import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';

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
  mixins: [widgetSettingsMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: cloneDeep(widget),
        widget_preferences: {
          itemsPerPage: PAGINATION_LIMIT,
          viewFilters: [],
          mainFilter: {},
        },
      },
    };
  },
  mounted() {
    const { widget_preferences: widgetPreference } = this.userPreference;

    this.settings.widget_preferences = {
      itemsPerPage: get(widgetPreference, 'itemsPerPage', PAGINATION_LIMIT),
      viewFilters: get(widgetPreference, 'viewFilters', []),
      mainFilter: get(widgetPreference, 'mainFilter', {}),
    };
  },
  methods: {
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
