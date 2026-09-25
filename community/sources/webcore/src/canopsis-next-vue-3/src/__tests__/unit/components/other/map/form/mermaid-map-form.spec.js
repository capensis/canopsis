import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import MermaidMapForm from '@/components/other/map/form/mermaid-map-form.vue';

const stubs = {
  'c-name-field': {
    name: 'CNameField',
    emits: ['update:modelValue'],
    template: '<div />',
  },
  'mermaid-editor': {
    name: 'MermaidEditor',
    emits: ['update:form'],
    template: '<div />',
  },
};

const selectNameField = wrapper => wrapper.findComponent({ name: 'CNameField' });
const selectMermaidEditor = wrapper => wrapper.findComponent({ name: 'MermaidEditor' });

describe('mermaid-map-form', () => {
  const factory = generateShallowRenderer(MermaidMapForm, { stubs });

  test('Name changed after trigger name field', async () => {
    const form = {
      name: '',
      parameters: {},
    };
    const wrapper = factory({
      props: {
        form,
      },
    });

    const newName = Faker.datatype.string();

    selectNameField(wrapper).vm.$emit('update:modelValue', newName);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:form')).toEqual([[{
      ...form,
      name: newName,
    }]]);
  });

  test('Parameters changed after trigger mermaid editor field', async () => {
    const form = {
      name: '',
      parameters: {},
    };
    const wrapper = factory({
      props: {
        form,
      },
    });

    const newParameters = {
      param: Faker.datatype.string(),
    };

    selectMermaidEditor(wrapper).vm.$emit('update:form', newParameters);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:form')).toEqual([[{
      ...form,
      parameters: newParameters,
    }]]);
  });
});
