<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    field-periodic-refresh(v-model="form.parameters.periodic_refresh")
    v-divider
    template(v-if="hasAccessToListFilters")
      field-filters(
        v-model="form.parameters.mainFilter",
        :filters.sync="form.filters",
        :widget-id="widget._id",
        :addable="hasAccessToAddFilter",
        :editable="hasAccessToEditFilter",
        :entity-types="[$constants.ENTITY_TYPES.service]",
        with-entity,
        with-service-weather
      )
      v-divider
    alarms-list-modal-form(
      v-model="form.parameters.alarmsList",
      :templates="preparedWidgetTemplates",
      :templates-pending="widgetTemplatesPending"
    )
    v-divider
    field-number(v-model="form.parameters.limit", :title="$t('settings.limit')")
    v-divider
    field-color-indicator(v-model="form.parameters.colorIndicator")
    v-divider
    field-columns(
      v-model="form.parameters.serviceDependenciesColumns",
      :template="form.parameters.serviceDependenciesColumnsTemplate",
      :templates="entityColumnsWidgetTemplates",
      :templates-pending="widgetTemplatesPending",
      :label="$t('settings.treeOfDependenciesColumnNames')",
      :type="$constants.ENTITIES_TYPES.entity",
      with-color-indicator,
      @input="updateWidgetColumnsTemplate"
    )
    v-divider
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-sort-column(
        v-model="form.parameters.sort",
        :columns="sortColumns"
      )
      v-divider
      field-default-elements-per-page(v-model="form.parameters.modalItemsPerPage", :sub-title="$t('settings.modal')")
      v-divider
      field-text-editor-with-template(
        :value="form.parameters.blockTemplate",
        :template="form.parameters.blockTemplateTemplate",
        :templates="weatherItemWidgetTemplates",
        :variables="blockTemplateVariables",
        :title="$t('settings.weatherTemplate')",
        @input="updateBlockTemplate"
      )
      v-divider
      field-text-editor-with-template(
        :value="form.parameters.modalTemplate",
        :template="form.parameters.modalTemplateTemplate",
        :templates="weatherModalWidgetTemplates",
        :variables="blockTemplateVariables",
        :title="$t('settings.modalTemplate')",
        @input="updateModalTemplate"
      )
      v-divider
      field-text-editor-with-template(
        :value="form.parameters.entityTemplate",
        :template="form.parameters.entityTemplateTemplate",
        :templates="weatherEntityWidgetTemplates",
        :variables="entityVariables",
        :title="$t('settings.entityTemplate')",
        @input="updateEntityTemplate"
      )
      v-divider
      field-grid-size(v-model="form.parameters.columnMobile", :title="$t('settings.columnMobile')", mobile)
      v-divider
      field-grid-size(v-model="form.parameters.columnTablet", :title="$t('settings.columnTablet')", tablet)
      v-divider
      field-grid-size(v-model="form.parameters.columnDesktop", :title="$t('settings.columnDesktop')")
      v-divider
      margins-form(v-model="form.parameters.margin")
      v-divider
      field-slider(
        v-model="form.parameters.heightFactor",
        :title="$t('settings.height')",
        :min="1",
        :max="20"
      )
      v-divider
      field-counters-selector(v-model="form.parameters.counters", :title="$t('settings.counters')")
      v-divider
      field-switcher(v-model="form.parameters.isPriorityEnabled", :title="$t('settings.isPriorityEnabled')")
      v-divider
      field-switcher(v-model="form.parameters.isHideGrayEnabled", :title="$t('settings.isHideGrayEnabled')")
      v-divider
      field-modal-type(v-model="form.parameters.modalType")
      v-divider
      field-switcher(v-model="form.parameters.entitiesActionsInQueue", :title="$t('settings.entitiesActionsInQueue')")
    v-divider
</template>

<script>
import { ENTITY_FIELDS, ENTITY_TEMPLATE_FIELDS, SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';
import { entityVariablesMixin } from '@/mixins/widget/variables';
import { permissionsWidgetsServiceWeatherFilters } from '@/mixins/permissions/widgets/service-weather/filters';

import FieldTitle from './fields/common/title.vue';
import FieldSortColumn from './fields/service-weather/sort-column.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldTextEditorWithTemplate from './fields/common/text-editor-with-template.vue';
import FieldGridSize from './fields/common/grid-size.vue';
import FieldSlider from './fields/common/slider.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import FieldModalType from './fields/service-weather/modal-type.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldNumber from './fields/common/number.vue';
import FieldCountersSelector from './fields/common/counters-selector.vue';
import FieldColorIndicator from './fields/common/color-indicator.vue';
import AlarmsListModalForm from './forms/alarms-list-modal.vue';
import MarginsForm from './forms/margins.vue';
import WidgetSettings from './partials/widget-settings.vue';
import WidgetSettingsGroup from './partials/widget-settings-group.vue';

export default {
  name: SIDE_BARS.serviceWeatherSettings,
  components: {
    FieldTitle,
    FieldSortColumn,
    FieldPeriodicRefresh,
    FieldFilters,
    FieldColumns,
    FieldDefaultSortColumn,
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
    WidgetSettings,
    WidgetSettingsGroup,
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
