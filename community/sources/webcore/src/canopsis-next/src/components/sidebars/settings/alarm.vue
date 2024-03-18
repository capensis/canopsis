<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    field-periodic-refresh(v-model="form.parameters.periodic_refresh")
    v-divider
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-default-sort-column(
        v-model="form.parameters.sort",
        :columns="sortablePreparedWidgetColumns",
        :columns-label="$t('settings.columnName')"
      )
      v-divider
      field-columns(
        v-model="form.parameters.widgetColumns",
        :template="form.parameters.widgetColumnsTemplate",
        :templates="alarmColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        :label="$t('settings.columnNames')",
        :type="$constants.ENTITIES_TYPES.alarm",
        with-template,
        with-html,
        with-color-indicator,
        @update:template="updateWidgetColumnsTemplate"
      )
      v-divider
      field-columns(
        v-model="form.parameters.widgetGroupColumns",
        :template="form.parameters.widgetGroupColumnsTemplate",
        :templates="alarmColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        :label="$t('settings.groupColumnNames')",
        :type="$constants.ENTITIES_TYPES.alarm",
        with-html,
        with-color-indicator,
        @update:template="updateWidgetGroupColumnsTemplate"
      )
      v-divider
      field-columns(
        v-model="form.parameters.serviceDependenciesColumns",
        :template="form.parameters.serviceDependenciesColumnsTemplate",
        :templates="entityColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        :label="$t('settings.trackColumnNames')",
        :type="$constants.ENTITIES_TYPES.entity",
        with-color-indicator,
        @update:template="updateServiceDependenciesColumnsTemplate"
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
          with-pbehavior
        )
        v-divider
      template(v-if="hasAccessToListRemediationInstructionsFilters")
        field-remediation-instructions-filters(
          v-model="form.parameters.remediationInstructionsFilters",
          :addable="hasAccessToAddRemediationInstructionsFilter",
          :editable="hasAccessToEditRemediationInstructionsFilter"
        )
        v-divider
      field-switcher(
        v-model="form.parameters.clearFilterDisabled",
        :title="$t('settings.clearFilterDisabled')"
      )
      v-divider
      field-live-reporting(v-model="form.parameters.liveReporting")
      v-divider
      field-info-popup(
        v-model="form.parameters.infoPopups",
        :columns="preparedWidgetColumns"
      )
      v-divider
      field-text-editor-with-template(
        :value="form.parameters.moreInfoTemplate",
        :template="form.parameters.moreInfoTemplateTemplate",
        :title="$t('settings.moreInfosModal')",
        :variables="alarmVariables",
        :templates="alarmMoreInfosWidgetTemplates",
        addable,
        removable,
        @input="updateMoreInfo"
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
      widget-settings-group(:title="$t('common.ack')")
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
        field-fast-action-output(
          v-model="form.parameters.fastAckOutput",
          :label="$t('settings.fastAckOutput')"
        )
      v-divider
      widget-settings-group(:title="$t('common.cancel')")
        field-fast-action-output(
          v-model="form.parameters.fastCancelOutput",
          :label="$t('settings.fastCancelOutput')"
        )
      v-divider
      field-switcher(
        v-model="form.parameters.isSnoozeNoteRequired",
        :title="$t('settings.isSnoozeNoteRequired')"
      )
      v-divider
      field-switcher(
        v-model="form.parameters.isUncancelAlarmsCommentRequired",
        :title="$t('settings.isUncancelAlarmsCommentRequired')"
      )
      v-divider
      field-switcher(
        v-model="form.parameters.isMultiDeclareTicketEnabled",
        :title="$t('settings.isMultiDeclareTicketEnabled')"
      )
      v-divider
      export-csv-form(
        v-model="form.parameters",
        :type="$constants.ENTITIES_TYPES.alarm",
        :templates="alarmColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        datetime-format,
        with-instructions
      )
      v-divider
      field-switcher(
        v-model="form.parameters.sticky_header",
        :title="$t('settings.stickyHeader')"
      )
      v-divider
      widget-settings-group(:title="$t('settings.kioskMode')")
        field-switcher(
          v-model="form.parameters.kiosk.hideActions",
          :title="$t('settings.kiosk.hideActions')"
        )
        v-divider
        field-switcher(
          v-model="form.parameters.kiosk.hideMassSelection",
          :title="$t('settings.kiosk.hideMassSelection')"
        )
        v-divider
        field-switcher(
          v-model="form.parameters.kiosk.hideToolbar",
          :title="$t('settings.kiosk.hideToolbar')"
        )
    v-divider
</template>

<script>
import { SIDE_BARS, ALARM_UNSORTABLE_FIELDS, ALARM_FIELDS_TO_LABELS_KEYS } from '@/constants';

import { formToWidgetColumns } from '@/helpers/forms/shared/widget-column';
import { getColumnLabel, getSortable } from '@/helpers/widgets';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { alarmVariablesMixin } from '@/mixins/widget/variables';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';
import { permissionsWidgetsAlarmsListFilters } from '@/mixins/permissions/widgets/alarms-list/filters';
import {
  permissionsWidgetsAlarmsListRemediationInstructionsFilters,
} from '@/mixins/permissions/widgets/alarms-list/remediation-instructions-filters';

import FieldTitle from './fields/common/title.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldLiveReporting from './fields/common/live-reporting.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldTextEditorWithTemplate from './fields/common/text-editor-with-template.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import FieldFastActionOutput from './fields/alarm/fast-action-output.vue';
import FieldGridRangeSize from './fields/common/grid-range-size.vue';
import FieldRemediationInstructionsFilters from './fields/common/remediation-instructions-filters.vue';
import FieldOpenedResolvedFilter from './fields/alarm/opened-resolved-filter.vue';
import FieldInfoPopup from './fields/alarm/info-popup.vue';
import FieldDensity from './fields/common/density.vue';
import FieldNumber from './fields/common/number.vue';
import ExportCsvForm from './forms/export-csv.vue';
import WidgetSettings from './partials/widget-settings.vue';
import WidgetSettingsGroup from './partials/widget-settings-group.vue';

/**
 * Component to regroup the alarms list settings fields
 */
export default {
  name: SIDE_BARS.alarmSettings,
  components: {
    WidgetSettingsGroup,
    WidgetSettings,
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldLiveReporting,
    FieldPeriodicRefresh,
    FieldDefaultElementsPerPage,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldTextEditorWithTemplate,
    FieldSwitcher,
    FieldFastActionOutput,
    FieldGridRangeSize,
    FieldRemediationInstructionsFilters,
    FieldInfoPopup,
    FieldDensity,
    FieldNumber,
    ExportCsvForm,
  },
  mixins: [
    widgetSettingsMixin,
    entitiesInfosMixin,
    alarmVariablesMixin,
    widgetTemplatesMixin,
    permissionsWidgetsAlarmsListFilters,
    permissionsWidgetsAlarmsListRemediationInstructionsFilters,
  ],
  computed: {
    preparedWidgetColumns() {
      return formToWidgetColumns(this.form.parameters.widgetColumns).map(column => ({
        ...column,

        text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
      }));
    },

    sortablePreparedWidgetColumns() {
      return this.preparedWidgetColumns.filter(column => getSortable(column, ALARM_UNSORTABLE_FIELDS));
    },
  },
  mounted() {
    this.fetchInfos();
  },
  methods: {
    updateWidgetColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'widgetColumnsTemplate', template);
      this.$set(this.form.parameters, 'widgetColumns', columns);
    },

    updateWidgetGroupColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'widgetGroupColumnsTemplate', template);
      this.$set(this.form.parameters, 'widgetGroupColumns', columns);
    },

    updateServiceDependenciesColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'serviceDependenciesColumnsTemplate', template);
      this.$set(this.form.parameters, 'serviceDependenciesColumns', columns);
    },

    updateMoreInfo(content, template) {
      this.$set(this.form.parameters, 'moreInfoTemplate', content);

      if (template && template !== this.form.parameters.moreInfoTemplateTemplate) {
        this.$set(this.form.parameters, 'moreInfoTemplateTemplate', template);
      }
    },
  },
};
</script>
