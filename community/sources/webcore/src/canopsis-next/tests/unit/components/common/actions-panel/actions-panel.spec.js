import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { deleteAction, editAction, fakeAction } from '@unit/data/actions-panel';
import { createButtonStub } from '@unit/stubs/button';
import ActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

const stubs = {
  'c-action-btn': createButtonStub('c-action-btn'),
  'v-list-item': createButtonStub('v-list-item'),
};

const snapshotStubs = {
  'c-action-btn': true,
};

describe('actions-panel', () => {
  const factory = generateShallowRenderer(ActionsPanel, { stubs });
  const snapshotFactory = generateRenderer(ActionsPanel, { stubs: snapshotStubs });

  it('Method into list called after trigger click on action item button. On the extra large size.', async () => {
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

  it('Method into dropdown called after trigger click on action item button. On the extra large size.', async () => {
    const inlineCount = 1;
    const actions = [
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

    const dropdownActionElements = wrapper.findAll('v-menu-stub .v-list-item');

    expect(dropdownActionElements).toHaveLength(actions.length - inlineCount);

    const firstDropdownActionElement = dropdownActionElements.at(0);

    firstDropdownActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  it('Renders `actions-panel` with default props correctly on the extra large size', async () => {
    const wrapper = snapshotFactory({
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();

    const dropdownContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with default props correctly on the large size', async () => {
    const wrapper = snapshotFactory({
      mocks: {
        $mq: 'l',
      },
    });

    await flushPromises();

    const dropdownContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with actions correctly on the extra large size', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction],
      },
      mocks: {
        $mq: 'xl',
      },
    });

    await flushPromises();

    const dropdownContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with actions correctly on the large size', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction],
      },
      mocks: {
        $mq: 'l',
      },
    });

    await flushPromises();

    const dropdownContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with dropdown actions correctly on the large size', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction],
      },
      mocks: {
        $mq: 'l',
      },
    });

    const dropdownContent = wrapper.findMenu();

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with dropdown actions correctly on the tablet size', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction],
      },
      mocks: {
        $mq: 't',
      },
    });

    await flushPromises();

    const dropdownContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with dropdown actions correctly on the mobile size', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        actions: [editAction, deleteAction],
      },
      mocks: {
        $mq: 'm',
      },
    });

    await flushPromises();

    const dropdownContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });
});
