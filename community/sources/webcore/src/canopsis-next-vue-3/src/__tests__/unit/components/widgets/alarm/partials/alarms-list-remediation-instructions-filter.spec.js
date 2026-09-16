import { createMockedStoreModules, createRemediationInstructionModule } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';

import AlarmsListRemediationInstructionsFilter from '@/components/widgets/alarm/partials/alarms-list-remediation-instructions-filter.vue';

const stubs = {
  'c-select-field': true,
  'alarms-list-remediation-instructions-filter-chip': true,
  'alarms-list-remediation-instructions-filter-fields': true,
};

describe('alarms-list-remediation-instructions-filter', () => {
  const updateFilter = vi.fn();

  const { remediationInstructionModule, fetchRemediationInstructionsListWithoutStore } =
    createRemediationInstructionModule();

  const store = createMockedStoreModules([remediationInstructionModule]);
  const factory = generateShallowRenderer(AlarmsListRemediationInstructionsFilter, {
    stubs,
    store,
    listeners: {
      input: updateFilter,
    },
  });

  afterEach(() => {
    updateFilter.mockReset();
    fetchRemediationInstructionsListWithoutStore.mockReset();
  });

  it('Emits input with empty object when cleared', async () => {
    fetchRemediationInstructionsListWithoutStore.mockResolvedValue({
      data: [
        { _id: 'id1', name: 'name1' },
        { _id: 'id2', name: 'name2' },
      ],
      meta: {
        total_count: 2,
      },
    });

    const wrapper = factory({
      propsData: {
        filter: {
          instruction_filter_type: 2,
          instruction_type: 1,
          instruction_statuses: [0, 1],
          instruction_ids: ['id1', 'id2'],
        },
      },
    });

    // Simulate clearing the filter (select none)
    await wrapper.vm.updateFilter([]);
    expect(updateFilter).toHaveBeenCalledWith({});
  });

  it('Does not emit input if value is not empty', async () => {
    fetchRemediationInstructionsListWithoutStore.mockResolvedValue({
      data: [
        { _id: 'id1', name: 'name1' },
        { _id: 'id2', name: 'name2' },
      ],
      meta: {
        total_count: 2,
      },
    });

    const wrapper = factory({
      propsData: {
        filter: {
          instruction_filter_type: 2,
        },
      },
    });
    await wrapper.vm.updateFilter(['instruction_filter_type']);
    expect(updateFilter).not.toHaveBeenCalled();
  });

  it('Computes selectedItems and selectedAllFields correctly', async () => {
    fetchRemediationInstructionsListWithoutStore.mockResolvedValue({
      data: [
        { _id: 'id1', name: 'name1' },
        { _id: 'id2', name: 'name2' },
      ],
      meta: {
        total_count: 2,
      },
    });

    const wrapper = factory({
      propsData: {
        filter: {
          instruction_filter_type: 2,
          instruction_type: 1,
          instruction_statuses: [0, 1],
          instruction_ids: ['id1', 'id2'],
        },
      },
    });
    expect(wrapper.vm.selectedItems).toEqual([
      'instruction_filter_type',
      'instruction_type',
      'instruction_statuses',
      'instruction_ids',
    ]);
    expect(wrapper.vm.selectedAllFields).toBe(true);
  });
});
