import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { ackAction, deleteAction, editAction, fakeAction } from '@unit/data/actions-panel';
import { createButtonStub } from '@unit/stubs/button';

import { MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP } from '@/constants';

import ActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

const stubs = {
  'c-action-btn': createButtonStub('c-action-btn'),
  'v-list-tile': createButtonStub('v-list-tile'),
};

const snapshotStubs = {
  'c-action-btn': true,
};

describe('actions-panel', () => {
  const factory = generateShallowRenderer(ActionsPanel, { stubs });
  const snapshotFactory = generateRenderer(ActionsPanel, { stubs: snapshotStubs });

  test('Method into list called after trigger click on action item button. Size \'xl\'', async () => {
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
    const actionElements = wrapper.findAll('.c-action-btn');

    expect(actionElements).toHaveLength(actions.length);

    const secondActionElement = actionElements.at(1);

    secondActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  test('Method into dropdown called after trigger click on action item button. Size \'xl\'', async () => {
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

    const dropdownActionElements = wrapper.findAll('v-menu-stub .v-list-tile');

    expect(dropdownActionElements).toHaveLength(actions.length - inlineCount + 1);

    const firstDropdownActionElement = dropdownActionElements.at(0);

    firstDropdownActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  test('Renders `actions-panel` with default props correctly', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `actions-panel` with actions correctly. Size \'xl\'', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction],
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it.each(
    Object.keys(MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP),
  )('Renders `actions-panel` with three actions and 3 inlineCount correctly. Size \'%s\'', async ($mq) => {
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
    expect(wrapper).toMatchMenuSnapshot();
  });

  it.each(
    Object.keys(MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP),
  )('Renders `actions-panel` with three actions, 3 inlineCount and ignoreMediaQuery correctly. Size %s', async ($mq) => {
    const wrapper = snapshotFactory({
      propsData: {
        inlineCount: 3,
        actions: [editAction, deleteAction, ackAction],
        ignoreMediaQuery: true,
      },
      mocks: {
        $mq,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `actions-panel` with three actions, 2 inlineCount and ignoreMediaQuery', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        inlineCount: 2,
        actions: [editAction, deleteAction, ackAction],
        ignoreMediaQuery: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `actions-panel` with three actions, 3 inlineCount and ignoreMediaQuery', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        inlineCount: 3,
        actions: [editAction, deleteAction, ackAction],
        ignoreMediaQuery: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
