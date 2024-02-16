<template>
  <state-settings-list
    :state-settings="stateSettingsWithStaticSetting"
    :pending="stateSettingsPending"
    @edit="showEditStateSettingModal"
  />
</template>

<script>
import { MAX_LIMIT, MODALS, STATE_SETTING_METHODS } from '@/constants';

import { entitiesStateSettingMixin } from '@/mixins/entities/state-setting';

import StateSettingsList from '@/components/other/state-setting/state-settings-list.vue';

const SERVICE_STATE_SETTING_ID = 'serviceState';

export default {
  components: { StateSettingsList },
  mixins: [entitiesStateSettingMixin],
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
        name: MODALS.stateSetting,
        config: {
          stateSetting,
          action: async (data) => {
            await this.updateStateSetting({ data, id: stateSetting._id });

            this.fetchList();
          },
        },
      });
    },

    async fetchList() {
      this.fetchStateSettingsList({ params: { page: 1, limit: MAX_LIMIT } });
    },
  },
};
</script>
