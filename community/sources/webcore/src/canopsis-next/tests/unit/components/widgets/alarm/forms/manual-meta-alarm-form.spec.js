import Faker from 'faker';
import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';
import { createMockedStoreModules } from '@unit/utils/store';

import ManualMetaAlarmForm from '@/components/widgets/alarm/forms/manual-meta-alarm-form.vue';

const stubs = {
  'v-combobox': createInputStub('v-combobox'),
  'v-text-field': createInputStub('v-text-field'),
  'c-enabled-field': true,
};

const snapshotStubs = {
  'c-enabled-field': true,
};

const selectTextField = wrapper => wrapper.find('.v-text-field');
const selectComboboxField = wrapper => wrapper.find('.v-combobox');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');

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

  const factory = generateShallowRenderer(ManualMetaAlarmForm, { stubs });
  const snapshotFactory = generateRenderer(ManualMetaAlarmForm, { stubs: snapshotStubs });

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
      comment: Faker.datatype.string(),
      auto_resolve: false,
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

  test('Comment changed after trigger description field', () => {
    const form = {
      metaAlarm: Faker.datatype.string(),
      comment: Faker.datatype.string(),
      auto_resolve: false,
    };
    const wrapper = factory({
      store,
      propsData: {
        form,
      },
    });

    const comment = Faker.datatype.string();

    selectTextField(wrapper).setValue(comment);

    expect(wrapper).toEmit('input', { ...form, comment });
  });

  test('Auto resolve changed after trigger enabled field', () => {
    const form = {
      metaAlarm: Faker.datatype.string(),
      comment: Faker.datatype.string(),
      auto_resolve: false,
    };
    const wrapper = factory({
      store,
      propsData: {
        form,
      },
    });

    const autoResolve = true;

    const enabledField = selectEnabledField(wrapper);

    enabledField.triggerCustomEvent('input', autoResolve);

    expect(wrapper).toEmit('input', { ...form, auto_resolve: autoResolve });
  });

  test('Renders `manual-meta-alarm-form` with default props', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        form: {
          metaAlarm: 'metaAlarm',
          comment: 'comment',
          auto_resolve: true,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
