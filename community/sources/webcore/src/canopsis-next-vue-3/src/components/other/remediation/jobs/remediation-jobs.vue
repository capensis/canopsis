<template>
  <v-card-text>
    <remediation-jobs-list
      :remediation-jobs="remediationJobs"
      :pending="remediationJobsPending"
      :total-items="remediationJobsMeta.total_count"
      :options="options"
      :removable="hasDeleteAnyRemediationJobAccess"
      :updatable="hasUpdateAnyRemediationJobAccess"
      :duplicable="hasCreateAnyRemediationJobAccess"
      :active="active"
      @remove="showRemoveRemediationJobModal"
      @duplicate="showDuplicateRemediationJobModal"
      @edit="showEditRemediationJobModal"
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
import { useRemediationJob } from '@/hooks/store/modules/remediation-job';

import { useRemediationJobsActions } from './hooks/remediation-jobs';
import RemediationJobsList from './remediation-jobs-list.vue';

export default {
  components: {
    RemediationJobsList,
  },

  props: {
    active: {
      type: Boolean,
      default: true,
    },
  },

  setup() {
    const {
      remediationJobs,
      remediationJobsMeta,
      remediationJobsPending,
      fetchRemediationJobsList,
    } = useRemediationJob();

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useFetchListWithOptions({
      fetchListHandler: ({ params }) =>
        fetchRemediationJobsList({
          params: {
            ...params,
            with_flags: true,
          },
        }),
    });

    const {
      hasCreateAccess: hasCreateAnyRemediationJobAccess,
      hasUpdateAccess: hasUpdateAnyRemediationJobAccess,
      hasDeleteAccess: hasDeleteAnyRemediationJobAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.exploitation.remediationJob);

    const {
      showEditRemediationJobModal,
      showDuplicateRemediationJobModal,
      showRemoveRemediationJobModal,
    } = useRemediationJobsActions(fetchList);

    onMounted(fetchList);

    return {
      remediationJobs,
      remediationJobsMeta,
      remediationJobsPending,
      options,
      updateOptions,
      hasCreateAnyRemediationJobAccess,
      hasUpdateAnyRemediationJobAccess,
      hasDeleteAnyRemediationJobAccess,
      showEditRemediationJobModal,
      showDuplicateRemediationJobModal,
      showRemoveRemediationJobModal,
      fetchList,
    };
  },
};
</script>
