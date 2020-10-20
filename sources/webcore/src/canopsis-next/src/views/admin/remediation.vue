<template lang="pug">
  div
    the-page-header {{ $t('common.instructions') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(v-model="activeTab", fixed-tabs, slider-color="primary")
            template(v-if="hasReadAnyRemediationInstructionAccess")
              v-tab(:href="`#${$constants.REMEDIATION_TABS.instructions}`") {{ $t('remediation.tabs.instructions') }}
              v-tab-item(:value="$constants.REMEDIATION_TABS.instructions")
                v-card-text
                  remediation-instructions
            template
              v-tab(
                :href="`#${$constants.REMEDIATION_TABS.configurations}`"
              ) {{ $t('remediation.tabs.configurations') }}
              v-tab-item(:value="$constants.REMEDIATION_TABS.configurations")
                v-card-text
                  span {{ $t('remediation.tabs.configurations') }}
            template(v-if="hasReadAnyRemediationJobAccess")
              v-tab(:href="`#${$constants.REMEDIATION_TABS.jobs}`") {{ $t('remediation.tabs.jobs') }}
              v-tab-item(:value="$constants.REMEDIATION_TABS.jobs")
                v-card-text
                  remediation-jobs
    fab-buttons(@create="create", @refresh="refresh", :has-access="hasCreateAccess")
      span {{ tooltipText }}
</template>

<script>
import { MODALS, REMEDIATION_TABS } from '@/constants';

import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';
import RemediationInstructions from '@/components/other/remediation/instructions/remediation-instructions.vue';
import RemediationJobs from '@/components/other/remediation/jobs/remediation-jobs.vue';

import entitiesRemediationInstructionMixin from '@/mixins/entities/remediation/instruction';
import entitiesRemediationJobMixin from '@/mixins/entities/remediation/jobs';
import rightsTechnicalRemediationInstructionMixin from '@/mixins/rights/technical/remediation-instruction';
import rightsTechnicalRemediationJobMixin from '@/mixins/rights/technical/remediation-job';

export default {
  components: {
    RemediationInstructions,
    RemediationJobs,
    FabButtons,
  },
  mixins: [
    entitiesRemediationInstructionMixin,
    entitiesRemediationJobMixin,
    rightsTechnicalRemediationInstructionMixin,
    rightsTechnicalRemediationJobMixin,
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
        [REMEDIATION_TABS.configurations]: this.$t('modals.createRemediationConfiguration.title'),
        [REMEDIATION_TABS.jobs]: this.$t('modals.createRemediationJob.title'),
      }[this.activeTab];
    },

    hasCreateAccess() {
      return {
        [REMEDIATION_TABS.instructions]: this.hasCreateAnyRemediationInstructionAccess,
        [REMEDIATION_TABS.configurations]: true,
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
    },

    showCreateJobModal() {
    },
  },
};
</script>
