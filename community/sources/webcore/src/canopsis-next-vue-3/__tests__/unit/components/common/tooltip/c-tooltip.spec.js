import { generateShallowRenderer } from '@unit/utils/vue';
import { h } from 'vue';

import CTooltip from '@/components/common/tooltip/c-tooltip.vue';

const { sanitizeHtmlMock } = vi.hoisted(() => ({
  sanitizeHtmlMock: vi.fn(content => `<strong>${content}</strong>`),
}));

vi.mock('@/helpers/html', () => ({
  sanitizeHtml: sanitizeHtmlMock,
}));

const VTooltipStub = {
  name: 'v-tooltip',
  props: ['activator', 'disabled', 'eager', 'location', 'maxWidth', 'modelValue'],
  emits: ['update:modelValue'],
  template: `
    <div class="v-tooltip">
      <slot name="activator" :props="{ marker: 'activator' }" />
      <slot />
    </div>
  `,
};

describe('c-tooltip', () => {
  const factory = generateShallowRenderer(CTooltip, {
    stubs: {
      'v-tooltip': VTooltipStub,
    },
  });

  beforeEach(() => {
    sanitizeHtmlMock.mockClear();
  });

  it('renders tooltip content lazily', () => {
    const wrapper = factory({ propsData: { content: 'tooltip content' } });

    expect(wrapper.findComponent(VTooltipStub).props('eager')).toBe(false);
  });

  it('keeps eager disabled when it is passed explicitly', () => {
    const wrapper = factory({ attrs: { eager: true } });

    expect(wrapper.findComponent(VTooltipStub).props('eager')).toBe(false);
  });

  it('forwards Vuetify props and default content', () => {
    const wrapper = factory({
      attrs: {
        activator: 'parent',
        maxWidth: 320,
      },
      slots: {
        default: '<span class="tooltip-content">Tooltip content</span>',
      },
    });
    const tooltip = wrapper.findComponent(VTooltipStub);

    expect(tooltip.props()).toMatchObject({
      activator: 'parent',
      maxWidth: 320,
    });
    expect(wrapper.find('.tooltip-content').text()).toBe('Tooltip content');
  });

  it('provides Vuetify activator props through props and legacy on aliases', () => {
    let receivedSlotProps;
    factory({
      slots: {
        activator: (slotProps) => {
          receivedSlotProps = slotProps;

          return h('button');
        },
      },
    });

    expect(receivedSlotProps.props).toEqual({ marker: 'activator' });
    expect(receivedSlotProps.on).toBe(receivedSlotProps.props);
  });

  it('maps legacy directional props to Vuetify location', () => {
    const wrapper = factory({ propsData: { content: 'tooltip content', right: true } });

    expect(wrapper.findComponent(VTooltipStub).props('location')).toBe('right');
  });

  it('disables an empty content tooltip', () => {
    const wrapper = factory({ propsData: { content: '' } });

    expect(wrapper.findComponent(VTooltipStub).props('disabled')).toBe(true);
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

    expect(sanitizeHtmlMock).toHaveBeenCalledExactlyOnceWith(content);
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
