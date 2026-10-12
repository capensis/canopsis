import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import ClickOutside from '@/services/click-outside';

import AddInfoPopup from '@/components/modals/alarm/info-popup-setting/add-info-popup.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'info-popup-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const modalId = wrapper => ({ id: wrapper.props().modal.id });

const withModalInject = (options = {}) => ({
  ...options,
  propsData: {
    ...options.propsData,
    modal: {
      id: 'test-modal-id',
      ...options.propsData?.modal,
    },
  },
});

describe('add-info-popup', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(AddInfoPopup, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  test('Default parameters applied to form', () => {
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    expect(wrapper.vm.form).toEqual({
      column: undefined,
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
    const wrapper = factory(
      withModalInject({
        propsData: {
          modal: {
            config,
          },
        },
        mocks: {
          $modals,
        },
      }),
    );

    expect(wrapper.vm.form).toEqual({
      column: columns[1].value,
      template,
    });
  });

  test('Form submitted after trigger submit button', async () => {
    const action = vi.fn();
    const popup = {
      column: Faker.datatype.string(),
      template: Faker.datatype.string(),
    };
    const config = {
      columns: [{ value: popup.column }],
      popup,
      action,
    };
    const wrapper = factory(
      withModalInject({
        propsData: {
          modal: {
            config,
          },
        },
        mocks: {
          $modals,
        },
      }),
    );

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(popup);
    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = vi.fn().mockRejectedValue({ column: 'Column error' });
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      column: undefined,
      template: '',
    });
    expect($modals.hide).not.toBeCalled();
  });

  test('Form submitted after trigger submit button without action', async () => {
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    await wrapper.vm.submit();

    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      column: 'Column error',
      template: 'Template error',
    };
    const action = vi.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
    const column = Faker.datatype.string();

    const wrapper = factory(
      withModalInject({
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
      }),
    );

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      column,
      template: '',
    });
    expect($modals.hide).not.toHaveBeenCalledWith();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = vi.fn().mockRejectedValue(errors);
    const column = Faker.datatype.string();
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    await wrapper.vm.submit();

    expect(consoleErrorSpy).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledWith({
      column,
      template: '',
    });
    expect($modals.hide).not.toHaveBeenCalledWith();

    consoleErrorSpy.mockRestore();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = vi.fn();
    const column = Faker.datatype.string();
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    const newForm = {
      column,
      template: Faker.datatype.string(),
    };

    wrapper.vm.form = newForm;
    await flushPromises();
    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(newForm);
    expect($modals.hide).toBeCalledWith(modalId(wrapper));
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    wrapper.vm.closeModal();
    await flushPromises();

    expect($modals.hide).toBeCalledWith(modalId(wrapper));
  });
});
