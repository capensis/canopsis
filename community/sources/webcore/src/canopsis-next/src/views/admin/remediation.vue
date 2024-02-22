<template>
  <div>
    <c-page-header />
    <v-layout wrap>
      <v-flex xs12>
        <v-card class="ma-4 mt-0">
          <v-tabs
            v-model="activeTab"
            slider-color="primary"
            centered
          >
            <template v-if="hasReadAnyRemediationInstructionAccess">
              <v-tab :href="`#${$constants.REMEDIATION_TABS.instructions}`">
                {{ $t('common.instructions') }}
              </v-tab>
              <v-tab-item :value="$constants.REMEDIATION_TABS.instructions">
                <v-card-text>
                  <remediation-instructions />
                </v-card-text>
              </v-tab-item>
            </template>
            <template v-if="hasReadAnyRemediationConfigurationAccess">
              <v-tab :href="`#${$constants.REMEDIATION_TABS.configurations}`">
                {{ $t('remediation.tabs.configurations') }}
              </v-tab>
              <v-tab-item :value="$constants.REMEDIATION_TABS.configurations">
                <v-card-text>
                  <remediation-configurations />
                </v-card-text>
              </v-tab-item>
            </template>
            <template v-if="hasReadAnyRemediationJobAccess">
              <v-tab :href="`#${$constants.REMEDIATION_TABS.jobs}`">
                {{ $t('remediation.tabs.jobs') }}
              </v-tab>
              <v-tab-item :value="$constants.REMEDIATION_TABS.jobs">
                <v-card-text>
                  <remediation-jobs />
                </v-card-text>
              </v-tab-item>
            </template>
            <template v-if="hasReadRemediationStatisticAccess">
              <v-tab :href="`#${$constants.REMEDIATION_TABS.statistics}`">
                {{ $t('remediation.tabs.statistics') }}
              </v-tab>
              <v-tab-item :value="$constants.REMEDIATION_TABS.statistics">
                <v-card-text>
                  <remediation-statistics />
                </v-card-text>
              </v-tab-item>
            </template>
          </v-tabs>
        </v-card>
      </v-flex>
    </v-layout>
    <c-fab-btn
      :has-access="hasCreateAccess"
      v-on="fabListeners"
    >
      <span>{{ tooltipText }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { MODALS, REMEDIATION_TABS } from '@/constants';

import { entitiesRemediationInstructionMixin } from '@/mixins/entities/remediation/instruction';
import { entitiesRemediationConfigurationMixin } from '@/mixins/entities/remediation/configuration';
import { entitiesRemediationJobMixin } from '@/mixins/entities/remediation/job';
import { entitiesRemediationStatisticMixin } from '@/mixins/entities/remediation/statistic';
import {
  permissionsTechnicalRemediationInstructionMixin,
} from '@/mixins/permissions/technical/remediation-instruction';
import {
  permissionsTechnicalRemediationConfigurationMixin,
} from '@/mixins/permissions/technical/remediation-configuration';
import { permissionsTechnicalRemediationJobMixin } from '@/mixins/permissions/technical/remediation-job';
import { permissionsTechnicalRemediationStatisticMixin } from '@/mixins/permissions/technical/remediation-statistic';

import RemediationStatistics from '@/components/other/remediation/statistics/remediation-statistics.vue';
import RemediationJobs from '@/components/other/remediation/jobs/remediation-jobs.vue';
import RemediationConfigurations from '@/components/other/remediation/configurations/remediation-configurations.vue';
import RemediationInstructions from '@/components/other/remediation/instructions/remediation-instructions.vue';

export default {
  components: {
    RemediationInstructions,
    RemediationConfigurations,
    RemediationJobs,
    RemediationStatistics,
  },
  mixins: [
    entitiesRemediationInstructionMixin,
    entitiesRemediationConfigurationMixin,
    entitiesRemediationJobMixin,
    entitiesRemediationStatisticMixin,
    permissionsTechnicalRemediationInstructionMixin,
    permissionsTechnicalRemediationConfigurationMixin,
    permissionsTechnicalRemediationJobMixin,
    permissionsTechnicalRemediationStatisticMixin,
  ],
  data() {
    return {
      activeTab: REMEDIATION_TABS.instructions,
    };
  },
  computed: {
    fabListeners() {
      const listeners = {
        refresh: this.refresh,
      };

      if (this.activeTab !== REMEDIATION_TABS.statistics) {
        listeners.create = this.create;
      }

      return listeners;
    },

    tooltipText() {
      return {
        [REMEDIATION_TABS.instructions]: this.$t('modals.createRemediationInstruction.create.title'),
        [REMEDIATION_TABS.configurations]: this.$t('modals.createRemediationConfiguration.create.title'),
        [REMEDIATION_TABS.jobs]: this.$t('modals.createRemediationJob.create.title'),
      }[this.activeTab];
    },

    hasCreateAccess() {
      return {
        [REMEDIATION_TABS.instructions]: this.hasCreateAnyRemediationInstructionAccess,
        [REMEDIATION_TABS.configurations]: this.hasCreateAnyRemediationConfigurationAccess,
        [REMEDIATION_TABS.jobs]: this.hasCreateAnyRemediationJobAccess,
      }[this.activeTab];
    },
  },
  methods: {
    refresh() {
      switch (this.activeTab) {
        case REMEDIATION_TABS.instructions:
          this.fetchInstructionsList();
          break;
        case REMEDIATION_TABS.configurations:
          this.fetchConfigurationsList();
          break;
        case REMEDIATION_TABS.jobs:
          this.fetchJobsList();
          break;
        case REMEDIATION_TABS.statistics:
          this.fetchStatisticsList();
          break;
      }
    },

    create() {
      switch (this.activeTab) {
        case REMEDIATION_TABS.instructions:
          this.showCreateInstructionModal();
          break;
        case REMEDIATION_TABS.configurations:
          this.showCreateConfigurationModal();
          break;
        case REMEDIATION_TABS.jobs:
          this.showCreateJobModal();
          break;
      }
    },

    fetchInstructionsList() {
      this.fetchRemediationInstructionsListWithPreviousParams();
    },

    fetchConfigurationsList() {
      this.fetchRemediationConfigurationsListWithPreviousParams();
    },

    fetchJobsList() {
      this.fetchRemediationJobsListWithPreviousParams();
    },

    fetchStatisticsList() {
      this.fetchRemediationMetricsListWithPreviousParams();
    },

    showCreateInstructionModal() {
      this.$modals.show({
        name: MODALS.createRemediationInstruction,
        config: {
          action: async (instruction) => {
            await this.createRemediationInstruction({ data: instruction });

            this.$popups.success({
              text: this.$t('modals.createRemediationInstruction.create.popups.success', {
                instructionName: instruction.name,
              }),
            });

            await this.fetchInstructionsList();
          },
        },
      });
    },

    showCreateConfigurationModal() {
      this.$modals.show({
        name: MODALS.createRemediationConfiguration,
        config: {
          action: async (remediationConfiguration) => {
            await this.createRemediationConfiguration({ data: remediationConfiguration });

            this.$popups.success({
              text: this.$t('modals.createRemediationConfiguration.create.popups.success', {
                configurationName: remediationConfiguration.name,
              }),
            });

            await this.fetchConfigurationsList();
          },
        },
      });
    },

    showCreateJobModal() {
      this.$modals.show({
        name: MODALS.createRemediationJob,
        config: {
          action: async (remediationJob) => {
            await this.createRemediationJob({ data: remediationJob });

            this.$popups.success({
              text: this.$t('modals.createRemediationJob.create.popups.success', {
                jobName: remediationJob.name,
              }),
            });

            await this.fetchJobsList();
          },
        },
      });
    },
  },
};
</script>
