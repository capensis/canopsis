import { generateShallowRenderer } from '@unit/utils/vue';

import JunitGanttChart from '@/components/common/chart/junit-gantt-chart.vue';

const newItemsPerPage = 20;

const stubs = {
  'c-table-pagination': {
    props: ['page'],
    template: `
      <div class="c-table-pagination">
        <button class="next-page" @click="$emit('update:page', page + 1)"></button>
        <button class="items-per-page" @click="$emit('update:items-per-page', ${newItemsPerPage})"></button>
      </div>
    `,
  },

  'horizontal-bar': true,
};

describe('junit-gantt-chart', () => {
  const query = { page: 1, itemsPerPage: 10 };

  const items = [
    {
      _id: '766484c9-e952-49f7-bc70-12829bb7f3b7',
      name: 'test_case_9_6_0',
      time: 0.7,
      from: 0,
      to: 0.7,
      status: 2,
      message: 'test_case_msg_9_6_0',
      avg_to: 0,
      avg_time: 0,
      avg_status: 0,
    },
    {
      _id: '0b85f69c-d95e-4eea-8c44-6de3ccafb86b',
      name: 'test_case_9_6_1',
      time: 0.74,
      from: 0.7,
      to: 1.44,
      status: 3,
      message: 'test_case_msg_9_6_1',
      avg_to: 0,
      avg_time: 0,
      avg_status: 0,
    },
    {
      _id: '2822c0df-dffc-4a0b-8964-f2a2dc08b842',
      name: 'test_case_9_6_2',
      time: 0.29,
      from: 1.44,
      to: 1.73,
      status: 2,
      message: 'test_case_msg_9_6_2',
      avg_to: 0,
      avg_time: 0,
      avg_status: 0,
    },
    {
      _id: '45350b1f-ca1b-4c49-ad00-daec6d7bc731',
      name: 'test_case_9_6_3',
      time: 0.08,
      from: 1.73,
      to: 1.8,
      status: 1,
      message: 'test_case_msg_9_6_3',
      avg_to: 0,
      avg_time: 0,
      avg_status: 0,
    },
    {
      // Without time and message
      _id: 'c3ec3087-c9f0-48a7-94d8-53a8dd0614ac',
      name: 'test_case_9_6_4',
      from: 1.8,
      to: 1.8,
      status: 1,
      avg_to: 0,
      avg_time: 0,
      avg_status: 0,
    },
  ];

  const factory = generateShallowRenderer(JunitGanttChart, { stubs });

  it('Pagination next page event', async () => {
    const wrapper = factory({
      propsData: {
        items,
        query,
      },
    });

    await wrapper.find('.c-table-pagination .next-page').trigger('click');

    expect(wrapper).toEmit('update:query', { ...query, page: query.page + 1 });
  });

  it('Pagination rows per page event', async () => {
    const wrapper = factory({
      propsData: {
        items,
        query,
      },
    });

    await wrapper.find('.c-table-pagination .items-per-page').trigger('click');

    expect(wrapper).toEmit('update:query', { ...query, itemsPerPage: newItemsPerPage, page: 1 });
  });
});
