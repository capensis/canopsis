import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { ackAction, deleteAction, editAction, fakeAction } from '@unit/data/actions-panel';
import { createButtonStub } from '@unit/stubs/button';

import { ALARM_ACTION_BUTTON_MARGINS, ALARM_ACTION_BUTTON_WIDTHS, ALARM_DENSE_TYPES } from '@/constants';

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
  const snapshotFactory = generateRenderer(ActionsPanel, { stubs: snapshotStubs, attachTo: document.body });

  const ACTION_WIDTHS = {
    large: ALARM_ACTION_BUTTON_WIDTHS[ALARM_DENSE_TYPES.large]
      + ALARM_ACTION_BUTTON_MARGINS[ALARM_DENSE_TYPES.large],
    medium: ALARM_ACTION_BUTTON_WIDTHS[ALARM_DENSE_TYPES.medium]
      + ALARM_ACTION_BUTTON_MARGINS[ALARM_DENSE_TYPES.medium],
    small: ALARM_ACTION_BUTTON_WIDTHS[ALARM_DENSE_TYPES.small]
      + ALARM_ACTION_BUTTON_MARGINS[ALARM_DENSE_TYPES.small],
  };

  test('Method into list called after trigger click on action item button. With width for two items. For large dense.', async () => {
    const actions = [
      fakeAction(),
      fakeAction(),
    ];

    const wrapper = factory({
      propsData: {
        actions,
      },
    });

    wrapper.setData({ width: ACTION_WIDTHS.large * 2 });

    await flushPromises();
    const actionElements = wrapper.findAll('.c-action-btn');

    expect(actionElements).toHaveLength(actions.length);

    const secondActionElement = actionElements.at(1);

    secondActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  test('Method into dropdown called after trigger click on action item button. With width for three items. For large dense.', async () => {
    const inlineCount = 1;
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
    });

    wrapper.setData({ width: ACTION_WIDTHS.large * 2 });

    await flushPromises();

    const dropdownActionElements = wrapper.findAll('v-menu-stub .v-list-tile');

    expect(dropdownActionElements).toHaveLength(actions.length - inlineCount);

    const firstDropdownActionElement = dropdownActionElements.at(0);

    firstDropdownActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  test('Renders `actions-panel` with default props correctly. With zero width.', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `actions-panel` with three actions correctly. With width for two items. For large dense.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction, ackAction],
      },
    });

    wrapper.setData({ width: ACTION_WIDTHS.large * 2 });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `actions-panel` with three actions correctly. With width for two items. For medium dense.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction, ackAction],
        medium: true,
      },
    });

    wrapper.setData({ width: ACTION_WIDTHS.medium * 2 });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `actions-panel` with three actions correctly. With width for two items. For small dense.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction, ackAction],
        small: true,
      },
    });

    wrapper.setData({ width: ACTION_WIDTHS.small * 2 });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `actions-panel` with three actions correctly. With width for one item.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction, ackAction],
      },
    });

    wrapper.setData({ width: ACTION_WIDTHS.large });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `actions-panel` with three actions correctly. With width for three items.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction, ackAction],
      },
    });

    wrapper.setData({ width: ACTION_WIDTHS.large * 3 });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `actions-panel` with three actions correctly. With width for three items and changing to one item.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction, ackAction],
      },
    });

    wrapper.setData({ width: ACTION_WIDTHS.large * 3 });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();

    wrapper.setData({ width: ACTION_WIDTHS.large });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
