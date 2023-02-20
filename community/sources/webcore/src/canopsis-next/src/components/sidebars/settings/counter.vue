<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    field-filters(
      :filters.sync="form.filters",
      addable,
      editable,
      with-alarm,
      with-entity,
      with-pbehavior,
      hide-selector
    )
    v-divider
    field-opened-resolved-filter(v-model="form.parameters.opened")
    v-divider
    alarms-list-modal-form(
      v-model="form.parameters.alarmsList",
      :templates="preparedWidgetTemplates",
      :templates-pending="widgetTemplatesPending"
    )
    v-divider
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-template(
        v-model="form.parameters.blockTemplate",
        :title="$t('settings.blockTemplate')"
      )
      v-divider
      field-grid-size(
        v-model="form.parameters.columnMobile",
        :title="$t('settings.columnMobile')",
        mobile
      )
      v-divider
      field-grid-size(
        v-model="form.parameters.columnTablet",
        :title="$t('settings.columnTablet')",
        tablet
      )
      v-divider
      field-grid-size(
        v-model="form.parameters.columnDesktop",
        :title="$t('settings.columnDesktop')"
      )
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
      counter-levels-form(v-model="form.parameters.levels")
      v-divider
      field-switcher(
        v-model="form.parameters.isCorrelationEnabled",
        :title="$t('settings.isCorrelationEnabled')"
      )
    v-divider
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';

import FieldTitle from './fields/common/title.vue';
import FieldOpenedResolvedFilter from './fields/alarm/opened-resolved-filter.vue';
import FieldTemplate from './fields/common/template.vue';
import FieldGridSize from './fields/common/grid-size.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldSlider from './fields/common/slider.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import AlarmsListModalForm from './forms/alarms-list-modal.vue';
import MarginsForm from './forms/margins.vue';
import CounterLevelsForm from './forms/counter-levels.vue';
import WidgetSettings from './partials/widget-settings.vue';
import WidgetSettingsGroup from './partials/widget-settings-group.vue';

export default {
  name: SIDE_BARS.counterSettings,
  components: {
    FieldTitle,
    FieldOpenedResolvedFilter,
    FieldTemplate,
    FieldGridSize,
    FieldFilters,
    FieldSlider,
    FieldSwitcher,
    AlarmsListModalForm,
    MarginsForm,
    CounterLevelsForm,
    WidgetSettings,
    WidgetSettingsGroup,
  },
  mixins: [
    widgetSettingsMixin,
    entitiesInfosMixin,
    widgetTemplatesMixin,
  ],
  mounted() {
    return this.fetchInfos();
  },
};
</script>
