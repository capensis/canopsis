import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createInputStub } from '@unit/stubs/input';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import ClickOutside from '@/services/click-outside';

import TextEditor from '@/components/modals/common/text-editor.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'text-editor-field': createInputStub('text-editor-field'),
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
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

  beforeAll(() => vi.useFakeTimers());

  test('Form submitted with empty string after trigger submit button', async () => {
    const action = vi.fn();
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

    vi.runOnlyPendingTimers();
    await flushPromises();

    expect(action).toHaveBeenCalledWith('');
    expect($modals.hide).toHaveBeenCalledWith();
  });

  test('Form submitted with correct value after trigger submit button', async () => {
    const action = vi.fn();
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

    vi.runOnlyPendingTimers();
    await flushPromises();

    expect(action).toHaveBeenCalledWith(text);
    expect($modals.hide).toHaveBeenCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = vi.fn();
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

    vi.runOnlyPendingTimers();
    await flushPromises();

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

    vi.runOnlyPendingTimers();
    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith();
  });

  test('Validation rules applied to form from config', async () => {
    const action = vi.fn();
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

    vi.runOnlyPendingTimers();
    await flushPromises();

    expect(wrapper.getValidator().errors.any()).toBeTruthy();

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      text: 'Text error',
    };
    const action = vi.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
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

    vi.runOnlyPendingTimers();
    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toHaveBeenCalledWith(text);
    expect($modals.hide).not.toHaveBeenCalledWith();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = vi.fn().mockRejectedValue(errors);
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

    vi.runOnlyPendingTimers();
    await flushPromises();

    expect(consoleErrorSpy).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledWith(text);
    expect($modals.hide).not.toHaveBeenCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    vi.useFakeTimers();

    const action = vi.fn();
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

    vi.runAllTimers();

    await flushPromises();

    expect(action).toHaveBeenCalledWith(newValue);
    expect($modals.hide).toBeCalled();

    vi.useRealTimers();
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
});
