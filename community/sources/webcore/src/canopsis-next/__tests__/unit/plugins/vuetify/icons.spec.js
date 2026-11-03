import { mount } from '@vue/test-utils';
import { h } from 'vue';

import { vuetify } from '@/plugins/vuetify';
import {
  canopsisIcons,
  customIconsRegistry,
  registerCustomIcon,
  unregisterCustomIcon,
} from '@/plugins/vuetify/icons';

import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CustomIcon from '@/components/other/icons/partials/custom-icon.vue';

const name = 'xss-regression-icon';
const maliciousSvg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" onload="window.iconXss = true">
  <foreignObject width="1" height="1">
    <img xmlns="http://www.w3.org/1999/xhtml" src="x" onerror="window.iconXss = true" />
  </foreignObject>
  <script>window.iconXss = true</script>
  <a href="javascript:window.iconXss = true"><path d="M0 0h24v24H0z" /></a>
</svg>`;

describe('custom icon sanitization', () => {
  afterEach(() => {
    unregisterCustomIcon(name);
  });

  it('sanitizes SVG before storing it in the registry and rendering it', () => {
    registerCustomIcon(name, maliciousSvg, { fillBorder: true });

    const stored = customIconsRegistry[name].component.template;
    const wrapper = mount({
      render: () => h(canopsisIcons.component, { tag: 'i', icon: name }),
    });

    expect(stored).not.toMatch(/foreignObject|script|onload|onerror|javascript:/i);
    expect(wrapper.find('svg').attributes('viewBox')).toBe('0 0 24 24');
    expect(wrapper.find('path').attributes('d')).toBe('M0 0h24v24H0z');
    expect(wrapper.find('i').classes()).toContain('v-icon--fill-border');
    wrapper.unmount();
  });

  it('applies the same SVG policy to the icon list preview', () => {
    const wrapper = mount(CustomIcon, { props: { content: maliciousSvg } });

    expect(wrapper.html()).not.toMatch(/foreignObject|script|onload|onerror|javascript:/i);
    expect(wrapper.find('path').exists()).toBe(true);
    wrapper.unmount();
  });

  it('safely renders an uploaded icon referenced by a text template', () => {
    registerCustomIcon(name, maliciousSvg);

    const wrapper = mount(CRuntimeTemplate, {
      props: { template: `<v-icon>$${name}</v-icon>` },
      global: { plugins: [vuetify] },
    });

    expect(wrapper.find('svg').exists()).toBe(true);
    expect(wrapper.html()).not.toMatch(/foreignObject|script|onload|onerror|javascript:/i);
    wrapper.unmount();
  });
});
