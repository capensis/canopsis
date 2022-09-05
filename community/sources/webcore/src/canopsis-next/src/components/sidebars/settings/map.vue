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
            v-model="form.parameters.entities_under_pbehavior_enabled",
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
            :title="$t('settings.entityInfoPopup')",
            :variables="variables"
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
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT, SIDE_BARS } from '@/constants';

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

const { mapActions: mapServiceActions } = createNamespacedHelpers('service');

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
  },
  mixins: [
    widgetSettingsMixin,
    permissionsWidgetsMapFilters,
  ],
  data() {
    return {
      infos: [],
    };
  },
  computed: {
    infosSubVariables() {
      return [
        {
          text: this.$t('common.value'),
          value: 'value',
        },
        {
          text: this.$t('common.description'),
          value: 'description',
        },
      ];
    },

    infosVariables() {
      return this.infos.map(({ value }) => ({
        text: value,
        value,
        variables: this.infosSubVariables,
      }));
    },

    variables() {
      return [
        {
          text: this.$t('common.id'),
          value: 'entity._id',
        },
        {
          text: this.$t('common.name'),
          value: 'entity.name',
        },
        {
          text: this.$t('common.infos'),
          value: 'entity.infos',
          variables: this.infosVariables,
        },
        {
          text: this.$t('common.connector'),
          value: 'entity.connector',
        },
        {
          text: this.$t('common.connectorName'),
          value: 'entity.connector_name',
        },
        {
          text: this.$t('common.component'),
          value: 'entity.component',
        },
        {
          text: this.$t('common.resource'),
          value: 'entity.resource',
        },
        {
          text: this.$t('common.state'),
          value: 'entity.state.val',
        },
        {
          text: this.$t('common.status'),
          value: 'entity.status.val',
        },
        {
          text: this.$t('common.snooze'),
          value: 'entity.snooze',
        },
        {
          text: this.$t('common.ack'),
          value: 'entity.ack',
        },
        {
          text: this.$t('common.updated'),
          value: 'entity.last_update_date',
        },
        {
          text: this.$t('common.impactLevel'),
          value: 'entity.impact_level',
        },
        {
          text: this.$t('common.impactState'),
          value: 'entity.impact_state',
        },
        {
          text: this.$t('common.category'),
          value: 'entity.category.name',
        },
        {
          text: this.$tc('common.link', 2),
          value: 'entity.links',
        },
      ];
    },
  },
  mounted() {
    this.fetchInfos();
  },
  methods: {
    ...mapServiceActions({ fetchEntityInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore' }),

    async fetchInfos() {
      const { data: infos } = await this.fetchEntityInfosKeysWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.infos = infos;
    },
  },
};
</script>
