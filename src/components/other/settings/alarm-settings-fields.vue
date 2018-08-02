<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.title")
      v-divider
      field-default-column-sort(v-model="settings.defaultSortColumn")
      v-divider
      field-columns(v-model="settings.widgetColumns")
      v-divider
      field-periodic-refresh(v-model="settings.periodicRefresh")
      v-divider
      field-default-elements-per-page(v-model="settings.itemsPerPage")
      v-divider
      field-opened-resolved-filter(v-model="settings.alarmStateFilter")
      v-divider
      field-filters(v-model="settings.selectedFilter", :filters="settings.filters")
      v-divider
      field-info-popup(v-model="settings.popups", :widget="widget")
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
  mixins: [
    entitiesWidgetMixin,
    entitiesUserPreferenceMixin,
  ],
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
        periodicRefresh: cloneDeep(this.widget.periodicRefresh) || {},
        alarmStateFilter: cloneDeep(this.widget.alarms_state_filter) || {},
        popups: cloneDeep(this.widget.popup) || [],
      },
    };
  },
  created() {
    this.settings.itemsPerPage = this.userPreference.widget_preferences.itemsPerPage;
    this.settings.filters = this.userPreference.widget_preferences.user_filters;
    this.settings.selectedFilter = this.userPreference.widget_preferences.selected_filter;
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const widget = {
          ...this.widget,
          title: this.settings.title,
          default_sort_column: this.settings.defaultSortColumn,
          widget_columns: this.settings.widgetColumns,
          periodicRefresh: this.settings.periodicRefresh,
          alarms_state_filter: this.settings.alarmStateFilter,
          popup: this.settings.popups,
        };

        const userPreference = {
          ...this.userPreference,
          widget_preferences: {
            ...this.userPreference.widget_preferences,
            itemsPerPage: this.settings.itemsPerPage,
            user_filters: this.settings.filters,
            selected_filter: this.settings.selectedFilter,
          },
        };

        const actions = [this.createUserPreference({ userPreference })];

        if (this.isNew) {
          actions.push(this.createWidget({ widget }));
        } else {
          actions.push(this.updateWidget({ widget }));
        }

        await this.$store.dispatch('query/startPending', { id: widget.id }); // TODO: fix it
        await Promise.all(actions);
        await this.$store.dispatch('query/stopPending', { id: widget.id }); // TODO: fix it

        this.$emit('closeSettings');
      }
    },
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
