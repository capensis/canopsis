import flushPromises from 'flush-promises';
import Faker from 'faker';
import { omit } from 'lodash';

import { createAuthModule, createMockedStoreModules, createViewModule } from '@unit/utils/store';
import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups, mockRouter } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';

import ClickOutside from '@/services/click-outside';
import {
  CRUD_ACTIONS,
  DEFAULT_PERIODIC_REFRESH,
  MODALS,
  ROUTES_NAMES,
  USERS_PERMISSIONS,
} from '@/constants';

import CreateView from '@/components/modals/view/create-view.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'view-form': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const selectSubmitButton = wrapper => wrapper.find('.v-btn[type="submit"]');
const selectCancelButton = wrapper => wrapper.find('.v-btn[depressed]');
const selectRemoveButton = wrapper => wrapper.find('.v-btn[color="error"]');
const selectViewForm = wrapper => wrapper.find('view-form-stub');

describe('create-view', () => {
  const $modals = mockModals();
  const $popups = mockPopups();
  const $router = mockRouter();

  const fakedView = {
    _id: Faker.datatype.string(),
    title: Faker.datatype.string(),
    description: Faker.datatype.string(),
    enabled: Faker.datatype.boolean(),
    tags: [Faker.datatype.string()],
    group: {
      _id: Faker.datatype.string(),
    },
    periodic_refresh: DEFAULT_PERIODIC_REFRESH,
  };
  const fakedViewWithoutId = omit(fakedView, ['_id']);

  const { authModule, currentUserPermissionsById } = createAuthModule();
  currentUserPermissionsById.mockReturnValue({
    [fakedView._id]: {
      actions: [CRUD_ACTIONS.update, CRUD_ACTIONS.delete],
    },
    [USERS_PERMISSIONS.technical.view]: {
      actions: [
        CRUD_ACTIONS.create,
        CRUD_ACTIONS.update,
        CRUD_ACTIONS.read,
        CRUD_ACTIONS.delete,
      ],
    },
  });

  const {
    viewModule,
    createGroup,
    updateView,
    createView,
    copyView,
    removeView,
  } = createViewModule();
  const store = createMockedStoreModules([
    authModule,
    viewModule,
  ]);

  const factory = generateShallowRenderer(CreateView, {
    stubs,
    attachTo: document.body,
    propsData: {
      modal: {
        config: {},
      },
    },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
        $system: {},
      },
    },
    mocks: {
      $modals,
      $popups,
      $router,
    },
  });
  const snapshotFactory = generateRenderer(CreateView, {
    stubs,
    attachTo: document.body,
    propsData: {
      modal: {
        config: {},
      },
    },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
        $system: {},
      },
    },
    mocks: {
      $modals,
    },
  });

  test('View created after trigger submit button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
      store,
    });

    await flushPromises();

    const newView = {
      ...fakedViewWithoutId,
      group: Faker.datatype.string(),
    };
    const newGroup = {
      _id: Faker.datatype.string(),
      title: newView.title,
    };

    createGroup.mockReturnValueOnce(newGroup);

    selectViewForm(wrapper).vm.$emit('input', newView);
    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(createView).toBeDispatchedWith(
      {
        data: {
          ...newView,
          group: newGroup._id,
        },
      },
    );
    expect($popups.success).toBeCalledWith({ text: 'New view created!' });
    expect($modals.hide).toBeCalled();
  });

  test('Create error popup showed after trigger submit button with error', async () => {
    createView.mockRejectedValueOnce({
      title: 'Title error',
    });
    const wrapper = factory({
      propsData: {
        modal: {
          config: {},
        },
      },
      store,
    });

    await flushPromises();

    selectViewForm(wrapper).vm.$emit('input', fakedViewWithoutId);
    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(createView).toBeCalled();
    expect($popups.error).toBeCalledWith({ text: 'View creation failed...' });
    expect($modals.hide).not.toBeCalled();
  });

  test('View updated after trigger submit button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            view: fakedView,
          },
        },
      },
      store,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(updateView).toBeDispatchedWith(
      {
        data: {
          ...fakedViewWithoutId,
          group: fakedView.group._id,
        },
        id: fakedView._id,
      },
    );
    expect($modals.hide).toBeCalled();
    expect($popups.success).toBeCalledWith({ text: 'View successfully edited!' });
  });

  test('Update error popup showed after trigger submit button with error', async () => {
    updateView.mockRejectedValueOnce({
      description: 'Description error',
    });
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            view: fakedView,
          },
        },
      },
      store,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(updateView).toBeCalled();
    expect($popups.error).toBeCalledWith({ text: 'View edition failed...' });
    expect($modals.hide).not.toBeCalled();
  });

  test('View duplicated after trigger submit button', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            duplicate: true,
            view: fakedViewWithoutId,
          },
        },
      },
      store,
    });

    await flushPromises();

    selectViewForm(wrapper).vm.$emit('input', fakedViewWithoutId);
    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(copyView).toBeDispatchedWith(
      {
        data: {
          ...fakedViewWithoutId,
          group: fakedView.group._id,
        },
      },
    );
    expect($modals.hide).toBeCalled();
    expect($popups.success).toBeCalledWith({ text: 'View successfully duplicated!' });
  });

  test('Duplicate error popup showed after trigger submit button with error', async () => {
    copyView.mockRejectedValueOnce({
      tags: 'Tags error',
    });
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            duplicate: true,
            view: fakedViewWithoutId,
          },
        },
      },
      store,
    });

    await flushPromises();

    selectViewForm(wrapper).vm.$emit('input', fakedViewWithoutId);
    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(copyView).toBeCalled();
    expect($popups.error).toBeCalledWith({ text: 'View duplication failed...' });
    expect($modals.hide).not.toBeCalled();
  });

  test('Action called after trigger submit button with action', async () => {
    const action = jest.fn();
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            view: fakedView,
            action,
          },
        },
      },
      store,
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(action).toBeCalledWith({
      ...fakedViewWithoutId,
      group: fakedView.group._id,
    });
    expect($modals.hide).toBeCalled();
  });

  test('View removed after trigger remove button with action', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            view: fakedView,
          },
        },
      },
      mocks: {
        $route: {},
      },
      store,
    });

    selectRemoveButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(removeView).toBeDispatchedWith({ id: fakedView._id });
    expect($modals.hide).toBeCalled();
  });

  test('Page redirected after trigger remove button with action on the view page', async () => {
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            view: fakedView,
          },
        },
      },
      mocks: {
        $route: {
          name: ROUTES_NAMES.view,
          params: {
            id: fakedView._id,
          },
        },
      },
      store,
    });

    selectRemoveButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(removeView).toBeCalled();
    expect($router.push).toBeCalledWith({ name: ROUTES_NAMES.home });
    expect($modals.hide).toBeCalled();
  });

  test('Remove error popup showed after trigger remove button with error', async () => {
    removeView.mockRejectedValueOnce({});
    const wrapper = factory({
      propsData: {
        modal: {
          config: {
            view: fakedView,
          },
        },
      },
      mocks: {
        $route: {},
      },
      store,
    });

    selectRemoveButton(wrapper).trigger('click');

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(removeView).toBeDispatchedWith({ id: fakedView._id });
    expect($popups.error).toBeCalledWith({
      text: 'View deletion failed...',
    });
    expect($modals.hide).not.toBeCalled();
  });

  test('Cancel action called after trigger cancel button', () => {
    const wrapper = factory({
      store,
    });

    selectCancelButton(wrapper).trigger('click');

    expect($modals.hide).toBeCalled();
  });

  test('Renders `create-view` with empty modal', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        modal: {
          config: {
            view: {},
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `create-view` with title', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        modal: {
          config: {
            title: 'Create view custom title',
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `create-view` without rights', async () => {
    currentUserPermissionsById.mockReturnValueOnce({});

    const wrapper = snapshotFactory({
      propsData: {
        modal: {
          config: {},
        },
      },
      store: createMockedStoreModules([
        authModule,
        viewModule,
      ]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});