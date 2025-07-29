import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';

import { ALARM_STATES } from '@/constants';

import ClickOutside from '@/services/click-outside';

import CreateChangeStateEvent from '@/components/modals/alarm/create-change-state-event.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-change-state-field': true,
  'alarms-list-general-table': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-change-state-field': true,
  'alarms-list-general-table': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectChangeStateField = wrapper => wrapper.find('c-change-state-field-stub');

describe('create-change-state-event', () => {
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
  const changeStateEventData = {
    state: alarm.v.state.val,
    comment: '',
  };
  const config = { items };

  const factory = generateShallowRenderer(CreateChangeStateEvent, {
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
  const snapshotFactory = generateRenderer(CreateChangeStateEvent, {
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
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    const changeStateEventField = selectChangeStateField(wrapper);

    expect(changeStateEventField.vm.$attrs.value).toEqual({
      comment: '',
      state: alarm.v.state.val,
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

    expect(action).toHaveBeenCalledWith(changeStateEventData);
    expect($modals.hide).toHaveBeenCalled();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = jest.fn();
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

    const changeStateEventField = selectChangeStateField(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => changeStateEventField.vm,
      vm: changeStateEventField.vm,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).not.toHaveBeenCalled();
    expect($modals.hide).not.toHaveBeenCalled();

    validator.detach('name');
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const action = jest.fn();
    const formErrors = {
      state: 'State error',
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
    expect(action).toHaveBeenCalledWith(changeStateEventData);
    expect($modals.hide).not.toHaveBeenCalledWith();
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

    expect(consoleErrorSpy).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledWith(changeStateEventData);
    expect($modals.hide).not.toHaveBeenCalledWith();

    consoleErrorSpy.mockClear();
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

    const changeStateEventField = selectChangeStateField(wrapper);

    const newForm = {
      state: Faker.datatype.number(),
      comment: 'comment',
    };

    changeStateEventField.triggerCustomEvent('input', newForm);

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledWith({
      ...changeStateEventData,
      comment: newForm.comment,
      state: newForm.state,
    });
    expect($modals.hide).toHaveBeenCalled();
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

    expect($modals.hide).toHaveBeenCalled();
  });

  test('Action called with single item state when config.items has one item', async () => {
    const action = jest.fn();
    const singleItemState = ALARM_STATES.critical;
    const singleItem = {
      ...alarm,
      v: {
        ...alarm.v,
        state: {
          val: singleItemState,
        },
      },
    };

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            items: [singleItem],
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

    expect(action).toHaveBeenCalledWith({
      comment: '',
      state: singleItemState,
    });
    expect($modals.hide).toHaveBeenCalled();
  });

  test('Action called with major state when config.items has multiple items', async () => {
    const action = jest.fn();
    const multipleItems = [
      {
        ...alarm,
        v: {
          ...alarm.v,
          state: {
            val: ALARM_STATES.critical,
          },
        },
      },
      {
        ...alarm,
        v: {
          ...alarm.v,
          state: {
            val: ALARM_STATES.minor,
          },
        },
      },
    ];

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            items: multipleItems,
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

    expect(action).toHaveBeenCalledWith({
      comment: '',
      state: ALARM_STATES.major,
    });
    expect($modals.hide).toHaveBeenCalled();
  });

  test('Renders `create-change-state-event`', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            items: [{
              v: {
                state: {
                  val: ALARM_STATES.major,
                },
              },
            }],
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
