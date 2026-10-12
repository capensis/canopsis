<template>
  <v-card-text>
    <remediation-configurations-list
      :remediation-configurations="remediationConfigurations"
      :pending="remediationConfigurationsPending"
      :total-items="remediationConfigurationsMeta.total_count"
      :options.sync="options"
      :removable="hasDeleteAnyRemediationConfigurationAccess"
      :duplicable="hasCreateAnyRemediationConfigurationAccess"
      :updatable="hasUpdateAnyRemediationConfigurationAccess"
      :active="active"
      @remove="showRemoveRemediationConfigurationModal"
      @duplicate="showDuplicateRemediationConfigurationModal"
      @edit="showEditRemediationConfigurationModal"
      @refresh="fetchList"
    />
  </v-card-text>
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { localQueryMixin } from '@/mixins/query/query';
import { entitiesRemediationConfigurationMixin } from '@/mixins/entities/remediation/configuration';
import {
  permissionsTechnicalRemediationConfigurationMixin,
} from '@/mixins/permissions/technical/remediation-configuration';

import RemediationConfigurationsList from './remediation-configurations-list.vue';

export default {
  components: { RemediationConfigurationsList },
  mixins: [
    localQueryMixin,
    entitiesRemediationConfigurationMixin,
    permissionsTechnicalRemediationConfigurationMixin,
  ],
  props: {
    active: {
      type: Boolean,
      default: true,
    },
  },
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

    showDuplicateRemediationConfigurationModal(remediationConfiguration) {
      this.$modals.show({
        name: MODALS.createRemediationConfiguration,
        config: {
          title: this.$t('modals.createRemediationConfiguration.duplicate.title'),
          remediationConfiguration: omit(remediationConfiguration, ['_id']),
          action: async (configuration) => {
            await this.createRemediationConfiguration({ data: configuration });

            this.$popups.success({
              text: this.$t('modals.createRemediationConfiguration.duplicate.popups.success', {
                configurationName: remediationConfiguration.name,
              }),
            });

            await this.fetchList();
          },
        },
      });
    },

  },
};
</script>
