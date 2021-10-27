import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';
import { fakeUsersForTreeview } from '@unit/data/treeview';

import CTreeviewDataTable from '@/components/common/table/c-treeview-data-table.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(CTreeviewDataTable, {
  localVue,
  stubs: {
    transition: true,
  },

  ...options,
});

describe('c-treeview-data-table', () => {
  const snapshotItems = fakeUsersForTreeview({ count: 3, fake: false, depths: 2 });
  const headers = [
    { text: 'Username', value: 'username' },
    { text: 'Firstname', value: 'firstname' },
    { text: 'Lastname', value: 'lastname' },
    { text: 'Email', value: 'email' },
  ];

  it('Renders `c-treeview-data-table` with default and required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with default and required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
        loading: true,
        dark: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with `openAll` prop', () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
        openAll: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with custom props and expand slots', () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
        light: true,
      },
      slots: {
        expand: '<div class="expand-slot" />',
        'expand-append': '<div class="expand-append-slot" />',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with default props and data-table `items` slot', () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
        light: true,
      },
      slots: {
        items: '<tr class="items-slot" />',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with default props and data-table values slots', () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
        light: true,
      },
      slots: {
        username: '<span class="username-slot" />',
        firstname: '<span class="firstname-slot" />',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with default and required props (changed headers)', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers: headers.slice(0, -1),
        items: snapshotItems,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with default and required props and trigger click', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
      },
    });

    const toggleIcon = wrapper.find('.v-treeview-node__toggle');

    await toggleIcon.trigger('click');
    await flushPromises();

    const secondToggleIcon = wrapper.find('.v-treeview-node__children .v-treeview-node__toggle');

    await secondToggleIcon.trigger('click');
    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
