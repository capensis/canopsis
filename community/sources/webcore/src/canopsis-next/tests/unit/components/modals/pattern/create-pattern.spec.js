import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';

import { PATTERN_TYPES } from '@/constants';

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

describe('create-pattern', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(CreatePattern, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(CreatePattern, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

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
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith({
      title: '',
      is_corporate: false,
      type: PATTERN_TYPES.alarm,
      alarm_pattern: [],
    });
    expect($modals.hide).toBeCalledWith();
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
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalledWith();
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
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toBeCalledWith({
      title: '',
      is_corporate: false,
      type: PATTERN_TYPES.alarm,
      alarm_pattern: [],
    });
    expect($modals.hide).not.toBeCalledWith();
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
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toBeCalledWith({
      entity_pattern: customPattern.entity_pattern,
      is_corporate: customPattern.is_corporate,
      type: customPattern.type,
      title: customPattern.title,
    });
    expect($modals.hide).not.toBeCalledWith();

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
      },
    });

    const patternForm = selectPatternForm(wrapper);

    const newForm = {
      all: false,
      auto: false,
      manual: false,
      with: false,
      instructions: [{}],
    };

    patternForm.triggerCustomEvent('input', newForm);

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith(newForm);
    expect($modals.hide).toBeCalled();
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
      },
    });

    selectCancelButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalled();
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
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
