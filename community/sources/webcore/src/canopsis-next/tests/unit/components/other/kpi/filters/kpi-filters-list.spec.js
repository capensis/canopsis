import { generateRenderer } from '@unit/utils/vue';

import KpiFiltersList from '@/components/other/kpi/filters/kpi-filters-list.vue';
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-table-pagination': true,
  'kpi-filters-expand-item': true,
};

const selectExpandButtonByRow = (wrapper, index) => wrapper
  .findAll('tbody > tr')
  .at(index)
  .find('.v-data-table__expand-icon');

describe('kpi-filters-list', () => {
  const filtersItems = [
    {
      _id: 'c0ed9d92-67eb-4dc7-a2ab-9a551d45b9bf',
      name: 'Filter',
      created: 1614861888,
      updated: 1614861888,
    },
    {
      _id: '441a2a17-9036-48a3-9ff7-f393487395a9',
      name: 'Filter 2',
      created: 1614863990,
      updated: 1614863990,
    },
    {
      _id: '1cae4b8a-f598-480a-ad0c-b0a89a5c2e93',
      name: 'Filter 3',
      created: 1614864049,
      updated: 1614864049,
    },
    {
      _id: 'c46bffd9-8f5a-4c6c-b045-416e23ab1d44',
      name: 'Filter 4',
      created: 1614857014,
      updated: 1614857014,
    },
    {
      _id: 'd2403af7-712d-4353-911e-376f7a8053a7',
      name: 'Filter 5',
      created: 1613620731,
      updated: 1613620731,
    },
    {
      _id: '9bbb623c-7537-4c3b-afc0-0ace4f25a76b',
      name: 'Filter 6',
      created: 1613620721,
      updated: 1613620721,
    },
    {
      _id: 'fd35fcc4-36b0-445d-85be-999cc939047a',
      name: 'Filter 7',
      created: 1614864251,
      updated: 1614864251,
    },
    {
      _id: '70bfae47-cfdf-4a2c-9b43-3427f6aabea2',
      name: 'Filter 8',
      created: 1622781697,
      updated: 1622781697,
    },
    {
      _id: 'b3f67a16-019a-4694-9d74-ed762affaa04',
      name: 'Filter 9',
      created: 1615435601,
      updated: 1615435601,
    },
    {
      _id: 'e1f3e64a-dc99-42ed-af72-d8678f2e62bf',
      name: 'Filter 10',
      created: 1615442556,
      updated: 1615442556,
    },
    {
      _id: '15094f5a-9472-4700-b0cd-52305f754754',
      name: 'Filter 11',
      created: 1615440560,
      updated: 1615440560,
    },
  ];

  const snapshotFactory = generateRenderer(KpiFiltersList, { stubs });

  it('Renders `kpi-filters-list` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters: [],
        options: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `kpi-filters-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters: filtersItems,
        options: {
          page: 2,
          itemsPerPage: 10,
          sortBy: ['created'],
          sortDesc: [],
        },
        totalItems: 50,
        pending: true,
        removable: true,
        updatable: true,
        duplicable: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `kpi-filters-list` with expanded panel', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters: filtersItems,
        options: {
          page: 1,
          itemsPerPage: 10,
          sortBy: [],
          sortDesc: [],
        },
        totalItems: 50,
      },
    });

    await selectExpandButtonByRow(wrapper, 0).triggerCustomEvent('expand');

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `kpi-filters-list` with updatable and old_entity_patterns', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters: filtersItems.map(item => ({ ...item, old_entity_patterns: true })),
        options: {
          page: 1,
          itemsPerPage: 10,
          sortBy: [],
          sortDesc: [],
        },
        updatable: true,
        totalItems: 50,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
