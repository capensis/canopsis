<template>
  <c-page
    :creatable="hasCreateAnyStateSettingAccess"
    :create-tooltip="$t('modals.createStateSetting.create.title')"
    @refresh="fetchList"
    @create="showCreateStateSettingModal"
  >
    <state-settings-list
      :options.sync="options"
      :state-settings="stateSettings"
      :total-items="stateSettingsMeta.total_count"
      :pending="stateSettingsPending"
      :addable="hasCreateAnyStateSettingAccess"
      :updatable="hasUpdateAnyStateSettingAccess"
      :removable="hasDeleteAnyStateSettingAccess"
      @edit="showEditStateSettingModal"
      @duplicate="showDuplicateStateSettingModal"
      @remove="showRemoveStateSettingModal"
    />
  </c-page>
</template>

<script>

import { JUNIT_STATE_SETTING_ID, MODALS } from '@/constants';

import { localQueryMixin } from '@/mixins/query/query';
import { entitiesStateSettingMixin } from '@/mixins/entities/state-setting';
import { permissionsTechnicalStateSettingMixin } from '@/mixins/permissions/technical/state-setting';

import StateSettingsList from '@/components/other/state-setting/state-settings-list.vue';

export default {
  components: { StateSettingsList },
  mixins: [
    localQueryMixin,
    entitiesStateSettingMixin,
    permissionsTechnicalStateSettingMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    showCreateStateSettingModal() {
      this.$modals.show({
        name: MODALS.createStateSetting,
        config: {
          title: this.$t('modals.createStateSetting.create.title'),
          action: async (newStateSetting) => {
            await this.createStateSetting({ data: newStateSetting });

            this.$popups.success({ text: this.$t('modals.createStateSetting.create.success') });
            this.fetchList();
          },
        },
      });
    },

    showEditJunitStateSettingModal(stateSetting) {
      this.$modals.show({
        name: MODALS.createJunitStateSetting,
        config: {
          stateSetting,
          title: this.$t('modals.createJunitStateSetting.edit.title'),
          action: async (data) => {
            await this.updateStateSetting({ data, id: stateSetting._id });

            this.$popups.success({ text: this.$t('modals.createJunitStateSetting.edit.success') });
            this.fetchList();
          },
        },
      });
    },

    showEditStateSettingModal(stateSetting) {
      if (stateSetting._id === JUNIT_STATE_SETTING_ID) {
        this.showEditJunitStateSettingModal(stateSetting);

        return;
      }

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
