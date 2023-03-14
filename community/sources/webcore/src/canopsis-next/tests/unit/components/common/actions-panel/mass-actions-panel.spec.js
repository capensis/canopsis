import { shallowMount, mount, createVueInstance } from '@unit/utils/vue';

import { deleteAction, editAction, fakeAction } from '@unit/data/actions-panel';
import { createButtonStub } from '@unit/stubs/button';

import MassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';

const localVue = createVueInstance();

const stubs = {
  'c-action-btn': createButtonStub('c-action-btn'),
  'v-list-tile': createButtonStub('v-list-tile'),
};

const snapshotStubs = {
  'c-action-btn': true,
};

const factory = (options = {}) => shallowMount(MassActionsPanel, {
  localVue,
  stubs,
  ...options,
});

const snapshotFactory = (options = {}) => mount(MassActionsPanel, {
  localVue,
  stubs: snapshotStubs,
  ...options,
});

describe('mass-actions-panel', () => {
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
        $windowSize: 'l+',
      },
    });

    const actionElements = wrapper.findAll('button.c-action-btn');

    expect(actionElements).toHaveLength(actions.length);

    const secondActionElement = actionElements.at(1);

    secondActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
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
        $windowSize: 't',
      },
    });

    const dropdownActionElements = wrapper.findAll('v-menu-stub button.v-list-tile');

    expect(dropdownActionElements).toHaveLength(actions.length);

    const secondActionElement = dropdownActionElements.at(1);
    secondActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  it('Renders `mass-actions-panel` with actions on the large size', () => {
    const wrapper = snapshotFactory({
      mocks: {
        $windowSize: 'l+',
      },
      propsData: {
        actions: [editAction, deleteAction],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `mass-actions-panel` with actions correctly on the tablet size', () => {
    const wrapper = snapshotFactory({
      mocks: {
        $windowSize: 't',
      },
      propsData: {
        actions: [editAction, deleteAction],
      },
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `mass-actions-panel` with actions correctly on the mobile size', () => {
    const wrapper = snapshotFactory({
      mocks: {
        $windowSize: 'm',
      },
      propsData: {
        actions: [editAction, deleteAction],
      },
    });

    const dropdownContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });
});
