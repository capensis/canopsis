<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />

    <field-periodic-refresh
      v-model="form.parameters.periodic_refresh"
      :live-watching="form.parameters.liveWatching"
      with-live-watching
      @update:live-watching="form.parameters.liveWatching = $event"
    />

    <widget-settings-group :title="t('settings.liveReporting.title')">
      <field-live-reporting v-model="form.parameters.liveReporting" />
    </widget-settings-group>

    <widget-settings-group :title="t('common.column', 2)">
      <field-default-sort-columns-with-template
        v-model="form.parameters.sort"
        :columns="sortablePreparedWidgetColumns"
        :template="form.parameters.sortTemplate"
        :templates="alarmSortColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        @update:template="updateSortTemplate"
      />
      <field-columns
        v-model="form.parameters.widgetColumns"
        :template="form.parameters.widgetColumnsTemplate"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.columnNames')"
        :type="ENTITIES_TYPES.alarm"
        with-template
        with-html
        with-color-indicator
        with-filter-on-click
        @update:template="updateWidgetColumnsTemplate"
      />
      <field-columns
        v-model="form.parameters.widgetGroupColumns"
        :template="form.parameters.widgetGroupColumnsTemplate"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.groupColumnNames')"
        :type="ENTITIES_TYPES.alarm"
        with-template
        with-html
        with-color-indicator
        with-filter-on-click
        @update:template="updateGroupColumnsTemplate"
      />
      <field-columns
        v-model="form.parameters.serviceDependenciesColumns"
        :template="form.parameters.serviceDependenciesColumnsTemplate"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.trackColumnNames')"
        :type="ENTITIES_TYPES.entity"
        with-color-indicator
        @update:template="updateServiceDependenciesColumnsTemplate"
      />
      <field-resize-column-behavior v-model="form.parameters.columns" draggable resizable />
      <field-root-cause-settings v-model="form.parameters" />
      <field-info-popup
        v-model="form.parameters.infoPopups"
        :columns="preparedWidgetColumns"
        :variables="alarmVariables"
      />
    </widget-settings-group>

    <widget-settings-group :title="t('common.filter', 2)">
      <field-filters
        v-model:model-value="form.parameters.mainFilter"
        v-model:filters="form.filters"
        :widget-id="widget._id"
        :widget-type="form.type"
        addable
        editable
        with-alarm
        with-entity
        with-pbehavior
      />
      <field-opened-resolved-filter v-model="form.parameters.opened" />
      <field-switcher
        v-model="form.parameters.isCorrelationEnabled"
        :title="t('common.correlation')"
      />
      <field-switcher
        v-model="form.parameters.clearFilterEnabled"
        :title="t('settings.clearFilterEnabled')"
      />
    </widget-settings-group>

    <widget-settings-group :title="t('common.view', 1)">
      <field-default-elements-per-page v-model="form.parameters.itemsPerPage" />
      <field-density v-model="form.parameters.dense" />
      <field-switcher v-model="form.parameters.sticky_header" :title="t('settings.stickyHeader')" />
      <field-switcher
        v-model="form.parameters.sticky_horizontal_scroll"
        :title="t('settings.stickyHorizontalScroll')"
      />
      <field-switcher
        v-model="form.parameters.isVirtualScrollEnabled"
        :title="t('settings.isVirtualScrollEnabled')"
      />
      <widget-settings-group :title="t('settings.kioskMode')">
        <field-switcher
          v-model="form.parameters.kiosk.hideActions"
          :title="t('settings.kiosk.hideActions')"
        />
        <field-switcher
          v-model="form.parameters.kiosk.hideMassSelection"
          :title="t('settings.kiosk.hideMassSelection')"
        />
        <field-switcher
          v-model="form.parameters.kiosk.hideToolbar"
          :title="t('settings.kiosk.hideToolbar')"
        />
      </widget-settings-group>
    </widget-settings-group>

    <widget-settings-group :title="t('common.actionsLabel')">
      <widget-settings-group :title="t('common.ack')">
        <field-switcher
          v-model="form.parameters.isAckNoteRequired"
          :title="t('settings.isAckNoteRequired')"
        />
        <field-switcher
          v-model="form.parameters.isMultiAckEnabled"
          :title="t('settings.isMultiAckEnabled')"
        />
        <field-fast-action-output
          v-model="form.parameters.fastAckOutput"
          :label="t('settings.fastAckOutput')"
        />
      </widget-settings-group>

      <widget-settings-group :title="t('common.cancel')">
        <field-fast-action-output
          v-model="form.parameters.fastCancelOutput"
          :label="t('settings.fastCancelOutput')"
        />
        <field-switcher
          v-model="form.parameters.isUncancelAlarmsCommentRequired"
          :title="t('settings.isUncancelAlarmsCommentRequired')"
        />
      </widget-settings-group>

      <field-comment-templates v-model="form.parameters.comment_templates" />

      <fast-pbehavior-form v-model="form.parameters.fast_pbehaviors" />

      <field-switcher
        v-model="form.parameters.isSnoozeNoteRequired"
        :title="t('settings.isSnoozeNoteRequired')"
      />
      <field-switcher
        v-model="form.parameters.isRemoveAlarmsFromMetaAlarmCommentRequired"
        :title="t('settings.isRemoveAlarmsFromMetaAlarmCommentRequired')"
      />
      <field-switcher
        v-model="form.parameters.isMultiDeclareTicketEnabled"
        :title="t('settings.isMultiDeclareTicketEnabled')"
      />
      <field-switcher
        v-model="form.parameters.isActionsAllowWithOkState"
        :title="t('settings.isActionsAllowWithOkState')"
      />
      <field-switcher
        v-model="form.parameters.keepSelectedAfterAction"
        :title="t('common.massActionsPanel.keepSelectedAfterAction')"
      />
      <field-text-editor-with-template
        :model-value="form.parameters.exportPdfTemplate"
        :template="form.parameters.exportPdfTemplateTemplate"
        :title="t('settings.exportPdfTemplate')"
        :widget-template-type="WIDGET_TEMPLATES_TYPES.alarmExportToPdf"
        :variables="exportPdfAlarmVariables"
        :templates="alarmExportToPdfWidgetTemplates"
        :default-value="defaultExportPdfTemplateValue"
        :dialog-props="{ maxWidth: 1070 }"
        addable
        removable
        @update:model-value="updateExportPdf"
        @update:template="updateExportPdfTemplate"
      />
      <widget-settings-group :title="t('settings.quickActions.title')">
        <field-quick-alarm-actions
          v-model="form.parameters.quickActions"
          :template="form.parameters.quickActionsTemplate"
          :templates="alarmQuickActionsWidgetTemplates"
          :templates-pending="widgetTemplatesPending"
          @update:template="updateQuickActionsTemplate"
        />
      </widget-settings-group>
      <widget-settings-group :title="t('settings.quickMassActions.title')">
        <field-switcher
          v-model="form.parameters.hideMassActions"
          :title="t('settings.quickMassActions.hideSwitcher')"
        />
        <field-quick-alarm-actions
          v-model="form.parameters.quickMassActions"
          :template="form.parameters.quickMassActionsTemplate"
          :templates="alarmMassQuickActionsWidgetTemplates"
          :templates-pending="widgetTemplatesPending"
          massive
          @update:template="updateQuickMassActionsTemplate"
        />
      </widget-settings-group>
    </widget-settings-group>

    <widget-settings-group :title="t('settings.expandPanel.title')">
      <widget-settings-group :title="t('alarm.tabs.moreInfos')">
        <field-text-editor-with-template
          :model-value="form.parameters.moreInfoTemplate"
          :template="form.parameters.moreInfoTemplateTemplate"
          :title="t('common.infos')"
          :widget-template-type="WIDGET_TEMPLATES_TYPES.alarmMoreInfos"
          :variables="alarmVariables"
          :templates="alarmMoreInfosWidgetTemplates"
          addable
          removable
          @update:model-value="updateMoreInfo"
          @update:template="updateMoreInfoTemplate"
        />
        <field-grid-range-size
          v-model="form.parameters.expandGridRangeSize"
          :title="t('common.width')"
          inline
        />
      </widget-settings-group>
      <field-availability-graph-settings v-model="form.parameters.availability" />
      <field-switcher
        v-model="form.parameters.isHtmlEnabledOnTimeLine"
        :title="t('settings.isHtmlEnabledOnTimeLine')"
      />
      <field-switcher
        v-model="form.parameters.openExpandPanelByDoubleClick"
        :title="t('settings.openExpandPanelByDoubleClick')"
      />
      <field-switcher
        v-model="form.parameters.pausePeriodicRefreshOnExpandPanel"
        :title="t('settings.pausePeriodicRefreshOnExpandPanel')"
      />
    </widget-settings-group>
    <field-export-csv
      v-model="form.parameters"
      :type="ENTITIES_TYPES.alarm"
      :templates="alarmColumnsWidgetTemplates"
      :templates-pending="widgetTemplatesPending"
      :variables="templateVars.column"
      datetime-format
      with-instructions
      with-template
      optional-infos-attributes
    />
    <charts-form v-model="form.parameters.charts" />
    <widget-colors-form v-model="form.parameters.colors" with-table />
  </widget-settings>
</template>

<script setup>
import { computed, onMounted, ref, toRef } from 'vue';

import {
  ALARM_EXPORT_PDF_FIELDS,
  ALARM_EXPORT_PDF_FIELDS_TO_ORIGINAL_FIELDS,
  ALARM_TEMPLATE_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_UNSORTABLE_FIELDS,
  ENTITIES_TYPES,
  ENTITY_TEMPLATE_FIELDS,
  MAX_LIMIT,
  SIDE_BARS,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { formToWidgetColumns } from '@/helpers/entities/widget/column/form';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '@/helpers/entities/widget/list';
import { prepareWidgetTemplates } from '@/helpers/entities/widget/template/prepare';

import { useI18n } from '@/hooks/i18n';
import { useInfos } from '@/hooks/store/modules/infos';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';
import { useWidgetTemplate } from '@/hooks/store/modules/widget-template';
import { useVarsPrepare } from '@/hooks/vars/shared';
import { useWidgetSettings } from '@/hooks/widget/settings';

import FastPbehaviorForm from '@/components/sidebars/alarm/form/fast-pbehavior-form.vue';
import FieldCommentTemplates from '@/components/sidebars/alarm/form/fields/comment-templates.vue';
import FieldFastActionOutput from '@/components/sidebars/alarm/form/fields/fast-action-output.vue';
import FieldInfoPopup from '@/components/sidebars/alarm/form/fields/info-popup.vue';
import FieldOpenedResolvedFilter from '@/components/sidebars/alarm/form/fields/opened-resolved-filter.vue';
import FieldQuickAlarmActions from '@/components/sidebars/alarm/form/fields/quick-alarm-actions.vue';
import FieldResizeColumnBehavior from '@/components/sidebars/alarm/form/fields/resize-column-behavior.vue';
import FieldAvailabilityGraphSettings from '@/components/sidebars/availability/form/fields/availability-graph-settings.vue';
import ChartsForm from '@/components/sidebars/chart/form/charts-form.vue';
import FieldExportCsv from '@/components/sidebars/form/export-csv.vue';
import FieldColumns from '@/components/sidebars/form/fields/columns.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/form/fields/default-elements-per-page.vue';
import FieldDefaultSortColumnsWithTemplate from '@/components/sidebars/form/fields/default-sort-columns-with-template.vue';
import FieldDensity from '@/components/sidebars/form/fields/density.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldGridRangeSize from '@/components/sidebars/form/fields/grid-range-size.vue';
import FieldLiveReporting from '@/components/sidebars/form/fields/live-reporting.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldRootCauseSettings from '@/components/sidebars/form/fields/root-cause-settings.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';
import FieldTextEditorWithTemplate from '@/components/sidebars/form/fields/text-editor-with-template.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import WidgetColorsForm from '@/components/sidebars/form/fields/widget-colors-form.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html?raw';

defineOptions({
  name: SIDE_BARS.alarmSettings,
});

const props = defineProps({
  sidebar: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const sidebar = toRef(props, 'sidebar');
const widgetTemplates = ref([]);
const widgetTemplatesPending = ref(false);
const { fetchInfos } = useInfos();
const { fetchWidgetsTemplateVarsWithoutStore } = useTemplateVars();
const { fetchWidgetTemplatesListWithoutStore } = useWidgetTemplate();
const { prepare: prepareTemplateVars } = useVarsPrepare('templateTesting.templateVars');

const { form, widget, hasChanges, submitting, close, submit } = useWidgetSettings({ sidebar });

const preparedWidgetTemplates = computed(() => prepareWidgetTemplates(widgetTemplates.value));

const alarmSortColumnsWidgetTemplates = computed(() => preparedWidgetTemplates.value.filter(
    ({ type }) => type === WIDGET_TEMPLATES_TYPES.alarmSortColumns,
  ),
);
const alarmColumnsWidgetTemplates = computed(() => preparedWidgetTemplates.value.filter(
    ({ type }) => type === WIDGET_TEMPLATES_TYPES.alarmColumns,
  ),
);
const entityColumnsWidgetTemplates = computed(() => preparedWidgetTemplates.value.filter(
    ({ type }) => type === WIDGET_TEMPLATES_TYPES.entityColumns,
  ),
);
const alarmMoreInfosWidgetTemplates = computed(() => preparedWidgetTemplates.value.filter(
    ({ type }) => type === WIDGET_TEMPLATES_TYPES.alarmMoreInfos,
  ),
);
const alarmExportToPdfWidgetTemplates = computed(() => preparedWidgetTemplates.value.filter(
    ({ type }) => type === WIDGET_TEMPLATES_TYPES.alarmExportToPdf,
  ),
);
const alarmQuickActionsWidgetTemplates = computed(() => preparedWidgetTemplates.value.filter(
    ({ type }) => type === WIDGET_TEMPLATES_TYPES.alarmQuickActions,
  ),
);
const alarmMassQuickActionsWidgetTemplates = computed(() => preparedWidgetTemplates.value.filter(
    ({ type }) => type === WIDGET_TEMPLATES_TYPES.alarmMassQuickActions,
  ),
);
const defaultExportPdfTemplateValue = ALARM_EXPORT_PDF_TEMPLATE;
const preparedWidgetColumns = computed(() => formToWidgetColumns(form.value.parameters.widgetColumns).map(column => ({
    ...column,
    text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
  })),
);
const sortablePreparedWidgetColumns = computed(() => preparedWidgetColumns.value.filter(
    column => getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
  ),
);
const alarmVariables = computed(() => [
    { value: ENTITY_TEMPLATE_FIELDS.id },
    { value: ALARM_TEMPLATE_FIELDS.ack },
    { value: ALARM_TEMPLATE_FIELDS.state },
    { value: ALARM_TEMPLATE_FIELDS.status },
    { value: ALARM_TEMPLATE_FIELDS.ticket },
    { value: ALARM_TEMPLATE_FIELDS.component },
    { value: ALARM_TEMPLATE_FIELDS.connector },
    { value: ALARM_TEMPLATE_FIELDS.connectorName },
    { value: ALARM_TEMPLATE_FIELDS.resource },
    { value: ALARM_TEMPLATE_FIELDS.creationDate },
    { value: ALARM_TEMPLATE_FIELDS.displayName },
    { value: ALARM_TEMPLATE_FIELDS.output },
    { value: ALARM_TEMPLATE_FIELDS.lastUpdateDate },
    { value: ALARM_TEMPLATE_FIELDS.lastEventDate },
    { value: ALARM_TEMPLATE_FIELDS.pbehaviorInfo },
    { value: ALARM_TEMPLATE_FIELDS.duration },
    { value: ALARM_TEMPLATE_FIELDS.eventsCount },
  ].map(variable => ({
    value: `{{ ${variable.value} }}`,
    text: t(ALARM_FIELDS_TO_LABELS_KEYS[variable.value.replace('alarm.', '')], 2),
  })),
);
const exportPdfAlarmVariables = computed(() => [
    {
      value: ALARM_EXPORT_PDF_FIELDS.currentDate,
      text: t('common.currentDate'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.displayName,
      text: t('alarm.fields.displayName'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.state,
      text: t('common.state'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.initialState,
      text: t('alarm.fields.initialState'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.maxState,
      text: t('alarm.fields.maxState'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.status,
      text: t('common.status'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.connector,
      text: t('common.connector'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.connectorName,
      text: t('common.connectorName'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.component,
      text: t('common.component'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.resource,
      text: t('common.resource'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.creationDate,
      text: t('common.created'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.duration,
      text: t('common.duration'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.infos,
      text: t('common.infos'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.pbehaviorInfo,
      text: t('pbehavior.pbehaviorInfo'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.initialOutput,
      text: t('alarm.fields.initialOutput'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.output,
      text: t('common.output'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.lastEventDate,
      text: t('common.lastEventDate'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.lastUpdateDate,
      text: t('common.updated'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.acknowledgeDate,
      text: t('alarm.fields.ackAt'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.resolved,
      text: t('alarm.fields.resolved'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.eventsCount,
      text: t('alarm.fields.eventsCount'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.ticket,
      text: t('common.ticket', 1),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.comments,
      text: t('common.comment', 2),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.activationDate,
      text: t('alarm.fields.activationDate'),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.tags,
      text: t('common.tag', 2),
    },
    {
      value: ALARM_EXPORT_PDF_FIELDS.links,
      text: t('common.link', 2),
    },
  ].map(variable => ({
    ...variable,
    value: `{{ ${ALARM_EXPORT_PDF_FIELDS_TO_ORIGINAL_FIELDS[variable.value]} }}`,
  })),
);
const templateVars = ref({});

const updateSortTemplate = (template, columns) => {
  form.value.parameters.sortTemplate = template;
  form.value.parameters.sort = columns;
};
const updateWidgetColumnsTemplate = (template, columns) => {
  form.value.parameters.widgetColumnsTemplate = template;
  form.value.parameters.widgetColumns = columns;
};
const updateGroupColumnsTemplate = (template, columns) => {
  form.value.parameters.widgetGroupColumnsTemplate = template;
  form.value.parameters.widgetGroupColumns = columns;
};
const updateServiceDependenciesColumnsTemplate = (template, columns) => {
  form.value.parameters.serviceDependenciesColumnsTemplate = template;
  form.value.parameters.serviceDependenciesColumns = columns;
};
const updateMoreInfo = (value) => {
  form.value.parameters.moreInfoTemplate = value;
};
const updateMoreInfoTemplate = (value) => {
  form.value.parameters.moreInfoTemplateTemplate = value;
};
const updateExportPdf = (value) => {
  form.value.parameters.exportPdfTemplate = value;
};
const updateExportPdfTemplate = (value) => {
  form.value.parameters.exportPdfTemplateTemplate = value;
};
const updateQuickActionsTemplate = (template, actions) => {
  form.value.parameters.quickActionsTemplate = template;
  form.value.parameters.quickActions = actions;
};
const updateQuickMassActionsTemplate = (template, actions) => {
  form.value.parameters.quickMassActionsTemplate = template;
  form.value.parameters.quickMassActions = actions;
};

const fetchWidgetTemplates = async () => {
  try {
    widgetTemplatesPending.value = true;
    const { data } = await fetchWidgetTemplatesListWithoutStore({ params: { limit: MAX_LIMIT } });

    widgetTemplates.value = data;
  } catch (err) {
    console.warn(err);
  } finally {
    widgetTemplatesPending.value = false;
  }
};
const fetchTemplateVars = async () => {
  try {
    const data = await fetchWidgetsTemplateVarsWithoutStore();

    templateVars.value = prepareTemplateVars(data);
  } catch (err) {
    console.error(err);
  }
};

onMounted(() => {
  fetchInfos();
  fetchTemplateVars();
  fetchWidgetTemplates();
});

defineExpose({ close });
</script>
