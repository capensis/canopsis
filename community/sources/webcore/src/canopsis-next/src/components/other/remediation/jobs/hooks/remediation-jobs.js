import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useRemediationJob } from '@/hooks/store/modules/remediation-job';

/**
 * Provides modal actions for remediation jobs.
 *
 * @param {Function} [refresh=() => {}] - Callback to refresh data after actions.
 * @returns {Object}
 */
export const useRemediationJobsActions = (refresh = () => {}) => {
  const { t } = useI18n();
  const modals = useModals();
  const popups = usePopups();

  const { createRemediationJob, updateRemediationJob, removeRemediationJob } = useRemediationJob();

  const showEditRemediationJobModal = remediationJob =>
    modals.show({
      name: MODALS.createRemediationJob,
      config: {
        remediationJob,
        title: t('modals.createRemediationJob.edit.title'),
        action: async (job) => {
          await updateRemediationJob({ id: remediationJob._id, data: job });

          popups.success({
            text: t('modals.createRemediationJob.edit.popups.success', {
              jobName: job.name,
            }),
          });

          await refresh();
        },
      },
    });

  const showDuplicateRemediationJobModal = remediationJob =>
    modals.show({
      name: MODALS.createRemediationJob,
      config: {
        remediationJob: omit(remediationJob, ['_id']),
        title: t('modals.createRemediationJob.duplicate.title'),
        action: async (job) => {
          const newRemediationJob = await createRemediationJob({ data: job });

          popups.success({
            text: t('modals.createRemediationJob.duplicate.popups.success', {
              jobName: remediationJob.name,
            }),
          });

          await refresh();

          return newRemediationJob;
        },
      },
    });

  const showRemoveRemediationJobModal = remediationJob =>
    modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeRemediationJob({ id: remediationJob._id });
          await refresh();
        },
      },
    });

  return {
    showEditRemediationJobModal,
    showDuplicateRemediationJobModal,
    showRemoveRemediationJobModal,
  };
};
