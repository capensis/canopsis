import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { vi } from 'vitest';

import { JOB_EXECUTION_STATUS, JOBS_TABS } from '@/constants';

import JobTableActions from '@/components/other/job-management/partials/job-table-actions.vue';

const jobManagement = vi.hoisted(() => ({
  stopJobs: vi.fn(),
  fetchJobsListWithPreviousParams: vi.fn(),
}));

vi.mock('@/components/other/job-management/hooks/job-management', () => ({
  useJobManagement: () => jobManagement,
}));

vi.mock('@/helpers/async', () => ({
  promisedWait: vi.fn(() => Promise.resolve()),
}));

const CActionBtnStub = {
  emits: ['click'],
  template: '<button class="stop" @click="$emit(\'click\')" />',
};

describe('job-table-actions', () => {
  const popups = {
    error: vi.fn(),
    info: vi.fn(),
  };
  const factory = generateShallowRenderer(JobTableActions, {
    stubs: {
      'c-action-btn': CActionBtnStub,
    },
  });

  beforeEach(() => {
    globalThis.__TEST_POPUPS__ = popups;
    jobManagement.fetchJobsListWithPreviousParams.mockResolvedValue();
  });

  afterEach(() => {
    delete globalThis.__TEST_POPUPS__;
  });

  it('includes the stopped job name in the confirmation toast', async () => {
    const job = {
      _id: 'job-id',
      rule_name: 'Daily cleanup',
      status: JOB_EXECUTION_STATUS.running,
      pending_jobs: 0,
    };
    jobManagement.stopJobs.mockResolvedValue([{
      id: job._id,
      item: { _id: job._id },
      status: 200,
    }]);

    const wrapper = factory({
      propsData: {
        item: job,
        jobType: JOBS_TABS.instructions,
      },
    });

    await wrapper.find('.stop').trigger('click');
    await flushPromises();

    expect(wrapper.props('item')).toEqual(job);
    expect(jobManagement.stopJobs).toHaveBeenCalledWith({
      data: [{ _id: job._id, stop_pending_jobs: false }],
    });
    expect(popups.info).toHaveBeenCalledWith({
      text: 'Job <strong>Daily cleanup</strong> stopped',
    });
  });

  it('does not show the stop action when a stop command is pending', () => {
    const wrapper = factory({
      propsData: {
        item: {
          _id: 'stopped-job-id',
          status: JOB_EXECUTION_STATUS.running,
          stopped: true,
        },
        jobType: JOBS_TABS.instructions,
      },
    });

    expect(wrapper.find('.stop').exists()).toBe(false);
  });

});
