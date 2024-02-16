import flushPromises from 'flush-promises';

import { generateRenderer } from '@unit/utils/vue';
import { fakeUsersForTreeview } from '@unit/data/treeview';

import CTreeviewDataTable from '@/components/common/table/c-treeview-data-table.vue';

describe('c-treeview-data-table', () => {
  const snapshotItems = fakeUsersForTreeview({ count: 3, fake: false, depths: 2 });
  const headers = [
    { text: 'Username', value: 'username' },
    { text: 'Firstname', value: 'firstname' },
    { text: 'Lastname', value: 'lastname' },
    { text: 'Email', value: 'email' },
  ];

  const snapshotFactory = generateRenderer(CTreeviewDataTable);

  it('Renders `c-treeview-data-table` with default and required props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with default and required props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
        loading: true,
        dark: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with `openAll` prop', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers,
        items: snapshotItems,
        openAll: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with custom props and expand slots', async () => {
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

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with default props and data-table `items` slot', async () => {
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

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with default props and data-table values slots', async () => {
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

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-treeview-data-table` with default and required props (changed headers)', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        headers: headers.slice(0, -1),
        items: snapshotItems,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
