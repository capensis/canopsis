<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title")
      v-divider
      field-default-column-sort(v-model="settings.widget.default_sort_column")
      v-divider
      field-columns(v-model="settings.widget.widget_columns")
      v-divider
      field-periodic-refresh(v-model="settings.widget.periodicRefresh")
      v-divider
      field-default-elements-per-page(v-model="settings.widget_preferences.itemsPerPage")
      v-divider
      field-opened-resolved-filter(v-model="settings.widget.alarms_state_filter")
      v-divider
      field-filters(
      v-model="settings.widget_preferences.selected_filter",
      :filters.sync="settings.widget_preferences.user_filters"
      )
      v-divider
      field-info-popup(v-model="settings.widget.popup", :widget="widget")
      v-divider
      field-more-info(v-model="settings.widget.more_infos_popup")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import { PAGINATION_LIMIT } from '@/config';
import { SIDE_BARS } from '@/constants';
import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldTitle from '../partial/fields/title.vue';
import FieldDefaultColumnSort from '../partial/fields/default-column-sort.vue';
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
    FieldTitle,
    FieldDefaultColumnSort,
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
    const { widget } = this.config;

    return {
      settings: {
        widget: {
          title: widget.title,
          default_sort_column: cloneDeep(widget.default_sort_column) || {},
          widget_columns: cloneDeep(widget.widget_columns) || [],
          periodicRefresh: cloneDeep(widget.periodicRefresh) || {},
          alarms_state_filter: {
            opened: widget.alarms_state_filter.opened || widget.alarms_state_filter.state === 'opened',
            resolved: widget.alarms_state_filter.resolved || widget.alarms_state_filter.state === 'resolved',
          },
          popup: cloneDeep(widget.popup) || [],
          more_infos_popup: widget.more_infos_popup || '',
        },
        widget_preferences: {
          itemsPerPage: PAGINATION_LIMIT,
          user_filters: [],
          selected_filter: {},
        },
      },
    };
  },
  mounted() {
    this.settings.widget_preferences = {
      itemsPerPage: this.userPreference.widget_preferences.itemsPerPage,
      user_filters: this.userPreference.widget_preferences.user_filters,
      selected_filter: this.userPreference.widget_preferences.selected_filter,
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
        widget_columns: widget.widget_columns.map(v => ({
          ...v,
          value: this.prefixFormatter(v.value),
        })),
        popup: widget.popup.map(v => ({
          ...v,
          column: this.prefixFormatter(v.column),
        })),
      };
    },
  },
};
</script>
