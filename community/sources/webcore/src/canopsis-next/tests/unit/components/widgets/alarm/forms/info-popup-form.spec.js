import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import InfoPopupForm from '@/components/widgets/alarm/forms/info-popup-form.vue';
import { ALARM_FIELDS } from '@/constants';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'text-editor-field': true,
};

const snapshotStubs = {
  'text-editor-field': true,
};

const factory = (options = {}) => shallowMount(InfoPopupForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(InfoPopupForm, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectColumnField = wrapper => wrapper.find('.v-select');
const selectTextEditorField = wrapper => wrapper.find('text-editor-field-stub');

describe('info-popup-form', () => {
  test('Column changed after trigger select field', () => {
    const columns = [
      { value: Faker.datatype.string() },
      { value: Faker.datatype.string() },
    ];
    const template = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        form: {
          column: columns[0].value,
          template,
        },
        columns,
      },
    });

    const columnField = selectColumnField(wrapper);

    columnField.setValue(columns[1].value);

    expect(wrapper).toEmit('input', {
      template,
      column: columns[1].value,
    });
  });

  test('Template changed after trigger text editor field', () => {
    const column = Faker.datatype.string();
    const template = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        form: {
          column,
          template,
        },
      },
    });

    const textEditorField = selectTextEditorField(wrapper);

    const newTemplate = Faker.datatype.string();

    textEditorField.vm.$emit('input', newTemplate);

    expect(wrapper).toEmit('input', {
      template: newTemplate,
      column,
    });
  });

  test('Renders `info-popup-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `info-popup-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          column: 'column',
          template: 'template',
        },
        columns: [
          {
            value: ALARM_FIELDS.entityName,
          },
        ],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
