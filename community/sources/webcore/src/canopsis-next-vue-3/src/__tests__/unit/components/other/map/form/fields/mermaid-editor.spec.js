import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { MERMAID_THEMES } from '@/constants';

import { mermaidPointToForm } from '@/helpers/entities/map/form';

import MermaidEditor from '@/components/other/map/form/fields/mermaid-editor.vue';

const stubs = {
  'mermaid-code-editor': {
    name: 'MermaidCodeEditor',
    emits: ['input'],
    template: '<div />',
  },
  'add-location-btn': {
    name: 'AddLocationBtn',
    emits: ['input'],
    template: '<div />',
  },
  'mermaid-theme-field': {
    name: 'MermaidThemeField',
    emits: ['input'],
    template: '<div />',
  },
  'mermaid-code-preview': true,
  'mermaid-points-editor': {
    name: 'MermaidPointsEditor',
    props: ['value', 'addOnClick'],
    emits: ['input'],
    template: '<div />',
  },
  'map-editor-required-field': {
    name: 'MapEditorRequiredField',
    props: ['valid', 'message', 'name'],
    template: '<div><slot :has-error="!valid" :messages="[message]" /></div>',
  },
  'v-messages': true,
};

const selectMermaidCodeEditor = wrapper => wrapper.findComponent({ name: 'MermaidCodeEditor' });
const selectAddLocationBtn = wrapper => wrapper.findComponent({ name: 'AddLocationBtn' });
const selectMermaidThemeField = wrapper => wrapper.findComponent({ name: 'MermaidThemeField' });
const selectMermaidPoints = wrapper => wrapper.findComponent({ name: 'MermaidPointsEditor' });
const selectRequiredField = wrapper => wrapper.findComponent({ name: 'MapEditorRequiredField' });

describe('mermaid-editor', () => {
  const initialForm = {
    code: '',
    theme: MERMAID_THEMES.base,
    points: [],
  };

  const factory = generateShallowRenderer(MermaidEditor, { stubs });

  test('Code changed after trigger code editor', async () => {
    const wrapper = factory({
      props: {
        form: initialForm,
      },
    });

    const newCode = Faker.datatype.string();

    selectMermaidCodeEditor(wrapper).vm.$emit('input', newCode);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:form')).toEqual([[{
      ...initialForm,
      code: newCode,
    }]]);
  });

  test('Theme changed after trigger theme field', async () => {
    const wrapper = factory({
      props: {
        form: initialForm,
      },
    });

    const newTheme = MERMAID_THEMES.canopsis;

    selectMermaidThemeField(wrapper).vm.$emit('input', newTheme);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:form')).toEqual([[{
      ...initialForm,
      theme: newTheme,
    }]]);
  });

  test('Points changed after trigger points editor', async () => {
    const wrapper = factory({
      props: {
        form: initialForm,
      },
    });

    const newPoints = [mermaidPointToForm()];

    selectMermaidPoints(wrapper).vm.$emit('input', newPoints);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('update:form')).toEqual([[{
      ...initialForm,
      points: newPoints,
    }]]);
  });

  test('Add on click mode enabled after trigger add location btn', async () => {
    const wrapper = factory({
      props: {
        form: initialForm,
      },
    });

    selectAddLocationBtn(wrapper).vm.$emit('input', true);
    await wrapper.vm.$nextTick();

    expect(selectMermaidPoints(wrapper).props('addOnClick')).toBe(true);
  });

  test('Form re-validated after change form with error', async () => {
    const wrapper = factory({
      props: {
        name: 're-validate-name',
        form: initialForm,
      },
    });

    expect(selectRequiredField(wrapper).props('valid')).toBe(false);

    await wrapper.setProps({
      form: {
        ...initialForm,
        code: 'mermaid-code',
        points: [mermaidPointToForm()],
      },
    });

    await flushPromises();

    expect(selectRequiredField(wrapper).props('valid')).toBe(true);
  });
});
