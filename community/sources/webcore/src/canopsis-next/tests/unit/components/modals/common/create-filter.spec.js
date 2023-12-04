import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { ALARM_PATTERN_FIELDS, PATTERN_CONDITIONS, PATTERN_CUSTOM_ITEM_VALUE, PATTERN_OPERATORS } from '@/constants';
import CreateFilter from '@/components/modals/common/create-filter.vue';
import ClickOutside from '@/services/click-outside';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'patterns-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'patterns-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectPatternsForm = wrapper => wrapper.find('patterns-form-stub');

describe('create-filter', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

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

  const factory = generateShallowRenderer(CreateFilter, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(CreateFilter, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  test('Form submitted without fields after trigger submit button', async () => {
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
      is_user_preference: false,
      title: '',
    });
    expect($modals.hide).toBeCalledWith();
  });

  test('Form submitted with all fields after trigger submit button', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            withTitle: true,
            withEntity: true,
            withPbehavior: true,
            withAlarm: true,
            withEvent: true,
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
      is_user_preference: false,
      title: '',
      alarm_pattern: [],
      entity_pattern: [],
      event_pattern: [],
      pbehavior_pattern: [],
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

    const patternsForm = selectPatternsForm(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => patternsForm.vm,
      vm: patternsForm.vm,
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
      alarm_pattern: 'Alarm pattern error',
      entity_pattern: 'Entity pattern error',
      pbehavior_pattern: 'PBehavior pattern error',
      weather_service_pattern: 'Weather service pattern error',
    };
    const action = jest.fn().mockRejectedValue({
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
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(addedErrors).toEqual(formErrors);
    expect(action).toBeCalledWith({
      is_user_preference: false,
      alarm_pattern: [],
      entity_pattern: [],
      pbehavior_pattern: [],
      weather_service_pattern: [],
      title: '',
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
    const customFilter = {
      title: 'Title',
      is_user_preference: true,
      alarm_pattern: [
        [customPattern],
      ],
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
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toBeCalledWith(customFilter);
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal submitted with correct data after trigger form', async () => {
    const action = jest.fn();
    const filter = {
      title: 'Title',
      is_user_preference: true,
      alarm_pattern: [
        [customPattern],
      ],
      entity_pattern: [
        [defaultPattern],
      ],
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
      },
    });

    const patternsForm = selectPatternsForm(wrapper);

    const newForm = {
      title: filter.title,
      is_user_preference: filter.is_user_preference,
      alarm_pattern: {
        id: PATTERN_CUSTOM_ITEM_VALUE,
        groups: [{
          rules: [{
            attribute: ALARM_PATTERN_FIELDS.ack,
            operator: PATTERN_OPERATORS.acked,
          }],
        }],
      },
      entity_pattern: {
        id: PATTERN_CUSTOM_ITEM_VALUE,
        groups: [],
      },
    };

    patternsForm.vm.$emit('input', newForm);

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith({
      alarm_pattern: [
        [
          {
            field: ALARM_PATTERN_FIELDS.ack,
            cond: {
              type: PATTERN_CONDITIONS.exist,
              value: true,
            },
          },
        ],
      ],
      entity_pattern: [],
      is_user_preference: filter.is_user_preference,
      title: filter.title,
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

  test('Renders `create-filter` with empty modal', () => {
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

  test('Renders `create-filter` with all parameters', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            withTitle: true,
            withEntity: true,
            withPbehavior: true,
            withAlarm: true,
            withEvent: true,
            title: 'Create filter title',
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
