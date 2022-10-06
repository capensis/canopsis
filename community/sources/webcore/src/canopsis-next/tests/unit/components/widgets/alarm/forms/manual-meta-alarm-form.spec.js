import Faker from 'faker';
import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';
import { createMockedStoreModules } from '@unit/utils/store';

import ManualMetaAlarmForm from '@/components/widgets/alarm/forms/manual-meta-alarm-form.vue';

const localVue = createVueInstance();

const stubs = {
  'v-combobox': createInputStub('v-combobox'),
  'v-text-field': createInputStub('v-text-field'),
};

const factory = (options = {}) => shallowMount(ManualMetaAlarmForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(ManualMetaAlarmForm, {
  localVue,

  ...options,
});

const selectTextField = wrapper => wrapper.find('.v-text-field');
const selectComboboxField = wrapper => wrapper.find('.v-combobox');

describe('manual-meta-alarm-form', () => {
  const fetchManualMetaAlarmsListWithoutStore = jest.fn().mockReturnValue([]);
  const alarmModule = {
    name: 'alarm',
    actions: {
      fetchManualMetaAlarmsListWithoutStore,
    },
  };
  const store = createMockedStoreModules([
    alarmModule,
  ]);

  test('Alarms fetched after mount', () => {
    factory({
      store,
      propsData: {
        form: {},
      },
    });

    expect(fetchManualMetaAlarmsListWithoutStore).toBeCalled();
  });

  test('Meta alarm changed after trigger text field', () => {
    const form = {
      metaAlarm: Faker.datatype.string(),
      output: Faker.datatype.string(),
    };
    const wrapper = factory({
      store,
      propsData: {
        form,
      },
    });

    const metaAlarm = Faker.datatype.string();

    const comboboxField = selectComboboxField(wrapper);

    comboboxField.setValue(metaAlarm);

    expect(wrapper).toEmit('input', { ...form, metaAlarm });
  });

  test('Output changed after trigger description field', () => {
    const form = {
      metaAlarm: Faker.datatype.string(),
      output: Faker.datatype.string(),
    };
    const wrapper = factory({
      store,
      propsData: {
        form,
      },
    });

    const output = Faker.datatype.string();

    const textField = selectTextField(wrapper);

    textField.setValue(output);

    expect(wrapper).toEmit('input', { ...form, output });
  });

  test('Renders `manual-meta-alarm-form` with default props', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        form: {
          metaAlarm: 'metaAlarm',
          output: 'output',
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `manual-meta-alarm-form` with alarms', async () => {
    fetchManualMetaAlarmsListWithoutStore.mockReturnValueOnce([
      { _id: 'entity-id', name: 'alarm-display-name' },
    ]);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        alarmModule,
      ]),
      propsData: {
        form: {},
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});