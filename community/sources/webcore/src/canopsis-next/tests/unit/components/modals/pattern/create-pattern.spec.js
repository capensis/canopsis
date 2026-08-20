import Faker from 'faker';
import { merge } from 'lodash';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAuthModule, createLlmModule, createMockedStoreModules, createUserModule } from '@unit/utils/store';
import { mockModals, mockPopups, mockSidebar } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';

import { PATTERN_TYPES, PATTERNS_FIELDS } from '@/constants';

import ClickOutside from '@/services/click-outside';

import CreatePattern from '@/components/modals/pattern/create-pattern.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pattern-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pattern-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectPatternForm = wrapper => wrapper
  .find('pattern-form-stub');

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
    updateModalConfig: jest.fn(),
    updateDialogProps: jest.fn(),
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

  const renderCreatePattern = generateRenderer(CreatePattern, {
    stubs: snapshotStubs,
    store,
  });
  const snapshotFactory = (options = {}) => renderCreatePattern(withModalInject(options));

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
      title: '',
      is_corporate: false,
      type: PATTERN_TYPES.alarm,
      alarm_pattern: [],
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

    const patternForm = selectPatternForm(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => patternForm.vm,
      vm: patternForm.vm,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();

    validator.detach('name');
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
      title: 'Title error',
      is_corporate: 'Is corporate error',
      type: 'Type error',
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
      title: '',
      is_corporate: false,
      type: PATTERN_TYPES.alarm,
      alarm_pattern: [],
    });

    expect($modals.hide).not.toHaveBeenCalledWith();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error')
      .mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    const action = jest.fn().mockRejectedValue(errors);
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

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

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

    expect($modals.hide).not.toHaveBeenCalledWith();

    consoleErrorSpy.mockClear();
  });

  test.each([
    { type: PATTERN_TYPES.alarm, expectedField: PATTERNS_FIELDS.alarm },
    { type: PATTERN_TYPES.entity, expectedField: PATTERNS_FIELDS.entity },
    { type: PATTERN_TYPES.pbehavior, expectedField: PATTERNS_FIELDS.pbehavior },
    { type: PATTERN_TYPES.serviceWeather, expectedField: PATTERNS_FIELDS.serviceWeather },
  ])('Modal submitted with type: "$type" data after trigger form', async ({ type, expectedField }) => {
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
      type,
      name: Faker.datatype.string(),
      groups: [],
    };

    selectPatternForm(wrapper).triggerCustomEvent('input', newForm);

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledWith({
      type: newForm.type,
      name: newForm.name,
      [expectedField]: [],
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
      mocks: {
        $modals,
        $sidebar,
      },
    });

    selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalledWith(wrapper.props().modal);
  });

  test('Renders `create-pattern` with empty modal', () => {
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

  test('Renders `create-pattern` with pattern', () => {
    const pattern = {
      title: 'Title',
      id: 'Id',
      type: PATTERN_TYPES.alarm,
      is_corporate: true,
      alarm_pattern: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            type: PATTERN_TYPES.alarm,
            pattern,
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

  test('Renders `create-pattern` with hidden title', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            hideTitle: true,
            text: 'create-pattern text',
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
