import { generateShallowRenderer } from '@unit/utils/vue';

import KpiRatingSettings from '@/components/other/kpi/rating-settings/kpi-rating-settings.vue';

const mocks = vi.hoisted(() => ({
  fetchRatingSettingsList: vi.fn(),
  bulkUpdateRatingSettings: vi.fn(),
}));

vi.mock('@/hooks/store/modules/rating-settings', () => ({
  useRatingSettings: () => ({
    ratingSettings: [],
    ratingSettingsMeta: { total_count: 0 },
    ratingSettingsPending: false,
    ...mocks,
  }),
}));

const stubs = { 'kpi-rating-settings-list': true };

describe('kpi-rating-settings', () => {
  const factory = generateShallowRenderer(KpiRatingSettings, { stubs });

  it('fetches rating settings with the initial pagination options', async () => {
    const wrapper = factory();
    await wrapper.vm.$nextTick();

    expect(mocks.fetchRatingSettingsList).toHaveBeenCalled();
    expect(wrapper.vm.options).toMatchObject({ page: 1, itemsPerPage: 10 });
  });

  it('exposes the rating settings query and selection handlers', () => {
    const wrapper = factory();

    expect(wrapper.vm.fetchList).toEqual(expect.any(Function));
    expect(wrapper.vm.updateOptions).toEqual(expect.any(Function));
    expect(wrapper.vm.changeSelectedRatingSettings).toEqual(expect.any(Function));
  });
});
