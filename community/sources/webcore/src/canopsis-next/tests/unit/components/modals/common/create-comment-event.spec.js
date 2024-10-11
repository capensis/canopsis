import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createInputStub } from '@unit/stubs/input';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { fakeAlarms } from '@unit/data/alarm';

import ClickOutside from '@/services/click-outside';

import CreateCommentEvent from '@/components/modals/common/create-comment-event.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-name-field': createInputStub('c-name-field'),
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectNameField = wrapper => wrapper.find('.c-name-field');

describe('create-comment-event', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(CreateCommentEvent, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(CreateCommentEvent, {
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

  const timestamp = 1386435600000;

  beforeAll(() => jest.useFakeTimers({ now: timestamp }));

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();
    const config = { action };

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

    const submitButton = selectSubmitButton(wrapper);
    const textField = selectNameField(wrapper);

    const comment = Faker.datatype.string();

    textField.setValue(comment);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledTimes(1);
    expect(action).toBeCalledWith({ comment });
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = jest.fn();
    const config = { action };
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

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => wrapper.vm,
      vm: wrapper.vm,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const action = jest.fn();
    const formErrors = {
      comment: 'Comment error',
    };
    action.mockRejectedValueOnce({ ...formErrors, unavailableField: 'Error' });
    const wrapper = factory({
      propsData: {
        modal: {
          config: { action },
        },
      },
      mocks: {
        $modals,
      },
    });

    const comment = Faker.datatype.string();

    selectNameField(wrapper).setValue(comment);
    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toBeCalledTimes(1);
    expect(action).toBeCalledWith({ comment });
    expect($modals.hide).not.toBeCalledWith();
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
          config: { action },
        },
      },
      mocks: {
        $modals,
        $popups,
      },
    });

    const submitButton = selectSubmitButton(wrapper);
    const textField = selectNameField(wrapper);

    const comment = Faker.datatype.string();

    textField.setValue(comment);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toBeCalledTimes(1);
    expect(action).toBeCalledWith({ comment });
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
      mocks: {
        $modals,
      },
    });

    const cancelButton = selectCancelButton(wrapper);

    cancelButton.trigger('click');

    await flushPromises(true);

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-comment-event` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `create-comment-event` with items', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            items: fakeAlarms(10),
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
