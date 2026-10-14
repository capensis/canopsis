import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateRenderer } from '@unit/utils/vue';

import ClickOutside from '@/services/click-outside';

import { saveJsonFile } from '@/helpers/file/files';
import { convertObjectToTreeview } from '@/helpers/treeview';

import VariablesHelp from '@/components/modals/common/variables-help.vue';

vi.mock('@/helpers/file/files', () => ({ saveJsonFile: vi.fn() }));

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'c-object-treeview': true,
  'c-ellipsis': true,
  'v-menu': { template: '<div><slot name="activator" /><slot /></div>' },
};

describe('variables-help', () => {
  mockDateNow(1386435600000);

  const $modals = mockModals();
  const $popups = mockPopups();
  const variablesObjectSecond = {
    number_prop: 1,
    string_prop: 'string',
    bool_prop: false,
    null_prop: null,
    undefined_prop: undefined,
    array_prop: [{ obj_prop: { deep_prop: 53, empty_obj: {} } }, []],
  };
  const variables = [
    convertObjectToTreeview({}, 'test'),
    convertObjectToTreeview(variablesObjectSecond, 'test2'),
  ];
  const factory = generateRenderer(VariablesHelp, {
    stubs,
    mocks: { $popups, $modals },
    parentComponent: { provide: { $clickOutside: new ClickOutside(), $system: {} } },
  });

  const createWrapper = () => factory({ propsData: { modal: { config: { variables } } } });

  test('copies a variable path successfully', async () => {
    const writeText = vi.fn().mockResolvedValue();
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText } });
    const wrapper = createWrapper();

    await wrapper.vm.copyPathToClipboard('test2.array_prop.[0].obj_prop.deep_prop');

    expect(writeText).toHaveBeenCalledWith('test2.array_prop.[0].obj_prop.deep_prop');
    expect($popups.success).toHaveBeenCalledWith({ text: 'Path copied to clipboard' });
  });

  test('shows an error when copying a variable path fails', async () => {
    const writeText = vi.fn().mockRejectedValue(new Error('clipboard unavailable'));
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText } });
    const wrapper = createWrapper();

    await wrapper.vm.copyPathToClipboard('test2.array_prop.[0].obj_prop.deep_prop');

    expect($popups.error).toHaveBeenCalledWith({ text: 'Something went wrong...' });
  });

  test('exports a tree item as JSON', () => {
    const wrapper = createWrapper();
    const arrayItem = variables[1].children.find(({ name }) => name === 'array_prop');

    wrapper.vm.exportAsJson(arrayItem);

    expect(saveJsonFile).toHaveBeenCalledWith(
      variablesObjectSecond.array_prop,
      'array_prop-07/12/2013 18:00:00',
    );
  });
});
