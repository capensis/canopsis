import { generateRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';

import ManageInfos from '@/components/widgets/context/manage-infos.vue';

describe('manage-infos', () => {
  const $modals = mockModals();
  const infos = [
    { name: 'charlie', value: 'third' },
    { name: 'alpha', value: 'first' },
    { name: 'bravo', value: 'second' },
  ];
  const factory = generateRenderer(ManageInfos, {
    propsData: { infos },
    mocks: { $modals },
    stubs: { 'c-action-btn': true },
  });

  test('Edits the selected info after sorting and pagination', async () => {
    const wrapper = factory();
    const footer = wrapper.findComponent({ name: 'VDataFooter' });

    await footer.triggerCustomEvent('update:options', {
      itemsPerPage: 1, sortBy: ['name'], sortDesc: [false],
    });
    await footer.triggerCustomEvent('update:options', { page: 2 });

    expect(wrapper.find('tbody tr td').text()).toBe('bravo');

    await wrapper.find('c-action-btn-stub[type="edit"]').triggerCustomEvent('click');

    const [{ config }] = $modals.show.mock.calls[0];
    const editedInfo = { name: 'renamed', value: 'updated' };

    expect(config.entityInfo).toEqual(infos[2]);

    config.action(editedInfo);

    expect(wrapper).toEmitInput([infos[0], infos[1], editedInfo]);
  });

  test('Deletes the selected info after sorting and pagination', async () => {
    const wrapper = factory();
    const footer = wrapper.findComponent({ name: 'VDataFooter' });

    await footer.triggerCustomEvent('update:options', {
      itemsPerPage: 1, sortBy: ['name'], sortDesc: [false],
    });
    await footer.triggerCustomEvent('update:options', { page: 2 });

    expect(wrapper.find('tbody tr td').text()).toBe('bravo');

    await wrapper.find('c-action-btn-stub[type="delete"]').triggerCustomEvent('click');

    expect(wrapper).toEmitInput([infos[0], infos[1]]);
  });
});
