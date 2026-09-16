import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import MermaidCodeEditor from '@/components/other/map/form/fields/mermaid-code-editor.vue';

const stubs = {
  'code-editor': true,
};

const selectCodeEditorNode = wrapper => wrapper.vm.$children[0];

describe('mermaid-code-editor', () => {
  const factory = generateShallowRenderer(MermaidCodeEditor, { stubs });

  test('Value changed after trigger code editor', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const newValue = Faker.datatype.string();

    const codeEditorNode = selectCodeEditorNode(wrapper);

    codeEditorNode.$emit('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });

  test('Error added after update value', async () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    await wrapper.setProps({
      value: 'invalid code',
    });

    const codeEditorNode = selectCodeEditorNode(wrapper);

    expect(codeEditorNode.errorMarkers).toEqual([
      {
        endColumn: 1,
        endLineNumber: 1,
        severity: 8,
        startColumn: 0,
        startLineNumber: 1,
        message: expect.any(String),
      },
    ]);
  });

  test('Error removed after input valid value', async () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    await wrapper.setProps({
      value: 'invalid code',
    });

    expect(selectCodeEditorNode(wrapper).errorMarkers).toHaveLength(1);

    await wrapper.setProps({
      value: 'graph TB\n  a-->b',
    });

    expect(selectCodeEditorNode(wrapper).errorMarkers).toHaveLength(0);
  });
});
