import { generateShallowRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import PeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';

const stubs = {
  'widget-settings-item': true,
  'periodic-refresh-field': {
    name: 'PeriodicRefreshFieldStub',
    emits: ['update:periodicRefresh'],
    template: '<div class="periodic-refresh-field-stub" />',
  },
  'live-watching-field': true,
};

const selectPeriodicRefreshField = wrapper =>
  wrapper.findComponent({ name: 'PeriodicRefreshFieldStub' });
const selectLiveWatchingField = wrapper => wrapper.find('live-watching-field-stub');

describe('periodic-refresh', () => {
  const factory = generateShallowRenderer(PeriodicRefresh, { stubs });

  it('Value changed after trigger periodic refresh field', async () => {
    const wrapper = factory({
      propsData: {
        value: {
          value: 1,
          unit: TIME_UNITS.day,
        },
      },
    });

    const newValue = {
      value: 2,
      unit: TIME_UNITS.week,
    };

    selectPeriodicRefreshField(wrapper).vm.$emit('update:periodicRefresh', newValue);
    await wrapper.vm.$nextTick();

    expect(wrapper).toEmitInput(newValue);
  });

  it('Live watching triggers input on changes', () => {
    const periodicRefresh = {
      value: 1,
      unit: TIME_UNITS.second,
    };

    const wrapper = factory({
      propsData: {
        value: periodicRefresh,
        withLiveWatching: true,
      },
    });

    const newLiveWatching = true;

    selectLiveWatchingField(wrapper).triggerCustomEvent('input', newLiveWatching);

    expect(wrapper).toEmit('update:liveWatching', newLiveWatching);
  });
});
