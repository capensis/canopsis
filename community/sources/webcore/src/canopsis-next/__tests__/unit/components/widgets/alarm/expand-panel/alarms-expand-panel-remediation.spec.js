import { fakeAlarm } from '@unit/data/alarm';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import AlarmsExpandPanelRemediation from '@/components/widgets/alarm/expand-panel/alarms-expand-panel-remediation.vue';

const { fetchExecutionsWithoutStore } = vi.hoisted(() => ({
  fetchExecutionsWithoutStore: vi.fn(),
}));

vi.mock('@/hooks/store/modules/alarm', () => ({
  useAlarm: () => ({
    fetchExecutionsWithoutStore,
  }),
}));

const stubs = {
  'remediation-instruction-executions-list': true,
};

const selectRemediationInstructionExecutionsList = wrapper =>
  wrapper.findComponent({ name: 'RemediationInstructionExecutionsList' });

describe('alarms-expand-panel-remediation', () => {
  const options = {
    page: 1,
    limit: 10,
  };
  const alarm = fakeAlarm();

  const factory = generateShallowRenderer(AlarmsExpandPanelRemediation, {
    stubs,
  });

  afterEach(() => {
    fetchExecutionsWithoutStore.mockReset();
  });

  test('Check fetchList calling on mounted', async () => {
    fetchExecutionsWithoutStore.mockResolvedValue({ data: [], meta: {} });

    factory({
      propsData: { alarm },
    });

    await flushPromises();

    expect(fetchExecutionsWithoutStore).toHaveBeenCalledWith({
      id: alarm._id,
      params: options,
    });
  });

  test('Check fetchList calling on options updating', async () => {
    fetchExecutionsWithoutStore.mockResolvedValue({ data: [], meta: {} });
    const newOptions = { ...options, itemsPerPage: PAGINATION_LIMIT, page: 2 };
    const wrapper = factory({
      propsData: { alarm },
    });

    const remediationInstructionExecutionsList =
      selectRemediationInstructionExecutionsList(wrapper);

    await remediationInstructionExecutionsList.triggerCustomEventForFirstChild(
      'update:options',
      newOptions,
    );

    expect(fetchExecutionsWithoutStore).toHaveBeenNthCalledWith(2, {
      id: alarm._id,
      params: omit(newOptions, ['itemsPerPage']),
    });
  });
});
