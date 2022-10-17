<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="form.title", :title="$t('common.title')")
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
      alarms-list-modal-form(v-model="form.parameters.alarmsList")
      v-divider
      v-list-group
        template(#activator="")
          v-list-tile {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-template(
            v-model="form.parameters.blockTemplate",
            :title="$t('settings.blockTemplate')"
          )
          v-divider
          field-grid-size(
            v-model="form.parameters.columnSM",
            :title="$t('settings.columnMobile')",
            mobile
          )
          v-divider
          field-grid-size(
            v-model="form.parameters.columnMD",
            :title="$t('settings.columnTablet')",
            tablet
          )
          v-divider
          field-grid-size(
            v-model="form.parameters.columnLG",
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
    v-btn.primary(
      :loading="submitting",
      :disabled="submitting",
      @click="submit"
    ) {{ $t('common.save') }}
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldOpenedResolvedFilter from '@/components/sidebars/settings/fields/alarm/opened-resolved-filter.vue';
import FieldTemplate from '@/components/sidebars/settings/fields/common/template.vue';
import FieldGridSize from '@/components/sidebars/settings/fields/common/grid-size.vue';
import FieldFilters from '@/components/sidebars/settings/fields/common/filters.vue';
import FieldSlider from '@/components/sidebars/settings/fields/common/slider.vue';
import FieldSwitcher from '@/components/sidebars/settings/fields/common/switcher.vue';
import AlarmsListModalForm from '@/components/sidebars/settings/forms/alarms-list-modal.vue';
import MarginsForm from '@/components/sidebars/settings/forms/margins.vue';
import CounterLevelsForm from '@/components/sidebars/settings/forms/counter-levels.vue';

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
  },
  mixins: [widgetSettingsMixin],
};
</script>
