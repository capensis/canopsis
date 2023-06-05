import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createMockedStoreModules } from '@unit/utils/store';
import ClickOutside from '@/services/click-outside';
import {
  EVENT_DEFAULT_ORIGIN,
  EVENT_ENTITY_TYPES,
  EVENT_INITIATORS,
  TIME_UNITS,
} from '@/constants';

import CreateSnoozeEvent from '@/components/modals/alarm/create-snooze-event.vue';

const localVue = createVueInstance();

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'snooze-event-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'snooze-event-form': true,
};

const factory = (options = {}) => shallowMount(CreateSnoozeEvent, {
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

const snapshotFactory = (options = {}) => mount(CreateSnoozeEvent, {
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
const selectSnoozeEventForm = wrapper => wrapper.find('snooze-event-form-stub');

describe('create-snooze-event', () => {
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
  const items = [alarm];
  const eventData = {
    id: alarm._id,
    component: alarm.v.component,
    connector: alarm.v.connector,
    connector_name: alarm.v.connector_name,
    resource: alarm.v.resource,
    crecord_type: EVENT_ENTITY_TYPES.snooze,
    event_type: EVENT_ENTITY_TYPES.snooze,
    initiator: EVENT_INITIATORS.user,
    origin: EVENT_DEFAULT_ORIGIN,
    ref_rk: `${alarm.v.resource}/${alarm.v.component}`,
    source_type: alarm.entity.type,
    state: alarm.v.state.val,
    timestamp: timestamp / 1000,
  };
  const snoozeEventData = {
    ...eventData,

    duration: 60,
    output: '',
  };
  const config = { items };

  const createEvent = jest.fn();
  const eventModule = {
    name: 'event',
    actions: {
      create: createEvent,
    },
  };
  const store = createMockedStoreModules([eventModule]);

  afterEach(() => {
    createEvent.mockClear();
  });

  test('Default parameters applied to form', () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    const snoozeEventForm = selectSnoozeEventForm(wrapper);

    expect(snoozeEventForm.vm.form).toEqual({
      output: '',
      duration: {
        unit: TIME_UNITS.minute,
        value: 1,
      },
    });
  });

  test('Form submitted after trigger submit button', async () => {
    const afterSubmit = jest.fn();

    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            items,
            afterSubmit,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(createEvent).toBeCalledTimes(1);
    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [snoozeEventData],
      },
      undefined,
    );
    expect(afterSubmit).toBeCalled();
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    const snoozeEventForm = selectSnoozeEventForm(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => snoozeEventForm.vm,
      vm: snoozeEventForm.vm,
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(createEvent).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();

    validator.detach('name');
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      duration: 'Ticket error',
      output: 'Output error',
    };
    createEvent.mockRejectedValueOnce({ ...formErrors, unavailableField: 'Error' });
    const wrapper = factory({
      store,
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

    submitButton.trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(createEvent).toBeCalledTimes(1);
    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [snoozeEventData],
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

    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
        $popups,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

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
        data: [snoozeEventData],
      },
      undefined,
    );
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    const snoozeEventForm = selectSnoozeEventForm(wrapper);

    const newForm = {
      duration: {
        unit: TIME_UNITS.hour,
        value: 1,
      },
      output: 'output',
    };

    snoozeEventForm.vm.$emit('input', newForm);

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(createEvent).toBeCalledTimes(1);
    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          ...snoozeEventData,
          duration: 3600,
          output: newForm.output,
        }],
      },
      undefined,
    );
    expect($modals.hide).toBeCalled();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    const cancelButton = selectCancelButton(wrapper);

    cancelButton.trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-snooze-event` with empty modal', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `create-snooze-event` with config data', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        modal: {
          config: {
            items,
            isNoteRequired: true,
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
