import { omit } from 'lodash';

import { flushPromises, generateShallowRenderer } from '@unit/utils/vue';
import { createAlarmModule, createMockedStoreModules } from '@unit/utils/store';
import { fakeAlarm } from '@unit/data/alarm';

import { PAGINATION_LIMIT } from '@/config';

import AlarmsExpandPanelRemediation from '@/components/widgets/alarm/expand-panel/alarms-expand-panel-remediation.vue';

const stubs = {
  'remediation-instruction-executions-list': true,
};

const selectRemediationInstructionExecutionsList = wrapper => (
  wrapper.find('remediation-instruction-executions-list-stub')
);

describe('alarms-expand-panel-remediation', () => {
  const options = {
    page: 1,
    limit: 10,
  };
  const alarm = fakeAlarm();
  const {
    alarmModule,
    fetchExecutionsWithoutStore,
  } = createAlarmModule();

  const store = createMockedStoreModules([alarmModule]);
  const factory = generateShallowRenderer(AlarmsExpandPanelRemediation, {
    stubs,
    store,
  });

  test('Check fetchList calling on mounted', async () => {
    fetchExecutionsWithoutStore.mockResolvedValue({ data: [], meta: {} });

    factory({
      propsData: { alarm },
    });

    await flushPromises();

    expect(fetchExecutionsWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: alarm._id,
        params: options,
      },
      undefined,
    );
  });

  test('Check fetchList calling on options updating', async () => {
    fetchExecutionsWithoutStore.mockResolvedValue({ data: [], meta: {} });
    const newOptions = { ...options, itemsPerPage: PAGINATION_LIMIT, page: 2 };
    const wrapper = factory({
      propsData: { alarm },
    });

    const remediationInstructionExecutionsList = selectRemediationInstructionExecutionsList(wrapper);

    await remediationInstructionExecutionsList.triggerCustomEventForFirstChild('update:options', newOptions);

    expect(fetchExecutionsWithoutStore).toHaveBeenNthCalledWith(
      2,
      expect.any(Object),
      {
        id: alarm._id,
        params: omit(newOptions, ['itemsPerPage']),
      },
      undefined,
    );
  });
});
