import { createInputStub } from '@unit/stubs/input';
import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { COLOR_INDICATOR_TYPES, SORT_ORDERS } from '@/constants';

import AlarmsListModal from '@/components/sidebars/alarm/form/alarms-list-modal.vue';

const stubs = {
  'widget-settings-group': true,
  'field-resize-column-behavior': true,
  'field-default-sort-column': true,
  'field-default-sort-columns-with-template': true,
  'field-columns': true,
  'field-default-elements-per-page': true,
  'field-info-popup': true,
  'field-text-editor-with-template': true,
  'field-root-cause-settings': createInputStub('field-root-cause-settings'),
};

const selectFieldDefaultSortColumnsWithTemplate = wrapper =>
  wrapper.findComponent({ name: 'field-default-sort-columns-with-template' });
const selectFieldColumns = wrapper => wrapper.findComponent({ name: 'field-columns' });
const selectFieldDefaultElementsPerPage = wrapper =>
  wrapper.findComponent({ name: 'field-default-elements-per-page' });
const selectFieldInfoPopup = wrapper => wrapper.findComponent({ name: 'field-info-popup' });
const selectFieldTextEditorWithTemplate = wrapper =>
  wrapper.findAllComponents({ name: 'field-text-editor-with-template' }).at(0);
const selectFieldRootCauseSettings = wrapper =>
  wrapper.findComponent({ name: 'field-root-cause-settings' });

describe('alarms-list-modal', () => {
  const getForm = () => ({
    widgetColumns: [],
    itemsPerPage: Faker.datatype.number(),
    infoPopups: [
      {
        column: Faker.datatype.string(),
        template: Faker.datatype.string(),
      },
    ],
    columns: {},
    moreInfoTemplate: Faker.datatype.string(),
  });

  const associativeTableModule = {
    name: 'associativeTable',
    actions: {
      fetch: vi.fn(() => ({})),
    },
  };

  const store = createMockedStoreModules([associativeTableModule]);

  const factory = generateShallowRenderer(AlarmsListModal, {
    store,
    stubs,
  });

  test('Sort columns changed after trigger sort columns field', () => {
    const form = getForm();
    const wrapper = factory({
      propsData: { value: form },
    });

    const newSort = [
      {
        sort_by: Faker.datatype.string(),
        sort: SORT_ORDERS.desc,
      },
    ];

    selectFieldDefaultSortColumnsWithTemplate(wrapper).vm.$emit('update:model-value', newSort);

    expect(wrapper).toEmitInput({ ...form, sort: newSort });
  });

  test('Sort template changed after trigger sort columns template field', () => {
    const form = getForm();
    const wrapper = factory({
      propsData: { value: form },
    });

    const template = Faker.datatype.uuid();
    const sort = [
      {
        sort_by: Faker.datatype.string(),
        sort: SORT_ORDERS.asc,
      },
    ];

    selectFieldDefaultSortColumnsWithTemplate(wrapper).vm.$emit('update:template', template, sort);

    expect(wrapper).toEmitInput({ ...form, sort, sortTemplate: template });
  });

  test('Columns changed after trigger columns field', () => {
    const form = getForm();
    const wrapper = factory({
      propsData: { value: form },
    });

    const newColumns = [
      {
        label: Faker.datatype.string(),
        value: Faker.datatype.string(),
      },
    ];

    selectFieldColumns(wrapper).vm.$emit('update:model-value', newColumns);

    expect(wrapper).toEmitInput({ ...form, widgetColumns: newColumns });
  });

  test('Items per page changed after trigger items per page field', () => {
    const form = getForm();
    const wrapper = factory({
      propsData: { value: form },
    });

    const newValue = Faker.datatype.number({
      min: form.itemsPerPage + 1,
    });

    selectFieldDefaultElementsPerPage(wrapper).vm.$emit('update:model-value', newValue);

    expect(wrapper).toEmitInput({ ...form, itemsPerPage: newValue });
  });

  test('Info popups changed after trigger info popup field', () => {
    const form = getForm();
    const wrapper = factory({
      propsData: { value: form },
    });

    const newInfoPopups = [
      {
        column: Faker.datatype.string(),
        value: Faker.datatype.string(),
      },
    ];

    selectFieldInfoPopup(wrapper).vm.$emit('update:model-value', newInfoPopups);

    expect(wrapper).toEmitInput({ ...form, infoPopups: newInfoPopups });
  });

  test('More info template changed after trigger text editor field', () => {
    const form = getForm();
    const wrapper = factory({
      propsData: { value: form },
    });

    const newTemplate = Faker.datatype.string();

    selectFieldTextEditorWithTemplate(wrapper).vm.$emit('update:model-value', newTemplate);

    expect(wrapper).toEmitInput({ ...form, moreInfoTemplate: newTemplate });
  });

  test('Show root cause by state click changed after trigger switcher field', () => {
    const form = getForm();
    const wrapper = factory({
      propsData: { value: form },
    });

    const newParameters = {
      ...form,
      showRootCauseByStateClick: false,
      rootCauseColorIndicator: COLOR_INDICATOR_TYPES.impactState,
    };
    selectFieldRootCauseSettings(wrapper).vm.$emit('update:model-value', newParameters);

    expect(wrapper).toEmitInput(newParameters);
  });
});
