import { generateShallowRenderer, generateRenderer, flushPromises } from '@unit/utils/vue';
import { deleteAction, editAction, fakeAction } from '@unit/data/actions-panel';

import MassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';

const actionsPanelBtnStub = {
  props: { action: { type: Object, required: true } },
  template: '<button class="actions-panel-btn" @click="action.method && action.method()"><slot /></button>',
};

const actionsPanelMenuStub = {
  props: { actions: { type: Array, default: () => [] } },
  template: `
    <div class="actions-panel-menu">
      <button
        v-for="(action, i) in actions"
        :key="i"
        class="actions-panel-menu-item"
        @click="action.method && action.method()"
      />
    </div>
  `,
};

const stubs = {
  'actions-panel-btn': actionsPanelBtnStub,
  'actions-panel-menu': actionsPanelMenuStub,
};

const snapshotStubs = {
  'c-action-btn': true,
  'c-list': true,
};

describe('mass-actions-panel', () => {
  const factory = generateShallowRenderer(MassActionsPanel, { stubs });
  const snapshotFactory = generateRenderer(MassActionsPanel, { stubs: snapshotStubs });

  it('Method into list called after trigger click on action item button. On the large size.', () => {
    const actions = [
      fakeAction(),
      fakeAction(),
    ];
    const wrapper = factory({
      propsData: {
        actions,
      },
      mocks: {
        $mq: 'l+',
      },
    });

    const actionElements = wrapper.findAll('button.actions-panel-btn');

    expect(actionElements).toHaveLength(actions.length);

    const secondActionElement = actionElements.at(1);

    secondActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toHaveBeenCalledTimes(1);
  });

  it('Method into dropdown called after trigger click on action item button. On the tablet size.', () => {
    const actions = [
      fakeAction(),
      fakeAction(),
    ];
    const wrapper = factory({
      propsData: {
        actions,
      },
      mocks: {
        $mq: 't',
      },
    });

    const dropdownActionElements = wrapper.findAll('button.actions-panel-menu-item');

    expect(dropdownActionElements).toHaveLength(actions.length);

    const secondActionElement = dropdownActionElements.at(1);
    secondActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  it('Renders `mass-actions-panel` with actions on the large size', async () => {
    const wrapper = snapshotFactory({
      mocks: {
        $mq: 'l+',
      },
      propsData: {
        actions: [editAction, deleteAction],
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `mass-actions-panel` with actions correctly on the tablet size', async () => {
    const wrapper = snapshotFactory({
      mocks: {
        $mq: 't',
      },
      propsData: {
        actions: [editAction, deleteAction],
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `mass-actions-panel` with actions correctly on the mobile size', async () => {
    const wrapper = snapshotFactory({
      mocks: {
        $mq: 'm',
      },
      propsData: {
        actions: [editAction, deleteAction],
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
