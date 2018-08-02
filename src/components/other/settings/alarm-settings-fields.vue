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
      field-default-elements-per-page(v-model="settings.userPreference.itemsPerPage")
      v-divider
      field-opened-resolved-filter(v-model="settings.widget.alarms_state_filter")
      v-divider
      field-filters(v-model="settings.userPreference.selected_filter", :filters="settings.userPreference.user_filters")
      v-divider
      field-info-popup(v-model="settings.widget.popup", :widget="widget")
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
        widget: {
          title: this.widget.title,
          default_sort_column: cloneDeep(this.widget.default_sort_column),
          widget_columns: cloneDeep(this.widget.widget_columns),
          periodicRefresh: cloneDeep(this.widget.periodicRefresh),
          alarms_state_filter: cloneDeep(this.widget.alarms_state_filter),
          popup: cloneDeep(this.widget.popup),
        },
        userPreference: {},
      },
    };
  },
  created() {
    this.settings.userPreference.itemsPerPage = this.userPreference.widget_preferences.itemsPerPage;
    this.settings.userPreference.user_filters = this.userPreference.widget_preferences.user_filters;
    this.settings.userPreference.selected_filter = cloneDeep(this.userPreference.widget_preferences.selected_filter);
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const widget = {
          ...this.widget,
          ...this.settings.widget,
        };

        const userPreference = {
          ...this.userPreference,
          widget_preferences: {
            ...this.userPreference.widget_preferences,
            ...this.settings.userPreference,
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
