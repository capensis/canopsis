import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import MermaidCodePreview from '@/components/other/map/partials/mermaid-code-preview.vue';

const { renderMermaidMock } = vi.hoisted(() => ({
  renderMermaidMock: vi.fn(),
}));

vi.mock('@/helpers/mermaid', () => ({
  renderMermaid: renderMermaidMock,
}));

describe('mermaid-code-preview', () => {
  const factory = generateShallowRenderer(MermaidCodePreview);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  test('Renders SVG after Mermaid promise resolves', async () => {
    renderMermaidMock.mockResolvedValue('<svg viewBox="0 0 100 200"></svg>');

    const wrapper = factory({
      propsData: {
        value: 'graph TB\n a-->b',
      },
    });

    expect(wrapper.text()).not.toContain('[object Promise]');

    await flushPromises();

    const svg = wrapper.find('svg');

    expect(svg.exists()).toBe(true);
    expect(svg.attributes('height')).toBe('200');
    expect(svg.attributes('style')).toContain('height: 200px');
    expect(wrapper.classes()).not.toContain('mermaid-preview--error');
  });

  test('Handles an asynchronous Mermaid render error', async () => {
    const error = new Error('Invalid diagram');
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});

    renderMermaidMock.mockRejectedValue(error);

    const wrapper = factory({
      propsData: {
        value: 'invalid',
      },
    });

    await flushPromises();

    expect(consoleError).toHaveBeenCalledWith(error);
    expect(wrapper.classes()).toContain('mermaid-preview--error');

    consoleError.mockRestore();
  });
});
