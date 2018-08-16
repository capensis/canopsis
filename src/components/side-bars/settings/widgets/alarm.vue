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
      field-more-info
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import pick from 'lodash/pick';
import cloneDeep from 'lodash/cloneDeep';

import { PAGINATION_LIMIT } from '@/config';
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
 *
 * @prop {Object} widget - active widget
 */
export default {
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
          alarms_state_filter: cloneDeep(widget.alarms_state_filter) || {},
          popup: cloneDeep(widget.popup) || [],
        },
        widget_preferences: {
          itemsPerPage: PAGINATION_LIMIT,
          user_filters: [],
          selected_filter: {},
        },
      },
    };
  },
  created() {
    this.settings.widget_preferences = pick(this.userPreference.widget_preferences, [
      'itemsPerPage',
      'user_filters',
      'selected_filter',
    ]);
  },
};
</script>
