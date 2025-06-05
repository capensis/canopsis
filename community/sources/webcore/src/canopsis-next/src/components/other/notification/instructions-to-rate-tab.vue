<template>
  <v-card-text>
    <div class="text-center pa-4">
      <h3 class="mb-4">
        {{ $t('notifications.tabs.instructionsToRate') }}
      </h3>
      <c-advanced-data-table
        :headers="headers"
        :items="items"
        :loading="pending"
        :total-items="meta.total_count"
        :options="options"
        advanced-pagination
        @update:options="$emit('update:options', $event)"
      >
        <template #actions="{ item }">
          <v-btn
            color="primary"
            small
            @click="showRateModal(item)"
          >
            {{ $t('common.rate') }}
          </v-btn>
        </template>
      </c-advanced-data-table>
    </div>
  </v-card-text>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';

export default {
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
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();
    const { rateInstruction } = useRemdeitionInstruction();

    // Computed
    const headers = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('remediation.instructionStat.lastExecutedOn'), value: 'last_executed_on' },
      { text: t('common.actions'), value: 'actions', sortable: false },
    ]);

    const showRateModal = (instruction) => {
      modals.show({
        name: MODALS.rate,
        config: {
          title: t('modals.rateInstruction.title', { name: instruction.name }),
          text: t('modals.rateInstruction.text'),
          action: async (data) => {
            await rateInstruction({ id: instruction._id, data });
            emit('refresh');
          },
        },
      });
    };

    return {
      headers,
      showRateModal,
    };
  },
};
</script>
