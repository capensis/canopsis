import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import MermaidMapForm from '@/components/other/map/form/mermaid-map-form.vue';

const stubs = {
  'c-name-field': true,
  'mermaid-editor': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectMermaidEditor = wrapper => wrapper.find('mermaid-editor-stub');

describe('mermaid-map-form', () => {
  const factory = generateShallowRenderer(MermaidMapForm, { stubs });
  const snapshotFactory = generateRenderer(MermaidMapForm, { stubs });

  test('Name changed after trigger name field', () => {
    const form = {
      name: '',
      parameters: {},
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newName = Faker.datatype.string();

    const nameField = selectNameField(wrapper);

    nameField.triggerCustomEvent('input', newName);

    expect(wrapper).toEmitInput({
      ...form,
      name: newName,
    });
  });

  test('Parameters changed after trigger mermaid editor field', () => {
    const form = {
      name: '',
      parameters: {},
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newParameters = {
      param: Faker.datatype.string(),
    };

    const mermaidEditor = selectMermaidEditor(wrapper);

    mermaidEditor.triggerCustomEvent('input', newParameters);

    expect(wrapper).toEmitInput({
      ...form,
      parameters: newParameters,
    });
  });

  test('Renders `mermaid-map-form` with form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          name: 'Mermaid',
          parameters: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
