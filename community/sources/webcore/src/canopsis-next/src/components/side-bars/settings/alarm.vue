<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-periodic-refresh(v-model="settings.widget.parameters.periodic_refresh")
      v-divider
      v-list-group
        template(#activator="")
          v-list-tile {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(
            v-model="settings.widget.parameters.sort",
            :columns="settings.widget.parameters.widgetColumns",
            :columns-label="$t('settings.columnName')"
          )
          v-divider
          field-columns(
            v-model="settings.widget.parameters.widgetColumns",
            :label="$t('settings.columnNames')",
            with-html,
            with-color-indicator
          )
          v-divider
          field-columns(
            v-model="settings.widget.parameters.widgetGroupColumns",
            :label="$t('settings.groupColumnNames')",
            with-html,
            with-color-indicator
          )
          v-divider
          field-columns(
            v-model="settings.widget.parameters.serviceDependenciesColumns",
            :label="$t('settings.trackColumnNames')",
            with-color-indicator
          )
          v-divider
          field-default-elements-per-page(v-model="settings.userPreferenceContent.itemsPerPage")
          v-divider
          field-opened-resolved-filter(v-model="settings.widget.parameters.opened")
          v-divider
          template(v-if="hasAccessToListFilters")
            field-filters(
              v-model="settings.widget.parameters.mainFilter",
              :entities-type="$constants.ENTITIES_TYPES.alarm",
              :filters.sync="settings.widget.parameters.viewFilters",
              :condition.sync="settings.widget.parameters.mainFilterCondition",
              :addable="hasAccessToAddFilter",
              :editable="hasAccessToEditFilter",
              @input="updateMainFilterUpdatedAt"
            )
            v-divider
          template(v-if="hasAccessToListRemediationInstructionsFilters")
            field-remediation-instructions-filters(
              v-model="settings.widget.parameters.remediationInstructionsFilters",
              :addable="hasAccessToAddRemediationInstructionsFilter",
              :editable="hasAccessToEditRemediationInstructionsFilter"
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
            v-model="settings.widget.parameters.isHtmlEnabledOnTimeLine",
            :title="$t('settings.isHtmlEnabledOnTimeLine')"
          )
          v-divider
          v-list-group
            v-list-tile(slot="activator") Ack
            v-list.grey.lighten-4.px-2.py-0(expand)
            field-switcher(
              v-model="settings.widget.parameters.isAckNoteRequired",
              :title="$t('settings.isAckNoteRequired')"
            )
            v-divider
            field-switcher(
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
          export-csv-form(v-model="settings.widget.parameters", datetime-format)
          v-divider
          field-switcher(
            v-model="settings.widget.parameters.sticky_header",
            :title="$t('settings.stickyHeader')"
          )
      v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { get } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { SIDE_BARS } from '@/constants';

import { alarmListWidgetToForm, formToAlarmListWidget } from '@/helpers/forms/widgets/alarm';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { permissionsWidgetsAlarmsListFilters } from '@/mixins/permissions/widgets/alarms-list/filters';
import { permissionsWidgetsAlarmsListRemediationInstructionsFilters }
  from '@/mixins/permissions/widgets/alarms-list/remediation-instructions-filters';

import FieldTitle from './fields/common/title.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldLiveReporting from './fields/common/live-reporting.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
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
import ExportCsvForm from './forms/export-csv.vue';

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
    ExportCsvForm,
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsAlarmsListFilters,
    permissionsWidgetsAlarmsListRemediationInstructionsFilters,
  ],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: alarmListWidgetToForm(widget),
        userPreferenceContent: { itemsPerPage: PAGINATION_LIMIT },
      },
    };
  },
  mounted() {
    const { content } = this.userPreference;

    this.settings.userPreferenceContent = {
      itemsPerPage: get(content, 'itemsPerPage', PAGINATION_LIMIT),
    };
  },
  methods: {
    updateMainFilterUpdatedAt() {
      this.settings.widget.parameters.mainFilterUpdatedAt = Date.now();
    },

    prepareWidgetSettings() {
      const { widget } = this.settings;

      return formToAlarmListWidget(widget);
    },
  },
};
</script>
