import { omit } from 'lodash';
import flushPromises from 'flush-promises';

export const createSettingsMocks = () => {
  const createWidget = jest.fn();
  const updateWidget = jest.fn();
  const copyWidget = jest.fn();
  const fetchActiveView = jest.fn();
  const fetchUserPreference = jest.fn();
  const getUserPreferenceByWidgetId = jest.fn()
    .mockReturnValue({ content: {} });

  return {
    createWidget,
    updateWidget,
    copyWidget,
    fetchActiveView,
    activeViewModule: {
      name: 'activeView',
      actions: {
        fetch: fetchActiveView,
      },
    },

    widgetModule: {
      name: 'view/widget',
      actions: {
        create: createWidget,
        update: updateWidget,
        copy: copyWidget,
      },
    },

    authModule: {
      name: 'auth',
      getters: {
        currentUserPermissionsById: {},
      },
    },

    userPreferenceModule: {
      name: 'userPreference',
      actions: {
        fetchItem: fetchUserPreference,
      },
      getters: {
        getItemByWidgetId: () => getUserPreferenceByWidgetId,
      },
    },
  };
};

export const getWidgetRequestWithNewProperty = (widget, key, value) => ({
  ...omit(widget, ['_id']),

  [key]: value,
});

export const getWidgetRequestWithNewParametersProperty = (widget, key, value) => ({
  ...omit(widget, ['_id']),

  parameters: {
    ...widget.parameters,

    [key]: value,
  },
});

export const submitWithExpects = async (wrapper, { fetchActiveView, hideSidebar, widgetMethod, expectData }) => {
  const submitButton = wrapper.find('button.v-btn');

  submitButton.trigger('click');

  await flushPromises();

  expect(widgetMethod).toHaveBeenCalledTimes(1);
  expect(widgetMethod).toHaveBeenLastCalledWith(
    expect.any(Object),
    expectData,
    undefined,
  );
  expect(fetchActiveView).toHaveBeenCalledTimes(1);
  expect(hideSidebar).toHaveBeenCalledTimes(1);
};
