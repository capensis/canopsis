<template>
  <widget-settings
    :submitting="submitting"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <field-periodic-refresh v-model="form.parameters" />
    <template v-if="hasAccessToListFilters">
      <field-filters
        v-model="form.parameters.mainFilter"
        :filters.sync="form.filters"
        :widget-id="widget._id"
        :addable="hasAccessToAddFilter"
        :editable="hasAccessToEditFilter"
        :entity-types="[$constants.ENTITY_TYPES.service]"
        with-entity
        with-service-weather
        entity-counters-type
      />
    </template>
    <alarms-list-modal-form
      v-model="form.parameters.alarmsList"
      :templates="preparedWidgetTemplates"
      :templates-pending="widgetTemplatesPending"
    />
    <field-number
      v-model="form.parameters.limit"
      :title="$t('settings.limit')"
    />
    <field-color-indicator v-model="form.parameters.colorIndicator" />
    <field-columns
      v-model="form.parameters.serviceDependenciesColumns"
      :template="form.parameters.serviceDependenciesColumnsTemplate"
      :templates="entityColumnsWidgetTemplates"
      :templates-pending="widgetTemplatesPending"
      :label="$t('settings.treeOfDependenciesColumnNames')"
      :type="$constants.ENTITIES_TYPES.entity"
      with-color-indicator
      @input="updateWidgetColumnsTemplate"
    />
    <field-tree-of-dependencies-settings v-model="form.parameters.treeOfDependenciesShowType" />
    <widget-settings-group :title="$t('settings.advancedSettings')">
      <field-sort-column
        v-model="form.parameters.sort"
        :columns="sortColumns"
      />
      <field-default-elements-per-page
        v-model="form.parameters.modalItemsPerPage"
        :sub-title="$t('settings.modal')"
      />
      <field-root-cause-settings v-model="form.parameters" />
      <field-text-editor-with-template
        :value="form.parameters.blockTemplate"
        :template="form.parameters.blockTemplateTemplate"
        :templates="weatherItemWidgetTemplates"
        :variables="blockTemplateVariables"
        :title="$t('settings.weatherTemplate')"
        @input="updateBlockTemplate"
      />
      <field-text-editor-with-template
        :value="form.parameters.modalTemplate"
        :template="form.parameters.modalTemplateTemplate"
        :templates="weatherModalWidgetTemplates"
        :variables="blockTemplateVariables"
        :title="$t('settings.modalTemplate')"
        @input="updateModalTemplate"
      />
      <field-text-editor-with-template
        :value="form.parameters.entityTemplate"
        :template="form.parameters.entityTemplateTemplate"
        :templates="weatherEntityWidgetTemplates"
        :variables="entityTemplateVariables"
        :title="$t('settings.entityTemplate')"
        @input="updateEntityTemplate"
      />
      <field-grid-size
        v-model="form.parameters.columnMobile"
        :title="$t('settings.columnMobile')"
        mobile
      />
      <field-grid-size
        v-model="form.parameters.columnTablet"
        :title="$t('settings.columnTablet')"
        tablet
      />
      <field-grid-size
        v-model="form.parameters.columnDesktop"
        :title="$t('settings.columnDesktop')"
      />
      <margins-form v-model="form.parameters.margin" />
      <field-slider
        v-model="form.parameters.heightFactor"
        :title="$t('settings.height')"
        :min="1"
        :max="20"
      />
      <field-counters-selector
        v-model="form.parameters.counters"
        :title="$t('settings.counters')"
      />
      <field-switcher
        v-model="form.parameters.isPriorityEnabled"
        :title="$t('settings.isPriorityEnabled')"
      />
      <field-switcher
        v-model="form.parameters.isHideGrayEnabled"
        :title="$t('settings.isHideGrayEnabled')"
      />
      <field-switcher
        v-model="form.parameters.isSecondaryIconEnabled"
        :title="$t('settings.isSecondaryIconEnabled')"
      />
      <field-modal-type v-model="form.parameters.modalType" />
      <widget-settings-group :title="$t('settings.stateSettings')">
        <field-action-required-settings v-model="form.parameters.actionRequiredSettings" />
        <field-no-action-required-settings v-model="form.parameters.noActionRequiredSettings" />
      </widget-settings-group>
      <v-divider />
      <field-switcher
        v-model="form.parameters.entitiesActionsInQueue"
        :title="$t('settings.entitiesActionsInQueue')"
      />
    </widget-settings-group>
  </widget-settings>
</template>

<script>
import { ENTITY_FIELDS, ENTITY_FIELDS_TO_LABELS_KEYS, ENTITY_TEMPLATE_FIELDS, SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';
import { entityVariablesMixin } from '@/mixins/widget/variables';
import { permissionsWidgetsServiceWeatherFilters } from '@/mixins/permissions/widgets/service-weather/filters';

import FieldTreeOfDependenciesSettings from '../form/fields/tree-of-dependencies-settings.vue';
import FieldTitle from '../form/fields/title.vue';
import FieldPeriodicRefresh from '../form/fields/periodic-refresh.vue';
import FieldFilters from '../form/fields/filters.vue';
import FieldColumns from '../form/fields/columns.vue';
import FieldTextEditorWithTemplate from '../form/fields/text-editor-with-template.vue';
import FieldGridSize from '../form/fields/grid-size.vue';
import FieldSlider from '../form/fields/slider.vue';
import FieldSwitcher from '../form/fields/switcher.vue';
import FieldDefaultElementsPerPage from '../form/fields/default-elements-per-page.vue';
import FieldNumber from '../form/fields/number.vue';
import FieldColorIndicator from '../form/fields/color-indicator.vue';
import AlarmsListModalForm from '../alarm/form/alarms-list-modal.vue';
import MarginsForm from '../form/margins.vue';
import WidgetSettings from '../partials/widget-settings.vue';
import WidgetSettingsGroup from '../partials/widget-settings-group.vue';
import FieldRootCauseSettings from '../form/fields/root-cause-settings.vue';

import FieldCountersSelector from './form/fields/counters-selector.vue';
import FieldSortColumn from './form/fields/sort-column.vue';
import FieldModalType from './form/fields/modal-type.vue';
import FieldActionRequiredSettings from './form/fields/field-action-required-settings.vue';
import FieldNoActionRequiredSettings from './form/fields/field-no-action-required-settings.vue';

export default {
  name: SIDE_BARS.serviceWeatherSettings,
  components: {
    FieldTreeOfDependenciesSettings,
    FieldTitle,
    FieldSortColumn,
    FieldPeriodicRefresh,
    FieldFilters,
    FieldColumns,
    FieldTextEditorWithTemplate,
    FieldGridSize,
    FieldSlider,
    FieldSwitcher,
    FieldModalType,
    FieldDefaultElementsPerPage,
    FieldNumber,
    FieldCountersSelector,
    FieldColorIndicator,
    AlarmsListModalForm,
    MarginsForm,
    FieldActionRequiredSettings,
    FieldNoActionRequiredSettings,
    WidgetSettings,
    WidgetSettingsGroup,
    FieldRootCauseSettings,
  },
  mixins: [
    widgetSettingsMixin,
    entitiesInfosMixin,
    widgetTemplatesMixin,
    entityVariablesMixin,
    permissionsWidgetsServiceWeatherFilters,
  ],
  computed: {
    sortColumns() {
      return [
        { label: this.$t('common.name'), value: ENTITY_FIELDS.name },
        { label: this.$t('common.state'), value: ENTITY_FIELDS.state },
      ];
    },

    blockTemplateVariables() {
      const excludeFields = [
        ENTITY_TEMPLATE_FIELDS.ticket,
        ENTITY_TEMPLATE_FIELDS.statsKo,
        ENTITY_TEMPLATE_FIELDS.statsOk,
        ENTITY_TEMPLATE_FIELDS.alarmDisplayName,
        ENTITY_TEMPLATE_FIELDS.alarmCreationDate,
      ];

      return this.entityVariables.filter(({ value }) => !excludeFields.includes(value));
    },

    entityTemplateVariables() {
      return [
        ...this.entityVariables,

        {
          value: ENTITY_TEMPLATE_FIELDS.sourceType,
          text: this.$t(ENTITY_FIELDS_TO_LABELS_KEYS[ENTITY_FIELDS.sourceType]),
        },
      ];
    },
  },
  mounted() {
    this.fetchInfos();
  },
  methods: {
    updateServiceDependenciesColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'serviceDependenciesColumnsTemplate', template);
      this.$set(this.form.parameters, 'serviceDependenciesColumns', columns);
    },

    updateWidgetColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'widgetColumnsTemplate', template);
      this.$set(this.form.parameters, 'widgetColumns', columns);
    },

    updateBlockTemplate(text, template) {
      this.$set(this.form.parameters, 'blockTemplate', text);

      if (template && template !== this.form.parameters.blockTemplateTemplate) {
        this.$set(this.form.parameters, 'blockTemplateTemplate', template);
      }
    },

    updateModalTemplate(text, template) {
      this.$set(this.form.parameters, 'modalTemplate', text);

      if (template && template !== this.form.parameters.modalTemplateTemplate) {
        this.$set(this.form.parameters, 'modalTemplateTemplate', template);
      }
    },

    updateEntityTemplate(text, template) {
      this.$set(this.form.parameters, 'entityTemplate', text);

      if (template && template !== this.form.parameters.entityTemplateTemplate) {
        this.$set(this.form.parameters, 'entityTemplateTemplate', template);
      }
    },
  },
};
</script>
