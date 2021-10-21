import { mount, createVueInstance } from '@unit/utils/vue';

import CEmptyDataTableColumns from '@/components/common/table/c-empty-data-table-columns.vue';

const localVue = createVueInstance();

describe('c-empty-data-table-columns', () => {
  it('Renders `c-empty-data-table-columns` correctly', () => {
    const wrapper = mount(CEmptyDataTableColumns, { localVue });

    expect(wrapper.element).toMatchSnapshot();
  });
});
