<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="form.title", :title="$t('common.title')")
      v-divider
      field-periodic-refresh(v-model="form.parameters.periodic_refresh")
      v-divider
      v-list-group
        template(#activator="")
          v-list-tile {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(
            v-model="form.parameters.sort",
            :columns="form.parameters.widgetColumns",
            :columns-label="$t('settings.columnName')"
          )
          v-divider
          field-columns(
            v-model="form.parameters.widgetColumns",
            :label="$t('settings.columnNames')",
            with-template,
            with-html,
            with-color-indicator
          )
          v-divider
          field-columns(
            v-model="form.parameters.widgetGroupColumns",
            :label="$t('settings.groupColumnNames')",
            with-html,
            with-color-indicator
          )
          v-divider
          field-columns(
            v-model="form.parameters.serviceDependenciesColumns",
            :label="$t('settings.trackColumnNames')",
            with-color-indicator
          )
          v-divider
          field-default-elements-per-page(v-model="form.parameters.itemsPerPage")
          v-divider
          field-density(v-model="form.parameters.dense")
          v-divider
          field-opened-resolved-filter(v-model="form.parameters.opened")
          v-divider
          template(v-if="hasAccessToListFilters")
            field-filters(
              v-model="form.parameters.mainFilter",
              :filters.sync="form.filters",
              :widget-id="widget._id",
              :addable="hasAccessToAddFilter",
              :editable="hasAccessToEditFilter",
              with-alarm,
              with-entity,
              with-pbehavior,
              @input="updateMainFilterUpdatedAt"
            )
            v-divider
          template(v-if="hasAccessToListRemediationInstructionsFilters")
            field-remediation-instructions-filters(
              v-model="form.parameters.remediationInstructionsFilters",
              :addable="hasAccessToAddRemediationInstructionsFilter",
              :editable="hasAccessToEditRemediationInstructionsFilter"
            )
            v-divider
          field-live-reporting(v-model="form.parameters.liveReporting")
          v-divider
          field-info-popup(
            v-model="form.parameters.infoPopups",
            :columns="form.parameters.widgetColumns"
          )
          v-divider
          field-text-editor(
            v-model="form.parameters.moreInfoTemplate",
            :title="$t('settings.moreInfosModal')"
          )
          v-divider
          field-grid-range-size(
            v-model="form.parameters.expandGridRangeSize",
            :title="$t('settings.expandGridRangeSize')"
          )
          v-divider
          field-switcher(
            v-model="form.parameters.isHtmlEnabledOnTimeLine",
            :title="$t('settings.isHtmlEnabledOnTimeLine')"
          )
          v-divider
          v-list-group
            template(#activator="")
              v-list-tile Ack
            v-list.grey.lighten-4.px-2.py-0(expand)
              field-switcher(
                v-model="form.parameters.isAckNoteRequired",
                :title="$t('settings.isAckNoteRequired')"
              )
              v-divider
              field-switcher(
                v-model="form.parameters.isMultiAckEnabled",
                :title="$t('settings.isMultiAckEnabled')"
              )
              v-divider
              field-fast-ack-output(v-model="form.parameters.fastAckOutput")
          v-divider
          field-switcher(
            v-model="form.parameters.isSnoozeNoteRequired",
            :title="$t('settings.isSnoozeNoteRequired')"
          )
          v-divider
          field-switcher(
            v-model="form.parameters.isMultiDeclareTicketEnabled",
            :title="$t('settings.isMultiDeclareTicketEnabled')"
          )
          v-divider
          field-enabled-limit(
            v-model="form.parameters.linksCategoriesAsList",
            :title="$t('settings.linksCategoriesAsList')",
            :label="$t('settings.linksCategoriesLimit')"
          )
          v-divider
          export-csv-form(v-model="form.parameters", datetime-format)
          v-divider
          field-switcher(
            v-model="form.parameters.sticky_header",
            :title="$t('settings.stickyHeader')"
          )
      v-divider
    v-btn.primary(
      :loading="submitting",
      :disabled="submitting",
      @click="submit"
    ) {{ $t('common.save') }}
</template>

<script>
import { SIDE_BARS } from '@/constants';

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
import FieldDensity from './fields/common/density.vue';
import ExportCsvForm from './forms/export-csv.vue';

/**
 * Component to regroup the alarms list settings fields
 */
export default {
  name: SIDE_BARS.alarmSettings,
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
    FieldDensity,
    ExportCsvForm,
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsAlarmsListFilters,
    permissionsWidgetsAlarmsListRemediationInstructionsFilters,
  ],
};
</script>
