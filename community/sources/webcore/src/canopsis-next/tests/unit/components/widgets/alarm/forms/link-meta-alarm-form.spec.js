import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';
import { createMockedStoreModules, createMetaAlarmModule } from '@unit/utils/store';

import LinkMetaAlarmForm from '@/components/widgets/alarm/forms/link-meta-alarm-form.vue';

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

describe('link-meta-alarm-form', () => {
  const { metaAlarmModule, fetchMetaAlarmsListWithoutStore } = createMetaAlarmModule();
  const store = createMockedStoreModules([
    metaAlarmModule,
  ]);

  const factory = generateShallowRenderer(LinkMetaAlarmForm, { stubs });
  const snapshotFactory = generateRenderer(LinkMetaAlarmForm, { stubs: snapshotStubs });

  test('Alarms fetched after mount', () => {
    factory({
      store,
      propsData: {
        form: {},
      },
    });

    expect(fetchMetaAlarmsListWithoutStore).toBeCalled();
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

    expect(wrapper).toEmitInput({ ...form, metaAlarm });
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

    expect(wrapper).toEmitInput({ ...form, comment });
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

    expect(wrapper).toEmitInput({ ...form, auto_resolve: autoResolve });
  });

  test('Renders `link-meta-alarm-form` with default props', async () => {
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
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `link-meta-alarm-form` with alarms', async () => {
    fetchMetaAlarmsListWithoutStore.mockReturnValueOnce([
      { _id: 'entity-id', name: 'alarm-display-name' },
    ]);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([metaAlarmModule]),
      propsData: {
        form: {},
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
