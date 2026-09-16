import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import InfoPopupForm from '@/components/widgets/alarm/forms/info-popup-form.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
  'text-editor-field': true,
};

const selectColumnField = wrapper => wrapper.find('.v-select');
const selectTextEditorField = wrapper => wrapper.find('text-editor-field-stub');

describe('info-popup-form', () => {
  const factory = generateShallowRenderer(InfoPopupForm, { stubs });

  test('Column changed after trigger select field', () => {
    const columns = [{ value: Faker.datatype.string() }, { value: Faker.datatype.string() }];
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
});
