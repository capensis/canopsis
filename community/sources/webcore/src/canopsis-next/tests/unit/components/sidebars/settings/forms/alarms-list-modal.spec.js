import Faker from 'faker';

import { createVueInstance, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import AlarmsListModal from '@/components/sidebars/settings/forms/alarms-list-modal.vue';

const localVue = createVueInstance();

const stubs = {
  'widget-settings-group': true,
  'field-columns': true,
  'field-default-elements-per-page': true,
  'field-info-popup': true,
  'field-text-editor': true,
};

const snapshotStubs = {
  'widget-settings-group': true,
  'field-columns': true,
  'field-default-elements-per-page': true,
  'field-info-popup': true,
  'field-text-editor': true,
};

const selectFieldColumns = wrapper => wrapper.find('field-columns-stub');
const selectFieldDefaultElementsPerPage = wrapper => wrapper.find('field-default-elements-per-page-stub');
const selectFieldInfoPopup = wrapper => wrapper.find('field-info-popup-stub');
const selectFieldTextEditor = wrapper => wrapper.find('field-text-editor-stub');

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

  const factory = generateShallowRenderer(AlarmsListModal, {
    localVue,
    stubs,
  });

  const snapshotFactory = generateRenderer(AlarmsListModal, {
    localVue,
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

    selectFieldColumns(wrapper).vm.$emit('input', newColumns);

    expect(wrapper).toEmit('input', { ...form, widgetColumns: newColumns });
  });

  test('Items per page changed after trigger items per page field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newValue = Faker.datatype.number({
      min: form.itemsPerPage + 1,
    });

    selectFieldDefaultElementsPerPage(wrapper).vm.$emit('input', newValue);

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

    selectFieldInfoPopup(wrapper).vm.$emit('input', newInfoPopups);

    expect(wrapper).toEmit('input', { ...form, infoPopups: newInfoPopups });
  });

  test('More info template changed after trigger text editor field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newTemplate = Faker.datatype.string();

    selectFieldTextEditor(wrapper).vm.$emit('input', newTemplate);

    expect(wrapper).toEmit('input', { ...form, moreInfoTemplate: newTemplate });
  });

  test('Renders `alarms-list-modal` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});