import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups, mockSidebar } from '@unit/utils/mock-hooks';
import {
  createAuthModule,
  createLlmModule,
  createMockedStoreModules,
  createUserModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { merge } from 'lodash';

import ClickOutside from '@/services/click-outside';

import CreatePbehavior from '@/components/modals/pbehavior/create-pbehavior.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pattern-progress': true,
  'pbehavior-form': true,
  'pbehavior-general-form': true,
  'c-enabled-field': true,
  'v-layout': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const selectPbehaviorForm = wrapper => wrapper.find('pbehavior-form-stub');
const modalId = wrapper => ({ id: wrapper.props().modal.id });

const withModalInject = (options = {}) => {
  const rawModal = options.propsData?.modal ?? { config: {} };
  const modalWithId = rawModal.id ? rawModal : { id: 'test-modal-id', ...rawModal };
  const modal = {
    ...modalWithId,
    config: {
      timezone: 'Europe/Paris',
      ...modalWithId.config,
    },
  };

  return merge({}, options, {
    propsData: {
      ...options.propsData,
      modal,
    },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
        $modal: modal,
      },
    },
  });
};

describe('create-pbehavior', () => {
  const $modals = {
    ...mockModals(),
    updateModalConfig: vi.fn(),
    updateDialogProps: vi.fn(),
  };
  const $popups = mockPopups();
  const $sidebar = mockSidebar();
  const defaultPbehavior = {
    _id: expect.any(String),
    name: '',
    color: '',
    comments: [],
    enabled: true,
    inherited: false,
    entity_pattern: [],
    exceptions: [],
    exdates: [],
    reason: undefined,
    type: undefined,
    rrule: '',
    tstart: null,
    tstop: null,
    timezone: 'Europe/Paris',
  };

  const defaultPbehaviorRequest = {
    ...defaultPbehavior,
    exec_pattern: true,
  };

  const { authModule } = createAuthModule();
  const { userModule } = createUserModule();
  const { llmModule } = createLlmModule();
  const store = createMockedStoreModules([authModule, userModule, llmModule]);

  const shallowCreatePbehavior = generateShallowRenderer(CreatePbehavior, {
    stubs,
    store,
    attachTo: document.body,
    mocks: { $modals, $popups, $sidebar },
  });
  const factory = (options = {}) => shallowCreatePbehavior(withModalInject(options));

  test('Form submitted after trigger submit button', async () => {
    const action = vi.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
          },
        },
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(defaultPbehaviorRequest);
    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = vi.fn().mockRejectedValue({ name: 'Name error' });
    const afterSubmit = vi.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
            afterSubmit,
          },
        },
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(defaultPbehaviorRequest);
    expect(afterSubmit).not.toHaveBeenCalled();
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Form submitted after trigger submit button without action', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
    });

    await wrapper.vm.submit();

    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      name: 'Name error',
      rrule: 'Rrule error',
      tstart: 'Tstart error',
    };
    const action = vi.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
          },
        },
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(defaultPbehaviorRequest);
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = vi.fn().mockRejectedValue(errors);
    const customPbehavior = {
      name: Faker.datatype.string(),
      _id: Faker.datatype.string(),
      entity_pattern: [],
    };
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            pbehavior: customPbehavior,
            action,
          },
        },
      },
    });

    await wrapper.vm.submit();

    expect(consoleErrorSpy).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledWith({
      ...defaultPbehaviorRequest,
      entity_pattern: customPbehavior.entity_pattern,
      name: customPbehavior.name,
      _id: customPbehavior._id,
    });
    expect($modals.hide).not.toHaveBeenCalled();

    consoleErrorSpy.mockRestore();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = vi.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
          },
        },
      },
    });

    const newForm = {
      name: Faker.datatype.string(),
      _id: Faker.datatype.string(),
      color: Faker.internet.color(),
      entity_pattern: [],
      rrule: null,
      timezone: 'Europe/Paris',
    };

    await selectPbehaviorForm(wrapper).triggerCustomEvent('input', {
      ...defaultPbehavior,
      ...newForm,
    });
    await flushPromises();
    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      ...defaultPbehaviorRequest,
      ...newForm,
    });
    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
    });

    wrapper.vm.closeModal();
    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });
});
