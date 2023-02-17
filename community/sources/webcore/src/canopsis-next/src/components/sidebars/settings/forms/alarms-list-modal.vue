<template lang="pug">
  widget-settings-group(:title="$t(`settings.titles.${$constants.SIDE_BARS.alarmSettings}`)")
    field-columns(
      v-field="form.widgetColumns",
      :template="form.widgetColumnsTemplate",
      :templates="alarmColumnsWidgetTemplates",
      :templates-pending="templatesPending",
      :label="$t('settings.columnNames')",
      :type="$constants.ENTITIES_TYPES.alarm",
      :alarm-infos="alarmInfos",
      :entity-infos="entityInfos",
      :infos-pending="infosPending",
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

import { WIDGET_TEMPLATES_TYPES } from '@/constants';

import { formBaseMixin } from '@/mixins/form';
import { alarmVariablesMixin } from '@/mixins/widget/variables/alarm';

import FieldColumns from '@/components/sidebars/settings/fields/common/columns.vue';
import FieldInfoPopup from '@/components/sidebars/settings/fields/alarm/info-popup.vue';
import FieldTextEditorWithTemplate from '@/components/sidebars/settings/fields/common/text-editor-with-template.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/settings/fields/common/default-elements-per-page.vue';
import WidgetSettingsGroup from '@/components/sidebars/settings/partials/widget-settings-group.vue';

export default {
  components: {
    WidgetSettingsGroup,
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
    alarmInfos: {
      type: Array,
      default: () => [],
    },
    entityInfos: {
      type: Array,
      default: () => [],
    },
    infosPending: {
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
