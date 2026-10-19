import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createInputStub } from '@unit/stubs/input';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import ClickOutside from '@/services/click-outside';

import TextFieldEditor from '@/components/modals/common/text-field-editor.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-alert': true,
  'c-text-field': createInputStub('c-text-field'),
  'v-text-field': createInputStub('v-text-field'),
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

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
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith('');
    expect($modals.hide).toHaveBeenCalledWith({ id: undefined });

    wrapper.destroy();
  });

  test('Form submitted with correct value after trigger submit button', async () => {
    const action = vi.fn();
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
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(value);
    expect($modals.hide).toHaveBeenCalledWith({ id: undefined });

    wrapper.destroy();
  });

  test('Form didn\'t close after trigger submit button with action error', async () => {
    const action = vi.fn().mockRejectedValue({ unavailableField: 'Error' });
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();
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
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect(action).toBeCalled();
    expect($modals.hide).not.toBeCalled();

    consoleErrorSpy.mockRestore();

    wrapper.destroy();
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
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect($modals.hide).toHaveBeenCalledWith({ id: undefined });

    wrapper.destroy();
  });

  test('Validation rules applied to field from config', async () => {
    const action = vi.fn();
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
        $popups,
      },
    });

    expect(wrapper.vm.field.validationRules).toBe('required');

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();

    wrapper.destroy();
  });

  test('Errors added after trigger submit button with field error', async () => {
    const name = Faker.lorem.word();
    const formErrors = {
      [name]: 'Text error',
    };
    const action = vi.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
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
        $popups,
      },
    });
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(value);
    expect($modals.hide).not.toHaveBeenCalled();

    consoleErrorSpy.mockRestore();

    wrapper.destroy();
  });

  test('Errors added after trigger submit button with common error', async () => {
    const error = Faker.lorem.words();
    const name = Faker.lorem.word();
    const action = vi.fn().mockRejectedValue({ [name]: error });
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
        $popups,
      },
    });
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith('');
    expect($modals.hide).not.toHaveBeenCalled();

    consoleErrorSpy.mockRestore();

    wrapper.destroy();
  });

  test('Errors added after trigger submit button with error message', async () => {
    const message = Faker.lorem.words();
    const name = Faker.lorem.word();
    const action = vi.fn().mockRejectedValue({ [name]: message });
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
        $popups,
      },
    });
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith('');
    expect($modals.hide).not.toHaveBeenCalled();

    consoleErrorSpy.mockRestore();

    wrapper.destroy();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = vi.fn();
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
        $popups,
      },
    });

    const newValue = Faker.datatype.string();

    wrapper.vm.text = newValue;

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(newValue);
    expect($modals.hide).toHaveBeenCalled();

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
        $popups,
      },
    });

    wrapper.vm.closeModal();
    await flushPromises();

    expect($modals.hide).toHaveBeenCalled();
  });
});
