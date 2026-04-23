<template>
  <c-page
    :creatable="hasCreateAnyLlmAccess"
    :create-tooltip="$t('modals.createLlm.create.title')"
    @refresh="refresh"
    @create="showCreateLlmModal"
  >
    <div class="pa-4 pb-2">
      <llms-important-notes-banner />
    </div>
    <llms-list
      :llms="llms"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      :updatable="hasUpdateAnyLlmAccess"
      :removable="hasDeleteAnyLlmAccess"
      @edit="showEditLlmModal"
      @remove="showRemoveLlmModal"
      @remove-selected="showRemoveSelectedLlmsModal"
    />
  </c-page>
</template>

<script>
import { onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { pickIds } from '@/helpers/array';

import { useCallActionWithPopup } from '@/hooks/actions/call';
import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useLlm } from '@/hooks/store/modules/llm';
import { useObserver } from '@/hooks/observer';

import LlmsImportantNotesBanner from '@/components/other/llm/partials/llms-important-notes-banner.vue';
import LlmsList from '@/components/other/llm/llms-list.vue';

export default {
  components: { LlmsImportantNotesBanner, LlmsList },
  setup() {
    const { t, tc } = useI18n();
    const modals = useModals();
    const { callActionWithPopup } = useCallActionWithPopup();
    const {
      fetchLlmsListWithoutStore,
      createLlm,
      updateLlm,
      removeLlm,
      bulkRemoveLlms,
    } = useLlm();
    const {
      hasCreateAccess: hasCreateAnyLlmAccess,
      hasUpdateAccess: hasUpdateAnyLlmAccess,
      hasDeleteAccess: hasDeleteAnyLlmAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.llm);

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
    const showCreateLlmModal = () => modals.show({
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
    const showEditLlmModal = llm => modals.show({
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
          action: () => callActionWithPopup(
            () => removeLlm({ id: llm._id }),
            fetchList,
          ),
        },
      });
    };

    /**
     * Opens confirmation phrase for bulk delete, removes selected LLMs, then refreshes the list.
     *
     * @param {Object[]} [selected=[]] - Selected rows; each item must include `_id`.
     */
    const showRemoveSelectedLlmsModal = (selected = []) => {
      const count = selected.length;
      const countValues = { count };

      modals.show({
        name: MODALS.confirmationPhrase,
        config: {
          title: tc('modals.confirmationPhrase.deleteSelectedLlms.title', count, countValues),
          text: tc('modals.confirmationPhrase.deleteSelectedLlms.text', count, countValues),
          phraseText: t('modals.confirmationPhrase.deleteSelectedLlms.phraseText'),
          phrase: t('modals.confirmationPhrase.deleteSelectedLlms.phrase'),
          action: () => callActionWithPopup(
            () => bulkRemoveLlms({ data: pickIds(selected) }),
            fetchList,
          ),
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
      hasCreateAnyLlmAccess,
      hasUpdateAnyLlmAccess,
      hasDeleteAnyLlmAccess,

      llms,
      meta,
      pending,
      options,
      updateOptions,

      refresh,
      showCreateLlmModal,
      showEditLlmModal,
      showRemoveLlmModal,
      showRemoveSelectedLlmsModal,
    };
  },
};
</script>
