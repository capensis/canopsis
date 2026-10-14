import { generateShallowRenderer } from '@unit/utils/vue';

import ActionsPanelMenu from '@/components/common/actions-panel/actions-panel-menu.vue';

describe('actions-panel-menu', () => {
  const factory = generateShallowRenderer(ActionsPanelMenu);

  test('keeps menu open for an action with disabled close on click', () => {
    const action = {
      title: 'Export PDF',
      closeOnClick: false,
      method: vi.fn(),
    };
    const wrapper = factory({ propsData: { actions: [action] } });

    wrapper.vm.menu = true;
    wrapper.vm.selectAction(action);

    expect(action.method).toHaveBeenCalledOnce();
    expect(wrapper.vm.menu).toBe(true);
  });

  test('closes menu for an action by default', () => {
    const action = {
      title: 'Open modal',
      method: vi.fn(),
    };
    const wrapper = factory({ propsData: { actions: [action] } });

    wrapper.vm.menu = true;
    wrapper.vm.selectAction(action);

    expect(action.method).toHaveBeenCalledOnce();
    expect(wrapper.vm.menu).toBe(false);
  });
});
