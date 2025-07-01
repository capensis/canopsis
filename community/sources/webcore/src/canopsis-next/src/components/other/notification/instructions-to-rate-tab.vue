<template>
  <v-card-text>
    <div class="text-center pa-4">
      <h3 class="mb-4">
        {{ $t('notifications.tabs.instructionsToRate') }}
      </h3>
      <remediation-instruction-stats-list
        :remediation-instruction-stats="items"
        :pending="pending"
        :options="options"
        :total-items="meta.total_count"
        :accumulated-before="meta.accumulated_before"
        :interval="interval"
        @rate="showRateInstructionModal"
        @update:options="$emit('update:options', $event)"
      />
    </div>
  </v-card-text>
</template>

<script>
import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';

import RemediationInstructionStatsList from '@/components/other/remediation/instruction-stats/remediation-instruction-stats-list.vue';

export default {
  components: {
    RemediationInstructionStatsList,
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    meta: {
      type: Object,
      default: () => {},
    },
    options: {
      type: Object,
      default: () => {},
    },
    interval: {
      type: Object,
      default: () => {},
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();
    const { rateRemediationInstruction } = useRemdeitionInstruction();

    const refresh = () => emit('refresh');

    /**
     * @todo: MAY BE REFACTORED TO USE A HOOK BECAUSE WE HAVE A DOUBLICATED MODAL
     */
    const showRateInstructionModal = (instruction = {}) => modals.show({
      name: MODALS.rate,
      config: {
        title: t('modals.rateInstruction.title', { name: instruction.name }),
        text: t('modals.rateInstruction.text'),
        action: async (data) => {
          await rateRemediationInstruction({ id: instruction._id, data });

          return refresh();
        },
      },
    });

    return {
      showRateInstructionModal,
    };
  },
};
</script>
