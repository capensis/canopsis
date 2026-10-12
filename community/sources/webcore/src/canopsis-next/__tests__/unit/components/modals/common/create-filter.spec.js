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
import { merge } from 'lodash';

import { ALARM_PATTERN_FIELDS, PATTERN_CONDITIONS } from '@/constants';

import ClickOutside from '@/services/click-outside';

import CreateFilter from '@/components/modals/common/create-filter.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pattern-progress': true,
  'patterns-form': true,
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

describe('create-filter', () => {
  const $modals = {
    ...mockModals(),
    updateModalConfig: vi.fn(),
    updateDialogProps: vi.fn(),
  };
  const $sidebar = mockSidebar();
  const $popups = mockPopups();

  const { authModule } = createAuthModule();
  const { userModule } = createUserModule();
  const { llmModule } = createLlmModule();
  const store = createMockedStoreModules([authModule, userModule, llmModule]);

  const defaultPattern = {
    field: '',
    cond: {
      type: PATTERN_CONDITIONS.equal,
      value: '',
    },
  };
  const customPattern = {
    field: ALARM_PATTERN_FIELDS.component,
    cond: {
      type: PATTERN_CONDITIONS.notEqual,
      value: 'component',
    },
  };

  const shallowCreateFilter = generateShallowRenderer(CreateFilter, {
    stubs,
    store,
    attachTo: document.body,
  });
  const factory = (options = {}) => shallowCreateFilter(withModalInject(options));

  test('Form submitted without fields after trigger submit button', async () => {
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
      mocks: {
        $modals,
        $sidebar,
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      is_user_preference: false,
      title: '',
    });
    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Form submitted with all fields after trigger submit button', async () => {
    const action = vi.fn();
    const modal = {
      config: {
        withTitle: true,
        withEntity: true,
        withPbehavior: true,
        withAlarm: true,
        withEvent: true,
        action,
      },
    };
    const wrapper = factory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
        $sidebar,
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      is_user_preference: false,
      title: '',
      alarm_pattern: [],
      entity_pattern: [],
      pbehavior_pattern: [],
      event_pattern: [],
    });
    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Form didn\'t close after trigger submit button with action error', async () => {
    const action = vi.fn().mockRejectedValue({ unavailableField: 'Error' });
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
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalled();
    expect($modals.hide).not.toHaveBeenCalled();

    consoleErrorSpy.mockRestore();
  });

  test('Form submitted after trigger submit button without action', async () => {
    const modal = {
      config: {},
    };
    const wrapper = factory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
        $sidebar,
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      title: 'Title error',
      alarm_pattern: 'Alarm pattern error',
      entity_pattern: 'Entity pattern error',
      pbehavior_pattern: 'PBehavior pattern error',
    };
    const action = vi.fn().mockRejectedValue({
      ...formErrors,
      unavailableField: 'Error',
    });
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            withAlarm: true,
            withEntity: true,
            withPbehavior: true,
            withServiceWeather: true,

            action,
          },
        },
      },
      mocks: {
        $modals,
        $sidebar,
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      is_user_preference: false,
      alarm_pattern: [],
      entity_pattern: [],
      pbehavior_pattern: [],
      weather_service_pattern: [],
      title: '',
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
    const customFilter = {
      title: 'Title',
      is_user_preference: true,
      alarm_pattern: [[customPattern]],
    };
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            filter: customFilter,
            withAlarm: true,
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
    expect(action).toHaveBeenCalledWith(customFilter);
    expect($modals.hide).not.toHaveBeenCalled();

    consoleErrorSpy.mockRestore();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = vi.fn();
    const filter = {
      title: 'Title',
      is_user_preference: true,
      alarm_pattern: [[customPattern]],
      entity_pattern: [[defaultPattern]],
    };
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            withTitle: true,
            withEntity: true,
            withAlarm: true,
            filter,

            action,
          },
        },
      },
      mocks: {
        $modals,
        $sidebar,
        $popups,
      },
    });

    await wrapper.vm.submit();

    expect(action).toHaveBeenCalledWith({
      alarm_pattern: filter.alarm_pattern,
      entity_pattern: filter.entity_pattern,
      is_user_preference: filter.is_user_preference,
      title: filter.title,
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
      mocks: {
        $modals,
        $sidebar,
        $popups,
      },
    });

    wrapper.vm.closeModal();
    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith(modalId(wrapper));
  });
});
