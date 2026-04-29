import { generateShallowRenderer, generateRenderer, flushPromises } from '@unit/utils/vue';
import { installResizeObserver, uninstallResizeObserver } from '@unit/utils/resize-observer';
import { ackAction, deleteAction, editAction, fakeAction } from '@unit/data/actions-panel';

import MassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';

const cActionBtnStub = {
  template: '<button class="c-action-btn-stub" type="button" v-on="$listeners" />',
};

const stubs = {
  'c-action-btn': cActionBtnStub,
};

const snapshotStubs = {
  'c-action-btn': true,
};

const threeActionsWithTypes = [
  { ...editAction, type: 'edit' },
  { ...deleteAction, type: 'delete' },
  { ...ackAction, type: 'ack' },
];

describe('mass-actions-panel', () => {
  const factory = generateShallowRenderer(MassActionsPanel, { stubs });
  const snapshotFactory = generateRenderer(MassActionsPanel, { stubs: snapshotStubs });

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

    const wrapper = snapshotFactory({
      propsData: { actions },
    });

    await flushPromises();
    await wrapper.activateAllMenus();

    const menuItems = wrapper.findAll('.v-list-item');

    expect(menuItems.length).toBe(2);

    menuItems.at(1).trigger('click');

    expect(actions[1].method).toHaveBeenCalledTimes(1);
  });

  it('Renders `mass-actions-panel` with actions on the large size', async () => {
    installResizeObserver(400);

    const wrapper = snapshotFactory({
      propsData: {
        actions: threeActionsWithTypes,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `mass-actions-panel` with actions correctly on the tablet size', async () => {
    installResizeObserver(100);

    const wrapper = snapshotFactory({
      propsData: {
        actions: threeActionsWithTypes,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `mass-actions-panel` with actions correctly on the mobile size', async () => {
    installResizeObserver(36);

    const wrapper = snapshotFactory({
      propsData: {
        actions: threeActionsWithTypes,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
