import Faker from 'faker';
import { merge } from 'lodash';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  createAuthModule,
  createLlmModule,
  createMockedStoreModules,
  createTemplateVarsModule,
  createUserModule,
} from '@unit/utils/store';
import { mockConsole, mockModals, mockPopups, mockSidebar } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';

import ClickOutside from '@/services/click-outside';

import { formToService, serviceToForm } from '@/helpers/entities/service/form';

import CreateService from '@/components/modals/service/create-service.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pattern-progress': true,
  'service-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pattern-progress': true,
  'service-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectServiceForm = wrapper => wrapper.find('service-form-stub');

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
    updateModalConfig: jest.fn(),
    updateDialogProps: jest.fn(),
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

  const renderCreateService = generateRenderer(CreateService, {
    stubs: snapshotStubs,
    store,
    mocks: { $modals, $popups, $sidebar },
  });
  const snapshotFactory = (options = {}) => renderCreateService(withModalInject(options));

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledWith(defaultService);
    expect($modals.hide).toHaveBeenCalledWith(wrapper.props().modal);
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
    });

    const serviceForm = selectServiceForm(wrapper);
    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      vm: serviceForm.vm,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).not.toHaveBeenCalled();
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith(wrapper.props().modal);
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      name: 'Name error',
      category: 'Category error',
      enabled: 'Category error',
      infos: 'Category error',
    };
    const action = jest.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toHaveBeenCalledWith(defaultService);
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = jest.fn().mockRejectedValue(errors);
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

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
    const action = jest.fn();
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
    await selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledWith({
      ...defaultService,
      name: newForm.name,
      category: newForm.category._id,
    });
    expect($modals.hide).toBeCalledWith(wrapper.props().modal);
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
    });

    selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalledWith(wrapper.props().modal);
  });

  test('Renders `create-service` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `create-service` with pbehavior', () => {
    const service = {
      ...defaultService,
      name: 'Service name',
    };
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            item: service,
            title: 'Custom create service title',
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
