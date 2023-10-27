<template lang="pug">
  state-settings-list(
    :pagination.sync="pagination",
    :state-settings="stateSettingsWithStaticSetting",
    :total-items="stateSettingsMeta.total_count",
    :pending="stateSettingsPending",
    :updatable="hasUpdateAnyStateSettingAccess",
    :removable="hasDeleteAnyStateSettingAccess",
    @edit="showEditStateSettingModal",
    @duplicate="showDuplicateStateSettingModal",
    @remove="showRemoveStateSettingModal"
  )
</template>

<script>
import { MODALS, STATE_SETTING_METHODS } from '@/constants';

import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesStateSettingMixin } from '@/mixins/entities/state-setting';
import { permissionsTechnicalStateSettingMixin } from '@/mixins/permissions/technical/state-setting';

import StateSettingsList from '@/components/other/state-setting/state-settings-list.vue';

const SERVICE_STATE_SETTING_ID = 'serviceState';

export default {
  components: { StateSettingsList },
  mixins: [
    localQueryMixin,
    entitiesStateSettingMixin,
    permissionsTechnicalStateSettingMixin,
  ],
  computed: {
    stateSettingsWithStaticSetting() {
      const stateSettings = this.stateSettings.map(stateSetting => ({ editable: true, ...stateSetting }));

      stateSettings.unshift({
        type: this.$t('stateSetting.serviceState'),
        _id: SERVICE_STATE_SETTING_ID,
        method: STATE_SETTING_METHODS.worst,
        editable: false,
      });

      return stateSettings;
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    showEditStateSettingModal(stateSetting) {
      this.$modals.show({
        name: MODALS.createStateSetting,
        config: {
          stateSetting,
          title: this.$t('modals.createStateSetting.edit.title'),
          action: async (data) => {
            await this.updateStateSetting({ data, id: stateSetting._id });

            this.$popups.success({ text: this.$t('modals.createStateSetting.edit.success') });
            this.fetchList();
          },
        },
      });
    },

    showDuplicateStateSettingModal(stateSetting) {
      this.$modals.show({
        name: MODALS.createStateSetting,
        config: {
          stateSetting,
          title: this.$t('modals.createStateSetting.duplicate.title'),
          action: async (data) => {
            await this.createStateSetting({ data });

            this.$popups.success({ text: this.$t('modals.createStateSetting.duplicate.success') });
            this.fetchList();
          },
        },
      });
    },

    showRemoveStateSettingModal(stateSetting) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeStateSetting({ id: stateSetting._id });

            this.$popups.success({ text: this.$t('modals.createStateSetting.remove.success') });
            this.fetchList();
          },
        },
      });
    },

    fetchList() {
      return this.fetchStateSettingsList({ params: this.getQuery() });
    },
  },
};
</script>
