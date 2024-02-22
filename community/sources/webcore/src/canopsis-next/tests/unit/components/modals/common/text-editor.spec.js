import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';

import ClickOutside from '@/services/click-outside';

import TextEditor from '@/components/modals/common/text-editor.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'text-editor-field': true,
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
const selectTextEditorField = wrapper => wrapper.find('text-editor-field-stub');

describe('text-editor', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(TextEditor, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
      provide: {
        $clickOutside: new ClickOutside(),
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(TextEditor, {
    stubs: snapshotStubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
      provide: {
        $clickOutside: new ClickOutside(),
        $system: {},
      },
    },
  });

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

    await flushPromises();

    expect(action).toBeCalledWith('');
    expect($modals.hide).toBeCalledWith();

    wrapper.destroy();
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

    await flushPromises();

    expect(action).toBeCalledWith(text);
    expect($modals.hide).toBeCalledWith();

    wrapper.destroy();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
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

    const textEditorField = selectTextEditorField(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => textEditorField.vm,
      vm: textEditorField.vm,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();

    validator.detach('name');
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

    await flushPromises();

    expect($modals.hide).toBeCalledWith();

    wrapper.destroy();
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

    await flushPromises();

    expect(wrapper.getValidator().errors.any()).toBeTruthy();

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();

    wrapper.destroy();
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

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toBeCalledWith(text);
    expect($modals.hide).not.toBeCalledWith();

    wrapper.destroy();
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

    await flushPromises();

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toBeCalledWith(text);
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();

    wrapper.destroy();
  });

  test('Modal submitted with correct data after trigger form', async () => {
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

    const textEditorField = selectTextEditorField(wrapper);

    const newValue = Faker.datatype.string();

    textEditorField.triggerCustomEvent('input', newValue);

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith(newValue);
    expect($modals.hide).toBeCalled();

    wrapper.destroy();
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

    selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalled();

    wrapper.destroy();
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

    wrapper.destroy();
  });
});
