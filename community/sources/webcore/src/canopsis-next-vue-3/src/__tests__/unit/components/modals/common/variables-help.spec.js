import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { generateRenderer } from '@unit/utils/vue';

import ClickOutside from '@/services/click-outside';

import { saveJsonFile } from '@/helpers/file/files';
import { convertObjectToTreeview } from '@/helpers/treeview';

import VariablesHelp from '@/components/modals/common/variables-help.vue';

vi.mock('@/helpers/file/files', () => ({
  saveJsonFile: vi.fn(),
}));

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'patterns-form': true,
  'c-ellipsis': true,
  'v-menu': {
    template: `
      <div class="v-menu">
        <slot name="activator" />
        <slot />
      </div>
    `,
  },
};

const selectTreeviewNodes = wrapper => wrapper.findAll('.v-treeview-node');
const selectTreeviewNodeByIndex = (wrapper, index) => selectTreeviewNodes(wrapper).at(index);
const selectListTileByIndex = (wrapper, index) => wrapper.findAll('.v-list-item').at(index);

describe('variables-help', () => {
  mockDateNow(1386435600000);

  const $modals = mockModals();
  const $popups = mockPopups();
  const variablesObjectFirst = {};
  const variablesObjectSecond = {
    number_prop: 1,
    string_prop: 'string',
    bool_prop: false,
    null_prop: null,
    undefined_prop: undefined,
    array_prop: [
      {
        obj_prop: {
          deep_prop: 53,
          empty_obj: {},
        },
      },
      [],
    ],
  };
  const variables = [
    convertObjectToTreeview(variablesObjectFirst, 'test'),
    convertObjectToTreeview(variablesObjectSecond, 'test2'),
  ];
  const modal = {
    config: {
      variables,
    },
  };

  const factory = generateRenderer(VariablesHelp, {
    stubs,
    mocks: { $popups },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
        $system: {},
      },
    },
  });

  test('Path success copied after trigger copy button', async () => {
    const copyText = vi.fn().mockResolvedValue();
    const wrapper = factory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
        $copyText: copyText,
      },
    });

    await wrapper.openAllTreeviewNodes(wrapper);

    await selectListTileByIndex(selectTreeviewNodeByIndex(wrapper, 10), 0).trigger('click');

    expect(copyText).toHaveBeenCalledWith('test2.array_prop.[0].obj_prop.deep_prop');
    expect($popups.success).toHaveBeenCalledWith({ text: 'Path copied to clipboard' });
  });

  test('Path error copied after trigger copy button', async () => {
    const copyText = vi.fn().mockRejectedValue();
    const wrapper = factory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
        $copyText: copyText,
      },
    });

    await wrapper.openAllTreeviewNodes(wrapper);

    await selectListTileByIndex(selectTreeviewNodeByIndex(wrapper, 10), 0).trigger('click');

    expect(copyText).toHaveBeenCalledWith('test2.array_prop.[0].obj_prop.deep_prop');
    expect($popups.error).toHaveBeenCalledWith({ text: 'Something went wrong...' });
  });

  test('Object exported after trigger action button', async () => {
    const wrapper = factory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
      },
    });

    await wrapper.openAllTreeviewNodes(wrapper);

    await selectListTileByIndex(selectTreeviewNodeByIndex(wrapper, 7), 0).trigger('click');

    expect(saveJsonFile).toHaveBeenCalledWith(
      variablesObjectSecond.array_prop,
      'array_prop-07/12/2013 18:00:00',
    );
  });
});
