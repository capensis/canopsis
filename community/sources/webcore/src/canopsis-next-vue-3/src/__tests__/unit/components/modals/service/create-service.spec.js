import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockConsole, mockModals, mockPopups, mockSidebar } from '@unit/utils/mock-hooks';
import {
  createAuthModule,
  createLlmModule,
  createMockedStoreModules,
  createTemplateVarsModule,
  createUserModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { merge } from 'lodash';

import ClickOutside from '@/services/click-outside';

import { formToService, serviceToForm } from '@/helpers/entities/service/form';

import CreateService from '@/components/modals/service/create-service.vue';

vi.mock('@/hooks/vars/template', async () => {
  const { ref } = await vi.importActual('vue');

  return {
    useTemplateVarsList: () => ({
      vars: ref({}),
      pending: ref(false),
      fetchList: vi.fn(),
    }),
  };
});

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-progress-overlay': true,
  'pattern-progress': true,
  'service-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const selectServiceForm = wrapper => wrapper.find('service-form-stub');
const modalId = wrapper => ({ id: wrapper.props().modal.id });

const withModalInject = (options = {}) => {
  const rawModal = options.propsData?.modal ?? { config: {} };
  const modal = rawModal.id ? rawModal : { id: 'test-modal-id', ...rawModal };

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

describe('create-service', () => {
  const $modals = {
    ...mockModals(),
    updateModalConfig: vi.fn(),
    updateDialogProps: vi.fn(),
  };
  const $popups = mockPopups();
  const $sidebar = mockSidebar();
  const consoleMock = mockConsole();

  const { templateVarsModule } = createTemplateVarsModule();
  const { authModule } = createAuthModule();
  const { userModule } = createUserModule();
  const { llmModule } = createLlmModule();
  const store = createMockedStoreModules([templateVarsModule, authModule, userModule, llmModule]);

  const defaultServiceForm = serviceToForm();
  const defaultService = formToService(defaultServiceForm);

  const shallowCreateService = generateShallowRenderer(CreateService, {
    stubs,
    attachTo: document.body,
    store,
    mocks: { $modals, $popups, $sidebar },
  });
  const factory = (options = {}) => shallowCreateService(withModalInject(options));

  test('Form submitted after trigger submit button', async () => {
    const action = vi.fn();
    const modal = {
      config: {
        action,
      },
    };
    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(defaultService);
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

    expect(action).toHaveBeenCalledWith(defaultService);
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
      category: 'Category error',
      enabled: 'Category error',
      infos: 'Category error',
    };
    const action = vi.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
    const modal = {
      config: {
        action,
      },
    };
    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(defaultService);
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = vi.fn().mockRejectedValue(errors);
    const customService = {
      ...defaultService,
      name: 'Custom name',
      category: {
        _id: 'custom-category',
      },
    };
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            item: customService,
            action,
          },
        },
      },
    });

    await wrapper.vm.submit();

    expect(consoleMock.error).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledWith({
      ...defaultService,
      name: customService.name,
      category: customService.category._id,
    });
    expect($modals.hide).not.toHaveBeenCalled();
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
      ...defaultServiceForm,
      name: Faker.datatype.string(),
      category: {
        _id: Faker.datatype.string(),
      },
    };

    await selectServiceForm(wrapper).triggerCustomEvent('input', newForm);
    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      ...defaultService,
      name: newForm.name,
      category: newForm.category._id,
    });
    expect($modals.hide).toBeCalledWith(modalId(wrapper));
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

    expect($modals.hide).toBeCalledWith(modalId(wrapper));
  });
});
