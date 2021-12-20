<template lang="pug">
  v-card-text
    remediation-configurations-list(
      :remediation-configurations="remediationConfigurations",
      :pending="remediationConfigurationsPending",
      :total-items="remediationConfigurationsMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedRemediationConfigurationsModal",
      @remove="showRemoveRemediationConfigurationModal",
      @edit="showEditRemediationConfigurationModal"
    )
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesRemediationConfigurationMixin } from '@/mixins/entities/remediation/configuration';
import { localQueryMixin } from '@/mixins/query-local/query';

import RemediationConfigurationsList from './remediation-configurations-list.vue';

export default {
  components: { RemediationConfigurationsList },
  mixins: [
    localQueryMixin,
    entitiesRemediationConfigurationMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      const params = this.getQuery();
      params.with_flags = true;

      this.fetchRemediationConfigurationsList({ params });
    },

    showEditRemediationConfigurationModal(remediationConfiguration) {
      this.$modals.show({
        name: MODALS.createRemediationConfiguration,
        config: {
          title: this.$t('modals.createRemediationConfiguration.edit.title'),
          remediationConfiguration,
          action: async (configuration) => {
            await this.updateRemediationConfiguration({ id: remediationConfiguration._id, data: configuration });

            this.$popups.success({
              text: this.$t('modals.createRemediationConfiguration.edit.popups.success', {
                configurationName: configuration.name,
              }),
            });

            await this.fetchList();
          },
        },
      });
    },

    showRemoveRemediationConfigurationModal(remediationConfiguration) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeRemediationConfiguration({ id: remediationConfiguration._id });
            await this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedRemediationConfigurationsModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id: id }) => this.removeRemediationConfiguration({ id })));
            await this.fetchList();
          },
        },
      });
    },
  },
};
</script>
