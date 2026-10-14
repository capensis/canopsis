import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import ClickOutside from '@/services/click-outside';

import EditLiveReporting from '@/components/modals/alarm/edit-live-reporting.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'date-interval-selector': true,
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

describe('edit-live-reporting', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(EditLiveReporting, {
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
            config: {},
          },
        },
        mocks: {
          $modals,
        },
      }),
    );

    expect(wrapper.vm.values).toEqual({
      tstart: '',
      tstop: '',
      time_field: '',
    });
  });

  test('Config parameters applied to form', () => {
    const config = {
      tstart: 'now-1d',
      tstop: 'now',
      time_field: 't',
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

    expect(wrapper.vm.values).toEqual(config);
  });

  test('Form submitted after trigger submit button', async () => {
    const action = vi.fn();
    const config = {
      tstart: 'now-1d',
      tstop: 'now',
      time_field: 't',
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

    expect(action).toHaveBeenCalledWith({
      tstart: config.tstart,
      tstop: config.tstop,
      time_field: config.time_field,
    });
    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = vi.fn().mockRejectedValue({ tstart: 'Tstart error' });
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      tstart: '',
      tstop: '',
      time_field: '',
    });
    expect($modals.hide).not.toBeCalled();
  });

  test('Form submitted after trigger submit button without action', async () => {
    const wrapper = factory(
      withModalInject({
        propsData: {
          modal: {
            config: {},
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
      tstart: 'Tstart error',
      tstop: 'Tstop error',
      time_field: 'Time field error',
    };
    const action = vi.fn().mockRejectedValue({
      ...formErrors,
      unavailableField: 'Error',
    });
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      tstart: '',
      tstop: '',
      time_field: '',
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
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    await wrapper.vm.submit();

    expect(consoleErrorSpy).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledWith({
      tstart: '',
      tstop: '',
      time_field: '',
    });
    expect($modals.hide).not.toHaveBeenCalledWith();

    consoleErrorSpy.mockRestore();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = vi.fn();
    const wrapper = factory(
      withModalInject({
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
      }),
    );

    const newForm = {
      tstart: 'now-1d',
      tstop: 'now',
      time_field: 't',
    };

    wrapper.vm.updateForm(newForm);
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
            config: {},
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
