<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-filters(
        :entities-type="$constants.ENTITIES_TYPES.alarm",
        :filters.sync="settings.widget.parameters.viewFilters",
        hide-select,
        addable,
        editable
      )
      v-divider
      field-opened-resolved-filter(v-model="settings.widget.parameters.opened")
      v-divider
      alarms-list-modal-form(v-model="settings.widget.parameters.alarmsList")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-template(
            v-model="settings.widget.parameters.blockTemplate",
            :title="$t('settings.blockTemplate')"
          )
          v-divider
          field-grid-size(
            v-model="settings.widget.parameters.columnSM",
            :title="$t('settings.columnSM')"
          )
          v-divider
          field-grid-size(
            v-model="settings.widget.parameters.columnMD",
            :title="$t('settings.columnMD')"
          )
          v-divider
          field-grid-size(
            v-model="settings.widget.parameters.columnLG",
            :title="$t('settings.columnLG')"
          )
          v-divider
          margins-form(v-model="settings.widget.parameters.margin")
          v-divider
          field-slider(
            data-test="widgetHeightFactory",
            v-model="settings.widget.parameters.heightFactor",
            :title="$t('settings.height')",
            :min="1",
            :max="20"
          )
          v-divider
          counter-levels-form(v-model="settings.widget.parameters.levels")
          v-divider
          field-switcher(
            v-model="settings.widget.parameters.isCorrelationEnabled",
            :title="$t('settings.isCorrelationEnabled')"
          )
      v-divider
    v-btn.primary(data-test="submitWeather", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { sideBarSettingsWidgetAlarmMixin } from '@/mixins/side-bar/settings/widgets/alarm';

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

export default {
  name: SIDE_BARS.counterSettings,
  $_veeValidate: {
    validator: 'new',
  },
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
  },
  mixins: [widgetSettingsMixin, sideBarSettingsWidgetAlarmMixin],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: this.prepareWidgetWithAlarmParametersSettings(cloneDeep(widget), true),
      },
    };
  },
  methods: {
    prepareWidgetSettings() {
      const { widget } = this.settings;

      return this.prepareWidgetWithAlarmParametersSettings(widget);
    },
  },
};
</script>
