import { shallowMount, mount, createVueInstance } from '@unit/utils/vue';

import { deleteAction, editAction, fakeAction } from '@unit/data/actions-panel';
import { actionsPanelItem } from '@unit/stubs/actions-panel';
import ActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

const localVue = createVueInstance();

const stubs = {
  'actions-panel-item': actionsPanelItem,
};

const snapshotStubs = {
  'c-action-btn': true,
};

const factory = (options = {}) => shallowMount(ActionsPanel, {
  localVue,
  stubs,
  ...options,
});

describe('actions-panel', () => {
  it('Method into list called after trigger click on action item button. On the extra large size.', () => {
    const actions = [
      fakeAction(),
      fakeAction(),
    ];
    const wrapper = factory({
      propsData: {
        actions,
      },
      mocks: {
        $windowSize: 'xl',
      },
    });

    const actionElements = wrapper.findAll('button.actions-panel-item');

    expect(actionElements).toHaveLength(actions.length);

    const secondActionElement = actionElements.at(1);

    secondActionElement.trigger('click');

    const [, secondAction] = actions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  it('Method into dropdown called after trigger click on action item button. On the extra large size.', () => {
    const dropDownActions = [
      fakeAction(),
      fakeAction(),
    ];
    const wrapper = factory({
      propsData: {
        dropDownActions,
      },
      mocks: {
        $windowSize: 'xl',
      },
    });

    const dropdownActionElements = wrapper.findAll('v-menu-stub button.actions-panel-item');

    expect(dropdownActionElements).toHaveLength(dropDownActions.length);

    const secondActionElement = dropdownActionElements.at(1);

    secondActionElement.trigger('click');

    const [, secondAction] = dropDownActions;
    expect(secondAction.method).toBeCalledTimes(1);
  });

  it('Method into list called after trigger click on action item button. On the large size.', () => {
    const action = fakeAction();
    const dropdownAction = fakeAction();
    const wrapper = factory({
      propsData: {
        actions: [action],
        dropDownActions: [dropdownAction],
      },
      mocks: {
        $windowSize: 'l',
      },
    });

    const actionElements = wrapper.findAll('button.actions-panel-item');

    expect(actionElements).toHaveLength(2);

    const secondActionElement = actionElements.at(1);

    secondActionElement.trigger('click');

    expect(dropdownAction.method).toBeCalledTimes(1);
  });

  it('Renders `actions-panel` with default props correctly on the extra large size', () => {
    const wrapper = mount(ActionsPanel, {
      localVue,
      stubs: snapshotStubs,
      mocks: {
        $windowSize: 'xl',
      },
    });

    const dropdownContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with default props correctly on the large size', () => {
    const wrapper = mount(ActionsPanel, {
      localVue,
      stubs: snapshotStubs,
      mocks: {
        $windowSize: 'l',
      },
    });

    const dropdownContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with actions correctly on the extra large size', () => {
    const wrapper = mount(ActionsPanel, {
      localVue,
      stubs: snapshotStubs,
      mocks: {
        $windowSize: 'xl',
      },
      propsData: {
        actions: [editAction, deleteAction],
      },
    });

    const dropdownContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with actions correctly on the large size', () => {
    const wrapper = mount(ActionsPanel, {
      localVue,
      stubs: snapshotStubs,
      mocks: {
        $windowSize: 'l',
      },
      propsData: {
        actions: [editAction, deleteAction],
      },
    });

    const dropdownContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with dropdown actions correctly on the extra large size', () => {
    const wrapper = mount(ActionsPanel, {
      localVue,
      stubs: snapshotStubs,
      mocks: {
        $windowSize: 'xl',
      },
      propsData: {
        dropDownActions: [editAction, deleteAction],
      },
    });

    const dropdownContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });

  it('Renders `actions-panel` with dropdown actions correctly on the large size', () => {
    const wrapper = mount(ActionsPanel, {
      localVue,
      stubs: snapshotStubs,
      mocks: {
        $windowSize: 'l',
      },
      propsData: {
        dropDownActions: [editAction, deleteAction],
      },
    });

    const dropdownContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(dropdownContent.element).toMatchSnapshot();
  });
});
