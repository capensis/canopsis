<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title")
      v-divider
      field-default-sort-column(v-model="settings.widget.default_sort_column")
      v-divider
      field-columns(v-model="settings.widget.widget_columns")
      v-divider
      field-context-entities-types-filter(v-model="settings.widget_preferences.selectedTypes")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import pick from 'lodash/pick';
import cloneDeep from 'lodash/cloneDeep';

import { SIDE_BARS } from '@/constants';
import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldTitle from '../partial/fields/title.vue';
import FieldDefaultSortColumn from '../partial/fields/default-sort-column.vue';
import FieldColumns from '../partial/fields/columns.vue';
import FieldContextEntitiesTypesFilter from '../partial/fields/context-entities-types-filter.vue';

/**
 * Component to regroup the entities list settings fields
 */
export default {
  name: SIDE_BARS.contextSettings,
  components: {
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldContextEntitiesTypesFilter,
  },
  mixins: [
    widgetSettingsMixin,
  ],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: {
          title: widget.title,
          default_sort_column: cloneDeep(widget.default_sort_column),
          widget_columns: cloneDeep(widget.widget_columns),
        },
        widget_preferences: {
          selectedTypes: [],
        },
      },
    };
  },
  created() {
    this.settings.widget_preferences = pick(this.userPreference.widget_preferences, 'selectedTypes');
  },
};
</script>
