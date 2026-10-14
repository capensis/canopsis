import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import MermaidCodeEditor from '@/components/other/map/form/fields/mermaid-code-editor.vue';

const stubs = {
  'code-editor': {
    name: 'CodeEditor',
    props: ['value', 'errorMarkers'],
    emits: ['update:value'],
    template: '<div />',
  },
};

const selectCodeEditor = wrapper => wrapper.findComponent({ name: 'CodeEditor' });

describe('mermaid-code-editor', () => {
  const factory = generateShallowRenderer(MermaidCodeEditor, { stubs });

  test('Value changed after trigger code editor', async () => {
    const wrapper = factory({
      props: {
        value: '',
      },
    });

    const newValue = Faker.datatype.string();

    selectCodeEditor(wrapper).vm.$emit('update:value', newValue);
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted('input')).toEqual([[newValue]]);
  });

  test('Error added after update value', async () => {
    const wrapper = factory({
      props: {
        value: '',
      },
    });

    await wrapper.setProps({
      value: 'graph TB\n a--',
    });
    await vi.waitFor(() => {
      expect(selectCodeEditor(wrapper).props('errorMarkers')).toEqual([
        {
          endColumn: 1,
          endLineNumber: 3,
          severity: 8,
          startColumn: 2,
          startLineNumber: 2,
          message: expect.any(String),
        },
      ]);
    });
  });

  test('Error removed after input valid value', async () => {
    const wrapper = factory({
      props: {
        value: '',
      },
    });

    await wrapper.setProps({
      value: 'graph TB\n a--',
    });
    await vi.waitFor(() => {
      expect(selectCodeEditor(wrapper).props('errorMarkers')).toHaveLength(1);
    });

    await wrapper.setProps({
      value: 'graph TB\n  a-->b',
    });
    await vi.waitFor(() => {
      expect(selectCodeEditor(wrapper).props('errorMarkers')).toHaveLength(0);
    });
  });
});
