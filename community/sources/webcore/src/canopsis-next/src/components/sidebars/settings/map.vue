<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="form.title", :title="$t('common.title')")
      v-divider
      field-periodic-refresh(v-model="form.parameters.periodic_refresh")
      v-divider
      field-map(v-model="form.parameters.map")
      v-divider
      v-list-group
        template(#activator="")
          v-list-tile {{ $t('settings.entityDisplaySettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-color-indicator(v-model="form.parameters.color_indicator")
          v-divider
          field-switcher(
            v-field="form.parameters.entities_under_pbehavior_enabled",
            :title="$t('settings.entitiesUnderPbehaviorEnabled')"
          )
      v-divider
      v-list-group
        template(#activator="")
          v-list-tile {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          template(v-if="hasAccessToListFilters")
            field-filters(
              v-model="form.parameters.mainFilter",
              :filters.sync="form.filters",
              :widget-id="widget._id",
              :addable="hasAccessToAddFilter",
              :editable="hasAccessToEditFilter",
              with-alarm,
              with-entity,
              with-pbehavior,
              @input="updateMainFilterUpdatedAt"
            )
            v-divider
            field-text-editor(
              v-model="form.parameters.entity_info_template",
              :title="$t('settings.entityInfoPopup')"
            )
            v-divider

          field-columns(
            v-model="form.parameters.alarms_columns",
            :label="$t('settings.alarmsColumns')",
            with-template,
            with-html,
            with-color-indicator
          )
          v-divider
          field-columns(
            v-model="form.parameters.entities_columns",
            :label="$t('settings.entitiesColumns')",
            with-html,
            with-color-indicator
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
import { permissionsWidgetsMapFilters } from '@/mixins/permissions/widgets/map/filters';

import FieldTitle from './fields/common/title.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldMap from './fields/map/map.vue';
import FieldColorIndicator from './fields/common/color-indicator.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldTextEditor from './fields/common/text-editor.vue';
import FieldColumns from './fields/common/columns.vue';

/**
 * Component to regroup the map settings fields
 */
export default {
  name: SIDE_BARS.alarmSettings,
  components: {
    FieldTitle,
    FieldPeriodicRefresh,
    FieldMap,
    FieldColorIndicator,
    FieldSwitcher,
    FieldFilters,
    FieldTextEditor,
    FieldColumns,
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsMapFilters,
  ],
};
</script>
