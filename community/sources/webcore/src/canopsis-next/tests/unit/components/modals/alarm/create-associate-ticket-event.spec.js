import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createMockedStoreModules } from '@unit/utils/store';
import ClickOutside from '@/services/click-outside';

import CreateAssociateTicketEvent from '@/components/modals/declare-ticket/create-associate-ticket-event.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-alert': true,
  'alarm-general-table': true,
  'associate-ticket-event-form': true,
  'c-description-field': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-alert': true,
  'alarm-general-table': true,
  'associate-ticket-event-form': true,
  'c-description-field': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectAssociateTicketEventForm = wrapper => wrapper.find('associate-ticket-event-form-stub');

describe('create-associate-ticket-event', () => {
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
  const assocTicketEventData = {
    comment: '',
    ticket: '',
    data: {},
    system_name: '',
    ticket_resources: false,
    url: '',
  };

  const createEvent = jest.fn();
  const eventModule = {
    name: 'event',
    actions: {
      create: createEvent,
    },
  };
  const store = createMockedStoreModules([eventModule]);

  const factory = generateShallowRenderer(CreateAssociateTicketEvent, {

    stubs,
    attachTo: document.body,
    propsData: {
      modal: {
        config: {
          items: [],
        },
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

  const snapshotFactory = generateRenderer(CreateAssociateTicketEvent, {

    stubs: snapshotStubs,
    propsData: {
      modal: {
        config: {
          items: [],
        },
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

  afterEach(() => {
    createEvent.mockClear();
  });

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            items: [{
              ...alarm,
              v: {
                ...alarm.v,
                ack: {},
              },
            }],
            action,
          },
        },
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith({
      comment: '',
      ticket: '',
      data: {},
      system_name: '',
      ticket_resources: false,
      url: '',
    });
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            items,
          },
        },
      },
    });

    const associateTicketEventForm = selectAssociateTicketEventForm(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => associateTicketEventForm.vm,
      vm: associateTicketEventForm.vm,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(createEvent).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      ticket_id: 'Ticket id error',
      ticket_url: 'Ticket url error',
      system_name: 'System name error',
    };
    const action = jest.fn().mockRejectedValueOnce({ ...formErrors, unavailableField: 'Error' });
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            action,
            items: [{
              ...alarm,
              v: {
                ...alarm.v,
                ack: null,
              },
            }],
          },
        },
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(addedErrors).toEqual(formErrors);
    expect(action).toBeCalledWith(assocTicketEventData);
    expect($modals.hide).not.toBeCalledWith();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = jest.fn().mockRejectedValueOnce(errors);

    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            action,
            items: [{
              ...alarm,
              v: {
                ...alarm.v,
                ack: {},
              },
            }],
          },
        },
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toBeCalledWith(assocTicketEventData);
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal submitted with correct data after trigger form with ticket', async () => {
    const action = jest.fn();
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            action,
            items,
          },
        },
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(action).toHaveBeenCalledWith(assocTicketEventData);
    expect($modals.hide).toBeCalled();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            items,
          },
        },
      },
    });

    const cancelButton = selectCancelButton(wrapper);

    cancelButton.trigger('click');

    await flushPromises(true);

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-associate-ticket-event` with empty modal', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        modal: {
          config: {
            items,
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `create-associate-ticket-event` with config data', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        modal: {
          config: {
            items,
            fastAckOutput: {
              enabled: true,
              value: 'Test',
            },
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
