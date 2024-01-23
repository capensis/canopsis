import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import FlowchartMapForm from '@/components/other/map/form/flowchart-map-form.vue';

const stubs = {
  'c-name-field': true,
  'flowchart-editor': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectFlowchartEditor = wrapper => wrapper.find('flowchart-editor-stub');

describe('flowchart-map-form', () => {
  const factory = generateShallowRenderer(FlowchartMapForm, { stubs });
  const snapshotFactory = generateRenderer(FlowchartMapForm, { stubs });

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

  test('Parameters changed after trigger flowchart editor field', () => {
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

    const flowchartEditor = selectFlowchartEditor(wrapper);

    flowchartEditor.triggerCustomEvent('input', newParameters);

    expect(wrapper).toEmitInput({
      ...form,
      parameters: newParameters,
    });
  });

  test('Renders `flowchart-map-form` with form', () => {
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
