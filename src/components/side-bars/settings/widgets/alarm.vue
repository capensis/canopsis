<template lang="pug">
  div
    v-list.pt-0(expand)
      v-combobox(
      :value="row",
      @change="updateRow"
      :items="availableRows",
      label="test",
      :search-input.sync="search"
      :error-messages="errors.collect('group')",
      item-text="title",
      item-value="title"
      )
        template(slot="no-data")
          v-list-tile
            v-list-tile-content
              v-list-tile-title(v-html="$t('modals.createView.noData')")
      v-container(v-if="row")
        v-slider(
        :value="3",
        :max="row.availableColumns.sm"
        :min="0"
        ticks="always"
        @input="$emit('input', $event)"
        v-validate="'min_value:3'",
        data-vv-name="row.sm",
        :error-messages="errors.collect('row.sm')",
        always-dirty,
        thumb-label,
        )
        v-slider(
        :value="3",
        :max="row.availableColumns.md"
        :min="0"
        ticks="always"
        @input="$emit('input', $event)"
        always-dirty,
        thumb-label,
        )
        v-slider(
        :value="3",
        :max="row.availableColumns.lg"
        :min="0"
        ticks="always"
        @input="$emit('input', $event)"
        always-dirty,
        thumb-label,
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
import cloneDeep from 'lodash/cloneDeep';

import { PAGINATION_LIMIT } from '@/config';
import { SIDE_BARS } from '@/constants';
import widgetSettingsMixin from '@/mixins/widget/settings';
import entitiesViewMixin from '@/mixins/entities/view';

import FieldGridSize from '../partial/fields/grid-size.vue';
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
    FieldGridSize,
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
  mixins: [widgetSettingsMixin, entitiesViewMixin],
  data() {
    const { widget } = this.config;

    return {
      row: null,
      search: null,
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
  mounted() {
    const { itemsPerPage, viewFilters, mainFilter } = this.userPreference.widget_preferences;

    this.settings.widget_preferences = {
      itemsPerPage,
      viewFilters,
      mainFilter,
    };
  },
  methods: {
    updateRow(value) {
      if (value !== this.row) {
        if (typeof value === 'string') {
          let newRow = this.availableRows.find(v => v.title === value);

          if (!newRow) {
            newRow = { title: value, _id: 'asdasd', availableColumns: { sm: 12, md: 12, lg: 12 } };
          }

          this.row = newRow;
        } else {
          this.row = value;
        }
      }
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
