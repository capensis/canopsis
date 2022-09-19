import Faker from 'faker';
import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import MermaidEditor from '@/components/other/map/form/partials/mermaid-editor.vue';
import { MERMAID_THEMES } from '@/constants';
import { mermaidPointToForm } from '@/helpers/forms/map';

const localVue = createVueInstance();

const stubs = {
  'mermaid-code-editor': true,
  'add-location-btn': true,
  'mermaid-theme-field': true,
  'mermaid-code-preview': true,
  'mermaid-points': true,
};

const factory = (options = {}) => shallowMount(MermaidEditor, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(MermaidEditor, {
  localVue,
  stubs,

  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },

  ...options,
});

const selectMermaidCodeEditor = wrapper => wrapper.find('mermaid-code-editor-stub');
const selectAddLocationBtn = wrapper => wrapper.find('add-location-btn-stub');
const selectMermaidThemeField = wrapper => wrapper.find('mermaid-theme-field-stub');
const selectMermaidPoints = wrapper => wrapper.find('mermaid-points-stub');

describe('mermaid-editor', () => {
  const initialForm = {
    code: '',
    theme: MERMAID_THEMES.base,
    points: [],
  };

  test('Code changed after trigger code editor', () => {
    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const newCode = Faker.datatype.string();

    const mermaidCodeEditor = selectMermaidCodeEditor(wrapper);

    mermaidCodeEditor.vm.$emit('input', newCode);

    expect(wrapper).toEmit('input', {
      ...initialForm,
      code: newCode,
    });
  });

  test('Theme changed after trigger theme field', () => {
    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const newTheme = MERMAID_THEMES.canopsis;

    const mermaidThemeField = selectMermaidThemeField(wrapper);

    mermaidThemeField.vm.$emit('input', newTheme);

    expect(wrapper).toEmit('input', {
      ...initialForm,
      theme: newTheme,
    });
  });

  test('Points changed after trigger points editor', () => {
    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const newPoints = [mermaidPointToForm()];

    const mermaidPoints = selectMermaidPoints(wrapper);

    mermaidPoints.vm.$emit('input', newPoints);

    expect(wrapper).toEmit('input', {
      ...initialForm,
      points: newPoints,
    });
  });

  test('Add on click mode enabled after trigger add location btn', async () => {
    const wrapper = factory({
      propsData: {
        form: initialForm,
      },
    });

    const addLocationBtn = selectAddLocationBtn(wrapper);

    await addLocationBtn.vm.$emit('input', true);

    const mermaidPoints = selectMermaidPoints(wrapper);

    expect(mermaidPoints.vm.$attrs['add-on-click']).toBeTruthy();
  });

  test('Form re-validated after change form with error', async () => {
    const wrapper = factory({
      propsData: {
        name: 're-validate-name',
        form: initialForm,
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(validator.errors.items).toHaveLength(1);

    wrapper.setProps({
      form: {
        ...initialForm,
        code: 'mermaid-code',
        points: [mermaidPointToForm()],
      },
    });

    await flushPromises();

    expect(validator.errors.items).toHaveLength(0);
  });

  test('Renders `mermaid-editor` with form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: initialForm,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `mermaid-editor` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: initialForm,
        minHeight: 500,
        name: 'custom_name',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `mermaid-editor` with validation errors ', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: initialForm,
        minHeight: 500,
        name: 'custom_name',
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });
});
