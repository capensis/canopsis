<template>
  <widget-settings-group :title="$t(`settings.titles.${SIDE_BARS.alarmSettings}`)">
    <field-default-sort-columns-with-template
      :model-value="form.sort"
      :columns="sortablePreparedWidgetColumns"
      :template="form.sortTemplate"
      :templates="alarmSortColumnsWidgetTemplates"
      :templates-pending="templatesPending"
      @update:model-value="updateField('sort', $event)"
      @update:template="updateSortTemplate"
    />
    <field-columns
      :model-value="form.widgetColumns"
      :template="form.widgetColumnsTemplate"
      :templates="alarmColumnsWidgetTemplates"
      :templates-pending="templatesPending"
      :label="$t('settings.columnNames')"
      :type="ENTITIES_TYPES.alarm"
      with-html
      with-template
      with-color-indicator
      @update:model-value="updateField('widgetColumns', $event)"
      @update:template="updateColumnsTemplate"
    />
    <field-resize-column-behavior
      :model-value="form.columns"
      resizable
      @update:model-value="updateField('columns', $event)"
    />
    <field-default-elements-per-page
      :model-value="form.itemsPerPage"
      @update:model-value="updateField('itemsPerPage', $event)"
    />
    <field-info-popup
      :model-value="form.infoPopups"
      :columns="preparedWidgetColumns"
      @update:model-value="updateField('infoPopups', $event)"
    />
    <field-text-editor-with-template
      :model-value="form.moreInfoTemplate"
      :template="form.moreInfoTemplateTemplate"
      :title="$t('settings.moreInfosModal')"
      :widget-template-type="WIDGET_TEMPLATES_TYPES.alarmMoreInfos"
      :variables="alarmVariables"
      :templates="alarmMoreInfosWidgetTemplates"
      @update:model-value="updateMoreInfo"
      @update:template="updateMoreInfoTemplate"
    />
    <field-text-editor-with-template
      :model-value="form.exportPdfTemplate"
      :template="form.exportPdfTemplateTemplate"
      :title="$t('settings.exportPdfTemplate')"
      :widget-template-type="WIDGET_TEMPLATES_TYPES.alarmExportToPdf"
      :variables="exportPdfAlarmVariables"
      :default-value="defaultExportPdfTemplateValue"
      :dialog-props="{ maxWidth: 1070 }"
      :templates="alarmExportToPdfTemplates"
      addable
      removable
      @update:model-value="updateExportPdf"
      @update:template="updateExportPdfTemplate"
    />
    <field-root-cause-settings
      :model-value="form"
      @update:model-value="updateModel"
    />
  </widget-settings-group>
</template>

<script>
import { filter } from 'lodash';

import {
  ALARM_EXPORT_PDF_FIELDS,
  ALARM_EXPORT_PDF_FIELDS_TO_ORIGINAL_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_TEMPLATE_FIELDS,
  ALARM_UNSORTABLE_FIELDS,
  ENTITIES_TYPES,
  ENTITY_TEMPLATE_FIELDS,
  SIDE_BARS,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { formToWidgetColumns } from '@/helpers/entities/widget/column/form';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '@/helpers/entities/widget/list';

import { formMixin } from '@/mixins/form';

import FieldInfoPopup from '@/components/sidebars/alarm/form/fields/info-popup.vue';
import FieldResizeColumnBehavior from '@/components/sidebars/alarm/form/fields/resize-column-behavior.vue';
import FieldColumns from '@/components/sidebars/form/fields/columns.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/form/fields/default-elements-per-page.vue';
import FieldDefaultSortColumnsWithTemplate from '@/components/sidebars/form/fields/default-sort-columns-with-template.vue';
import FieldRootCauseSettings from '@/components/sidebars/form/fields/root-cause-settings.vue';
import FieldTextEditorWithTemplate from '@/components/sidebars/form/fields/text-editor-with-template.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html?raw';

export default {
  components: {
    FieldRootCauseSettings,
    WidgetSettingsGroup,
    FieldDefaultSortColumnsWithTemplate,
    FieldColumns,
    FieldResizeColumnBehavior,
    FieldInfoPopup,
    FieldTextEditorWithTemplate,
    FieldDefaultElementsPerPage,
  },
  mixins: [formMixin],

  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    templates: {
      type: Array,
      default: () => [],
    },
    templatesPending: {
      type: Boolean,
      default: false,
    },
    datetimeFormat: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:value'],

  data() {
    return {
      ENTITIES_TYPES,
      SIDE_BARS,
      WIDGET_TEMPLATES_TYPES,
    };
  },

  computed: {
    form() {
      return this.value;
    },

    defaultExportPdfTemplateValue() {
      return ALARM_EXPORT_PDF_TEMPLATE;
    },

    alarmSortColumnsWidgetTemplates() {
      return filter(this.templates, { type: WIDGET_TEMPLATES_TYPES.alarmSortColumns });
    },

    alarmColumnsWidgetTemplates() {
      return filter(this.templates, { type: WIDGET_TEMPLATES_TYPES.alarmColumns });
    },

    alarmMoreInfosWidgetTemplates() {
      return filter(this.templates, { type: WIDGET_TEMPLATES_TYPES.alarmMoreInfos });
    },

    alarmExportToPdfTemplates() {
      return filter(this.templates, { type: WIDGET_TEMPLATES_TYPES.alarmExportToPdf });
    },

    preparedWidgetColumns() {
      return formToWidgetColumns(this.form.widgetColumns).map(column => ({
        ...column,
        text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
      }));
    },

    sortablePreparedWidgetColumns() {
      return this.preparedWidgetColumns.filter(column => getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
      );
    },

    alarmVariables() {
      return [
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
        text: this.$t(ALARM_FIELDS_TO_LABELS_KEYS[variable.value.replace('alarm.', '')], 2),
      }));
    },

    exportPdfAlarmVariables() {
      return [
        {
          value: ALARM_EXPORT_PDF_FIELDS.currentDate,
          text: this.$t('common.currentDate'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.displayName,
          text: this.$t('alarm.fields.displayName'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.state,
          text: this.$t('common.state'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.initialState,
          text: this.$t('alarm.fields.initialState'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.maxState,
          text: this.$t('alarm.fields.maxState'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.status,
          text: this.$t('common.status'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.connector,
          text: this.$t('common.connector'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.connectorName,
          text: this.$t('common.connectorName'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.component,
          text: this.$t('common.component'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.resource,
          text: this.$t('common.resource'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.creationDate,
          text: this.$t('common.created'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.duration,
          text: this.$t('common.duration'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.infos,
          text: this.$t('common.infos'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.pbehaviorInfo,
          text: this.$t('pbehavior.pbehaviorInfo'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.initialOutput,
          text: this.$t('alarm.fields.initialOutput'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.output,
          text: this.$t('common.output'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.lastEventDate,
          text: this.$t('common.lastEventDate'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.lastUpdateDate,
          text: this.$t('common.updated'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.acknowledgeDate,
          text: this.$t('alarm.fields.ackAt'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.resolved,
          text: this.$t('alarm.fields.resolved'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.eventsCount,
          text: this.$t('alarm.fields.eventsCount'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.ticket,
          text: this.$t('common.ticket', 1),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.comments,
          text: this.$t('common.comment', 2),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.activationDate,
          text: this.$t('alarm.fields.activationDate'),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.tags,
          text: this.$t('common.tag', 2),
        },
        {
          value: ALARM_EXPORT_PDF_FIELDS.links,
          text: this.$t('common.link', 2),
        },
      ].map(variable => ({
        ...variable,
        value: `{{ ${ALARM_EXPORT_PDF_FIELDS_TO_ORIGINAL_FIELDS[variable.value]} }}`,
      }));
    },
  },

  methods: {
    updateTemplate(field, template, value) {
      this.updateModel({
        ...this.form,
        [field]: value,
        [`${field}Template`]: template,
      });
    },

    updateSortTemplate(template = '', sort = []) {
      this.updateTemplate('sort', template, sort);
    },

    updateColumnsTemplate(template, columns) {
      this.updateTemplate('widgetColumns', template, columns);
    },

    updateMoreInfo(content) {
      this.updateModel({
        ...this.form,
        moreInfoTemplate: content,
      });
    },

    updateMoreInfoTemplate(template, content = this.form.moreInfoTemplate) {
      this.updateTemplate('moreInfoTemplate', template, content);
    },

    updateExportPdf(content) {
      this.updateModel({
        ...this.form,
        exportPdfTemplate: content,
      });
    },

    updateExportPdfTemplate(template, content = this.form.exportPdfTemplate) {
      this.updateTemplate('exportPdfTemplate', template, content);
    },
  },
};
</script>
