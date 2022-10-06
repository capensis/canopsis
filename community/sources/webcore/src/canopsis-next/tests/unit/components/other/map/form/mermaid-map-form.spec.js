import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import MermaidMapForm from '@/components/other/map/form/mermaid-map-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-name-field': true,
  'mermaid-editor': true,
};

const factory = (options = {}) => shallowMount(MermaidMapForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(MermaidMapForm, {
  localVue,
  stubs,

  ...options,
});

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectMermaidEditor = wrapper => wrapper.find('mermaid-editor-stub');

describe('mermaid-map-form', () => {
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

    mermaidEditor.vm.$emit('input', newParameters);

    expect(wrapper).toEmit('input', {
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

    expect(wrapper.element).toMatchSnapshot();
  });
});