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

import { COLORS } from '@/config';

import ClickOutside from '@/services/click-outside';

import CreateTag from '@/components/modals/tag/create-tag.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pattern-progress': true,
  'tag-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const selectTagForm = wrapper => wrapper.find('tag-form-stub');
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

describe('create-tag', () => {
  const $modals = {
    ...mockModals(),
    updateModalConfig: vi.fn(),
    updateDialogProps: vi.fn(),
  };
  const $popups = mockPopups();
  const $sidebar = mockSidebar();

  const { authModule } = createAuthModule();
  const { userModule } = createUserModule();
  const { llmModule } = createLlmModule();
  const store = createMockedStoreModules([authModule, userModule, llmModule]);

  const shallowCreateTag = generateShallowRenderer(CreateTag, {
    stubs,
    store,
    attachTo: document.body,
  });
  const factory = (options = {}) => shallowCreateTag(withModalInject(options));

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
      mocks: {
        $modals,
        $sidebar,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      value: '',
      alarm_pattern: [],
      entity_pattern: [],
      color: COLORS.secondary,
    });
    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = vi.fn().mockRejectedValue({ value: 'Value error' });
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
      mocks: {
        $modals,
        $sidebar,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      value: '',
      color: COLORS.secondary,
      alarm_pattern: [],
      entity_pattern: [],
    });
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
      mocks: {
        $modals,
        $sidebar,
      },
    });

    await wrapper.vm.submit();

    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      value: 'Value error',
      color: 'Color error',
      patterns: 'Patterns error',
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
      mocks: {
        $modals,
        $sidebar,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      value: '',
      color: COLORS.secondary,
      alarm_pattern: [],
      entity_pattern: [],
    });
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = vi.fn().mockRejectedValue(errors);
    const customTag = {
      value: Faker.datatype.string(),
      color: Faker.internet.color(),
      alarm_pattern: [],
      entity_pattern: [],
    };
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            tag: customTag,
            action,
          },
        },
      },
      mocks: {
        $modals,
        $popups,
        $sidebar,
      },
    });

    await wrapper.vm.submit();

    expect(consoleErrorSpy).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledWith({
      entity_pattern: customTag.entity_pattern,
      alarm_pattern: customTag.alarm_pattern,
      color: customTag.color,
      value: customTag.value,
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
      mocks: {
        $modals,
        $sidebar,
      },
    });

    const newForm = {
      value: Faker.datatype.string(),
      color: Faker.internet.color(),
    };

    await selectTagForm(wrapper).triggerCustomEvent('input', newForm);
    await flushPromises();
    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith(newForm);
    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
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
        $sidebar,
      },
    });

    wrapper.vm.closeModal();
    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });
});
