<template>
  <pbehavior-exceptions-list
    :pbehavior-exceptions="pbehaviorExceptions"
    :pending="pbehaviorExceptionsPending"
    :total-items="pbehaviorExceptionsMeta.total_count"
    :options="options"
    :removable="hasDeleteAccess"
    :updatable="hasUpdateAccess"
    :active="active"
    @edit="showEditPbehaviorExceptionModal"
    @remove="showRemovePbehaviorExceptionModal"
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
import { usePbehaviorException } from '@/hooks/store/modules/pbehavior-exception';

import PbehaviorExceptionsList from './pbehavior-exceptions-list.vue';

export default {
  components: {
    PbehaviorExceptionsList,
  },

  props: {
    active: {
      type: Boolean,
      default: true,
    },
  },

  setup() {
    const {
      pbehaviorExceptions,
      pbehaviorExceptionsPending,
      pbehaviorExceptionsMeta,
      fetchPbehaviorExceptionsList,
      updatePbehaviorException,
      removePbehaviorException,
    } = usePbehaviorException();

    const { hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.planningExceptions,
    );

    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();

    const {
      options,
      handler: fetchList,
      updateOptions,
    } = useLocalQueryWithOptions({
      onUpdate: fetchQuery => fetchPbehaviorExceptionsList({
          params: {
            ...convertQueryToRequest(fetchQuery),
            with_flags: true,
            with_hidden: true,
          },
        }),
    });

    const tryRemovePbehaviorException = async (id) => {
      try {
        await removePbehaviorException({ id });
      } catch (err) {
        popups.error({ text: err.error || t('errors.default') });
      }
    };

    const showRemovePbehaviorExceptionModal = id => modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await tryRemovePbehaviorException(id);
            await fetchList();
          },
        },
      });

    const showEditPbehaviorExceptionModal = pbehaviorException => modals.show({
        name: MODALS.createPbehaviorException,
        config: {
          pbehaviorException,
          action: async (data) => {
            await updatePbehaviorException({
              id: pbehaviorException._id,
              data,
            });
            await fetchList();
          },
        },
      });

    onMounted(fetchList);

    return {
      pbehaviorExceptions,
      pbehaviorExceptionsPending,
      pbehaviorExceptionsMeta,
      options,
      hasUpdateAccess,
      hasDeleteAccess,
      fetchList,
      updateOptions,
      showEditPbehaviorExceptionModal,
      showRemovePbehaviorExceptionModal,
    };
  },
};
</script>
