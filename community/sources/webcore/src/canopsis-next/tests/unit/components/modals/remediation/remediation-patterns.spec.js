import Faker from 'faker';
import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import ClickOutside from '@/services/click-outside';

import RemediationPatterns from '@/components/modals/remediation/remediation-patterns.vue';
import { PATTERN_CUSTOM_ITEM_VALUE } from '@/constants';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
  'c-patterns-field': true,
  'c-collapse-panel': true,
  'remediation-patterns-pbehavior-types-form': true,
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-patterns-field': true,
  'c-collapse-panel': true,
  'remediation-patterns-pbehavior-types-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');
const selectRemediationPatternsPbehaviorTypesForm = wrapper => wrapper
  .find('remediation-patterns-pbehavior-types-form-stub');

describe('remediation-patterns', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const factory = generateShallowRenderer(RemediationPatterns, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(RemediationPatterns, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  const alarmPattern = {
    id: PATTERN_CUSTOM_ITEM_VALUE,
    groups: [],
  };

  const entityPattern = {
    id: 'entity-pattern',
    is_corporate: true,
    groups: [],
  };

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
      active_on_pbh: [],
      disabled_on_pbh: [],
      alarm_pattern: [],
      entity_pattern: [],
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

    const patternsField = selectPatternsField(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => patternsField.vm,
      vm: patternsField.vm,
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
      alarm_pattern: 'Alarm pattern error',
      entity_pattern: 'Entity pattern error',
      active_on_pbh: 'Active on pbh error',
      disabled_on_pbh: 'Disabled on pbh error',
    };
    const action = jest.fn().mockRejectedValue({
      ...formErrors,
      unavailableField: 'Error',
    });
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
      active_on_pbh: [],
      disabled_on_pbh: [],
      alarm_pattern: [],
      entity_pattern: [],
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
    const customInstruction = {
      active_on_pbh: ['active-type'],
      disabled_on_pbh: ['disabled-type'],
      alarm_pattern: [],
      entity_pattern: [],
    };
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            instruction: customInstruction,
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
    expect(action).toBeCalledWith(customInstruction);
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

    const remediationPatternsPbehaviorTypesForm = selectRemediationPatternsPbehaviorTypesForm(wrapper);

    const newForm = {
      active_on_pbh: [Faker.datatype.string()],
      disabled_on_pbh: [Faker.datatype.string()],
      alarm_pattern: alarmPattern,
      entity_pattern: entityPattern,
    };

    remediationPatternsPbehaviorTypesForm.triggerCustomEvent('input', newForm);

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith({
      active_on_pbh: newForm.active_on_pbh,
      disabled_on_pbh: newForm.disabled_on_pbh,
      corporate_entity_pattern: newForm.entity_pattern.id,
      entity_pattern: newForm.entity_pattern.groups,
      alarm_pattern: newForm.alarm_pattern.groups,
    });
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

  test('Renders `remediation-patterns` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
      mocks: {
        $modals,
        $popups,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `remediation-patterns` with instruction', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            instruction: {
              active_on_pbh: ['active-type'],
              disabled_on_pbh: ['disabled-type'],
              alarm_pattern: [],
              entity_pattern: [],
            },
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `remediation-patterns` with title', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            title: 'Remediation patterns title',
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
