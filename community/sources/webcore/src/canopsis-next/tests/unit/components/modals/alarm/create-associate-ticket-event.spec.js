import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createInputStub } from '@unit/stubs/input';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createMockedStoreModules } from '@unit/utils/store';
import ClickOutside from '@/services/click-outside';
import {
  ENTITIES_TYPES,
  EVENT_DEFAULT_ORIGIN,
  EVENT_ENTITY_TYPES,
  EVENT_INITIATORS,
} from '@/constants';

import CreateAssociateTicketEvent from '@/components/modals/alarm/create-associate-ticket-event.vue';

const localVue = createVueInstance();

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
  'v-text-field': createInputStub('v-text-field'),
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
};

const factory = (options = {}) => shallowMount(CreateAssociateTicketEvent, {
  localVue,
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

  ...options,
});

const snapshotFactory = (options = {}) => mount(CreateAssociateTicketEvent, {
  localVue,
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

  ...options,
});

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectTextField = wrapper => wrapper.find('.v-text-field');

describe('create-associate-ticket-event', () => {
  const timestamp = 1386435600000;

  mockDateNow(timestamp);
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
  const itemsType = ENTITIES_TYPES.alarm;
  const itemsIds = [alarm._id];
  const eventData = {
    id: alarm._id,
    component: alarm.v.component,
    connector: alarm.v.connector,
    connector_name: alarm.v.connector_name,
    resource: alarm.v.resource,
    crecord_type: EVENT_ENTITY_TYPES.assocTicket,
    event_type: EVENT_ENTITY_TYPES.assocTicket,
    initiator: EVENT_INITIATORS.user,
    origin: EVENT_DEFAULT_ORIGIN,
    ref_rk: `${alarm.v.resource}/${alarm.v.component}`,
    source_type: alarm.entity.type,
    state: alarm.v.state.val,
    state_type: alarm.v.status.val,
    timestamp: timestamp / 1000,
  };
  const assocTicketEventData = {
    ...eventData,

    output: 'Associated ticket number',
    ticket: '',
  };

  const getEntitiesList = jest.fn().mockReturnValue([alarm]);
  const entitiesModule = {
    name: 'entities',
    getters: {
      getList: () => getEntitiesList,
    },
  };

  const createEvent = jest.fn();
  const eventModule = {
    name: 'event',
    actions: {
      create: createEvent,
    },
  };
  const store = createMockedStoreModules([entitiesModule, eventModule]);

  afterEach(() => {
    createEvent.mockClear();
    getEntitiesList.mockClear();
  });

  test('Default parameters applied to form', () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    const textField = selectTextField(wrapper);

    expect(textField.vm.value).toBe('');
  });

  test('Form submitted after trigger submit button', async () => {
    const afterSubmit = jest.fn();
    const config = {
      itemsType,
      itemsIds,
      afterSubmit,
    };
    getEntitiesList.mockReturnValueOnce([{
      ...alarm,
      v: {
        ...alarm.v,
        ack: {},
      },
    }]);
    const wrapper = factory({
      store: createMockedStoreModules([
        entitiesModule,
        eventModule,
      ]),
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);
    const textField = selectTextField(wrapper);

    const ticket = Faker.datatype.string();

    textField.setValue(ticket);

    submitButton.trigger('click');

    await flushPromises();

    expect(createEvent).toBeCalledTimes(1);
    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          ...assocTicketEventData,
          ticket,
        }],
      },
      undefined,
    );
    expect(afterSubmit).toBeCalled();
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(createEvent).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      ticket: 'Ticket error',
    };
    createEvent.mockRejectedValueOnce({ ...formErrors, unavailableField: 'Error' });
    getEntitiesList.mockReturnValueOnce([{
      ...alarm,
      v: {
        ...alarm.v,
        ack: null,
      },
    }]);
    const fastAckOutput = {
      enabled: true,
      value: Faker.datatype.string(),
    };
    const wrapper = factory({
      store: createMockedStoreModules([
        entitiesModule,
        eventModule,
      ]),
      propsData: {
        modal: {
          config: {
            fastAckOutput,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);
    const textField = selectTextField(wrapper);

    const ticket = Faker.datatype.string();

    textField.setValue(ticket);

    submitButton.trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(createEvent).toBeCalledTimes(1);
    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          ...eventData,
          crecord_type: EVENT_ENTITY_TYPES.ack,
          event_type: EVENT_ENTITY_TYPES.ack,
          output: fastAckOutput.value,
          ticket,
        }],
      },
      undefined,
    );
    expect($modals.hide).not.toBeCalledWith();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    createEvent.mockRejectedValueOnce(errors);
    getEntitiesList.mockReturnValueOnce([{
      ...alarm,
      v: {
        ...alarm.v,
        ack: {},
      },
    }]);

    const wrapper = factory({
      store: createMockedStoreModules([
        entitiesModule,
        eventModule,
      ]),
      mocks: {
        $modals,
        $popups,
      },
    });

    const submitButton = selectSubmitButton(wrapper);
    const textField = selectTextField(wrapper);

    const ticket = Faker.datatype.string();

    textField.setValue(ticket);

    submitButton.trigger('click');

    await flushPromises();

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(createEvent).toBeCalledTimes(1);
    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          ...assocTicketEventData,
          ticket,
        }],
      },
      undefined,
    );
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal submitted with correct data after trigger form with ticket', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    const textField = selectTextField(wrapper);

    const ticket = Faker.datatype.string();

    textField.vm.$emit('input', ticket);

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(createEvent).toBeCalledTimes(2);
    expect(createEvent).toHaveBeenNthCalledWith(
      1,
      expect.any(Object),
      {
        data: [{
          ...eventData,
          crecord_type: EVENT_ENTITY_TYPES.ack,
          event_type: EVENT_ENTITY_TYPES.ack,
          output: '',
          ticket,
        }],
      },
      undefined,
    );
    expect(createEvent).toHaveBeenNthCalledWith(
      2,
      expect.any(Object),
      {
        data: [{
          ...assocTicketEventData,
          ticket,
        }],
      },
      undefined,
    );
    expect($modals.hide).toBeCalled();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    const cancelButton = selectCancelButton(wrapper);

    cancelButton.trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-associate-ticket-event` with empty modal', () => {
    const wrapper = snapshotFactory({
      store,
      mocks: {
        $modals,
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
            fastAckOutput: {
              enabled: true,
              value: 'Test',
            },
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
