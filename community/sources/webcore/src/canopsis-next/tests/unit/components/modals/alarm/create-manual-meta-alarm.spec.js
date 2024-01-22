import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import ClickOutside from '@/services/click-outside';
import { ENTITIES_STATES } from '@/constants';

import CreateManualMetaAlarm from '@/components/modals/alarm/create-manual-meta-alarm.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
  'manual-meta-alarm-form': true,
  'c-name-field': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
  'manual-meta-alarm-form': true,
};

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

  const factory = generateShallowRenderer(CreateManualMetaAlarm, {
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
  });
  const snapshotFactory = generateRenderer(CreateManualMetaAlarm, {
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
  });

  test('Default parameters applied to form', () => {
    const wrapper = factory({
      mocks: {
        $modals,
      },
    });

    const manualMetaAlarmForm = selectManualMetaAlarmForm(wrapper);

    expect(manualMetaAlarmForm.vm.form).toEqual({
      metaAlarm: null,
      comment: '',
      auto_resolve: false,
    });
  });

  test('Form submitted after trigger submit button', async () => {
    const action = jest.fn();

    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            action,
            items,
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
      auto_resolve: true,
    };

    manualMetaAlarmForm.vm.$emit('input', newData);

    submitButton.trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith({
      name: newData.metaAlarm,
      alarms: [alarm._id],
      comment: newData.comment,
      auto_resolve: newData.auto_resolve,
    });
    expect($modals.hide).toBeCalledWith();
  });

  test('Form didn\'t submitted after trigger submit button with error', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            ...config,
            action,
          },
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

    expect(action).not.toBeCalled();
    expect($modals.hide).not.toBeCalled();

    validator.detach('name');
  });

  test('Modal hidden after trigger cancel button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config,
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

  test('Renders `create-manual-meta-alarm` with empty modal', () => {
    const action = jest.fn();
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            ...config,
            action,
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
