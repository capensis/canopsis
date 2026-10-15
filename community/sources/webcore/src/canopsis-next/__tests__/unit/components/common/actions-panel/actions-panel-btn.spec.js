import { generateShallowRenderer } from '@unit/utils/vue';

import ActionsPanelBtn from '@/components/common/actions-panel/actions-panel-btn.vue';

const stubs = {
  'c-action-btn': true,
};

const selectActionButton = wrapper => wrapper.find('c-action-btn-stub');

describe('actions-panel-btn', () => {
  const factory = generateShallowRenderer(ActionsPanelBtn, { stubs });

  test('keeps tooltip enabled for a disabled button when disabledButton is used', () => {
    const wrapper = factory({
      props: {
        action: {
          title: 'Inherited pause',
          icon: 'play_arrow',
          disabledButton: true,
        },
      },
    });
    const button = selectActionButton(wrapper);

    expect(button.attributes('tooltip')).toBe('Inherited pause');
    expect(button.attributes('disabled-button')).toBe('true');
    expect(button.attributes('disabled')).toBeUndefined();
  });
});
