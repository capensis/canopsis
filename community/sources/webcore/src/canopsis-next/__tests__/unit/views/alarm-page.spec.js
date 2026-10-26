import { generateShallowRenderer } from '@unit/utils/vue';

import AlarmPage from '@/views/alarm-page.vue';

const {
  alarm,
  fetchAlarmItem,
  fetchWidgetWithoutStore,
  getAlarmItem,
  route,
} = vi.hoisted(() => ({
  alarm: { _id: 'alarm-id' },
  fetchAlarmItem: vi.fn().mockResolvedValue(),
  fetchWidgetWithoutStore: vi.fn().mockResolvedValue(),
  getAlarmItem: vi.fn(),
  route: { query: {} },
}));

vi.mock('vue-router', () => ({
  useRoute: () => route,
}));

vi.mock('@/hooks/store/modules/alarm', () => ({
  useAlarm: () => ({
    getItem: getAlarmItem,
    fetchItem: fetchAlarmItem,
  }),
}));

vi.mock('@/hooks/store/modules/widget', () => ({
  useWidget: () => ({
    fetchWidgetWithoutStore,
  }),
}));

const stubs = {
  'alarms-list-table': true,
};

describe('alarm-page', () => {
  const factory = generateShallowRenderer(AlarmPage, { stubs });

  beforeEach(() => {
    getAlarmItem.mockReturnValue(alarm);
  });

  test('shows an alarm returned by the store getter', () => {
    const wrapper = factory({
      propsData: {
        id: alarm._id,
      },
    });

    expect(getAlarmItem).toHaveBeenCalledWith(alarm._id);
    expect(wrapper.vm.alarmItems).toEqual([alarm]);
  });
});
