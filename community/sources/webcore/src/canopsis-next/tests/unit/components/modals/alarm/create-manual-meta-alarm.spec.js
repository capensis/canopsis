import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createMockedStoreModules } from '@unit/utils/store';
import ClickOutside from '@/services/click-outside';
import {
  ENTITIES_STATES,
  ENTITIES_TYPES,
  EVENT_ENTITY_TYPES,
  MANUAL_META_ALARM_EVENT_DEFAULT_FIELDS,
} from '@/constants';

import CreateManualMetaAlarm from '@/components/modals/alarm/create-manual-meta-alarm.vue';

const localVue = createVueInstance();

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
  'manual-meta-alarm-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
  'manual-meta-alarm-form': true,
};

const factory = (options = {}) => shallowMount(CreateManualMetaAlarm, {
  localVue,
  stubs,
  attachTo: document.body,
  propsData: {
    modal: {
      config: {},
    },
  },

  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(CreateManualMetaAlarm, {
  localVue,
  stubs: snapshotStubs,
  propsData: {
    modal: {
      config: {},
    },
  },

  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectManualMetaAlarmForm = wrapper => wrapper.find('manual-meta-alarm-form-stub');

describe('create-manual-meta-alarm', () => {
  const timestamp = 1386435600000;

  mockDateNow(timestamp);
  const $modals = mockModals();
  const $popups = mockPopups();

  const alarm = {
    _id: Faker.datatype.string(),
    v: {
      connector: Faker.datatype.string(),
      connector_name: Faker.datatype.string(),
      component: Faker.datatype.string(),
      resource: Faker.datatype.string(),
      state: {
        val: Faker.datatype.number(),
      },
      status: {
        val: ENTITIES_STATES.ok,
      },
    },
    entity: {
      _id: Faker.datatype.number(),
      type: Faker.datatype.number(),
    },
  };
  const itemsType = ENTITIES_TYPES.alarm;
  const itemsIds = [alarm._id];
  const manualMetaAlarmEventData = {
    ...MANUAL_META_ALARM_EVENT_DEFAULT_FIELDS,

    event_type: EVENT_ENTITY_TYPES.manualMetaAlarmUpdate,
    state: ENTITIES_STATES.minor,

    output: '',
  };

  const getEntitiesList = jest.fn().mockReturnValue([alarm]);
  const entitiesModule = {
    name: 'entities',
    getters: {
      getList: () => getEntitiesList,
    },
  };

  const createEvent = jest.fn();
  const eventModule = {
    name: 'event',
    actions: {
      create: createEvent,
    },
  };
  const store = createMockedStoreModules([entitiesModule, eventModule]);

  afterEach(() => {
    createEvent.mockClear();
    getEntitiesList.mockClear();
  });

  test('Default parameters applied to form', () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    const manualMetaAlarmForm = selectManualMetaAlarmForm(wrapper);

    expect(manualMetaAlarmForm.vm.form).toEqual({
      metaAlarm: null,
      output: '',
    });
  });

  test('Form submitted after trigger submit button', async () => {
    const afterSubmit = jest.fn();
    const config = {
      itemsType,
      itemsIds,
      afterSubmit,
    };
    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config,
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);
    const manualMetaAlarmForm = selectManualMetaAlarmForm(wrapper);

    const newData = {
      output: Faker.datatype.string(),
      metaAlarm: Faker.datatype.string(),
    };

    manualMetaAlarmForm.vm.$emit('input', newData);

    submitButton.trigger('click');

    await flushPromises();

    expect(createEvent).toBeCalledTimes(1);
    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          ...manualMetaAlarmEventData,
          display_name: newData.metaAlarm,
          output: newData.output,
          ma_children: [alarm.entity._id],
          event_type: EVENT_ENTITY_TYPES.manualMetaAlarmGroup,
        }],
      },
      undefined,
    );
    expect(afterSubmit).toBeCalled();
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);
    const manualMetaAlarmForm = selectManualMetaAlarmForm(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => manualMetaAlarmForm.vm,
      vm: manualMetaAlarmForm.vm,
    });

    submitButton.trigger('click');

    await flushPromises();

    expect(createEvent).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();

    validator.detach('name');
  });

  test('Errors added after trigger submit button with action errors', async () => {
    const formErrors = {
      output: 'Output error',
      metaAlarm: 'Meta alarm error',
    };
    createEvent.mockRejectedValueOnce({ ...formErrors, unavailableField: 'Error' });
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);
    const manualMetaAlarmForm = selectManualMetaAlarmForm(wrapper);

    const newData = {
      output: Faker.datatype.string(),
      metaAlarm: {
        _id: Faker.datatype.string(),
      },
    };

    manualMetaAlarmForm.vm.$emit('input', newData);

    submitButton.trigger('click');

    await flushPromises();

    const addedErrors = wrapper.getValidatorErrorsObject();

    expect(formErrors).toEqual(addedErrors);
    expect(createEvent).toBeCalledTimes(1);
    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          ...manualMetaAlarmEventData,
          ma_parents: [newData.metaAlarm._id],
          ma_children: [alarm.entity._id],
          output: newData.output,
          event_type: EVENT_ENTITY_TYPES.manualMetaAlarmUpdate,
        }],
      },
      undefined,
    );
    expect($modals.hide).not.toBeCalledWith();
  });

  test('Error popup showed after trigger submit button with action errors', async () => {
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    const errors = {
      unavailableField: 'Error',
      anotherUnavailableField: 'Second error',
    };
    createEvent.mockRejectedValueOnce(errors);

    const wrapper = factory({
      store: createMockedStoreModules([
        entitiesModule,
        eventModule,
      ]),
      mocks: {
        $modals,
        $popups,
      },
    });

    const submitButton = selectSubmitButton(wrapper);
    const manualMetaAlarmForm = selectManualMetaAlarmForm(wrapper);

    const newData = {
      output: Faker.datatype.string(),
      metaAlarm: Faker.datatype.string(),
    };

    manualMetaAlarmForm.vm.$emit('input', newData);

    submitButton.trigger('click');

    await flushPromises();

    expect(consoleErrorSpy).toBeCalledWith(errors);
    expect($popups.error).toBeCalledWith({
      text: `${errors.unavailableField}\n${errors.anotherUnavailableField}`,
    });
    expect(createEvent).toBeCalledTimes(1);
    expect(createEvent).toBeCalledWith(
      expect.any(Object),
      {
        data: [{
          ...manualMetaAlarmEventData,
          display_name: newData.metaAlarm,
          ma_children: [alarm.entity._id],
          output: newData.output,
          event_type: EVENT_ENTITY_TYPES.manualMetaAlarmGroup,
        }],
      },
      undefined,
    );
    expect($modals.hide).not.toBeCalledWith();

    consoleErrorSpy.mockClear();
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    const cancelButton = selectCancelButton(wrapper);

    cancelButton.trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-manual-meta-alarm` with empty modal', () => {
    const wrapper = snapshotFactory({
      store,
      mocks: {
        $modals,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
