<template lang="pug">
  div
    c-page-header
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(v-model="activeTab", slider-color="primary", fixed-tabs)
            template(v-if="hasReadAnyRemediationInstructionAccess")
              v-tab(:href="`#${$constants.REMEDIATION_TABS.instructions}`") {{ $t('remediation.tabs.instructions') }}
              v-tab-item(:value="$constants.REMEDIATION_TABS.instructions", lazy)
                v-card-text
                  remediation-instructions
            template(v-if="hasReadAnyRemediationConfigurationAccess")
              v-tab(
                :href="`#${$constants.REMEDIATION_TABS.configurations}`"
              ) {{ $t('remediation.tabs.configurations') }}
              v-tab-item(:value="$constants.REMEDIATION_TABS.configurations", lazy)
                v-card-text
                  remediation-configurations
            template(v-if="hasReadAnyRemediationJobAccess")
              v-tab(:href="`#${$constants.REMEDIATION_TABS.jobs}`") {{ $t('remediation.tabs.jobs') }}
              v-tab-item(:value="$constants.REMEDIATION_TABS.jobs", lazy)
                v-card-text
                  remediation-jobs
    c-fab-btn(@create="create", @refresh="refresh", :has-access="hasCreateAccess")
      span {{ tooltipText }}
</template>

<script>
import { MODALS, REMEDIATION_TABS } from '@/constants';

import RemediationInstructions from '@/components/other/remediation/instructions/remediation-instructions.vue';
import RemediationJobs from '@/components/other/remediation/jobs/remediation-jobs.vue';
import RemediationConfigurations from '@/components/other/remediation/configurations/remediation-configurations.vue';

import { entitiesRemediationInstructionMixin } from '@/mixins/entities/remediation/instruction';
import { entitiesRemediationConfigurationMixin } from '@/mixins/entities/remediation/configuration';
import { entitiesRemediationJobMixin } from '@/mixins/entities/remediation/job';
import {
  permissionsTechnicalRemediationInstructionMixin,
} from '@/mixins/permissions/technical/remediation-instruction';
import {
  permissionsTechnicalRemediationConfigurationMixin,
} from '@/mixins/permissions/technical/remediation-configuration';
import { permissionsTechnicalRemediationJobMixin } from '@/mixins/permissions/technical/remediation-job';

export default {
  components: {
    RemediationInstructions,
    RemediationConfigurations,
    RemediationJobs,
  },
  mixins: [
    entitiesRemediationInstructionMixin,
    entitiesRemediationConfigurationMixin,
    entitiesRemediationJobMixin,
    permissionsTechnicalRemediationInstructionMixin,
    permissionsTechnicalRemediationConfigurationMixin,
    permissionsTechnicalRemediationJobMixin,
  ],
  data() {
    return {
      activeTab: REMEDIATION_TABS.instructions,
    };
  },
  computed: {
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
