<template lang="pug">
  widget-settings-group(:title="$t(`settings.titles.${$constants.SIDE_BARS.alarmSettings}`)")
    field-default-sort-column(
      v-model="form.sort",
      :columns="sortablePreparedWidgetColumns",
      :columns-label="$t('settings.columnName')"
    )
    v-divider
    field-columns(
      v-field="form.widgetColumns",
      :template="form.widgetColumnsTemplate",
      :templates="alarmColumnsWidgetTemplates",
      :templates-pending="templatesPending",
      :label="$t('settings.columnNames')",
      :type="$constants.ENTITIES_TYPES.alarm",
      with-html,
      with-template,
      with-color-indicator,
      @update:template="updateColumnsTemplate"
    )
    v-divider
    field-default-elements-per-page(v-field="form.itemsPerPage")
    v-divider
    field-info-popup(
      v-field="form.infoPopups",
      :columns="form.widgetColumns"
    )
    v-divider
    field-text-editor-with-template(
      :value="form.moreInfoTemplate",
      :template="form.moreInfoTemplateTemplate",
      :title="$t('settings.moreInfosModal')",
      :variables="alarmVariables",
      :templates="alarmMoreInfosWidgetTemplates",
      @input="updateMoreInfo"
    )
</template>

<script>
import { filter } from 'lodash';

import { ALARM_FIELDS_TO_LABELS_KEYS, ALARM_UNSORTABLE_FIELDS, WIDGET_TEMPLATES_TYPES } from '@/constants';

import { formToWidgetColumns } from '@/helpers/forms/shared/widget-column';
import { getColumnLabel, getSortable } from '@/helpers/widgets';

import { formBaseMixin } from '@/mixins/form';
import { alarmVariablesMixin } from '@/mixins/widget/variables/alarm';

import FieldDefaultSortColumn from '@/components/sidebars/settings/fields/common/default-sort-column.vue';
import FieldColumns from '@/components/sidebars/settings/fields/common/columns.vue';
import FieldInfoPopup from '@/components/sidebars/settings/fields/alarm/info-popup.vue';
import FieldTextEditorWithTemplate from '@/components/sidebars/settings/fields/common/text-editor-with-template.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/settings/fields/common/default-elements-per-page.vue';
import WidgetSettingsGroup from '@/components/sidebars/settings/partials/widget-settings-group.vue';

export default {
  components: {
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
    alarmColumnsWidgetTemplates() {
      return filter(this.templates, { type: WIDGET_TEMPLATES_TYPES.alarmColumns });
    },

    alarmMoreInfosWidgetTemplates() {
      return filter(this.templates, { type: WIDGET_TEMPLATES_TYPES.alarmMoreInfos });
    },

    preparedWidgetColumns() {
      return formToWidgetColumns(this.form.widgetColumns).map(column => ({
        ...column,

        text: getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
      }));
    },

    sortablePreparedWidgetColumns() {
      return this.preparedWidgetColumns.filter(column => getSortable(column, ALARM_UNSORTABLE_FIELDS));
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
  },
};
</script>
