import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import InfoPopupForm from '@/components/widgets/alarm/forms/info-popup-form.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'text-editor-field': true,
};

const selectColumnField = wrapper => wrapper.findComponent('.v-select');
const selectTextEditorField = wrapper => wrapper.find('text-editor-field-stub');

describe('info-popup-form', () => {
  const factory = generateShallowRenderer(InfoPopupForm, { stubs });

  test('Prepared column names and values are passed to the select', () => {
    const columns = [
      { text: Faker.datatype.string(), value: Faker.datatype.string() },
      { text: Faker.datatype.string(), value: Faker.datatype.string() },
    ];
    const wrapper = factory({
      propsData: {
        columns,
      },
    });

    expect(selectColumnField(wrapper).props('items')).toEqual(columns);
  });

  test('Column changed after trigger select field', () => {
    const columns = [{ value: 'alarm.v.state.val' }, { value: 'alarm.v.status.val' }];
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

    columnField.triggerCustomEvent('update:modelValue', columns[1].value);

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
});
