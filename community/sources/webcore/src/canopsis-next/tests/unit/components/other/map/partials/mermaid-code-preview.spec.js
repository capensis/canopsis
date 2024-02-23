import { generateRenderer } from '@unit/utils/vue';

import { MERMAID_THEMES } from '@/constants';

import { renderMermaid } from '@/helpers/mermaid';

import MermaidCodePreview from '@/components/other/map/partials/mermaid-code-preview.vue';

jest.mock('@/helpers/mermaid', () => ({
  renderMermaid: jest.fn(
    (code, config) => `<svg><text>${code}</text><text>${JSON.stringify(config)}</text></svg>`,
  ),
}));

describe('mermaid-code-preview', () => {
  const snapshotFactory = generateRenderer(MermaidCodePreview);

  test('Renders `mermaid-code-preview` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `mermaid-code-preview` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'graph TB\n  a-->b',
        theme: MERMAID_THEMES.canopsis,
        name: 'custom_name',
        config: {
          requirement: {
            useMaxWidth: true,
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `mermaid-code-preview` with error', () => {
    renderMermaid.mockImplementation(() => {
      throw new Error();
    });
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });
});
