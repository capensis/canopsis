import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createInputStub } from '@unit/stubs/input';
import ClickOutside from '@/services/click-outside';

import TextFieldEditor from '@/components/modals/common/text-field-editor.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'v-text-field': createInputStub('v-text-field'),
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectTextField = wrapper => wrapper.find('.v-text-field');

describe('text-field-editor', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(TextFieldEditor, {
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
  const snapshotFactory = generateRenderer(TextFieldEditor, {
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

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith('');
    expect($modals.hide).toBeCalledWith();
  });

  test('Form submitted with correct value after trigger submit button', async () => {
    const action = jest.fn();
    const value = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            field: {
              value,
              name: 'field',
            },
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith(value);
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            field: {
              name: 'name',
            },
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const textField = selectTextField(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => textField.vm,
      vm: textField.vm,
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

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

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect($modals.hide).toBeCalledWith();
  });

  test('Validation rules applied to form from config', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            field: {
              validationRules: 'required',
              name: 'field',
            },
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(wrapper.getValidator().errors.any()).toBeTruthy();

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();
  });

  test('Errors added after trigger submit button with field error', async () => {
    const name = Faker.lorem.word();
    const formErrors = {
      [name]: 'Text error',
    };
    const action = jest.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
    const value = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            field: {
              name,
              value,
            },

            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });
    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toBeCalledWith(value);
    expect($modals.hide).not.toBeCalledWith();
  });

  test('Errors added after trigger submit button with common error', async () => {
    const error = Faker.lorem.words();
    const name = Faker.lorem.word();
    const action = jest.fn().mockRejectedValue({ error });
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            field: {
              name,
              value: '',
            },

            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });
    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect({ [name]: error }).toEqual(addedErrors);
    expect(action).toBeCalledWith('');
    expect($modals.hide).not.toBeCalledWith();
  });

  test('Errors added after trigger submit button with error message', async () => {
    const message = Faker.lorem.words();
    const name = Faker.lorem.word();
    const action = jest.fn().mockRejectedValue({ message });
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            field: {
              name,
              value: '',
            },

            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });
    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect({ [name]: message }).toEqual(addedErrors);
    expect(action).toBeCalledWith('');
    expect($modals.hide).not.toBeCalledWith();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            field: {
              value: '',
              name: 'field',
            },

            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const textField = selectTextField(wrapper);

    const newValue = Faker.datatype.string();

    textField.setValue(newValue);

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises(true);

    expect(action).toBeCalledWith(newValue);
    expect($modals.hide).toBeCalled();
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

  test('Renders `text-field-editor` with empty modal', () => {
    const wrapper = snapshotFactory({
      mocks: {
        $modals,
        $popups,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `text-field-editor` with all parameters', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            title: 'Text field editor title',
            label: 'Text field editor label',
            field: {
              name: 'field-name',
              value: 'field-value',
              label: 'field-label',
              validationRules: 'required',
            },
          },
        },
      },
      mocks: {
        $modals,
        $popups,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
