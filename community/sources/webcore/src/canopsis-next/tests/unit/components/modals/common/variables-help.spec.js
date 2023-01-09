import { createVueInstance, generateRenderer } from '@unit/utils/vue';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { convertObjectToTreeview } from '@/helpers/treeview';
import { saveJsonFile } from '@/helpers/file/files';
import ClickOutside from '@/services/click-outside';

import VariablesHelp from '@/components/modals/common/variables-help.vue';

jest.mock('@/helpers/file/files', () => ({
  saveJsonFile: jest.fn(),
}));

const localVue = createVueInstance();

const snapshotStubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'patterns-form': true,
  'c-action-btn': true,
  'c-copy-btn': true,
  'c-ellipsis': true,
};

const selectTreeviewNodes = wrapper => wrapper.findAll('.v-treeview-node');
const selectTreeviewNodeByIndex = (wrapper, index) => selectTreeviewNodes(wrapper).at(index);
const selectTreeviewChildren = wrapper => wrapper.find('.v-treeview-node__children');
const selectCopyButton = wrapper => wrapper.find('c-copy-btn-stub');
const selectActionButton = wrapper => wrapper.find('c-action-btn-stub');
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
    localVue,
    stubs: snapshotStubs,
    mocks: { $popups },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  test('Path copied after trigger copy button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        modal,
      },
      mocks: {
        $modals,
      },
    });

    await openAllNodes(wrapper);

    const copyButton = selectCopyButton(selectTreeviewNodeByIndex(wrapper, 10));

    expect(copyButton.vm.$attrs.value).toEqual('test2.array_prop.[0].obj_prop.deep_prop');

    await copyButton.vm.$emit('success');
    expect($popups.success).toBeCalledWith({ text: 'Path copied to clipboard' });

    await copyButton.vm.$emit('error');
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

    const actionButton = selectActionButton(selectTreeviewNodeByIndex(wrapper, 7));

    await actionButton.vm.$emit('click');

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
});
