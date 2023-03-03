import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createMockedStoreModules, createManualMetaAlarmModule } from '@unit/utils/store';
import ClickOutside from '@/services/click-outside';
import { ENTITIES_STATES } from '@/constants';

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
  const items = [alarm];
  const config = { items };

  const {
    manualMetaAlarmModule,
    createManualMetaAlarm,
    addAlarmsIntoManualMetaAlarm,
  } = createManualMetaAlarmModule();

  const store = createMockedStoreModules([manualMetaAlarmModule]);

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
      comment: '',
    });
  });

  test('Form submitted after trigger submit button', async () => {
    const afterSubmit = jest.fn();

    const wrapper = factory({
      store,
      propsData: {
        modal: {
          config: {
            items,
            afterSubmit,
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    const submitButton = selectSubmitButton(wrapper);
    const manualMetaAlarmForm = selectManualMetaAlarmForm(wrapper);

    const newData = {
      comment: Faker.datatype.string(),
      metaAlarm: Faker.datatype.string(),
    };

    manualMetaAlarmForm.vm.$emit('input', newData);

    submitButton.trigger('click');

    await flushPromises();

    expect(createManualMetaAlarm).toBeCalledTimes(1);
    expect(createManualMetaAlarm).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          name: newData.metaAlarm,
          alarms: [alarm._id],
          comment: newData.comment,
        },
      },
      undefined,
    );
    expect(afterSubmit).toBeCalled();
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
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

    expect(createManualMetaAlarm).not.toBeCalled();
    expect(addAlarmsIntoManualMetaAlarm).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();

    validator.detach('name');
  });

  test('Modal hidden after trigger cancel button', async () => {
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

    const cancelButton = selectCancelButton(wrapper);

    cancelButton.trigger('click');

    await flushPromises();

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-manual-meta-alarm` with empty modal', () => {
    const wrapper = snapshotFactory({
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
