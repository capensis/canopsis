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
      field-filters(:filters.sync="settings.widget.parameters.filters", hideSelect)
      v-divider
      field-opened-resolved-filter(v-model="settings.widget.parameters.alarmsStateFilter")
      v-divider
      field-switcher(
      v-model="settings.widget.parameters.considerPbehaviors",
      :title="$t('settings.considerPbehaviors.title')"
      )
      v-divider
      field-criticity-levels(v-model="settings.widget.parameters.criticityLevels")
      v-divider
      field-levels-colors-selector(v-model="settings.widget.parameters.criticityLevelsColors", hideSuffix)
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.titles.alarmListSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-columns(v-model="settings.widget.parameters.alarmsList.widgetColumns")
          v-divider
          field-default-elements-per-page(v-model="settings.widget.parameters.alarmsList.itemsPerPage")
          v-divider
          field-info-popup(v-model="settings.widget.parameters.alarmsList.infoPopups")
          v-divider
          field-more-info(v-model="settings.widget.parameters.alarmsList.moreInfoTemplate")
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
import FieldCriticityLevels from '../partial/fields/criticity-levels.vue';
import FieldLevelsColorsSelector from '../partial/fields/levels-colors-selector.vue';
import FieldColumns from '../partial/fields/columns.vue';
import FieldDefaultElementsPerPage from '../partial/fields/default-elements-per-page.vue';
import FieldInfoPopup from '../partial/fields/info-popup.vue';
import FieldMoreInfo from '../partial/fields/more-info.vue';

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
    FieldCriticityLevels,
    FieldLevelsColorsSelector,
    FieldColumns,
    FieldDefaultElementsPerPage,
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
      },
    };
  },
  computed: {
    entitiesType() {
      return ENTITIES_TYPES.entity;
    },
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

          alarmsList: {
            ...widget.parameters.alarmsList,

            widgetColumns: widget.parameters.alarmsList.widgetColumns.map(v => ({
              ...v,
              value: this.prefixFormatter(v.value),
            })),

            infoPopups: widget.parameters.alarmsList.infoPopups.map(v => ({
              ...v,
              column: this.prefixFormatter(v.column),
            })),
          },
        },
      };
    },
  },
};
</script>
