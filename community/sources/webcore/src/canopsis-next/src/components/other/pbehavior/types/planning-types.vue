<template>
  <pbehavior-types-list
    :pbehavior-types="pbehaviorTypes"
    :pending="pbehaviorTypesPending"
    :total-items="pbehaviorTypesMeta.total_count"
    :options="options"
    :removable="hasDeleteAccess"
    :updatable="hasUpdateAccess"
    :active="active"
    @edit="showEditPbehaviorTypeModal"
    @remove="showRemovePbehaviorTypeModal"
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
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

import PbehaviorTypesList from './pbehavior-types-list.vue';

export default {
  components: {
    PbehaviorTypesList,
  },

  props: {
    active: {
      type: Boolean,
      default: true,
    },
  },

  setup() {
    const {
      pbehaviorTypes,
      pbehaviorTypesPending,
      pbehaviorTypesMeta,
      fetchPbehaviorTypesList,
      updatePbehaviorType,
      removePbehaviorType,
    } = usePbehaviorType();

    const { hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.planningType,
    );

    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();

    const {
      options,
      handler: fetchList,
      updateOptions,
    } = useLocalQueryWithOptions({
      onUpdate: fetchQuery => fetchPbehaviorTypesList({
          params: {
            ...convertQueryToRequest(fetchQuery),
            with_flags: true,
            with_hidden: true,
          },
        }),
    });

    const tryRemovePbehaviorType = async (id) => {
      try {
        await removePbehaviorType({ id });
      } catch (err) {
        popups.error({ text: err.error || t('errors.default') });
      }
    };

    const showRemovePbehaviorTypeModal = id => modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await tryRemovePbehaviorType(id);
            await fetchList();
          },
        },
      });

    const showEditPbehaviorTypeModal = pbehaviorType => modals.show({
        name: MODALS.createPbehaviorType,
        config: {
          pbehaviorType,
          action: async (data) => {
            await updatePbehaviorType({
              id: pbehaviorType._id,
              data,
            });
            await fetchList();
          },
        },
      });

    onMounted(fetchList);

    return {
      pbehaviorTypes,
      pbehaviorTypesPending,
      pbehaviorTypesMeta,
      options,
      hasUpdateAccess,
      hasDeleteAccess,
      fetchList,
      updateOptions,
      showEditPbehaviorTypeModal,
      showRemovePbehaviorTypeModal,
    };
  },
};
</script>
