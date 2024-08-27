import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';

import { ALARM_STATES } from '@/constants';

import ClickOutside from '@/services/click-outside';

import LinkToMetaAlarm from '@/components/modals/alarm/link-to-meta-alarm.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
  'link-meta-alarm-form': true,
  'c-name-field': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'alarm-general-table': true,
  'link-meta-alarm-form': true,
};

const selectButtons = wrapper => wrapper.findAll('button.v-btn');
const selectSubmitButton = wrapper => selectButtons(wrapper).at(1);
const selectCancelButton = wrapper => selectButtons(wrapper).at(0);
const selectLinkMetaAlarmForm = wrapper => wrapper.find('link-meta-alarm-form-stub');

describe('link-to-meta-alarm', () => {
  const timestamp = 1386435600000;
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
        val: ALARM_STATES.ok,
      },
    },
    entity: {
      _id: Faker.datatype.number(),
      type: Faker.datatype.number(),
    },
  };
  const items = [alarm];
  const config = { items };

  const factory = generateShallowRenderer(LinkToMetaAlarm, {
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
  const snapshotFactory = generateRenderer(LinkToMetaAlarm, {
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

  beforeAll(() => jest.useFakeTimers({ now: timestamp }));

  test('Default parameters applied to form', () => {
    const wrapper = factory({
      mocks: {
        $modals,
      },
    });

    const manualMetaAlarmForm = selectLinkMetaAlarmForm(wrapper);

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
    const manualMetaAlarmForm = selectLinkMetaAlarmForm(wrapper);

    const newData = {
      comment: Faker.datatype.string(),
      metaAlarm: Faker.datatype.string(),
      auto_resolve: true,
    };

    manualMetaAlarmForm.triggerCustomEvent('input', newData);

    submitButton.trigger('click');

    await flushPromises(true);

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
    const manualMetaAlarmForm = selectLinkMetaAlarmForm(wrapper);

    const validator = wrapper.getValidator();

    validator.attach({
      name: 'name',
      rules: 'required:true',
      getter: () => false,
      context: () => manualMetaAlarmForm.vm,
      vm: manualMetaAlarmForm.vm,
    });

    submitButton.trigger('click');

    await flushPromises(true);

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

    await flushPromises(true);

    expect($modals.hide).toBeCalled();
  });

  test('Renders `link-to-meta-alarm` with empty modal', () => {
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
