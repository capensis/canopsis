import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';

import ClickOutside from '@/services/click-outside';

import CreatePbehavior from '@/components/modals/pbehavior/create-pbehavior.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pbehavior-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'pbehavior-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectPbehaviorForm = wrapper => wrapper.find('pbehavior-form-stub');

describe('create-pbehavior', () => {
  const $modals = mockModals();
  const $popups = mockPopups();
  const defaultPbehavior = {
    _id: expect.any(String),
    name: '',
    color: '',
    comments: [],
    enabled: true,
    inherited: false,
    entity_pattern: [],
    exceptions: [],
    exdates: [],
    reason: undefined,
    type: undefined,
    rrule: null,
    tstart: null,
    tstop: null,
    timezone: 'Europe/Paris',
  };

  const defaultPbehaviorRequest = {
    ...defaultPbehavior,
    exec_pattern: true,
  };

  const factory = generateShallowRenderer(CreatePbehavior, {

    stubs,
    attachTo: document.body,
    mocks: { $modals, $popups },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(CreatePbehavior, {

    stubs: snapshotStubs,
    mocks: { $modals, $popups },
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
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledWith(defaultPbehaviorRequest);
    expect($modals.hide).toHaveBeenCalledWith();
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

    const pbehaviorForm = selectPbehaviorForm(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => pbehaviorForm.vm,
      vm: pbehaviorForm.vm,
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

    expect($modals.hide).toHaveBeenCalledWith();
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      name: 'Name error',
      rrule: 'Rrule error',
      tstart: 'Tstart error',
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
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(action).toHaveBeenCalledWith(defaultPbehaviorRequest);
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
    const customPbehavior = {
      name: Faker.datatype.string(),
      _id: Faker.datatype.string(),
      entity_pattern: [],
    };
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            pbehavior: customPbehavior,
            action,
          },
        },
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(consoleErrorSpy).toHaveBeenCalledWith(errors);
    expect($popups.error).toHaveBeenCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(action).toHaveBeenCalledWith({
      ...defaultPbehaviorRequest,
      entity_pattern: customPbehavior.entity_pattern,
      name: customPbehavior.name,
      _id: customPbehavior._id,
    });
    expect($modals.hide).not.toHaveBeenCalledWith();

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
    });

    const newForm = {
      name: Faker.datatype.string(),
      _id: Faker.datatype.string(),
      color: Faker.internet.color(),
      entity_pattern: [],
      rrule: null,
      timezone: 'Europe/Paris',
    };

    selectPbehaviorForm(wrapper).triggerCustomEvent('input', { ...defaultPbehavior, ...newForm });
    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toHaveBeenCalledWith({
      ...defaultPbehaviorRequest,
      ...newForm,
    });
    expect($modals.hide).toHaveBeenCalled();
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

    expect($modals.hide).toHaveBeenCalled();
  });

  test('Renders `create-pbehavior` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `create-pbehavior` with pbehavior', () => {
    const pbehavior = {
      name: 'Pbehavior name',
      rrule: 'FREQ=DAILY',
      entity_pattern: [],
    };
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            pbehavior,
            title: 'Custom create pbehavior title',
            noPattern: true,
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
