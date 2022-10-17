<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="form.title", :title="$t('common.title')")
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
      alarms-list-modal-form(v-model="form.parameters.alarmsList")
      v-divider
      field-number(
        v-model="form.parameters.limit",
        :title="$t('settings.limit')"
      )
      v-divider
      field-color-indicator(v-model="form.parameters.colorIndicator")
      v-divider
      field-columns(
        v-model="form.parameters.serviceDependenciesColumns",
        :label="$t('settings.treeOfDependenciesColumnNames')",
        with-color-indicator
      )
      v-divider
      v-list-group
        template(#activator="")
          v-list-tile {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-sort-column(
            v-model="form.parameters.sort",
            :columns="sortColumns",
            :columns-label="$t('settings.orderBy')"
          )
          v-divider
          field-default-elements-per-page(v-model="form.parameters.modalItemsPerPage")
            template(#title="")
              span {{ $t('settings.defaultNumberOfElementsPerPage') }}
                span.font-italic.caption.ml-1 (Modal)
          v-divider
          field-template(
            v-model="form.parameters.blockTemplate",
            :title="$t('settings.weatherTemplate')"
          )
          v-divider
          field-template(
            v-model="form.parameters.modalTemplate",
            :title="$t('settings.modalTemplate')"
          )
          v-divider
          field-template(
            v-model="form.parameters.entityTemplate",
            :title="$t('settings.entityTemplate')"
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
          field-counters-selector(
            v-model="form.parameters.counters",
            :title="$t('settings.counters')"
          )
          v-divider
          field-switcher(
            v-model="form.parameters.isPriorityEnabled",
            :title="$t('settings.isPriorityEnabled')"
          )
          v-divider
          field-modal-type(v-model="form.parameters.modalType")
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
import { permissionsWidgetsServiceWeatherFilters } from '@/mixins/permissions/widgets/service-weather/filters';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldSortColumn from '@/components/sidebars/settings/fields/service-weather/sort-column.vue';
import FieldPeriodicRefresh from '@/components/sidebars/settings/fields/common/periodic-refresh.vue';
import FieldFilters from '@/components/sidebars/settings/fields/common/filters.vue';
import FieldColumns from '@/components/sidebars/settings/fields/common/columns.vue';
import FieldDefaultSortColumn from '@/components/sidebars/settings/fields/common/default-sort-column.vue';
import FieldTemplate from '@/components/sidebars/settings/fields/common/template.vue';
import FieldGridSize from '@/components/sidebars/settings/fields/common/grid-size.vue';
import FieldSlider from '@/components/sidebars/settings/fields/common/slider.vue';
import FieldSwitcher from '@/components/sidebars/settings/fields/common/switcher.vue';
import FieldModalType from '@/components/sidebars/settings/fields/service-weather/modal-type.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/settings/fields/common/default-elements-per-page.vue';
import FieldNumber from '@/components/sidebars/settings/fields/common/number.vue';
import FieldCountersSelector from '@/components/sidebars/settings/fields/common/counters-selector.vue';
import FieldColorIndicator from '@/components/sidebars/settings/fields/common/color-indicator.vue';
import AlarmsListModalForm from '@/components/sidebars/settings/forms/alarms-list-modal.vue';
import MarginsForm from '@/components/sidebars/settings/forms/margins.vue';

export default {
  name: SIDE_BARS.serviceWeatherSettings,
  components: {
    FieldTitle,
    FieldSortColumn,
    FieldPeriodicRefresh,
    FieldFilters,
    FieldColumns,
    FieldDefaultSortColumn,
    FieldTemplate,
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
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsServiceWeatherFilters,
  ],
  computed: {
    sortColumns() {
      return [
        { label: this.$t('common.name'), value: 'name' },
        { label: this.$t('common.state'), value: 'state' },
      ];
    },
  },
};
</script>
