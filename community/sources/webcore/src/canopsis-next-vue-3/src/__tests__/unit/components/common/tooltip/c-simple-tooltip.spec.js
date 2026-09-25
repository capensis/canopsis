import { generateShallowRenderer } from '@unit/utils/vue';

import CSimpleTooltip from '@/components/common/tooltip/c-simple-tooltip.vue';

const { sanitizeHtmlMock } = vi.hoisted(() => ({
  sanitizeHtmlMock: vi.fn(content => `<strong>${content}</strong>`),
}));

vi.mock('@/helpers/html', () => ({
  sanitizeHtml: sanitizeHtmlMock,
}));

const VTooltipStub = {
  name: 'v-tooltip',
  props: ['modelValue'],
  emits: ['update:modelValue'],
  template: `
    <div class="v-tooltip">
      <slot name="activator" :props="{}" />
      <slot />
    </div>
  `,
};

describe('c-simple-tooltip', () => {
  const factory = generateShallowRenderer(CSimpleTooltip, {
    stubs: {
      'v-tooltip': VTooltipStub,
    },
  });

  beforeEach(() => {
    sanitizeHtmlMock.mockClear();
  });

  it('sanitizes HTML content only after the tooltip is opened', async () => {
    const content = '<em>tooltip content</em>';
    const wrapper = factory({
      propsData: {
        content,
        html: true,
      },
    });
    const tooltip = wrapper.findComponent(VTooltipStub);

    expect(sanitizeHtmlMock).not.toHaveBeenCalled();

    tooltip.vm.$emit('update:modelValue', true);
    await wrapper.vm.$nextTick();

    expect(sanitizeHtmlMock).toHaveBeenCalledOnce();
    expect(sanitizeHtmlMock).toHaveBeenCalledWith(content);
  });

  it('does not sanitize plain text content when the tooltip is opened', async () => {
    const wrapper = factory({
      propsData: {
        content: 'plain tooltip content',
      },
    });
    const tooltip = wrapper.findComponent(VTooltipStub);

    tooltip.vm.$emit('update:modelValue', true);
    await wrapper.vm.$nextTick();

    expect(sanitizeHtmlMock).not.toHaveBeenCalled();
  });
});
