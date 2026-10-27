import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import CopyValue from '@/components/modals/common/copy-value.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-text-field': true,
  'c-copy-btn': true,
  'v-alert': {
    name: 'VAlert',
    props: ['type'],
    template: '<div><slot /></div>',
  },
  'v-btn': true,
};

describe('copy-value', () => {
  const $modals = mockModals();
  const $popups = mockPopups();
  const factory = generateShallowRenderer(CopyValue, { stubs });

  test('shows the configured warning', () => {
    const warningText = 'Copy this value now.';
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            value: 'value',
            warningText,
          },
        },
      },
      mocks: {
        $modals,
        $popups,
      },
    });

    const warning = wrapper.findComponent({ name: 'VAlert' });

    expect(warning.props('type')).toBe('warning');
    expect(warning.text()).toBe(warningText);
  });
});
