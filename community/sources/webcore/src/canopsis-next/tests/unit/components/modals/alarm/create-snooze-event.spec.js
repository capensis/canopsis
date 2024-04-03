import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';

import { TIME_UNITS } from '@/constants';

import ClickOutside from '@/services/click-outside';

import CreateSnoozeEvent from '@/components/modals/alarm/create-snooze-event.vue';

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
  const snoozeEventData = {
    duration: {
      unit: TIME_UNITS.minute,
      value: 1,
    },
    comment: '',
  };
  const config = { items };

  const factory = generateShallowRenderer(CreateSnoozeEvent, {
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
  });
  const snapshotFactory = generateRenderer(CreateSnoozeEvent, {
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
  });

  test('Default parameters applied to form', () => {
    const wrapper = factory({
      mocks: {
        $modals,
      },
    });

    const snoozeEventForm = selectSnoozeEventForm(wrapper);

    expect(snoozeEventForm.vm.form).toEqual({
      comment: '',
      duration: {
        unit: TIME_UNITS.minute,
        value: 1,
      },
    });
  });

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            items,
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith(snoozeEventData);
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            ...config,
            action,
          },
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).not.toBeCalled();

    validator.detach('name');
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const action = jest.fn();
    const formErrors = {
      duration: 'Ticket error',
      comment: 'Comment error',
    };
    action.mockRejectedValueOnce({ ...formErrors, unavailableField: 'Error' });
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            ...config,
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toBeCalledWith(snoozeEventData);
    expect($modals.hide).not.toBeCalledWith();

    action.mockClear();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const action = jest.fn();
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    action.mockRejectedValueOnce(errors);

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            ...config,
            action,
          },
        },
      },
      mocks: {
        $modals,
        $popups,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toBeCalledWith(snoozeEventData);
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
    action.mockClear();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            ...config,
            action,
          },
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
        value: 2,
      },
      comment: 'comment',
    };

    snoozeEventForm.triggerCustomEvent('input', newForm);

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith(newForm);
    expect($modals.hide).toBeCalled();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-snooze-event` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `create-snooze-event` with config data', () => {
    const wrapper = snapshotFactory({
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

    expect(wrapper).toMatchSnapshot();
  });
});
