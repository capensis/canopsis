import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createInputStub } from '@unit/stubs/input';

import ClickOutside from '@/services/click-outside';

import TextEditor from '@/components/modals/common/text-editor.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'text-editor-field': createInputStub('text-editor-field'),
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'text-editor-field': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectTextEditorField = wrapper => wrapper.find('.text-editor-field');

describe('text-editor', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(TextEditor, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(TextEditor, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
        $system: {},
      },
    },
  });

  beforeAll(() => jest.useFakeTimers());

  test('Form submitted with empty string after trigger submit button', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith('');
    expect($modals.hide).toBeCalledWith();
  });

  test('Form submitted with correct value after trigger submit button', async () => {
    const action = jest.fn();
    const text = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            text,
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith(text);
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
            rules: {
              required: true,
            },
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();
  });

  test('Form submitted after trigger submit button without action', async () => {
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect($modals.hide).toBeCalledWith();
  });

  test('Validation rules applied to form from config', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            text: '',
            rules: {
              required: true,
            },
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    expect(wrapper.getValidator().errors.any()).toBeTruthy();

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      text: 'Text error',
    };
    const action = jest.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
    const text = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            text,

            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises(true);

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toBeCalledWith(text);
    expect($modals.hide).not.toBeCalledWith();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error')
      .mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = jest.fn().mockRejectedValue(errors);
    const text = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            text,
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

    await flushPromises(true);

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toBeCalledWith(text);
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    jest.useFakeTimers();

    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            text: '',

            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const newValue = Faker.datatype.string();

    await selectTextEditorField(wrapper).triggerCustomEvent('input', newValue);
    await selectSubmitButton(wrapper).trigger('click');

    jest.runAllTimers();

    await flushPromises();

    expect(action).toBeCalledWith(newValue);
    expect($modals.hide).toBeCalled();

    jest.useRealTimers();
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

    await selectCancelButton(wrapper).trigger('click');

    expect($modals.hide).toBeCalled();
  });

  test('Renders `text-editor` with empty modal', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
      mocks: {
        $modals,
        $popups,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
