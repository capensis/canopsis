import Faker from 'faker';
import { shallowMount, mount, createVueInstance } from '@unit/utils/vue';

import { createButtonStub } from '@unit/stubs/button';
import { createListItemStub } from '@unit/stubs/list';
import ActionsPanelItem from '@/components/common/actions-panel/actions-panel-item.vue';

const localVue = createVueInstance();

const stubs = {
  'v-list-tile': createListItemStub('v-list-tile'),
  'c-action-btn': createButtonStub('c-action-btn'),
};

const snapshotStubs = {
  'c-action-btn': true,
};

const factory = (options = {}) => shallowMount(ActionsPanelItem, {
  localVue,
  stubs,
  ...options,
});

describe('actions-panel-item', () => {
  it('Method called after trigger click on the list item', () => {
    const method = jest.fn();
    const wrapper = factory({
      propsData: {
        title: Faker.datatype.string(),
        icon: Faker.datatype.string(),
        isDropDown: true,
        method,
      },
    });

    const listItemElement = wrapper.find('li.v-list-tile');

    listItemElement.trigger('click');

    expect(method).toBeCalledTimes(1);
  });

  it('Method called after trigger click on the action button', () => {
    const method = jest.fn();
    const wrapper = factory({
      propsData: {
        title: Faker.datatype.string(),
        icon: Faker.datatype.string(),
        method,
      },
    });

    const actionButtonElement = wrapper.find('button.c-action-btn');

    actionButtonElement.trigger('click');

    expect(method).toBeCalledTimes(1);
  });

  it('Renders `actions-panel-item` with default and required props correctly', () => {
    const wrapper = mount(ActionsPanelItem, {
      localVue,
      stubs: snapshotStubs,
      propsData: {
        title: 'Custom title',
        icon: 'customIcon',
        method: jest.fn(),
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders disabled `actions-panel-item` with default and required props correctly', () => {
    const wrapper = mount(ActionsPanelItem, {
      localVue,
      stubs: snapshotStubs,
      propsData: {
        title: 'Disabled custom title',
        icon: 'disabledCustomIcon',
        iconColor: 'disabledCustomIconColor',
        disabled: true,
        method: jest.fn(),
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `actions-panel-item` with custom props as dropdown item correctly', () => {
    const wrapper = mount(ActionsPanelItem, {
      localVue,
      stubs: snapshotStubs,
      propsData: {
        title: 'Custom dropdown title',
        icon: 'customDropdownIcon',
        iconColor: 'customDropdownIconColor',
        isDropDown: true,
        method: jest.fn(),
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `actions-panel-item` with custom props as dropdown item correctly 2', () => {
    const wrapper = mount(ActionsPanelItem, {
      localVue,
      stubs: snapshotStubs,
      propsData: {
        title: 'Disabled custom dropdown title',
        icon: 'disabledCustomDropdownIcon',
        iconColor: 'disabledCustomDropdownIconColor',
        isDropDown: true,
        disabled: true,
        method: jest.fn(),
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
