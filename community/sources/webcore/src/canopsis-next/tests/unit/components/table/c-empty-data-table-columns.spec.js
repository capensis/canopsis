import { generateRenderer } from '@unit/utils/vue';

import CEmptyDataTableColumns from '@/components/common/table/c-empty-data-table-columns.vue';

describe('c-empty-data-table-columns', () => {
  const snapshotFactory = generateRenderer(CEmptyDataTableColumns);

  it('Renders `c-empty-data-table-columns` correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });
});
