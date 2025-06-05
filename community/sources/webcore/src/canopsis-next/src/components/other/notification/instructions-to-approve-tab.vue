<template>
  <v-card-text>
    <div class="text-center pa-4">
      <h3 class="mb-4">
        {{ $t('notifications.tabs.instructionsToApprove') }}
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
            class="mr-2"
            color="success"
            small
            @click="showApproveModal(item)"
          >
            {{ $t('common.approve') }}
          </v-btn>
          <v-btn
            color="error"
            small
            @click="showDismissModal(item)"
          >
            {{ $t('common.dismiss') }}
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
    const { updateInstructionApproval } = useRemdeitionInstruction();

    // Computed
    const headers = computed(() => [
      { text: t('common.name'), value: 'name' },
      { text: t('common.description'), value: 'description' },
      { text: t('common.author'), value: 'author' },
      { text: t('common.actions'), value: 'actions', sortable: false },
    ]);

    const showApproveModal = (instruction) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('modals.approveInstruction.title'),
          text: t('modals.approveInstruction.text', { name: instruction.name }),
          action: async () => {
            await updateInstructionApproval({
              id: instruction._id,
              data: { approved: true },
            });
            emit('refresh');
          },
        },
      });
    };

    const showDismissModal = (instruction) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('modals.dismissInstruction.title'),
          text: t('modals.dismissInstruction.text', { name: instruction.name }),
          action: async () => {
            await updateInstructionApproval({
              id: instruction._id,
              data: { approved: false },
            });
            emit('refresh');
          },
        },
      });
    };

    return {
      headers,
      showApproveModal,
      showDismissModal,
    };
  },
};
</script>
