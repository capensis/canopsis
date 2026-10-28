import { generateShallowRenderer } from '@unit/utils/vue';

import { JOB_EXECUTION_STATUS, JOBS_TABS } from '@/constants';

import JobManagementTable from '@/components/other/job-management/job-management-table.vue';

describe('job-management-table', () => {
  const factory = generateShallowRenderer(JobManagementTable, {
    stubs: {
      'c-advanced-data-table': true,
      'job-table-actions': true,
    },
  });

  it('disables selection for an instruction with a pending stop command', () => {
    const wrapper = factory({
      propsData: {
        jobType: JOBS_TABS.instructions,
        headers: [],
        items: [],
        options: {},
      },
    });

    expect(wrapper.vm.isDisabledItem({
      status: JOB_EXECUTION_STATUS.running,
      stopped: true,
    })).toBe(true);
    expect(wrapper.vm.isDisabledItem({
      status: JOB_EXECUTION_STATUS.running,
      stopped: false,
    })).toBe(false);
  });
});
