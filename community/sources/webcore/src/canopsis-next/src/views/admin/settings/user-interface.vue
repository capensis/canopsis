<template>
  <c-page @refresh="fetchInfo">
    <template #header>
      {{ $t('common.userInterface') }}
    </template>
    <user-interface :disabled="!hasUpdateParametersAccess" />
  </c-page>
</template>

<script>
import { USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useInfo } from '@/hooks/store/modules/info';

import UserInterface from '@/components/other/user-interface/user-interface.vue';

export default {
  components: { UserInterface },
  setup() {
    const { hasUpdateAccess: hasUpdateParametersAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.parameters);
    const { fetchInfo } = useInfo();

    return {
      hasUpdateParametersAccess,
      fetchInfo,
    };
  },
};
</script>
