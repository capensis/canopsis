import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { ALARM_FIELDS } from '@/constants';

import InfoPopupForm from '@/components/widgets/alarm/forms/info-popup-form.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'text-editor-field': true,
};

const snapshotStubs = {
  'text-editor-field': true,
};

const selectColumnField = wrapper => wrapper.find('.v-select');
const selectTextEditorField = wrapper => wrapper.find('text-editor-field-stub');

describe('info-popup-form', () => {
  const factory = generateShallowRenderer(InfoPopupForm, { stubs });
  const snapshotFactory = generateRenderer(InfoPopupForm, { stubs: snapshotStubs });

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

    expect(wrapper).toEmitInput({
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

    textEditorField.triggerCustomEvent('input', newTemplate);

    expect(wrapper).toEmitInput({
      template: newTemplate,
      column,
    });
  });

  test('Renders `info-popup-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
