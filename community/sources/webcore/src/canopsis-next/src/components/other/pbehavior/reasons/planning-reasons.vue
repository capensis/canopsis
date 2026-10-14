<template>
  <pbehavior-reasons-list
    :pbehavior-reasons="pbehaviorReasons"
    :pending="pbehaviorReasonsPending"
    :total-items="pbehaviorReasonsMeta.total_count"
    :options="options"
    :removable="hasDeleteAccess"
    :updatable="hasUpdateAccess"
    :active="active"
    @edit="showEditPbehaviorReasonModal"
    @remove="showRemovePbehaviorReasonModal"
    @refresh="fetchList"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason';

import PbehaviorReasonsList from './pbehavior-reasons-list.vue';

export default {
  components: {
    PbehaviorReasonsList,
  },

  props: {
    active: {
      type: Boolean,
      default: true,
    },
  },

  setup() {
    const {
      pbehaviorReasons,
      pbehaviorReasonsPending,
      pbehaviorReasonsMeta,
      fetchPbehaviorReasonsList,
      updatePbehaviorReason,
      removePbehaviorReason,
    } = usePbehaviorReason();

    const { hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.planningReason,
    );

    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();

    const {
      options,
      handler: fetchList,
      updateOptions,
    } = useLocalQueryWithOptions({
      onUpdate: fetchQuery => fetchPbehaviorReasonsList({
          params: {
            ...convertQueryToRequest(fetchQuery),
            with_flags: true,
            with_hidden: true,
          },
        }),
    });

    const tryRemovePbehaviorReason = async (id) => {
      try {
        await removePbehaviorReason({ id });
      } catch (err) {
        popups.error({ text: err.error || t('errors.default') });
      }
    };

    const showRemovePbehaviorReasonModal = id => modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await tryRemovePbehaviorReason(id);
            await fetchList();
          },
        },
      });

    const showEditPbehaviorReasonModal = pbehaviorReason => modals.show({
        name: MODALS.createPbehaviorReason,
        config: {
          pbehaviorReason,
          action: async (data) => {
            await updatePbehaviorReason({
              id: pbehaviorReason._id,
              data,
            });
            await fetchList();
          },
        },
      });

    onMounted(fetchList);

    return {
      pbehaviorReasons,
      pbehaviorReasonsPending,
      pbehaviorReasonsMeta,
      options,
      hasUpdateAccess,
      hasDeleteAccess,
      fetchList,
      updateOptions,
      showEditPbehaviorReasonModal,
      showRemovePbehaviorReasonModal,
    };
  },
};
</script>
