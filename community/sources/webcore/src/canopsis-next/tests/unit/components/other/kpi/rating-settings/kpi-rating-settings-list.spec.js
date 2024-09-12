import { flushPromises, generateRenderer } from '@unit/utils/vue';

import KpiRatingSettingsList from '@/components/other/kpi/rating-settings/kpi-rating-settings-list.vue';
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';

const stubs = {
  'c-advanced-data-table': CAdvancedDataTable,
  'c-search-field': true,
  'c-expand-btn': true,
  'c-action-btn': true,
  'c-table-pagination': true,
};

describe('kpi-rating-settings-list', () => {
  const ratingSettingsItems = [
    {
      id: 'c0ed9d92-67eb-4dc7-a2ab-9a551d45b9bf',
      label: 'Rating setting',
      enabled: true,
    },
    {
      id: '441a2a17-9036-48a3-9ff7-f393487395a9',
      label: 'Rating setting 2',
      enabled: true,
    },
    {
      id: '1cae4b8a-f598-480a-ad0c-b0a89a5c2e93',
      label: 'Rating setting 3',
      enabled: true,
    },
    {
      id: 'c46bffd9-8f5a-4c6c-b045-416e23ab1d44',
      label: 'Rating setting 4',
      enabled: false,
    },
    {
      id: 'd2403af7-712d-4353-911e-376f7a8053a7',
      label: 'Rating setting 5',
      enabled: false,
    },
    {
      id: '9bbb623c-7537-4c3b-afc0-0ace4f25a76b',
      label: 'Rating setting 6',
      enabled: false,
    },
    {
      id: 'fd35fcc4-36b0-445d-85be-999cc939047a',
      label: 'Rating setting 7',
      enabled: false,
    },
    {
      id: '70bfae47-cfdf-4a2c-9b43-3427f6aabea2',
      label: 'Rating setting 8',
      enabled: false,
    },
    {
      id: 'b3f67a16-019a-4694-9d74-ed762affaa04',
      label: 'Rating setting 9',
      enabled: false,
    },
    {
      id: 'e1f3e64a-dc99-42ed-af72-d8678f2e62bf',
      label: 'Rating setting 10',
      enabled: true,
    },
    {
      id: '15094f5a-9472-4700-b0cd-52305f754754',
      label: 'Rating setting 11',
      enabled: true,
    },
  ];

  const snapshotFactory = generateRenderer(KpiRatingSettingsList, { stubs });

  it('Rating settings changed after enable rows', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        ratingSettings: ratingSettingsItems,
        pagination: {
          page: 1,
          rowsPerPage: 10,
          search: '',
          sortBy: '',
          descending: false,
        },
        totalItems: 50,
        updatable: true,
      },
    });

    const rows = wrapper.findAll('tr > td');

    const enableButton = rows.at(0).find('input');

    await enableButton.trigger('click');

    const submitButton = wrapper.findAll('.v-btn').at(1);

    submitButton.trigger('click');

    const changeSelectedEvents = wrapper.emitted('change-selected');

    expect(changeSelectedEvents).toHaveLength(1);

    const [eventData] = changeSelectedEvents[0];

    const [firstRatingSetting] = ratingSettingsItems;
    expect(eventData).toEqual([{
      ...firstRatingSetting,
      enabled: !firstRatingSetting.enabled,
    }]);
  });

  it('Renders `kpi-rating-settings-list` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        ratingSettings: [],
        pagination: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `kpi-rating-settings-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        ratingSettings: ratingSettingsItems,
        pagination: {
          page: 2,
          rowsPerPage: 10,
          search: 'Rating setting',
          sortBy: '',
          descending: false,
        },
        totalItems: 50,
        pending: true,
        updatable: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `kpi-rating-settings-list` with enable rating', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        ratingSettings: ratingSettingsItems,
        pagination: {
          page: 1,
          rowsPerPage: 10,
          search: '',
          sortBy: '',
          descending: false,
        },
        totalItems: 50,
        updatable: true,
      },
    });

    const enableButton = wrapper
      .findAll('tr > td')
      .at(0)
      .find('input');

    enableButton.trigger('click');

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();

    enableButton.trigger('click');

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `kpi-rating-settings-list` after updated rating settings prop', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        ratingSettings: [],
        pagination: {
          page: 1,
          rowsPerPage: 10,
          search: '',
          sortBy: '',
          descending: false,
        },
        totalItems: 50,
        updatable: true,
      },
    });

    await wrapper.setProps({
      ratingSettings: ratingSettingsItems,
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
