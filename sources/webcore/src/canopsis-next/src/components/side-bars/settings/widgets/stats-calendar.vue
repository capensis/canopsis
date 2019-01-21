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
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-filters(
          v-model="settings.widget.parameters.filters",
          :filters.sync="settings.widget.parameters.filters",
          hideSelect
          )
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
          field-levels-colors-selector(
          v-model="settings.widget.parameters.criticityLevelsColors",
          colorType="hex",
          hideSuffix
          )
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';
import sideBarSettingsWidgetAlarmMixin from '@/mixins/side-bar/settings/widgets/alarm';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldOpenedResolvedFilter from './fields/alarm/opened-resolved-filter.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import FieldCriticityLevels from './fields/stats/criticity-levels.vue';
import FieldLevelsColorsSelector from './fields/stats/levels-colors-selector.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldInfoPopup from './fields/alarm/info-popup.vue';
import FieldMoreInfo from './fields/alarm/more-info.vue';

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
  mixins: [widgetSettingsMixin, sideBarSettingsWidgetAlarmMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: this.prepareWidgetWithAlarmParametersSettings(cloneDeep(widget), true),
      },
    };
  },
  methods: {
    prepareWidgetSettings() {
      const { widget } = this.settings;

      return this.prepareWidgetWithAlarmParametersSettings(widget);
    },
  },
};
</script>
