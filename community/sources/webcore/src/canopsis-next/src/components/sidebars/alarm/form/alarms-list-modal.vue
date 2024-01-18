<template>
  <widget-settings-group :title="$t(`settings.titles.${$constants.SIDE_BARS.alarmSettings}`)">
    <field-default-sort-column
      v-field="form.sort"
      :columns="sortablePreparedWidgetColumns"
      :columns-label="$t('settings.columnName')"
    />
    <field-columns
      v-field="form.widgetColumns"
      :template="form.widgetColumnsTemplate"
      :templates="alarmColumnsWidgetTemplates"
      :templates-pending="templatesPending"
      :label="$t('settings.columnNames')"
      :type="$constants.ENTITIES_TYPES.alarm"
      with-html
      with-template
      with-color-indicator
      @update:template="updateColumnsTemplate"
    />
    <field-default-elements-per-page v-field="form.itemsPerPage" />
    <field-info-popup
      v-field="form.infoPopups"
      :columns="form.widgetColumns"
    />
    <field-text-editor-with-template
      :value="form.moreInfoTemplate"
      :template="form.moreInfoTemplateTemplate"
      :title="$t('settings.moreInfosModal')"
      :variables="alarmVariables"
      :templates="alarmMoreInfosWidgetTemplates"
      @input="updateMoreInfo"
    />
    <field-text-editor-with-template
      :value="form.exportPdfTemplate"
      :template="form.exportPdfTemplateTemplate"
      :title="$t('settings.exportPdfTemplate')"
      :variables="exportPdfAlarmVariables"
      :default-value="defaultExportPdfTemplateValue"
      :dialog-props="{ maxWidth: 1070 }"
      :templates="alarmExportToPdfTemplates"
      addable
      removable
      @input="updateExportPdf"
    />
    <field-switcher
      v-field="form.showRootCauseByStateClick"
      :title="$t('settings.showRootCauseByStateClick')"
    />
  </widget-settings-group>
</template>

<script>
import { filter } from 'lodash';

import { ALARM_FIELDS_TO_LABELS_KEYS, ALARM_UNSORTABLE_FIELDS, WIDGET_TEMPLATES_TYPES } from '@/constants';

import { formToWidgetColumns } from '@/helpers/entities/widget/column/form';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '@/helpers/entities/widget/list';

import { formBaseMixin } from '@/mixins/form';
import { alarmVariablesMixin } from '@/mixins/widget/variables/alarm';

import FieldDefaultSortColumn from '@/components/sidebars/form/fields/default-sort-column.vue';
import FieldColumns from '@/components/sidebars/form/fields/columns.vue';
import FieldInfoPopup from '@/components/sidebars/alarm/form/fields/info-popup.vue';
import FieldTextEditorWithTemplate from '@/components/sidebars/form/fields/text-editor-with-template.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/form/fields/default-elements-per-page.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html';

export default {
  components: {
    FieldSwitcher,
    WidgetSettingsGroup,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldInfoPopup,
    FieldTextEditorWithTemplate,
    FieldDefaultElementsPerPage,
  },
  mixins: [
    formBaseMixin,
    alarmVariablesMixin,
  ],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
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
  computed: {
    defaultExportPdfTemplateValue() {
      return ALARM_EXPORT_PDF_TEMPLATE;
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
      return this.preparedWidgetColumns.filter(column => getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS));
    },
  },
  methods: {
    updateColumnsTemplate(template, columns) {
      this.updateModel({
        ...this.form,

        widgetColumnsTemplate: template,
        widgetColumns: columns,
      });
    },

    updateMoreInfo(content, template) {
      this.updateModel({
        ...this.form,

        moreInfoTemplate: content,
        moreInfoTemplateTemplate: template,
      });
    },

    updateExportPdf(content, template) {
      this.updateModel({
        ...this.form,

        exportPdfTemplate: content,
        exportPdfTemplateTemplate: template,
      });
    },
  },
};
</script>
