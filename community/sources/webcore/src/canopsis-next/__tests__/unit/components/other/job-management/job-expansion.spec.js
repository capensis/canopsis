import { mount, flushPromises } from '@vue/test-utils';
import { createPinia } from 'pinia';
import { h } from 'vue';

import { JOB_EXECUTION_STATUS, JOB_RULE_TYPES, JOBS_TABS } from '@/constants';

import { vuetify } from '@/plugins/vuetify';

import { i18n } from '@/i18n';

import CCollapsePanel from '@/components/common/block/c-collapse-panel.vue';
import CExpandPanelCard from '@/components/common/block/c-expand-panel-card.vue';
import CWindowExpandItem from '@/components/common/block/c-window-expand-item.vue';
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import CEllipsis from '@/components/common/table/c-ellipsis.vue';
import CTableExpandRow from '@/components/common/table/c-table-expand-row.vue';
import CObjectTreeview from '@/components/common/treeview/c-object-treeview.vue';
import CTreeview from '@/components/common/treeview/c-treeview.vue';
import JobDetailsExpandPanel from '@/components/other/job-management/partials/job-details-expand-panel.vue';

describe('job expansion', () => {
  it('opens and closes an event-filter webhook with a JSON response', async () => {
    const item = {
      _id: 'event-filter-job',
      rule_type: JOB_RULE_TYPES.eventFilter,
      status: JOB_EXECUTION_STATUS.succeeded,
      raw_request: 'POST /hook HTTP/1.1\r\nHost: example.test\r\n\r\n{}',
      raw_response: 'HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{"result":{"ok":true}}',
    };
    const wrapper = mount(CAdvancedDataTable, {
      props: {
        headers: [{ value: 'rule_type', text: 'Rule type' }],
        items: [item],
        options: { page: 1, itemsPerPage: 10 },
        expand: true,
        selectAll: true,
      },
      slots: {
        expand: ({ item: job }) => h(JobDetailsExpandPanel, { item: job, jobType: JOBS_TABS.webhooks }),
      },
      global: {
        plugins: [createPinia(), i18n, vuetify],
        components: {
          CCollapsePanel,
          CExpandPanelCard,
          CWindowExpandItem,
          CTableExpandRow,
          CObjectTreeview,
          CTreeview,
          CEllipsis,
        },
        stubs: {
          'c-icon': { template: '<button><slot /></button>' },
          'c-information-block-row': true,
          'c-copy-btn': true,
        },
      },
    });

    await wrapper.get('.v-data-table__expand-icon').trigger('click');
    await flushPromises();

    expect(wrapper.text()).toContain('ok');
    expect(wrapper.text()).toContain('true');

    await wrapper.get('.v-data-table__expand-icon').trigger('click');
    await new Promise(resolve => setTimeout(resolve, 300));
    await flushPromises();

    expect(wrapper.findComponent(JobDetailsExpandPanel).exists()).toBe(false);
    wrapper.unmount();
  });
});
