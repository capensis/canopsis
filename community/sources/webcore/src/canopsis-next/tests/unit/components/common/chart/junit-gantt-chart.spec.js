import { mount, createVueInstance } from '@unit/utils/vue';

import flushPromises from 'flush-promises';
import JunitGanttChart from '@/components/common/chart/junit-gantt-chart.vue';

const localVue = createVueInstance();

const stubs = {
  'c-table-pagination': {
    template: '<div class="c-table-pagination" />',
  },
};

const snapshotFactory = (options = {}) => mount(JunitGanttChart, {
  localVue,
  stubs,
  attachTo: document.body,

  ...options,
});

/**
 *
 *
 * junit-gantt-chart(
 :items="ganttIntervals",
 :historical="historical",
 :total-items="meta.total_count",
 :query.sync="query",
 :width="840"
 )
 *
 */

describe('junit-gantt-chart', () => {
  const historicalItems = [
    {
      _id: '5d72a420-ca3d-429f-a765-276af1d4cd55',
      name: 'test_case_9_1_0',
      time: 0.97,
      from: 0.00,
      to: 0.97,
      status: 0,
      message: '',
      avg_to: 0.97,
      avg_time: 0.97,
      avg_status: 0,
    },
    {
      _id: 'efe385e7-849a-49a8-aba9-ad8d5ff60b39',
      name: 'test_case_9_1_1',
      time: 0.77,
      from: 0.97,
      to: 1.73,
      status: 3,
      message: 'test_case_msg_9_1_1',
      avg_to: 1.73,
      avg_time: 0.77,
      avg_status: 3,
    },
    {
      _id: '671915bd-596c-4299-8565-efe9fc9d7376',
      name: 'test_case_9_1_2',
      time: 0.85,
      from: 1.73,
      to: 2.58,
      status: 2,
      message: 'test_case_msg_9_1_2',
      avg_to: 2.58,
      avg_time: 0.85,
      avg_status: 2,
    },
    {
      _id: 'ea7a3f46-7edd-4dbb-ba04-6429588552db',
      name: 'test_case_9_1_3',
      time: 0.91,
      from: 2.58,
      to: 3.49,
      status: 0,
      message: '',
      avg_to: 3.49,
      avg_time: 0.91,
      avg_status: 0,
    },
    {
      _id: 'b39bae6e-cf01-4706-804d-a3f5624d4c06',
      name: 'test_case_9_1_4',
      time: 0.89,
      from: 3.49,
      to: 4.38,
      status: 1,
      message: 'test_case_msg_9_1_4',
      avg_to: 3.49,
      avg_time: 0.00,
      avg_status: 1,
    },
  ];

  const items = [{
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
  }, {
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
  }, {
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
  }, {
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
  }, {
    _id: 'c3ec3087-c9f0-48a7-94d8-53a8dd0614ac',
    name: 'test_case_9_6_4',
    time: 0.91,
    from: 1.8,
    to: 2.71,
    status: 1,
    message: 'test_case_msg_9_6_4',
    avg_to: 0,
    avg_time: 0,
    avg_status: 0,
  }];

  const updatedItems = [{
    _id: '63c9388a-e734-4675-96f6-c71e01314910',
    name: 'test_case_9_6_5',
    time: 0.9,
    from: 2.71,
    to: 3.61,
    status: 0,
    message: '',
    avg_to: 0,
    avg_time: 0,
    avg_status: 0,
  }, {
    _id: '42463873-a0a9-469d-909a-766c6a4b5fcd',
    name: 'test_case_9_6_6',
    time: 0.93,
    from: 3.61,
    to: 4.54,
    status: 2,
    message: 'test_case_msg_9_6_6',
    avg_to: 0,
    avg_time: 0,
    avg_status: 0,
  }, {
    _id: '59d09586-bedd-4d83-a1b6-b50541bd042b',
    name: 'test_case_9_6_7',
    time: 0.19,
    from: 4.54,
    to: 4.73,
    status: 3,
    message: 'test_case_msg_9_6_7',
    avg_to: 0,
    avg_time: 0,
    avg_status: 0,
  }];

  it('Render `junit-gantt-chart` tooltip.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        items,
      },
    });

    const tooltip = {
      opacity: 1,
      caretX: 100,
      caretY: 100,
      dataPoints: [{ dataIndex: 0 }],
    };

    await flushPromises();

    wrapper.vm.getTooltip({ tooltip });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Render `junit-gantt-chart` tooltip with opacity 0.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        items,
      },
    });

    const tooltip = {
      opacity: 0,
      caretX: 100,
      caretY: 100,
      dataPoints: [{ dataIndex: 0 }],
    };

    await flushPromises();

    wrapper.vm.getTooltip({ tooltip });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Render `junit-gantt-chart` historical tooltip.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: historicalItems,
        historical: true,
      },
    });

    const tooltip = {
      opacity: 1,
      caretX: 100,
      caretY: 100,
      dataPoints: [{ dataIndex: 0 }],
    };

    await flushPromises();

    wrapper.vm.getTooltip({ tooltip });

    expect(wrapper.element).toMatchSnapshot(); // TODO: check icons
  });

  it('Renders `junit-gantt-chart` with default and required props.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        items,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `junit-gantt-chart` with historical data and prop.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        items: historicalItems,
        historical: true,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `junit-gantt-chart` with default, required props and updated items.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        items,
      },
    });

    await wrapper.setProps({ items: updatedItems });

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });
});
