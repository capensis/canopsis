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
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(v-model="settings.widget.parameters.sort")
          v-divider
          field-columns(v-model="settings.widget.parameters.widgetColumns")
          v-divider
          field-context-entities-types-filter(v-model="settings.widget_preferences.selectedTypes")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import get from 'lodash/get';
import cloneDeep from 'lodash/cloneDeep';

import { SIDE_BARS } from '@/constants';
import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldContextEntitiesTypesFilter from './fields/context/context-entities-types-filter.vue';

/**
 * Component to regroup the entities list settings fields
 */
export default {
  name: SIDE_BARS.contextSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldContextEntitiesTypesFilter,
  },
  mixins: [
    widgetSettingsMixin,
  ],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: cloneDeep(widget),
        widget_preferences: {
          selectedTypes: [],
        },
      },
    };
  },
  created() {
    const { widget_preferences: widgetPreference } = this.userPreference;

    this.settings.widget_preferences = {
      selectedTypes: get(widgetPreference, 'selectedTypes', []),
    };
  },
};
</script>
