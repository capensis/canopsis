<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    field-periodic-refresh(v-model="form.parameters.periodic_refresh")
    v-divider
    field-map(v-model="form.parameters.map")
    v-divider
    widget-settings-group(:title="$t('settings.entityDisplaySettings')")
      field-color-indicator(v-model="form.parameters.color_indicator")
      v-divider
      field-switcher(
        v-model="form.parameters.entities_under_pbehavior_enabled",
        :title="$t('settings.entitiesUnderPbehaviorEnabled')"
      )
    v-divider
    widget-settings-group(:title="$t('settings.advancedSettings')")
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
      field-text-editor(
        v-model="form.parameters.entity_info_template",
        :title="$t('settings.entityInfoPopup')",
        :variables="entityVariables"
      )
      v-divider

      field-columns(
        v-model="form.parameters.alarmsColumns",
        :template="form.parameters.alarmsColumnsTemplate",
        :templates="alarmColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        :label="$t('settings.alarmsColumns')",
        :type="$constants.ENTITIES_TYPES.alarm",
        with-template,
        with-html,
        @update:template="updateAlarmsColumnsTemplate"
      )
      v-divider
      field-columns(
        v-model="form.parameters.entitiesColumns",
        :template="form.parameters.entitiesColumnsTemplate",
        :templates="entityColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        :label="$t('settings.entitiesColumns')",
        :type="$constants.ENTITIES_TYPES.entity",
        with-html,
        with-color-indicator,
        @update:template="updateEntitiesColumnsTemplate"
      )
    v-divider
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entityVariablesMixin } from '@/mixins/widget/variables';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';
import { permissionsWidgetsMapFilters } from '@/mixins/permissions/widgets/map/filters';

import FieldTitle from './fields/common/title.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldMap from './fields/map/map.vue';
import FieldColorIndicator from './fields/common/color-indicator.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldTextEditor from './fields/common/text-editor.vue';
import FieldColumns from './fields/common/columns.vue';
import WidgetSettings from './partials/widget-settings.vue';
import WidgetSettingsGroup from './partials/widget-settings-group.vue';

/**
 * Component to regroup the map settings fields
 */
export default {
  name: SIDE_BARS.mapSettings,
  components: {
    FieldTitle,
    FieldPeriodicRefresh,
    FieldMap,
    FieldColorIndicator,
    FieldSwitcher,
    FieldFilters,
    FieldTextEditor,
    FieldColumns,
    WidgetSettings,
    WidgetSettingsGroup,
  },
  mixins: [
    widgetSettingsMixin,
    entityVariablesMixin,
    entitiesInfosMixin,
    widgetTemplatesMixin,
    permissionsWidgetsMapFilters,
  ],
  mounted() {
    this.fetchInfos();
  },
  methods: {
    updateAlarmsColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'alarmsColumnsTemplate', template);
      this.$set(this.form.parameters, 'alarmsColumns', columns);
    },

    updateEntitiesColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'entitiesColumnsTemplate', template);
      this.$set(this.form.parameters, 'entitiesColumns', columns);
    },
  },
};
</script>
