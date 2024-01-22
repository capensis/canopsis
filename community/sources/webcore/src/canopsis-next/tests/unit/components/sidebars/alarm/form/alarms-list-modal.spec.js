import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createMockedStoreModules } from '@unit/utils/store';
import AlarmsListModal from '@/components/sidebars/alarm/form/alarms-list-modal.vue';

const stubs = {
  'widget-settings-group': true,
  'field-default-sort-column': true,
  'field-columns': true,
  'field-default-elements-per-page': true,
  'field-info-popup': true,
  'field-text-editor-with-template': true,
};

const snapshotStubs = {
  'widget-settings-group': true,
  'field-default-sort-column': true,
  'field-columns': true,
  'field-default-elements-per-page': true,
  'field-info-popup': true,
  'field-text-editor-with-template': true,
  'field-switcher': true,
};

const selectFieldColumns = wrapper => wrapper.find('field-columns-stub');
const selectFieldDefaultElementsPerPage = wrapper => wrapper.find('field-default-elements-per-page-stub');
const selectFieldInfoPopup = wrapper => wrapper.find('field-info-popup-stub');
const selectFieldTextEditorWithTemplate = wrapper => wrapper.find('field-text-editor-with-template-stub');
const selectFieldSwitcher = wrapper => wrapper.find('field-switcher-stub');

describe('alarms-list-modal', () => {
  const form = {
    widgetColumns: [],
    itemsPerPage: Faker.datatype.number(),
    infoPopups: [{
      column: Faker.datatype.string(),
      template: Faker.datatype.string(),
    }],
    moreInfoTemplate: Faker.datatype.string(),
  };

  const associativeTableModule = {
    name: 'associativeTable',
    actions: {
      fetch: jest.fn(() => ({})),
    },
  };

  const store = createMockedStoreModules([
    associativeTableModule,
  ]);

  const factory = generateShallowRenderer(AlarmsListModal, {
    store,
    stubs,
  });

  const snapshotFactory = generateRenderer(AlarmsListModal, {
    store,
    stubs: snapshotStubs,
  });

  test('Columns changed after trigger columns field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newColumns = [
      {
        label: Faker.datatype.string(),
        value: Faker.datatype.string(),
      },
    ];

    selectFieldColumns(wrapper).triggerCustomEvent('input', newColumns);

    expect(wrapper).toEmit('input', { ...form, widgetColumns: newColumns });
  });

  test('Items per page changed after trigger items per page field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newValue = Faker.datatype.number({
      min: form.itemsPerPage + 1,
    });

    selectFieldDefaultElementsPerPage(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('input', { ...form, itemsPerPage: newValue });
  });

  test('Info popups changed after trigger info popup field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newInfoPopups = [
      {
        column: Faker.datatype.string(),
        value: Faker.datatype.string(),
      },
    ];

    selectFieldInfoPopup(wrapper).triggerCustomEvent('input', newInfoPopups);

    expect(wrapper).toEmit('input', { ...form, infoPopups: newInfoPopups });
  });

  test('More info template changed after trigger text editor field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newTemplate = Faker.datatype.string();

    selectFieldTextEditorWithTemplate(wrapper).triggerCustomEvent('input', newTemplate);

    expect(wrapper).toEmit('input', { ...form, moreInfoTemplate: newTemplate });
  });

  test('Show root cause by state click changed after trigger switcher field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newShowRootCauseByStateClick = Faker.datatype.boolean();

    selectFieldSwitcher(wrapper).triggerCustomEvent('input', newShowRootCauseByStateClick);

    expect(wrapper).toEmit('input', { ...form, showRootCauseByStateClick: newShowRootCauseByStateClick });
  });

  test('Renders `alarms-list-modal` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-modal` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          widgetColumns: [{}, {}],
          itemsPerPage: 11,
          infoPopups: [{}],
          moreInfoTemplate: '<div>more-info-template</div>',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
