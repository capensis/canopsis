<template>
  <c-page
    :creatable="hasCreateAccess"
    :create-tooltip="t('modals.createLlm.create.title')"
    @refresh="refresh"
    @create="showCreateLlmModal"
  >
    <div class="pa-4 pb-2">
      <llms-important-notes-banner />
    </div>

    <llms-list
      :llms="llms"
      :pending="pending"
      :options="options"
      :total-items="meta.total_count"
      :updatable="hasUpdateAccess"
      :removable="hasDeleteAccess"
      @edit="showEditLlmModal"
      @remove="showRemoveLlmModal"
      @refresh="refresh"
      @update:options="updateOptions"
    />
  </c-page>
</template>

<script>
import { onMounted } from 'vue';

import LlmsList from '@/components/other/llm/llms-list.vue';
import LlmsImportantNotesBanner from '@/components/other/llm/partials/llms-important-notes-banner.vue';
import { MODALS, USER_PERMISSIONS } from '@/constants';
import { useCallActionWithPopup } from '@/hooks/actions/call';
import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useObserver } from '@/hooks/observer';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useLlm } from '@/hooks/store/modules/llm';

export default {
  components: {
    LlmsImportantNotesBanner,
    LlmsList,
  },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const { callActionWithPopup } = useCallActionWithPopup();
    const { fetchLlmsListWithoutStore, createLlm, updateLlm, removeLlm } = useLlm();
    const { hasCreateAccess, hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.llm,
    );

    const {
      data: llms,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: fetchLlmsListWithoutStore,
    });

    /**
     * Opens the create LLM modal and refreshes the list after a successful create.
     */
    const showCreateLlmModal = () =>
      modals.show({
        name: MODALS.createLlm,
        config: {
          action: async (data) => {
            await createLlm({ data });

            return fetchList();
          },
        },
      });

    /**
     * Opens the edit LLM modal prefilled with the given row and refreshes the list after a successful update.
     *
     * @param {Object} llm - LLM entity from the table (must include `_id`).
     */
    const showEditLlmModal = llm =>
      modals.show({
        name: MODALS.createLlm,
        config: {
          llm,
          action: async (data) => {
            await updateLlm({ id: llm._id, data });

            return fetchList();
          },
        },
      });

    /**
     * Opens delete confirmation (type LLM name), removes one LLM by id, then refreshes the list.
     *
     * @param {Object} llm - LLM entity to delete (must include `_id` and `name`).
     */
    const showRemoveLlmModal = (llm) => {
      modals.show({
        name: MODALS.confirmationPhrase,
        config: {
          title: t('modals.confirmationPhrase.deleteLlm.title'),
          text: t('modals.confirmationPhrase.deleteLlm.text'),
          phraseText: t('modals.confirmationPhrase.deleteLlm.phraseText'),
          phrase: llm.name,
          action: () => callActionWithPopup(() => removeLlm({ id: llm._id }), fetchList),
        },
      });
    };

    const { observer } = useObserver({ key: '$refresh' });

    const refresh = () => observer.notify();

    onMounted(() => {
      observer.register(fetchList);
      fetchList();
    });

    return {
      t,
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,
      llms,
      meta,
      pending,
      options,

      refresh,
      updateOptions,
      showCreateLlmModal,
      showEditLlmModal,
      showRemoveLlmModal,
    };
  },
};
</script>
