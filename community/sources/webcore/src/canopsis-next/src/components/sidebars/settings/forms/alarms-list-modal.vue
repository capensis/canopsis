<template lang="pug">
  widget-settings-group(:title="$t(`settings.titles.${$constants.SIDE_BARS.alarmSettings}`)")
    field-columns(
      v-field="form.widgetColumns",
      :template="form.widgetColumnsTemplate",
      :templates="templates",
      :templates-pending="templatesPending",
      :label="$t('settings.columnNames')",
      :type="$constants.ENTITIES_TYPES.alarm",
      :alarm-infos="alarmInfos",
      :entity-infos="entityInfos",
      :infos-pending="infosPending",
      with-html,
      with-state,
      @update:template="updateTemplate"
    )
    v-divider
    field-default-elements-per-page(v-field="form.itemsPerPage")
    v-divider
    field-info-popup(
      v-field="form.infoPopups",
      :columns="form.widgetColumns"
    )
    v-divider
    field-text-editor(
      v-field="form.moreInfoTemplate",
      :title="$t('settings.moreInfosModal')"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import FieldColumns from '@/components/sidebars/settings/fields/common/columns.vue';
import FieldInfoPopup from '@/components/sidebars/settings/fields/alarm/info-popup.vue';
import FieldTextEditor from '@/components/sidebars/settings/fields/common/text-editor.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/settings/fields/common/default-elements-per-page.vue';
import WidgetSettingsGroup from '@/components/sidebars/settings/partials/widget-settings-group.vue';

export default {
  components: {
    WidgetSettingsGroup,
    FieldColumns,
    FieldInfoPopup,
    FieldTextEditor,
    FieldDefaultElementsPerPage,
  },
  mixins: [formBaseMixin],
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
  methods: {
    updateTemplate(template, columns) {
      this.updateModel({
        ...this.form,

        widgetColumnsTemplate: template,
        widgetColumns: columns,
      });
    },
  },
};
</script>
