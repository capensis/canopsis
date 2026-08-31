<template>
  <v-card-text>
    <remediation-configurations-list
      :remediation-configurations="remediationConfigurations"
      :pending="remediationConfigurationsPending"
      :total-items="remediationConfigurationsMeta.total_count"
      :options="options"
      :removable="hasDeleteAnyRemediationConfigurationAccess"
      :duplicable="hasCreateAnyRemediationConfigurationAccess"
      :updatable="hasUpdateAnyRemediationConfigurationAccess"
      :active="active"
      @remove="showRemoveRemediationConfigurationModal"
      @duplicate="showDuplicateRemediationConfigurationModal"
      @edit="showEditRemediationConfigurationModal"
      @refresh="fetchList"
      @update:options="updateOptions"
    />
  </v-card-text>
</template>

<script>
import { onMounted } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useFetchListWithOptions } from '@/hooks/query/shared';
import { useRemediationConfiguration } from '@/hooks/store/modules/remediation-configuration';

import { useRemediationConfigurationsActions } from './hooks/remediation-configurations';
import RemediationConfigurationsList from './remediation-configurations-list.vue';

export default {
  components: {
    RemediationConfigurationsList,
  },

  props: {
    active: {
      type: Boolean,
      default: true,
    },
  },

  setup() {
    const {
      remediationConfigurations,
      remediationConfigurationsMeta,
      remediationConfigurationsPending,
      fetchRemediationConfigurationsList,
    } = useRemediationConfiguration();

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useFetchListWithOptions({
      fetchListHandler: ({ params }) => fetchRemediationConfigurationsList({
          params: {
            ...params,
            with_flags: true,
          },
        }),
    });

    const {
      hasCreateAccess: hasCreateAnyRemediationConfigurationAccess,
      hasUpdateAccess: hasUpdateAnyRemediationConfigurationAccess,
      hasDeleteAccess: hasDeleteAnyRemediationConfigurationAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.exploitation.remediationConfiguration);

    const {
      showEditRemediationConfigurationModal,
      showRemoveRemediationConfigurationModal,
      showDuplicateRemediationConfigurationModal,
    } = useRemediationConfigurationsActions(fetchList);

    onMounted(fetchList);

    return {
      remediationConfigurations,
      remediationConfigurationsMeta,
      remediationConfigurationsPending,
      options,
      updateOptions,
      hasCreateAnyRemediationConfigurationAccess,
      hasUpdateAnyRemediationConfigurationAccess,
      hasDeleteAnyRemediationConfigurationAccess,
      showEditRemediationConfigurationModal,
      showRemoveRemediationConfigurationModal,
      showDuplicateRemediationConfigurationModal,
      fetchList,
    };
  },
};
</script>
