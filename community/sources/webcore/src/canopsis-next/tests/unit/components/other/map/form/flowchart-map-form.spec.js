import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import FlowchartMapForm from '@/components/other/map/form/flowchart-map-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-name-field': true,
  'flowchart-editor': true,
};

const factory = (options = {}) => shallowMount(FlowchartMapForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(FlowchartMapForm, {
  localVue,
  stubs,

  ...options,
});

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectFlowchartEditor = wrapper => wrapper.find('flowchart-editor-stub');

describe('flowchart-map-form', () => {
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

    nameField.vm.$emit('input', newName);

    expect(wrapper).toEmit('input', {
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

    flowchartEditor.vm.$emit('input', newParameters);

    expect(wrapper).toEmit('input', {
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
