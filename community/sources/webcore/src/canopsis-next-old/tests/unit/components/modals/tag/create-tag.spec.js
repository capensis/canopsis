import Faker from 'faker';
import { merge } from 'lodash';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAuthModule, createLlmModule, createMockedStoreModules, createUserModule } from '@unit/utils/store';
import { mockModals, mockPopups, mockSidebar } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';

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

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pattern-progress': true,
  'tag-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectTagForm = wrapper => wrapper
  .find('tag-form-stub');

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
    updateModalConfig: jest.fn(),
    updateDialogProps: jest.fn(),
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

  const renderCreateTag = generateRenderer(CreateTag, {
    stubs: snapshotStubs,
    store,
  });
  const snapshotFactory = (options = {}) => renderCreateTag(withModalInject(options));

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledWith({
      value: '',
      alarm_pattern: [],
      entity_pattern: [],
      color: COLORS.secondary,
    });
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
      mocks: {
        $modals,
        $sidebar,
      },
    });

    const tagForm = selectTagForm(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'value',
      rules: 'required:true',
      getter: () => false,
      context: () => tagForm.vm,
      vm: tagForm.vm,
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
      mocks: {
        $modals,
        $sidebar,
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith(wrapper.props().modal);
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      value: 'Value error',
      color: 'Color error',
      patterns: 'Patterns error',
    };
    const action = jest.fn().mockRejectedValue({ ...formErrors, unavailableField: 'Error' });
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toHaveBeenCalledWith({
      value: '',
      color: COLORS.secondary,
      alarm_pattern: [],
      entity_pattern: [],
    });
    expect($modals.hide).not.toHaveBeenCalled();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error')
      .mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = jest.fn().mockRejectedValue(errors);
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

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

    consoleErrorSpy.mockClear();
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
      mocks: {
        $modals,
        $sidebar,
      },
    });

    const newForm = {
      value: Faker.datatype.string(),
      color: Faker.internet.color(),
    };

    selectTagForm(wrapper).triggerCustomEvent('input', newForm);
    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledWith(newForm);
    expect($modals.hide).toHaveBeenCalledWith(wrapper.props().modal);
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

    selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toHaveBeenCalledWith(wrapper.props().modal);
  });

  test('Renders `create-tag` with empty modal', () => {
    const wrapper = snapshotFactory({
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

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `create-tag` with tag', () => {
    const tag = {
      value: 'Value',
      color: COLORS.primary,
      alarm_pattern: [],
      entity_pattern: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            title: 'create-tag title',
            tag,
          },
        },
      },
      mocks: {
        $modals,
        $sidebar,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
