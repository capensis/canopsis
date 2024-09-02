import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { fakeAlarm } from '@unit/data/alarm';
import { convertObjectToTreeview } from '@/helpers/treeview';
import { saveJsonFile } from '@/helpers/file/files';
import ClickOutside from '@/services/click-outside';

import VariablesHelp from '@/components/modals/common/variables-help.vue';

jest.mock('@/helpers/file/files', () => ({
  saveJsonFile: jest.fn(),
}));

const snapshotStubs = {
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
const selectTreeviewChildren = wrapper => wrapper.find('.v-treeview-node__children');
const selectListTileByIndex = (wrapper, index) => wrapper.findAll('.v-list__tile').at(index);
const selectTreeviewNodeToggle = wrapper => wrapper
  .find('.v-treeview-node__toggle:not(.v-treeview-node__toggle--open)');

const openAllNodes = async (wrapper) => {
  const { wrappers } = selectTreeviewNodes(wrapper);

  if (!wrappers.length) {
    return;
  }

  await Promise.all(
    wrappers.map(
      (node) => {
        const toggle = selectTreeviewNodeToggle(node);

        if (!toggle.element) {
          return Promise.resolve();
        }

        return toggle.trigger('click');
      },
    ),
  );
  await flushPromises();
  await selectTreeviewNodes(wrapper).wrappers.reduce(
    (acc, node) => {
      const children = selectTreeviewChildren(node);

      if (children.element) {
        return acc.then(
          () => openAllNodes(selectTreeviewChildren(node)),
        );
      }

      return acc;
    },
    Promise.resolve(),
  );
};

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
    array_prop: [{
      obj_prop: {
        deep_prop: 53,
        empty_obj: {},
      },
    }, []],
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

  const snapshotFactory = generateRenderer(VariablesHelp, {
    stubs: snapshotStubs,
    mocks: { $popups },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
        $system: {},
      },
    },
  });

  test('Path success copied after trigger copy button', async () => {
    const copyText = jest.fn().mockResolvedValue();
    const wrapper = snapshotFactory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
        $copyText: copyText,
      },
    });

    await openAllNodes(wrapper);

    await selectListTileByIndex(selectTreeviewNodeByIndex(wrapper, 10), 0).trigger('click');

    expect(copyText).toBeCalledWith('test2.array_prop.[0].obj_prop.deep_prop');
    expect($popups.success).toBeCalledWith({ text: 'Path copied to clipboard' });
  });

  test('Path error copied after trigger copy button', async () => {
    const copyText = jest.fn().mockRejectedValue();
    const wrapper = snapshotFactory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
        $copyText: copyText,
      },
    });

    await openAllNodes(wrapper);

    await selectListTileByIndex(selectTreeviewNodeByIndex(wrapper, 10), 0).trigger('click');

    expect(copyText).toBeCalledWith('test2.array_prop.[0].obj_prop.deep_prop');
    expect($popups.error).toBeCalledWith({ text: 'Something went wrong...' });
  });

  test('Object exported after trigger action button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
      },
    });

    await openAllNodes(wrapper);

    await selectListTileByIndex(selectTreeviewNodeByIndex(wrapper, 7), 0).trigger('click');

    expect(saveJsonFile).toBeCalledWith(
      variablesObjectSecond.array_prop,
      'array_prop-07/12/2013 18:00:00',
    );
  });

  test('Renders `variables-help` with empty modal', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `variables-help` with all parameters', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
      },
    });

    await openAllNodes(wrapper);

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `variables-help` with original flag', () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {
            variables: [
              { ...convertObjectToTreeview(variablesObjectFirst, 'alarm'), original: fakeAlarm() },
              convertObjectToTreeview(variablesObjectSecond, 'test2'),
            ],
          },
        },
      },
      mocks: {
        $modals,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
