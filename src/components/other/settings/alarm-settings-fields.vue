<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.title")
      v-divider
      field-default-column-sort(v-model="settings.defaultSortColumn")
      v-divider
      field-columns(v-model="settings.widgetColumns")
      v-divider
      field-periodic-refresh
      v-divider
      field-default-elements-per-page
      v-divider
      field-opened-resolved-filter(v-model="settings.alarmStateFilter")
      v-divider
      field-filters
      v-divider
      field-info-popup(:widget="widget")
      v-divider
      field-more-info
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldDefaultColumnSort from '@/components/other/settings/fields/default-column-sort.vue';
import FieldColumns from '@/components/other/settings/fields/columns.vue';
import FieldPeriodicRefresh from '@/components/other/settings/fields/periodic-refresh.vue';
import FieldDefaultElementsPerPage from '@/components/other/settings/fields/default-elements-per-page.vue';
import FieldOpenedResolvedFilter from '@/components/other/settings/fields/opened-resolved-filter.vue';
import FieldFilters from '@/components/other/settings/fields/filters.vue';
import FieldInfoPopup from '@/components/other/settings/fields/info-popup.vue';
import FieldMoreInfo from '@/components/other/settings/fields/more-info.vue';

import entitiesWidgetMixin from '@/mixins/entities/widget';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import { ALARM_FILTER_STATES } from '@/constants';

/**
 * Component to regroup the alarms list settings fields
 *
 * @prop {Object} widget - active widget
 * @prop {bool} isNew - is widget new
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
  mixins: {
    entitiesWidgetMixin,
    entitiesUserPreferenceMixin,
  },
  props: {
    widget: {
      type: Object,
      required: true,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      settings: {
        title: this.widget.title,
        defaultSortColumn: cloneDeep(this.widget.default_sort_column),
        widgetColumns: cloneDeep(this.widget.widget_columns),
        alarmStateFilter: this.widget.alarms_state_filter ?
          this.widget.alarms_state_filter.state : ALARM_FILTER_STATES.opened,
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();
      console.warn(isFormValid, this.settings);
    },
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
