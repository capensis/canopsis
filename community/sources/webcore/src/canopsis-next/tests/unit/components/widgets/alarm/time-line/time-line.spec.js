import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createMockedStoreModules } from '@unit/utils/store';

import TimeLine from '@/components/widgets/alarm/time-line/time-line.vue';

const localVue = createVueInstance();

const stubs = {
  'time-line-flag': true,
  'time-line-card': true,
};

const factory = (options = {}) => shallowMount(TimeLine, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(TimeLine, {
  localVue,
  stubs,

  ...options,
});

describe('time-line', () => {
  const fetchItem = jest.fn();
  const updateQuery = jest.fn();
  const getQueryById = jest.fn(() => ({}));
  const alarmModule = {
    name: 'alarm',
    getters: {},
    actions: {
      fetchItem,
    },
  };
  const queryModule = {
    name: 'query',
    getters: { getQueryById: () => getQueryById },
    actions: {
      update: updateQuery,
    },
  };
  const widget = {
    _id: 'widget-id',
  };
  const alarm = {
    _id: 'alarm-id',
    v: {},
  };
  const steps = [
    {
      _t: 'stateinc',
      t: 1626159262,
      a: 'root',
      m: 'Idle rule Test all resource',
      val: 3,
      initiator: 'system',
    },
    {
      _t: 'statusinc',
      t: 1626159262,
      a: 'root',
      m: 'Idle rule Test all resource',
      val: 5,
      initiator: 'system',
    },
    {
      _t: 'pbhenter',
      t: 1627641985,
      a: 'system',
      m: 'Pbehavior Name pbh. Type: Default pause. Reason: Test reason',
      val: 0,
      initiator: 'external',
    },
    {
      _t: 'pbhleave',
      t: 1632723441,
      a: 'system',
      m: 'Pbehavior Name pbh. Type: Default pause. Reason: Test reason',
      val: 0,
      initiator: 'external',
    },
    {
      _t: 'ack',
      t: 1632725253,
      a: 'root',
      m: '',
      val: 0,
      initiator: 'user',
    },
    {
      _t: 'pbhenter',
      t: 1632977650,
      a: 'system',
      m: 'Pbehavior Name. Type: Default maintenance. Reason: Test reason',
      val: 0,
      initiator: 'external',
    },
    {
      _t: 'pbhleave',
      t: 1634553310,
      a: 'system',
      m: 'Pbehavior Name. Type: Default maintenance. Reason: Test reason',
      val: 0,
      initiator: 'external',
    },
  ];

  const store = createMockedStoreModules([
    alarmModule,
    queryModule,
  ]);

  beforeEach(() => {
    fetchItem.mockClear();
    updateQuery.mockClear();
    getQueryById.mockClear();
  });

  it('Alarm fetched after mount, if alarm hasn\'t steps', () => {
    factory({
      store,
      propsData: {
        alarm,
        widget,
      },
    });

    expect(getQueryById).toBeCalledWith(widget._id);

    expect(fetchItem).toBeCalledWith(
      expect.any(Object),
      {
        id: alarm._id,
        params: {
          correlation: false,
          limit: 1,
          sort_dir: 'desc',
          sort_key: 't',
          with_instructions: true,
          with_steps: true,
        },
        dataPreparer: expect.any(Function),
      },
      undefined,
    );

    const [, options] = fetchItem.mock.calls[0];

    const fetchedAlarms = [alarm];

    expect(options.dataPreparer({ data: fetchedAlarms })).toBe(fetchedAlarms);
  });

  it('Renders `time-line` without steps', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm,
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `time-line` with steps', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          ...alarm,
          v: { steps },
        },
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `time-line` with updated steps', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          ...alarm,
          v: { steps },
        },
        widget,
      },
    });

    await wrapper.setProps({
      alarm: {
        ...alarm,
        v: {
          steps: steps.slice(0, 2),
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
