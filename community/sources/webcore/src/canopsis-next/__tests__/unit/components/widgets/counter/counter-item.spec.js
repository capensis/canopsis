import { mockModals } from '@unit/utils/mock-hooks';
import { createAuthModule, createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';

import { MODALS } from '@/constants';

import CounterItem from '@/components/widgets/counter/counter-item.vue';

const stubs = {
  'card-with-see-alarms-btn': true,
  'c-compiled-template': true,
};

const selectCard = wrapper => wrapper.find('card-with-see-alarms-btn-stub');

describe('counter-item', () => {
  const $modals = mockModals();
  const fetchAlarmsListWithoutStore = vi.fn();
  const alarmModule = {
    name: 'alarm',
    actions: {
      fetchListWithoutStore: fetchAlarmsListWithoutStore,
    },
  };
  const { authModule } = createAuthModule();
  const store = createMockedStoreModules([authModule, alarmModule]);

  const counter = {
    total: 5,
    filter: {
      _id: 'filter-id',
      title: 'Filter title',
    },
  };
  const widget = {
    _id: 'counter-widget-id',
    parameters: {
      alarmsList: {
        itemsPerPage: 25,
      },
      blockTemplate: '',
      heightFactor: 0,
      levels: {
        colors: {},
        counter: 'total',
        values: {},
      },
      margin: {
        top: 0,
        right: 0,
        bottom: 0,
        left: 0,
      },
      opened: true,
    },
  };
  const query = {
    correlation: true,
    opened: true,
  };

  const factory = generateShallowRenderer(CounterItem, { stubs });

  afterEach(() => {
    fetchAlarmsListWithoutStore.mockReset();
  });

  test('Alarms list modal showed after trigger see alarms button', async () => {
    const wrapper = factory({
      store,
      propsData: {
        counter,
        query,
        widget,
      },
      mocks: {
        $modals,
      },
    });

    selectCard(wrapper).triggerCustomEvent('show:alarms');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.alarmsList,
      config: {
        widget: expect.objectContaining({
          _id: widget._id,
          parameters: expect.objectContaining({
            itemsPerPage: widget.parameters.alarmsList.itemsPerPage,
            opened: widget.parameters.opened,
          }),
        }),
        title: 'modals.alarmsList.prefixTitle',
        fetchList: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];

    await config.fetchList({
      page: 2,
      tstart: 100,
      tstop: 200,
    });

    expect(fetchAlarmsListWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      params: {
        ...query,
        page: 2,
        filters: [counter.filter._id],
      },
    });
  });
});
