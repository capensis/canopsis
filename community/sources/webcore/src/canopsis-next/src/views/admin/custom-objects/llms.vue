<template>
  <c-page
    :creatable="hasCreateAnyLlmAccess"
    :create-tooltip="$t('modals.createLlm.create.title')"
    @refresh="fetchList"
    @create="showCreateLlmModal"
  >
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

import { useCRUDPermissions } from '@/hooks/auth';
import { useModals } from '@/hooks/modals';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useLlm } from '@/hooks/store/modules/llm';

import LlmsList from '@/components/other/llm/llms-list.vue';

export default {
  components: { LlmsList },
  setup() {
    const modals = useModals();
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
     * Opens delete confirmation, removes one LLM by id, then refreshes the list.
     *
     * @param {Object} llm - LLM entity to delete (must include `_id` and `name`).
     */
    const showRemoveLlmModal = llm => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeLlm({ id: llm._id });

          return fetchList();
        },
      },
    });

    /**
     * Opens confirmation for bulk delete, removes selected LLMs, then refreshes the list.
     *
     * @param {Object[]} [selected=[]] - Selected rows; each item must include `_id`.
     */
    const showRemoveSelectedLlmsModal = (selected = []) => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await bulkRemoveLlms({ data: pickIds(selected) });

          return fetchList();
        },
      },
    });

    onMounted(fetchList);

    return {
      hasCreateAnyLlmAccess,
      hasUpdateAnyLlmAccess,
      hasDeleteAnyLlmAccess,

      llms,
      meta,
      pending,
      options,
      updateOptions,

      fetchList,
      showCreateLlmModal,
      showEditLlmModal,
      showRemoveLlmModal,
      showRemoveSelectedLlmsModal,
    };
  },
};
</script>
