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

import { PATTERN_TYPES, PATTERNS_FIELDS } from '@/constants';

import ClickOutside from '@/services/click-outside';

import CreatePattern from '@/components/modals/pattern/create-pattern.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pattern-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

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

describe('create-pattern', () => {
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

  const shallowCreatePattern = generateShallowRenderer(CreatePattern, {
    stubs,
    store,
    attachTo: document.body,
  });
  const factory = (options = {}) => shallowCreatePattern(withModalInject(options));

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
      title: '',
      is_corporate: false,
      type: PATTERN_TYPES.alarm,
      alarm_pattern: [],
    });

    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Modal didn\'t close after trigger submit button with action error', async () => {
    const error = { unavailableField: 'Error' };
    const action = vi.fn().mockRejectedValue(error);
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();
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

    expect(action).toBeCalled();
    expect($modals.hide).not.toBeCalled();

    consoleErrorSpy.mockRestore();
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

  test('Modal didn\'t close after trigger submit button with form action errors', async () => {
    const formErrors = {
      title: 'Title error',
      is_corporate: 'Is corporate error',
      type: 'Type error',
    };
    const action = vi.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
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
        $popups,
        $sidebar,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      title: '',
      is_corporate: false,
      type: PATTERN_TYPES.alarm,
      alarm_pattern: [],
    });
    expect(afterSubmit).not.toHaveBeenCalled();

    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Error popup showed after trigger submit button with unavailable action errors', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = vi.fn().mockRejectedValue(errors);
    const customPattern = {
      title: Faker.datatype.string(),
      id: Faker.datatype.string(),
      type: PATTERN_TYPES.entity,
      is_corporate: true,
      entity_pattern: [],
    };
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            pattern: customPattern,
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
      entity_pattern: customPattern.entity_pattern,
      is_corporate: customPattern.is_corporate,
      type: customPattern.type,
      title: customPattern.title,
    });

    expect($modals.hide).not.toHaveBeenCalled();

    consoleErrorSpy.mockRestore();
  });

  test.each([
    { type: PATTERN_TYPES.alarm, expectedField: PATTERNS_FIELDS.alarm },
    { type: PATTERN_TYPES.entity, expectedField: PATTERNS_FIELDS.entity },
    { type: PATTERN_TYPES.pbehavior, expectedField: PATTERNS_FIELDS.pbehavior },
    { type: PATTERN_TYPES.serviceWeather, expectedField: PATTERNS_FIELDS.serviceWeather },
  ])(
    'Modal submitted with type: "$type" data after trigger form',
    async ({ type, expectedField }) => {
      const action = vi.fn();
      const newForm = {
        type,
        name: Faker.datatype.string(),
        groups: [],
      };
      const wrapper = factory({
        propsData: {
          modal: {
            config: {
              action,
              pattern: {
                title: newForm.name,
                type: newForm.type,
                [expectedField]: [],
              },
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
        type: newForm.type,
        title: newForm.name,
        is_corporate: false,
        [expectedField]: [],
      });
      expect($modals.hide).toBeCalledWith(modalId(wrapper));
    },
  );

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

    expect($modals.hide).toBeCalledWith(modalId(wrapper));
  });
});
