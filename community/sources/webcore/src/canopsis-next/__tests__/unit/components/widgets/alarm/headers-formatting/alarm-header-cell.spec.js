import { generateShallowRenderer } from '@unit/utils/vue';

import AlarmHeaderCell from '@/components/widgets/alarm/headers-formatting/alarm-header-cell.vue';

describe('alarm-header-cell', () => {
  const factory = generateShallowRenderer(AlarmHeaderCell);
  const header = {
    value: 'long-column',
    text: 'Very long column title',
  };

  test('Truncates the header and exposes its full text in a native tooltip', () => {
    const wrapper = factory({
      propsData: {
        header,
        ellipsisHeaders: true,
      },
    });

    const headerText = wrapper.get('span');

    expect(headerText.classes()).toContain('v-data-table-header__span--ellipsis');
    expect(headerText.attributes('title')).toBe(header.text);
    expect(headerText.text()).toBe(header.text);
  });

  test('Allows the header to wrap when truncation is disabled', () => {
    const wrapper = factory({ propsData: { header } });
    const headerText = wrapper.get('span');

    expect(headerText.classes()).not.toContain('v-data-table-header__span--ellipsis');
    expect(headerText.attributes('title')).toBeUndefined();
    expect(headerText.attributes('style')).toContain('white-space: normal');
  });
});
