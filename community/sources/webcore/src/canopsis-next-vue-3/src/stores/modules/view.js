import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { API_ROUTES } from '@/config';

import request, { useRequestCancelling } from '@/services/request';

export const useViewStore = defineStore('view', () => {
  const items = ref([]);
  const meta = ref({});
  const pending = ref(false);

  const getGroupById = computed(() => id => items.value.find(group => group._id === id));
  const getViewById = computed(() => (id) => {
    let view;

    items.value.some((group) => {
      view = group.views?.find(({ _id }) => _id === id);

      return view;
    });

    return view;
  });
  const getViewTabById = computed(() => (id) => {
    let tab;

    items.value.some((group) => {
      group.views?.some((view) => {
        tab = view.tabs?.find(({ _id }) => _id === id);

        return tab;
      });

      return tab;
    });

    return tab;
  });

  const mergeViewEntities = (view) => {
    const groupId = view?.group?._id;

    if (!groupId) {
      return;
    }

    const groupIndex = items.value.findIndex(group => group._id === groupId);

    if (groupIndex === -1) {
      return;
    }

    const group = items.value[groupIndex];
    const viewIndex = group.views?.findIndex(({ _id }) => _id === view._id) ?? -1;
    const views = [...(group.views ?? [])];

    if (viewIndex === -1) {
      views.push(view);
    } else {
      views[viewIndex] = view;
    }

    items.value[groupIndex] = {
      ...group,
      views,
    };
  };

  const mergeViewTabEntity = (viewTab) => {
    const tabViewId = viewTab?.view?._id ?? viewTab?.view;

    items.value = items.value.map(group => ({
      ...group,
      views: group.views?.map((view) => {
        if (tabViewId && view._id !== tabViewId) {
          return view;
        }

        const tabIndex = view.tabs?.findIndex(({ _id }) => _id === viewTab._id) ?? -1;

        if (tabIndex === -1 && view._id !== tabViewId) {
          return view;
        }

        const tabs = [...(view.tabs ?? [])];

        if (tabIndex === -1) {
          tabs.push(viewTab);
        } else {
          tabs[tabIndex] = viewTab;
        }

        return {
          ...view,
          tabs,
        };
      }),
    }));
  };

  const fetchList = async ({ params } = {}) => {
    try {
      pending.value = true;

      const response = await request.get(API_ROUTES.view.groups, { params });

      items.value = response?.data ?? response ?? [];
      meta.value = response?.meta ?? {};

      return response;
    } finally {
      pending.value = false;
    }
  };

  const fetchListWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.view.groups, { params });

  const create = ({ data } = {}) => request.post(API_ROUTES.view.groups, data);
  const update = ({ id, data } = {}) => request.put(`${API_ROUTES.view.groups}/${id}`, data);
  const remove = ({ id } = {}) => request.delete(`${API_ROUTES.view.groups}/${id}`);

  const createPrivateGroup = ({ data } = {}) => request.post(API_ROUTES.privateView.groups, data);
  const updatePrivateGroup = ({ id, data } = {}) =>
    request.put(`${API_ROUTES.privateView.groups}/${id}`, data);
  const removePrivateGroup = ({ id } = {}) =>
    request.delete(`${API_ROUTES.privateView.groups}/${id}`);

  const createView = ({ data } = {}) => request.post(API_ROUTES.view.list, data);
  const updateView = ({ id, data } = {}) => request.put(`${API_ROUTES.view.list}/${id}`, data);
  const updateWithoutStoreView = ({ id, data } = {}) =>
    request.put(`${API_ROUTES.view.list}/${id}`, data);
  const updateViewPositions = ({ data } = {}) => request.put(API_ROUTES.view.positions, data);
  const removeView = ({ id } = {}) => request.delete(`${API_ROUTES.view.list}/${id}`);
  const copyView = ({ id, data } = {}) => request.post(`${API_ROUTES.view.copy}/${id}`, data);
  const exportViewWithoutStore = ({ data } = {}) => request.post(API_ROUTES.view.export, data);
  const importViewWithoutStore = ({ data } = {}) => request.post(API_ROUTES.view.import, data);
  const createViewTab = ({ data } = {}) => request.post(API_ROUTES.view.tabs, data);
  const updateViewTab = ({ id, data } = {}) => request.put(`${API_ROUTES.view.tabs}/${id}`, data);
  const removeViewTab = ({ id } = {}) => request.delete(`${API_ROUTES.view.tabs}/${id}`);
  const copyViewTab = ({ id, data } = {}) => request.post(`${API_ROUTES.view.tabCopy}/${id}`, data);
  const updateViewTabPositions = ({ data } = {}) => request.put(API_ROUTES.view.tabPositions, data);
  const fetchViewTab = async ({ id, params } = {}) => {
    const viewTab = await request.get(`${API_ROUTES.view.tabs}/${id}`, { params });

    mergeViewTabEntity(viewTab);

    return viewTab;
  };

  const fetchView = ({ id }) =>
    useRequestCancelling(async (source) => {
      const view = await request.get(`${API_ROUTES.view.list}/${id}`, {
        cancelToken: source.token,
      });

      mergeViewEntities(view);

      return view;
    }, `view_${id}`);

  return {
    items,
    meta,
    pending,

    getGroupById,
    getViewById,
    getViewTabById,

    fetchList,
    fetchListWithoutStore,
    create,
    update,
    remove,
    createPrivateGroup,
    updatePrivateGroup,
    removePrivateGroup,
    createView,
    updateView,
    updateWithoutStoreView,
    updateViewPositions,
    removeView,
    copyView,
    exportViewWithoutStore,
    importViewWithoutStore,
    createViewTab,
    updateViewTab,
    removeViewTab,
    copyViewTab,
    updateViewTabPositions,
    fetchViewTab,
    fetchView,
  };
});
