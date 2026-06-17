<template>
  <v-container>
    <c-page @refresh="fetchAppInfo">
      <v-card-text>
        <user-interface :disabled="!hasUpdateParametersAccess" />
      </v-card-text>
    </c-page>
  </v-container>
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
    const { fetchAppInfo } = useInfo();

    return {
      hasUpdateParametersAccess,
      fetchAppInfo,
    };
  },
};
</script>
