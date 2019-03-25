import Vuetify from 'vuetify';
import { shallowMount, createLocalVue } from '@vue/test-utils';

import {
  WEATHER_ICONS,
  WATCHER_STATES_COLORS,
  ENTITIES_STATES,
  WATCHER_PBEHAVIOR_COLOR,
  PBEHAVIOR_TYPES,
} from '@/constants';

import WeatherItem from '@/components/other/service-weather/weather-item.vue';

const localVue = createLocalVue();

localVue.use(Vuetify);

const componentProps = {
  isEditingMode: false,
  template: '<h1>{{ entity.display_name }}</h1>',
  watcher: {
    action_required: false,
    active_pb_all: false,
    active_pb_some: false,
    active_pb_watcher: false,
    alerts_not_ack: false,
    display_name: 'Watcher',
    entity_id: 'watcher',
    infos: {},
    pbehavior: [],
    watcher_pbehavior: [],
    state: {
      val: 0,
    },
  },
  widget: {
    _id: 'widget_weather',
    parameters: {
      alarmsList: {},
      blockTemplate: '',
      columnLG: 3,
      columnMD: 4,
      columnSM: 6,
      entityTemplate: '',
      heightFactor: 4,
      margin: {
        bottom: 1,
        left: 1,
        right: 1,
        top: 1,
      },
      mfilter: {
        filter: '',
        title: '',
      },
    },
    title: 'Weather title',
    type: 'ServiceWeather',
  },
};

describe('Weather Item', () => {
  let wrapper;

  beforeEach(() => {
    wrapper = shallowMount(
      WeatherItem,
      {
        localVue,
        propsData: componentProps,
      },
    );
  });

  it('has the right color', () => {
    expect(wrapper.vm.format.color).toBe(WATCHER_STATES_COLORS[ENTITIES_STATES.ok]);

    wrapper.setProps({ watcher: { state: { val: 1 } } });

    expect(wrapper.vm.format.color).toBe(WATCHER_STATES_COLORS[ENTITIES_STATES.minor]);

    wrapper.setProps({ watcher: { state: { val: 2 } } });

    expect(wrapper.vm.format.color).toBe(WATCHER_STATES_COLORS[ENTITIES_STATES.major]);

    wrapper.setProps({ watcher: { state: { val: 3 } } });

    expect(wrapper.vm.format.color).toBe(WATCHER_STATES_COLORS[ENTITIES_STATES.critical]);

    wrapper.setProps({ watcher: { state: { val: 4 } } });

    expect(wrapper.vm.format.color).toBe(WATCHER_STATES_COLORS.invalid);

    wrapper.setProps({ watcher: { state: {} } });

    expect(wrapper.vm.format.color).toBe(WATCHER_STATES_COLORS.invalid);

    wrapper.setProps({ watcher: { active_pb_all: true } });

    expect(wrapper.vm.format.color).toBe(WATCHER_PBEHAVIOR_COLOR);

    wrapper.setProps({
      watcher: {
        active_pb_all: false,
        active_pb_watcher: true,
      },
    });

    expect(wrapper.vm.format.color).toBe(WATCHER_PBEHAVIOR_COLOR);
  });

  it('has the right icon', () => {
    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS[ENTITIES_STATES.ok]);

    wrapper.setProps({ watcher: { state: { val: 1 } } });

    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS[ENTITIES_STATES.minor]);

    wrapper.setProps({ watcher: { state: { val: 2 } } });

    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS[ENTITIES_STATES.major]);

    wrapper.setProps({ watcher: { state: { val: 3 } } });

    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS[ENTITIES_STATES.critical]);

    wrapper.setProps({ watcher: { state: { val: 4 } } });

    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS.invalid);

    wrapper.setProps({ watcher: { state: {} } });

    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS.invalid);

    wrapper.setProps({
      watcher: {
        active_pb_watcher: true,
        watcher_pbehavior: [
          {
            type_: PBEHAVIOR_TYPES.maintenance,
          },
        ],
      },
    });

    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS.maintenance);

    wrapper.setProps({
      watcher: {
        active_pb_watcher: true,
        watcher_pbehavior: [
          {
            type_: PBEHAVIOR_TYPES.outOfSurveillance,
          },
        ],
      },
    });

    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS.outOfSurveillance);

    wrapper.setProps({
      watcher: {
        active_pb_watcher: true,
        watcher_pbehavior: [
          {
            type_: PBEHAVIOR_TYPES.pause,
          },
        ],
      },
    });

    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS.pause);

    wrapper.setProps({
      watcher: {
        active_pb_watcher: true,
        watcher_pbehavior: [
          {
            type_: PBEHAVIOR_TYPES.maintenance,
          },
          {
            type_: PBEHAVIOR_TYPES.outOfSurveillance,
          },
          {
            type_: PBEHAVIOR_TYPES.pause,
          },
        ],
      },
    });

    expect(wrapper.vm.format.icon).toBe(WEATHER_ICONS.maintenance);
  });
});
