const { mermaidMock } = vi.hoisted(() => ({
  mermaidMock: {
    initialize: vi.fn(),
    render: vi.fn(),
    mermaidAPI: {
      initialize: vi.fn(),
    },
  },
}));

vi.mock('mermaid', () => ({ default: mermaidMock }));

import { renderMermaid } from '@/helpers/mermaid';

describe('mermaid helper', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test('Returns SVG from the asynchronous Mermaid render result', async () => {
    const options = { theme: 'default' };

    mermaidMock.render.mockResolvedValue({ svg: '<svg />' });

    await expect(renderMermaid('graph TB\n a-->b', options)).resolves.toBe('<svg />');

    expect(mermaidMock.initialize).toHaveBeenCalledWith({ ...options, startOnLoad: false });
    expect(mermaidMock.render).toHaveBeenCalledWith(
      expect.stringMatching(/^id_/),
      'graph TB\n a-->b',
    );
    expect(mermaidMock.mermaidAPI.initialize).toHaveBeenCalledWith({
      ...options,
      startOnLoad: false,
    });
  });

  test('Removes temporary Mermaid element after an asynchronous render error', async () => {
    const error = new Error('Invalid diagram');

    mermaidMock.render.mockImplementation(function mockRender(id) {
      const element = document.createElement('div');

      element.id = `d${id}`;
      document.body.appendChild(element);

      return Promise.reject(error);
    });

    await expect(renderMermaid('invalid')).rejects.toBe(error);

    expect(document.querySelector('[id^="did_"]')).toBeNull();
  });
});
