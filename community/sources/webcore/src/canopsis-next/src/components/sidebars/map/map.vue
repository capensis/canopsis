<template>
  <widget-settings
    :submitting="submitting"
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <v-divider />
    <field-periodic-refresh v-model="form.parameters" />
    <v-divider />
    <field-map v-model="form.parameters.map" />
    <v-divider />
    <widget-settings-group :title="$t('settings.entityDisplaySettings')">
      <field-color-indicator v-model="form.parameters.color_indicator" />
      <v-divider />
      <field-switcher
        v-model="form.parameters.entities_under_pbehavior_enabled"
        :title="$t('settings.entitiesUnderPbehaviorEnabled')"
      />
    </widget-settings-group>
    <v-divider />
    <widget-settings-group :title="$t('settings.advancedSettings')">
      <template v-if="hasAccessToListFilters">
        <field-filters
          v-model="form.parameters.mainFilter"
          :filters.sync="form.filters"
          :widget-id="widget._id"
          :addable="hasAccessToAddFilter"
          :editable="hasAccessToEditFilter"
          with-alarm
          with-entity
          with-pbehavior
        />
        <v-divider />
      </template>
      <field-text-editor
        v-model="form.parameters.entity_info_template"
        :title="$t('settings.entityInfoPopup')"
        :variables="entityInfoTemplateVariables"
      />
      <v-divider />
      <field-columns
        v-model="form.parameters.alarmsColumns"
        :template="form.parameters.alarmsColumnsTemplate"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="$t('settings.alarmsColumns')"
        :type="$constants.ENTITIES_TYPES.alarm"
        with-template
        with-html
        @update:template="updateAlarmsColumnsTemplate"
      />
      <v-divider />
      <field-columns
        v-model="form.parameters.entitiesColumns"
        :template="form.parameters.entitiesColumnsTemplate"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="$t('settings.entitiesColumns')"
        :type="$constants.ENTITIES_TYPES.entity"
        with-color-indicator
        @update:template="updateEntitiesColumnsTemplate"
      />
    </widget-settings-group>
    <v-divider />
  </widget-settings>
</template>

<script>
import { ENTITY_TEMPLATE_FIELDS, SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entityVariablesMixin } from '@/mixins/widget/variables';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';
import { permissionsWidgetsMapFilters } from '@/mixins/permissions/widgets/map/filters';

import FieldTitle from '../form/fields/title.vue';
import FieldPeriodicRefresh from '../form/fields/periodic-refresh.vue';
import FieldColorIndicator from '../form/fields/color-indicator.vue';
import FieldSwitcher from '../form/fields/switcher.vue';
import FieldFilters from '../form/fields/filters.vue';
import FieldTextEditor from '../form/fields/text-editor.vue';
import FieldColumns from '../form/fields/columns.vue';
import WidgetSettings from '../partials/widget-settings.vue';
import WidgetSettingsGroup from '../partials/widget-settings-group.vue';

import FieldMap from './form/fields/map.vue';

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
  computed: {
    entityInfoTemplateVariables() {
      const excludeFields = [
        ENTITY_TEMPLATE_FIELDS.alarmLastComment,
      ];

      return this.entityVariables.filter(({ value }) => !excludeFields.includes(value));
    },
  },
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
