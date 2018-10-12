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
      field-filters(:filters.sync="settings.widget.parameters.filters")
      v-divider
      field-opened-resolved-filter(v-model="settings.widget.parameters.alarmsStateFilter")
      v-divider
      field-switcher(v-model="settings.widget.parameters.considerPbehaviors", title="Consider pbehaviors")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import { SIDE_BARS, ENTITIES_TYPES } from '@/constants';
import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldRowGridSize from '../partial/fields/row-grid-size.vue';
import FieldTitle from '../partial/fields/title.vue';
import FieldOpenedResolvedFilter from '../partial/fields/opened-resolved-filter.vue';
import FieldFilters from '../partial/fields/filters.vue';
import FieldSwitcher from '../partial/fields/switcher.vue';

/**
 * Component to regroup the entities list settings fields
 */
export default {
  name: SIDE_BARS.statsCalendarSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldSwitcher,
  },
  mixins: [widgetSettingsMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: cloneDeep(widget),
      },
    };
  },
  computed: {
    entitiesType() {
      return ENTITIES_TYPES.entity;
    },
  },
};
</script>
