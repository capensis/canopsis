import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createDeclareTicketModule, createMockedStoreModules } from '@unit/utils/store';

import { ALARM_LIST_STEPS, MODALS } from '@/constants';

import ClickOutside from '@/services/click-outside';

import CreateDeclareTicketEvent from '@/components/modals/declare-ticket/create-declare-ticket-event.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-progress-overlay': true,
  'declare-ticket-events-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-progress-overlay': true,
  'declare-ticket-events-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);

describe('create-declare-ticket-event', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const alarm = {
    _id: Faker.datatype.string(),
    v: {
      connector: Faker.datatype.string(),
      connector_name: Faker.datatype.string(),
      component: Faker.datatype.string(),
      resource: Faker.datatype.string(),
      state: {
        val: Faker.datatype.number(),
      },
      status: {
        val: Faker.datatype.number(),
      },
    },
    entity: {
      type: Faker.datatype.number(),
    },
  };
  const items = [alarm];
  const config = {
    items,
    alarmsByTickets: {},
    ticketsByAlarms: {},
  };

  const { declareTicketRuleModule } = createDeclareTicketModule();
  const store = createMockedStoreModules([declareTicketRuleModule]);

  const factory = generateShallowRenderer(CreateDeclareTicketEvent, {
    stubs,
    attachTo: document.body,
    propsData: {
      modal: {
        config: {},
      },
    },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
    mocks: {
      $modals,
      $popups,
    },
  });

  const snapshotFactory = generateRenderer(CreateDeclareTicketEvent, {
    stubs: snapshotStubs,
    propsData: {
      modal: {
        config: {},
      },
    },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
    mocks: {
      $modals,
      $popups,
    },
  });

  const timestamp = 1386435600000;

  beforeAll(() => jest.useFakeTimers({ now: timestamp }));

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();

    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            items,
            action,
            alarmsByTickets: {},
            ticketsByAlarms: {},
          },
        },
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledTimes(1);
    expect(action).toBeCalledWith([], false);
    expect($modals.hide).toBeCalledWith();
  });

  test('Confirmation modal showed after submit with exist tickets', async () => {
    const action = jest.fn();
    const ticketRuleId = Faker.datatype.string();
    const alarmWithTickets = {
      ...alarm,
      v: {
        ...alarm.v,
        tickets: [
          {
            ticket_rule_id: ticketRuleId,
            _t: ALARM_LIST_STEPS.declareTicket,
          },
        ],
      },
    };
    const itemsWithTickets = [alarmWithTickets];

    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            items: itemsWithTickets,
            action,
            alarmsByTickets: {
              [ticketRuleId]: {
                alarms: [alarmWithTickets._id],
              },
            },
            ticketsByAlarms: {},
          },
        },
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: expect.objectContaining({
          action: expect.any(Function),
        }),
      },
    );

    const [{ config: confirmationModalConfig }] = $modals.show.mock.calls[0];

    await confirmationModalConfig.action();

    expect(action).toBeCalledTimes(1);
    expect(action).toBeCalledWith([
      {
        _id: ticketRuleId,
        alarms: [alarmWithTickets._id],
        comment: '',
        ticket_resources: false,
      },
    ], false);
    expect($modals.hide).toBeCalledWith();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = jest.fn().mockRejectedValueOnce(errors);

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            ...config,
            action,
          },
        },
      },
      store,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toBeCalledWith([], false);
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config,
        },
      },
    });

    selectCancelButton(wrapper).trigger('click');

    await flushPromises(true);

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-declare-ticket-event`', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        modal: {
          config,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
