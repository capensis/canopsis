import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import ClickOutside from '@/services/click-outside';

import AddInfoPopup from '@/components/modals/alarm/info-popup-setting/add-info-popup.vue';

const localVue = createVueInstance();

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'info-popup-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'info-popup-form': true,
};

const factory = (options = {}) => shallowMount(AddInfoPopup, {
  localVue,
  stubs,
  attachTo: document.body,

  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(AddInfoPopup, {
  localVue,
  stubs: snapshotStubs,

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
const selectInfoPopupForm = wrapper => wrapper
  .find('info-popup-form-stub');

describe('add-info-popup', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  test('Default parameters applied to form', () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns: [],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const infoPopupForm = selectInfoPopupForm(wrapper);

    expect(infoPopupForm.vm.form).toEqual({
      selectedColumn: undefined,
      template: '',
    });
  });

  test('Config parameters applied to form', () => {
    const column = Faker.datatype.string();
    const template = Faker.datatype.string();
    const columns = [
      {
        value: Faker.datatype.string(),
      },
      {
        value: column,
      },
    ];
    const config = {
      popup: {
        column,
        template,
      },
      columns,
    };
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

    const infoPopupForm = selectInfoPopupForm(wrapper);

    expect(infoPopupForm.vm.form).toEqual({
      column: columns[1].value,
      template,
    });
  });

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();
    const popup = {
      column: Faker.datatype.string(),
      template: Faker.datatype.string(),
    };
    const config = {
      columns: [{ value: popup.column }],
      popup,
      action,
    };
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

    submitButton.trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith(popup);
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns: [],
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const infoPopupForm = selectInfoPopupForm(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => infoPopupForm.vm,
      vm: infoPopupForm.vm,
    });

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();
  });

  test('Form submitted after trigger submit button without action', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns: [],
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

    expect($modals.hide).toBeCalledWith();
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      column: 'Column error',
      template: 'Template error',
    };
    const action = jest.fn()
      .mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
    const column = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns: [
              {
                value: column,
              },
            ],
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

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toBeCalledWith({
      column,
      template: '',
    });
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
    const column = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns: [{ value: column }],
            action,
          },
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
    expect(action).toBeCalledWith({
      column,
      template: '',
    });
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = jest.fn();
    const column = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns: [{ value: column }],
            action,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const infoPopupForm = selectInfoPopupForm(wrapper);

    const newForm = {
      column,
      template: Faker.datatype.string(),
    };

    infoPopupForm.vm.$emit('input', newForm);

    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith(newForm);
    expect($modals.hide).toBeCalled();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            columns: [],
          },
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

  test('Renders `add-info-popup` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            columns: [],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `add-info-popup` with popup', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            columns: [{}],
            popup: {
              column: 'column',
              template: 'template',
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
