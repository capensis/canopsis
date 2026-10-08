<template>
  <widget-settings
    :submitting="submitting"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <field-filters
      :filters.sync="form.filters"
      addable
      editable
      with-alarm
      with-entity
      with-pbehavior
      hide-selector
    />
    <field-opened-resolved-filter v-model="form.parameters.opened" />
    <alarms-list-modal-form
      v-model="form.parameters.alarmsList"
      :templates="preparedWidgetTemplates"
      :templates-pending="widgetTemplatesPending"
    />
    <widget-settings-group :title="$t('settings.advancedSettings')">
      <field-template
        v-model="form.parameters.blockTemplate"
        :title="$t('settings.blockTemplate')"
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
      <counter-levels-form v-model="form.parameters.levels" />
      <field-switcher
        v-model="form.parameters.isCorrelationEnabled"
        :title="$t('settings.isCorrelationEnabled')"
      />
    </widget-settings-group>
  </widget-settings>
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';

import FieldOpenedResolvedFilter from '../alarm/form/fields/opened-resolved-filter.vue';
import FieldTitle from '../form/fields/title.vue';
import FieldGridSize from '../form/fields/grid-size.vue';
import FieldFilters from '../form/fields/filters.vue';
import FieldSlider from '../form/fields/slider.vue';
import FieldSwitcher from '../form/fields/switcher.vue';
import AlarmsListModalForm from '../alarm/form/alarms-list-modal.vue';
import MarginsForm from '../form/margins.vue';
import WidgetSettings from '../partials/widget-settings.vue';
import WidgetSettingsGroup from '../partials/widget-settings-group.vue';

import CounterLevelsForm from './form/counter-levels.vue';
import FieldTemplate from './form/fields/template.vue';

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
