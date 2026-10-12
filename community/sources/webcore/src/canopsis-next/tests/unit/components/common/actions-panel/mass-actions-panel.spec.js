import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  ackAction,
  deleteAction,
  editAction,
  fakeAction,
} from '@unit/data/actions-panel';

import { MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP } from '@/constants';

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

  it('Method into list called after trigger click on action item button. Size \'xl\'', async () => {
    const actions = [
      fakeAction(),
      fakeAction(),
    ];

    const wrapper = factory({
      propsData: {
        actions,
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();
    const actionElements = wrapper.findAll('button.actions-panel-btn');

    expect(actionElements).toHaveLength(actions.length);

    const secondActionElement = actionElements.at(1);

    secondActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  it('Method into dropdown called after trigger click on action item button. Size \'m\'', async () => {
    const actions = [
      fakeAction(),
      fakeAction(),
    ];
    const wrapper = factory({
      propsData: {
        actions,
      },
      mocks: {
        $mq: 'm',
      },
    });

    await flushPromises();

    const dropdownActionElements = wrapper.findAll('button.actions-panel-menu-item');

    expect(dropdownActionElements).toHaveLength(actions.length);

    const secondDropdownActionElement = dropdownActionElements.at(1);

    secondDropdownActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  it('Method into dropdown called after trigger click on action item button. Size \'xl\'', async () => {
    const inlineCount = 2;
    const actions = [
      fakeAction(),
      fakeAction(),
      fakeAction(),
    ];
    const wrapper = factory({
      propsData: {
        actions,
        inlineCount,
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();

    const dropdownActionElements = wrapper.findAll('button.actions-panel-menu-item');

    expect(dropdownActionElements).toHaveLength(actions.length - inlineCount + 1);

    const firstDropdownActionElement = dropdownActionElements.at(0);

    firstDropdownActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  it('Renders `mass-actions-panel` with actions correctly. Size \'xl\'', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction],
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it.each(
    Object.keys(MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP),
  )('Renders `mass-actions-panel` with three actions and 3 inlineCount correctly. Size \'%s\'', async ($mq) => {
    const wrapper = snapshotFactory({
      propsData: {
        inlineCount: 3,
        actions: [editAction, deleteAction, ackAction],
      },
      mocks: {
        $mq,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `mass-actions-panel` with three actions and 2 inlineCount. Size \'xl\'', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        inlineCount: 2,
        actions: [editAction, deleteAction, ackAction],
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
