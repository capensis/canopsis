import { mockModals } from '@unit/utils/mock-hooks';
import { generateRenderer } from '@unit/utils/vue';

import ManageInfos from '@/components/widgets/context/manage-infos.vue';

const selectRows = wrapper => wrapper.findAll('tbody > tr');
const selectAction = (row, type) => row.find(`c-action-btn-stub[type="${type}"]`);
const sortDescendingAndGoToSecondPage = async (wrapper) => {
  const nameHeader = wrapper.findAll('th')[0];

  await nameHeader.trigger('click');
  await nameHeader.trigger('click');
  wrapper.findComponent({ name: 'VPagination' }).vm.$emit('update:modelValue', 2);
  await wrapper.vm.$nextTick();
};

describe('manage-infos', () => {
  const modals = mockModals();
  const factory = generateRenderer(ManageInfos, {
    stubs: { 'c-action-btn': true },
  });
  const infos = Array.from({ length: 12 }, (_, index) => ({
    name: `info-${String(index + 1).padStart(2, '0')}`,
    description: `Description ${index + 1}`,
    value: index + 1,
  }));

  test.each([1, 2])('deletes the selected information on page %s', async (page) => {
    const wrapper = factory({ props: { infos } });

    wrapper.findComponent({ name: 'VPagination' }).vm.$emit('update:modelValue', page);
    await wrapper.vm.$nextTick();

    const selectedInfo = infos[(page - 1) * 10 + 1];
    const row = selectRows(wrapper)[1];

    expect(row.text()).toContain(selectedInfo.name);
    await selectAction(row, 'delete').trigger('click');

    expect(wrapper.emitted('update:infos')).toEqual([
      [infos.filter(info => info.name !== selectedInfo.name)],
    ]);
  });

  test('deletes the selected information after changing the page size', async () => {
    const wrapper = factory({ props: { infos } });

    wrapper.findComponent({ name: 'VSelect' }).vm.$emit('update:modelValue', 5);
    await wrapper.vm.$nextTick();
    wrapper.findComponent({ name: 'VPagination' }).vm.$emit('update:modelValue', 3);
    await wrapper.vm.$nextTick();

    const row = selectRows(wrapper)[0];

    expect(row.text()).toContain(infos[10].name);
    await selectAction(row, 'delete').trigger('click');

    expect(wrapper.emitted('update:infos')).toEqual([
      [infos.filter((_, index) => index !== 10)],
    ]);
  });

  test('deletes the selected information after sorting and changing pages', async () => {
    const wrapper = factory({ props: { infos } });

    await sortDescendingAndGoToSecondPage(wrapper);

    const row = selectRows(wrapper)[0];
    const selectedInfo = infos[1];

    expect(row.text()).toContain(selectedInfo.name);
    await selectAction(row, 'delete').trigger('click');

    expect(wrapper.emitted('update:infos')).toEqual([
      [infos.filter(info => info.name !== selectedInfo.name)],
    ]);
  });

  test('edits the selected information after sorting and changing pages', async () => {
    const wrapper = factory({ props: { infos } });

    await sortDescendingAndGoToSecondPage(wrapper);

    const row = selectRows(wrapper)[0];
    const selectedInfo = infos[1];

    expect(row.text()).toContain(selectedInfo.name);
    await selectAction(row, 'edit').trigger('click');

    const { config } = modals.show.mock.calls[0][0];
    const editedInfo = { ...selectedInfo, name: 'renamed', value: 'updated' };

    expect(config.entityInfo).toEqual(selectedInfo);
    config.action(editedInfo);

    expect(wrapper.emitted('update:infos')).toEqual([
      [infos.map(info => info.name === selectedInfo.name ? editedInfo : info)],
    ]);
  });
});
