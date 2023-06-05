import { mount, createVueInstance } from '@unit/utils/vue';
import { createModalWrapperStub } from '@unit/stubs/modal';

import { createMockedStoreModules } from '@unit/utils/store';
import AlarmsList from '@/components/modals/alarm/alarms-list.vue';

const localVue = createVueInstance();

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarms-list-table-with-pagination': true,
};

const snapshotFactory = (options = {}) => mount(AlarmsList, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

describe('alarms-list', () => {
  const associativeTableModule = {
    name: 'associativeTable',
    actions: {
      fetch: jest.fn(() => ({})),
    },
  };

  const store = createMockedStoreModules([
    associativeTableModule,
  ]);

  test('Renders `alarms-list`', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        modal: {
          config: {
            widget: {
              parameters: {
                widgetColumns: [],
              },
            },
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
