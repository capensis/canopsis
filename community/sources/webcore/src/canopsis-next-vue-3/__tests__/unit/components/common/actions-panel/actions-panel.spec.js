import {
  fakeAction,
  fastPbehaviorAddAction,
  fastPbehaviorRemoveAction,
} from '@unit/data/actions-panel';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import ActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

const actionsPanelBtnStub = {
  props: { action: { type: Object, required: true } },
  template:
    '<button class="actions-panel-btn" @click="action.method && action.method()"><slot /></button>',
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

describe('actions-panel', () => {
  const factory = generateShallowRenderer(ActionsPanel, { stubs });

  it('Method into list called after trigger click on action item button. Size \'xl\'', async () => {
    const actions = [fakeAction(), fakeAction()];

    const wrapper = factory({
      propsData: {
        actions,
        ignoreMediaQuery: true,
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();
    const actionElements = wrapper.findAll('button.actions-panel-btn');

    expect(actionElements).toHaveLength(actions.length);

    const secondActionElement = actionElements.at(1);

    await secondActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  it('Fast pbehavior add action method called after trigger click', async () => {
    const actions = [fastPbehaviorAddAction];

    const wrapper = factory({
      propsData: {
        actions,
        ignoreMediaQuery: true,
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();

    const addButton = wrapper.find('button.actions-panel-btn');
    await addButton.trigger('click');

    expect(fastPbehaviorAddAction.method).toHaveBeenCalledTimes(1);
  });

  it('Fast pbehavior remove action method called after trigger click', async () => {
    const actions = [fastPbehaviorRemoveAction];

    const wrapper = factory({
      propsData: {
        actions,
        ignoreMediaQuery: true,
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();

    const removeButton = wrapper.find('button.actions-panel-btn');
    await removeButton.trigger('click');

    expect(fastPbehaviorRemoveAction.method).toHaveBeenCalledTimes(1);
  });

  it('Method into dropdown called after trigger click on action item button. Size \'xl\'', async () => {
    const inlineCount = 2;
    const actions = [fakeAction(), fakeAction(), fakeAction()];
    const wrapper = factory({
      propsData: {
        actions,
        inlineCount,
        ignoreMediaQuery: true,
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();

    const dropdownActionElements = wrapper.findAll('button.actions-panel-menu-item');

    expect(dropdownActionElements).toHaveLength(actions.length - inlineCount + 1);

    const firstDropdownActionElement = dropdownActionElements.at(0);

    await firstDropdownActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });
});
