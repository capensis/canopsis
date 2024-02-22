import { generateRenderer } from '@unit/utils/vue';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createMockedStoreModules } from '@unit/utils/store';

import AlarmsList from '@/components/modals/alarm/alarms-list.vue';

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarms-list-table-with-pagination': true,
};

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

  const snapshotFactory = generateRenderer(AlarmsList, { stubs: snapshotStubs });

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

    expect(wrapper).toMatchSnapshot();
  });
});
