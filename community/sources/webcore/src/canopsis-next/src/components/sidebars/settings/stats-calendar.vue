<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="form.title", :title="$t('common.title')")
      v-divider
      alarms-list-modal-form(v-model="form.parameters.alarmsList")
      v-divider
      v-list-group
        template(#activator="")
          v-list-tile
            div(:class="validationHeaderClass") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-filters(
            :filters.sync="form.filters",
            :widget-id="widget._id",
            addable,
            editable,
            with-alarm,
            with-entity,
            with-pbehavior,
            hide-selector
          )
          v-divider
          field-opened-resolved-filter(v-field="form.parameters.opened")
          v-divider
          field-switcher(
            v-field="form.parameters.considerPbehaviors",
            :title="$t('settings.considerPbehaviors.title')"
          )
          v-divider
          field-criticity-levels(v-field="form.parameters.criticityLevels")
          v-divider
          field-levels-colors-selector(
            v-field="form.parameters.criticityLevelsColors",
            color-type="hex",
            hide-suffix
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
import { formValidationHeaderMixin } from '@/mixins/form';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldOpenedResolvedFilter from '@/components/sidebars/settings/fields/alarm/opened-resolved-filter.vue';
import FieldFilters from '@/components/sidebars/settings/fields/common/filters.vue';
import FieldSwitcher from '@/components/sidebars/settings/fields/common/switcher.vue';
import FieldCriticityLevels from '@/components/sidebars/settings/fields/stats/criticity-levels.vue';
import FieldLevelsColorsSelector from '@/components/sidebars/settings/fields/stats/levels-colors-selector.vue';
import AlarmsListModalForm from '@/components/sidebars/settings/forms/alarms-list-modal.vue';
import StatsCalendarAdvancedForm from '@/components/sidebars/settings/forms/stats-calendar-advanced.vue';

/**
 * Component to regroup the entities list settings fields
 */
export default {
  name: SIDE_BARS.statsCalendarSettings,
  components: {
    FieldTitle,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldSwitcher,
    FieldCriticityLevels,
    FieldLevelsColorsSelector,
    AlarmsListModalForm,
    StatsCalendarAdvancedForm,
  },
  mixins: [widgetSettingsMixin, formValidationHeaderMixin],
};
</script>
