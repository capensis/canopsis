<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-periodic-refresh(v-model="settings.widget.parameters.periodicRefresh")
      v-divider
      v-list-group(data-test="advancedSettings")
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(
            v-model="settings.widget.parameters.sort",
            :columns="settings.widget.parameters.widgetColumns",
            :columnsLabel="$t('settings.columnName')"
          )
          v-divider
          field-columns(
            v-model="settings.widget.parameters.widgetColumns",
            :label="$t('settings.columnNames')",
            with-html,
            with-state
          )
          v-divider
          field-columns(
            v-model="settings.widget.parameters.widgetGroupColumns",
            :label="$t('settings.groupColumnNames')",
            with-html,
            with-state
          )
          v-divider
          field-default-elements-per-page(v-model="settings.widget_preferences.itemsPerPage")
          v-divider
          field-opened-resolved-filter(v-model="settings.widget.parameters.alarmsStateFilter")
          v-divider
          template(v-if="hasAccessToListFilters")
            field-filters(
              v-model="settings.widget.parameters.mainFilter",
              :entities-type="$constants.ENTITIES_TYPES.alarm",
              :filters.sync="settings.widget.parameters.viewFilters",
              :condition.sync="settings.widget.parameters.mainFilterCondition",
              :has-access-to-add-filter="hasAccessToAddFilter",
              :has-access-to-edit-filter="hasAccessToEditFilter",
              @input="updateMainFilterUpdatedAt"
            )
            v-divider
          template(v-if="hasAccessToListRemediationInstructionsFilters")
            field-remediation-instructions-filters(
              v-model="settings.widget.parameters.remediationInstructionsFilters",
              :has-access-to-add-filter="hasAccessToAddRemediationInstructionsFilter",
              :has-access-to-edit-filter="hasAccessToEditRemediationInstructionsFilter"
            )
            v-divider
          field-live-reporting(v-model="settings.widget.parameters.liveReporting")
          v-divider
          field-info-popup(
            v-model="settings.widget.parameters.infoPopups",
            :columns="settings.widget.parameters.widgetColumns"
          )
          v-divider
          field-text-editor(
            data-test="widgetMoreInfoTemplate",
            v-model="settings.widget.parameters.moreInfoTemplate",
            :title="$t('settings.moreInfosModal')"
          )
          v-divider
          field-grid-range-size(
            v-model="settings.widget.parameters.expandGridRangeSize",
            :title="$t('settings.expandGridRangeSize')"
          )
          v-divider
          field-switcher(
            data-test="isHtmlEnabledOnTimeLine",
            v-model="settings.widget.parameters.isHtmlEnabledOnTimeLine",
            :title="$t('settings.isHtmlEnabledOnTimeLine')"
          )
          v-divider
          v-list-group(data-test="ackGroup")
            v-list-tile(slot="activator") Ack
            v-list.grey.lighten-4.px-2.py-0(expand)
            field-switcher(
              data-test="isAckNoteRequired",
              v-model="settings.widget.parameters.isAckNoteRequired",
              :title="$t('settings.isAckNoteRequired')"
            )
            v-divider
            field-switcher(
              data-test="isMultiAckEnabled",
              v-model="settings.widget.parameters.isMultiAckEnabled",
              :title="$t('settings.isMultiAckEnabled')"
            )
            v-divider
            field-fast-ack-output(v-model="settings.widget.parameters.fastAckOutput")
          v-divider
          field-switcher(
            v-model="settings.widget.parameters.isSnoozeNoteRequired",
            :title="$t('settings.isSnoozeNoteRequired')"
          )
          v-divider
          field-enabled-limit(
            v-model="settings.widget.parameters.linksCategoriesAsList",
            :title="$t('settings.linksCategoriesAsList')",
            :label="$t('settings.linksCategoriesLimit')"
          )
          v-divider
          field-export-csv-separator(
            v-model="settings.widget.parameters.exportCsvSeparator",
            :title="$t('settings.exportCsvSeparator')"
          )
      v-divider
    v-btn.primary(data-test="submitAlarms", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { get, cloneDeep } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';
import sideBarSettingsWidgetAlarmMixin from '@/mixins/side-bar/settings/widgets/alarm';
import rightsWidgetsAlarmsListFilters from '@/mixins/rights/widgets/alarms-list/filters';
import rightsWidgetsAlarmsListRemediationInstructionsFilters
  from '@/mixins/rights/widgets/alarms-list/remediation-instructions-filters';

import FieldTitle from './fields/common/title.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldLiveReporting from './fields/common/live-reporting.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldExportCsvSeparator from './fields/common/export-csv-separator.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldTextEditor from './fields/common/text-editor.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import FieldFastAckOutput from './fields/alarm/fast-ack-output.vue';
import FieldGridRangeSize from './fields/common/grid-range-size.vue';
import FieldRemediationInstructionsFilters from './fields/common/remediation-instructions-filters.vue';
import FieldOpenedResolvedFilter from './fields/alarm/opened-resolved-filter.vue';
import FieldInfoPopup from './fields/alarm/info-popup.vue';
import FieldEnabledLimit from './fields/common/enabled-limit.vue';

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
    FieldDefaultSortColumn,
    FieldColumns,
    FieldLiveReporting,
    FieldPeriodicRefresh,
    FieldDefaultElementsPerPage,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldTextEditor,
    FieldSwitcher,
    FieldFastAckOutput,
    FieldGridRangeSize,
    FieldRemediationInstructionsFilters,
    FieldInfoPopup,
    FieldEnabledLimit,
    FieldExportCsvSeparator,
  },
  mixins: [
    widgetSettingsMixin,
    sideBarSettingsWidgetAlarmMixin,
    rightsWidgetsAlarmsListFilters,
    rightsWidgetsAlarmsListRemediationInstructionsFilters,
  ],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: this.prepareAlarmWidgetSettings(cloneDeep(widget), true),
        widget_preferences: {
          itemsPerPage: PAGINATION_LIMIT,
        },
      },
    };
  },
  mounted() {
    const { widget_preferences: widgetPreference } = this.userPreference;

    this.settings.widget_preferences = {
      itemsPerPage: get(widgetPreference, 'itemsPerPage', PAGINATION_LIMIT),
    };
  },
  methods: {
    updateMainFilterUpdatedAt() {
      this.settings.widget.parameters.mainFilterUpdatedAt = Date.now();
    },

    prepareWidgetSettings() {
      const { widget } = this.settings;

      return this.prepareAlarmWidgetSettings(widget);
    },
  },
};
</script>
