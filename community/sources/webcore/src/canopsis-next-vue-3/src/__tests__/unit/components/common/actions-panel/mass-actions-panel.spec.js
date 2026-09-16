import { fakeAction } from '@unit/data/actions-panel';
import { installResizeObserver, uninstallResizeObserver } from '@unit/utils/resize-observer';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import MassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';

const cActionBtnStub = {
  template: '<button class="c-action-btn-stub" type="button" v-on="$attrs" />',
};

const stubs = {
  'c-action-btn': cActionBtnStub,
};

describe('mass-actions-panel', () => {
  const factory = generateShallowRenderer(MassActionsPanel, { stubs });

  afterEach(() => {
    uninstallResizeObserver();
  });

  it('Method into inline button called after click when layout is wide enough', async () => {
    installResizeObserver(400);

    const actions = [
      { ...fakeAction(), type: 'a' },
      { ...fakeAction(), type: 'b' },
    ];

    const wrapper = factory({
      propsData: { actions },
    });

    await flushPromises();

    const buttons = wrapper.findAll('button.c-action-btn-stub');

    expect(buttons.length).toBe(2);

    buttons.at(1).trigger('click');

    expect(actions[1].method).toHaveBeenCalledTimes(1);
  });

  it('Method into overflow menu called after click when layout is narrow', async () => {
    installResizeObserver(36);

    const actions = [
      { ...fakeAction(), type: 'a' },
      { ...fakeAction(), type: 'b' },
    ];

    const wrapper = factory({
      propsData: { actions },
    });

    await flushPromises();
    await wrapper.activateAllMenus();

    const menuItems = wrapper.findAll('.v-list-item');

    expect(menuItems.length).toBe(2);

    menuItems.at(1).trigger('click');

    expect(actions[1].method).toHaveBeenCalledTimes(1);
  });
});
